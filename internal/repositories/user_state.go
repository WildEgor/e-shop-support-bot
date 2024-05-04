package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/WildEgor/e-shop-support-bot/internal/db/redis"
	"github.com/WildEgor/e-shop-support-bot/internal/models"
	rediss "github.com/go-redis/redis"
	"time"
)

type IUserStateRepository interface {
	// AddUserToQueue add user to scored (by created_at timestamp) queue. Using for awaited users
	AddUserToQueue(ctx context.Context, id int64) error
	// DeleteUserFromQueue remove from queue
	DeleteUserFromQueue(ctx context.Context, id int64) error
	// AddUserMessageToBuffer save user message to temp buffer for future process
	AddUserMessageToBuffer(ctx context.Context, msg *models.UserMessage) error
	// GetUserMessagesFromBuffer get first 100 messages from buffer
	GetUserMessagesFromBuffer(ctx context.Context, id int64) []*models.UserMessage
	// CleanUserMessagesFromBuffer delete all buffered messages
	CleanUserMessagesFromBuffer(ctx context.Context, id int64) error
	// SaveUserOptions save user lang and chat id
	SaveUserOptions(ctx context.Context, data *models.UserOptions) error
	// GetUserOptions get user lang and chat id
	GetUserOptions(ctx context.Context, id int64) (*models.UserOptions, error)
	// CheckUserState get user state - default/queue/room/unknown
	CheckUserState(ctx context.Context, id int64) string
}

// UserStateRepository represents user current state in chatting
type UserStateRepository struct {
	redis *redis.RedisConnection
}

func NewUserStateRepository(
	redis *redis.RedisConnection,
) *UserStateRepository {
	return &UserStateRepository{
		redis,
	}
}

func (r *UserStateRepository) AddUserToQueue(ctx context.Context, id int64) error {
	if err := r.redis.Client.WithContext(ctx).ZAdd(r.prefixUserQueueKey(), rediss.Z{
		Member: fmt.Sprintf("%d", id),
		Score:  float64(time.Now().Unix()),
	}).Err(); err != nil {
		return err
	}

	return nil
}

// DeleteUserFromQueue remove user from queue. Remove from queue when user processed
func (r *UserStateRepository) DeleteUserFromQueue(ctx context.Context, id int64) error {
	if err := r.redis.Client.WithContext(ctx).ZRem(r.prefixUserQueueKey(), fmt.Sprintf("%d", id)).Err(); err != nil {
		return err
	}

	return nil
}

func (r *UserStateRepository) AddUserMessageToBuffer(ctx context.Context, msg *models.UserMessage) error {
	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = r.redis.Client.WithContext(ctx).RPush(r.prefixUserMessagesBufferKey(msg.TelegramUserId), payload).Err()
	if err != nil {
		return err
	}

	return nil
}

// GetUserState get user state. Check if user in queue or in rooms state or first time here
func (r *UserStateRepository) CheckUserState(ctx context.Context, id int64) string {
	userId := fmt.Sprintf("%d", id)
	res, err := r.redis.Client.WithContext(ctx).ZScore(r.prefixUserQueueKey(), userId).Result()
	if err != nil && !errors.Is(err, rediss.Nil) {
		return models.UnknownState
	}

	if res != 0 {
		return models.QueueState
	}

	inRoom, err := r.redis.Client.Get(r.prefixRoomsKey(id)).Result()
	if err != nil && !errors.Is(err, rediss.Nil) {
		return models.UnknownState
	}

	if len(inRoom) != 0 {
		return models.RoomState
	}

	return models.DefaultState
}

// CleanUserMessagesFromBuffer cleanup
func (r *UserStateRepository) CleanUserMessagesFromBuffer(ctx context.Context, id int64) error {

	r.redis.Client.WithContext(ctx).Del(r.prefixUserMessagesBufferKey(id))

	return nil
}

// GetUserMessagesFromBuffer return first 100 messages from buffer
func (r *UserStateRepository) GetUserMessagesFromBuffer(ctx context.Context, id int64) []*models.UserMessage {
	val := make([]*models.UserMessage, 0)

	// HINT: limit to 100 messages
	r.redis.Client.WithContext(ctx).LTrim(r.prefixUserMessagesBufferKey(id), 0, 99)
	res, err := r.redis.Client.WithContext(ctx).LRange(r.prefixUserMessagesBufferKey(id), 0, 99).Result()
	if err != nil {
		return val
	}

	for _, msg := range res {
		var result *models.UserMessage
		if err := json.Unmarshal([]byte(msg), &result); err != nil {
			continue
		}

		val = append(val, result)
	}

	return val
}

// GetUserOptions return saved user opts (ex. language). Return default
func (r *UserStateRepository) GetUserOptions(ctx context.Context, id int64) (result *models.UserOptions, err error) {
	val, err := r.redis.Client.WithContext(ctx).Get(r.prefixUserStateKey(id)).Result()
	if err != nil {
		return result, err
	}

	if err := json.Unmarshal([]byte(val), &result); err != nil {
		return result, err
	}

	return
}

// SaveUserOptions save user opts (ex. language)
func (r *UserStateRepository) SaveUserOptions(ctx context.Context, data *models.UserOptions) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = r.redis.Client.WithContext(ctx).Set(r.prefixUserStateKey(data.TelegramId), payload, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *UserStateRepository) prefixUserStateKey(id int64) string {
	return fmt.Sprintf("users:%d:state", id)
}

func (r *UserStateRepository) prefixUserQueueKey() string {
	return "users_queue"
}

func (r *UserStateRepository) prefixUserMessagesBufferKey(id int64) string {
	return fmt.Sprintf("buf_messages:%d", id)
}

func (r *UserStateRepository) prefixRoomsKey(id int64) string {
	return fmt.Sprintf("rooms:%d", id)
}
