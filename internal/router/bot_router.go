package router

import (
	"github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/adapters/telegram"
	"github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/configs"
	accept_callback_handler "github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/handlers/accept_callback"
	break_action_handler "github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/handlers/break_action"
	decline_callback_handler "github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/handlers/decline_callback"
	edit_message_handler "github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/handlers/edit_message"
	new_message_handler "github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/handlers/new_message"
	no_right_callback_handler "github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/handlers/not_right_callback"
	rating_callback_handler "github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/handlers/rating_callback"
	start_action_handler "github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/handlers/start_action"
	middlewares2 "github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/middlewares/auth"
	middlewares "github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/middlewares/group"
	"github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/models"
	"github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/repositories"
)

type BotRouter struct {
	tcfg *configs.TelegramConfig

	tl  *telegram.TelegramListener
	sah *start_action_handler.StartActionHandler
	bah *break_action_handler.BreakActionHandler
	emh *edit_message_handler.EditMessageHandler
	ah  *accept_callback_handler.AcceptCallbackHandler
	dh  *decline_callback_handler.DeclineCallbackHandler
	nrh *no_right_callback_handler.NoRightCallbackHandler
	nmh *new_message_handler.NewMessageHandler
	rh  *rating_callback_handler.RatingCallbackHandler
	uor repositories.IUserStateRepository
	gr  repositories.IGroupRepository
}

func NewBotRouter(
	tcfg *configs.TelegramConfig,
	tl *telegram.TelegramListener,
	sah *start_action_handler.StartActionHandler,
	bah *break_action_handler.BreakActionHandler,
	emh *edit_message_handler.EditMessageHandler,
	ah *accept_callback_handler.AcceptCallbackHandler,
	dh *decline_callback_handler.DeclineCallbackHandler,
	nrh *no_right_callback_handler.NoRightCallbackHandler,
	nmh *new_message_handler.NewMessageHandler,
	rh *rating_callback_handler.RatingCallbackHandler,
	uor repositories.IUserStateRepository,
	gr repositories.IGroupRepository,
) *BotRouter {
	return &BotRouter{
		tcfg,
		tl,
		sah,
		bah,
		emh,
		ah,
		dh,
		nrh,
		nmh,
		rh,
		uor,
		gr,
	}
}

func (r *BotRouter) SetupBotRouter() {
	gmw := middlewares.NewExtractGroupMiddleware(r.uor, r.gr, r.tcfg)
	amw := middlewares2.NewAuthMiddleware()

	r.tl.Use(gmw.Next)
	r.tl.Use(amw.Next)

	// Text handlers
	// HINT: catch all messages, except commands (ex. /start)
	r.tl.HandleMessage("^[^/]{1,255}$", r.nmh.Handle)

	// Text edit handler
	r.tl.HandleEditedMessage(r.emh.Handle)

	// Commands handlers
	r.tl.HandleMessage(models.StartBotCommand, r.sah.Handle)
	r.tl.HandleMessage(models.BreakTopicCommand, r.bah.Handle)

	// Callbacks handlers
	r.tl.RegisterCallbackHandler(models.AcceptHelp, r.ah.Handle)
	r.tl.RegisterCallbackHandler(models.DeclineHelp, r.dh.Handle)
	r.tl.RegisterCallbackHandler(models.NoRight, r.nrh.Handle)
	r.tl.RegisterCallbackHandler(models.Rating, r.rh.Handle)
}
