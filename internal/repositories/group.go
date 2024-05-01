package repositories

import (
	"context"
	"github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/db/redis"
	"strconv"
)

type IGroupRepository interface {
	SaveGroupId(ctx context.Context, id int64) error
	GetGroupId(ctx context.Context) (int64, error)
}

// GroupRepository can save chat id
type GroupRepository struct {
	redis *redis.RedisConnection
}

func NewGroupRepository(
	redis *redis.RedisConnection,
) *GroupRepository {
	return &GroupRepository{
		redis,
	}
}

// SaveGroupId Save chat id where bot added (probably, support group)
func (r *GroupRepository) SaveGroupId(ctx context.Context, id int64) error {
	err := r.redis.Client.WithContext(ctx).Set(r.cacheKey(), id, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

// GetGroupId Return support group id
func (r *GroupRepository) GetGroupId(ctx context.Context) (int64, error) {
	val, err := r.redis.Client.WithContext(ctx).Get(r.cacheKey()).Result()
	if err != nil {
		return 0, err
	}

	v, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, err
	}

	return v, nil
}

// cacheKey
func (r *GroupRepository) cacheKey() string {
	return "support_chat"
}
