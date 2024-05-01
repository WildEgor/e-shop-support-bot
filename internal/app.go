package pkg

import (
	"context"
	"fmt"
	"github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/adapters"
	"github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/adapters/telegram"
	"github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/configs"
	eh "github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/handlers/errors"
	nfm "github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/middlewares/not_found"
	"github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/router"
	slogger "github.com/WildEgor/e-shop-gopack/pkg/libs/logger/handlers"
	"github.com/WildEgor/e-shop-gopack/pkg/libs/logger/models"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/template/html/v2"
	"github.com/google/wire"
	"log/slog"
)

var AppSet = wire.NewSet(
	NewApp,
	configs.ConfigsSet,
	router.RouterSet,
	adapters.AdaptersSet,
)

// Server represents the main server configuration.
type Server struct {
	App       *fiber.App
	Bot       *telegram.TelegramListener
	AppConfig *configs.AppConfig
}

func (srv *Server) Run(ctx context.Context) {
	go func() {
		slog.Debug("server is listening")
		if err := srv.App.Listen(fmt.Sprintf(":%s", srv.AppConfig.Port)); err != nil {
			panic(err)
		}
	}()

	srv.Bot.ListenUpdates(ctx)
}

func (srv *Server) Shutdown() {
	slog.Debug("shutdown bot")
	srv.Bot.Stop()

	slog.Debug("shutdown service")
	if err := srv.App.Shutdown(); err != nil {
		slog.Error("unable to shutdown server.", models.LogEntryAttr(&models.LogEntry{
			Err: err,
		}))
	}
}

func NewApp(
	ac *configs.AppConfig,
	eh *eh.ErrorsHandler,
	prr *router.PrivateRouter,
	pbr *router.PublicRouter,
	sr *router.SwaggerRouter,
	bot *telegram.TelegramListener,
	pc *configs.PostgresConfig,
) *Server {
	logger := slogger.NewLogger(
		slogger.WithOrganization(ac.Name),
	)
	if ac.IsProduction() {
		logger = slogger.NewLogger(
			slogger.WithOrganization(ac.Name),
			slogger.WithLevel("info"),
			slogger.WithOutput("stdout"),
			slogger.WithFormat("json"),
		)
	}
	slog.SetDefault(logger)

	app := fiber.New(fiber.Config{
		ErrorHandler: eh.Handle,
		Views:        html.New("./views", ".html"),
	})

	app.Use(cors.New(cors.Config{
		AllowHeaders: "Origin, Content-Type, Accept, Content-Length, Accept-Language, Accept-Encoding, Connection, Access-Control-Allow-Origin",
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))
	app.Use(recover.New())

	prr.Setup(app)
	pbr.Setup(app)
	sr.Setup(app)

	// 404 handler
	app.Use(nfm.NewNotFound())

	if !ac.IsProduction() {
		RunMigrate(pc.MigrationURI())
	}

	return &Server{
		App:       app,
		Bot:       bot,
		AppConfig: ac,
	}
}
