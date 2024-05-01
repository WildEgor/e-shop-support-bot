package redis

import (
	"github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/configs"
	"github.com/WildEgor/e-shop-gopack/pkg/libs/logger/models"
	"github.com/go-redis/redis"
	"log/slog"
)

// RedisConnection holds db conn
type RedisConnection struct {
	Client *redis.Client
	cfg    *configs.RedisConfig
}

func NewRedisConnection(
	redisConfig *configs.RedisConfig,
) *RedisConnection {
	conn := &RedisConnection{
		nil,
		redisConfig,
	}

	conn.Connect()

	return conn
}

func (rc *RedisConnection) Connect() {
	opt, err := redis.ParseURL(rc.cfg.URI)
	if err != nil {
		slog.Error("fail parse url", models.LogEntryAttr(&models.LogEntry{
			Err: err,
		}))
		panic(err)
	}

	rc.Client = redis.NewClient(opt)

	if _, err := rc.Client.Ping().Result(); err != nil {
		slog.Error("fail connect to redis ", err)
		panic(err)
	}

	slog.Info("success connect to redis")
}

func (rc *RedisConnection) Close() {
	if rc.Client == nil {
		return
	}

	if err := rc.Client.Close(); err != nil {
		slog.Error("fail disconnect redis", models.LogEntryAttr(&models.LogEntry{
			Err: err,
		}))
		return
	}

	slog.Info("connection to redis closed success")
}
