package db

import (
	"github.com/WildEgor/e-shop-support-bot/internal/db/postgres"
	"github.com/WildEgor/e-shop-support-bot/internal/db/redis"
	"github.com/google/wire"
)

var DbSet = wire.NewSet(
	redis.NewRedisConnection,
	postgres.NewPostgresConnection,
)
