package repositories_test

import (
	"context"
	"fmt"
	pkg "github.com/WildEgor/e-shop-support-bot/internal"
	"github.com/WildEgor/e-shop-support-bot/internal/configs"
	postgres2 "github.com/WildEgor/e-shop-support-bot/internal/db/postgres"
	"github.com/WildEgor/e-shop-support-bot/internal/db/redis"
	"github.com/WildEgor/e-shop-support-bot/internal/repositories"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	redis2 "github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"os"
	"strconv"
	"testing"
	"time"
)

var (
	RedisContainer *redis2.RedisContainer
	RedisConn      *redis.RedisConnection
)

var (
	PostgresContainer *postgres.PostgresContainer
	PostgresConn      *postgres2.PostgresConnection
)

var (
	TopicRepository     repositories.ITopicsRepository
	GroupRepository     repositories.IGroupRepository
	UserStateRepository repositories.IUserStateRepository
)

func TestMain(m *testing.M) {
	if os.Getenv("SKIP_E2E_TESTS") != "" {
		return
	}

	if err := setup(); err != nil {
		os.Exit(1)
	}

	exitCode := m.Run()

	if err := tearDown(); err != nil {
		os.Exit(1)
	}

	os.Exit(exitCode)
}

func setup() error {
	ctx := context.Background()

	rc, err := redis2.RunContainer(ctx,
		testcontainers.WithImage("docker.io/redis:7"),
		redis2.WithSnapshotting(10, 1),
		redis2.WithLogLevel(redis2.LogLevelVerbose),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	RedisContainer = rc

	port, err := rc.MappedPort(ctx, "6379")
	if err != nil {
		return err
	}

	redisConfig := &configs.RedisConfig{
		URI: fmt.Sprintf("redis://127.0.0.1:%s/0", port.Port()),
	}
	RedisConn = redis.NewRedisConnection(redisConfig)

	pc, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("docker.io/postgres:14-alpine"),
		postgres.WithDatabase("ds_db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	PostgresContainer = pc

	host, err := pc.Host(ctx)
	if err != nil {
		log.Fatalf("Could not get PostgreSQL container host: %s", err)
	}
	port, err = pc.MappedPort(ctx, "5432")
	if err != nil {
		log.Fatalf("Could not get PostgreSQL container port: %s", err)
	}

	parseInt, err := strconv.ParseInt(port.Port(), 10, 64)
	if err != nil {
		return err
	}

	postgresConfig := &configs.PostgresConfig{
		Host:     host,
		Port:     uint16(parseInt),
		User:     "postgres",
		Password: "postgres",
		Name:     "ds_db",
	}

	PostgresConn = postgres2.NewPostgresConnection(postgresConfig)

	TopicRepository = repositories.NewTopicsRepository(RedisConn, PostgresConn)
	GroupRepository = repositories.NewGroupRepository(RedisConn)
	UserStateRepository = repositories.NewUserStateRepository(RedisConn)

	pkg.RunMigrate(postgresConfig.MigrationURI())

	return nil
}

func tearDown() error {
	if err := RedisContainer.Terminate(context.Background()); err != nil {
		log.Fatalf("failed to terminate container: %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := PostgresContainer.Terminate(ctx); err != nil {
		log.Fatalf("failed to terminate container: %s", err)
	}

	RedisConn.Close()
	PostgresConn.Close()

	return nil
}
