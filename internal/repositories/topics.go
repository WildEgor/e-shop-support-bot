package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	logModels "github.com/WildEgor/e-shop-gopack/pkg/libs/logger/models"
	"github.com/WildEgor/e-shop-support-bot/internal/adapters/publisher"
	"github.com/WildEgor/e-shop-support-bot/internal/configs"
	"github.com/WildEgor/e-shop-support-bot/internal/db/postgres"
	"github.com/WildEgor/e-shop-support-bot/internal/db/redis"
	"github.com/WildEgor/e-shop-support-bot/internal/mappers"
	"github.com/WildEgor/e-shop-support-bot/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"log/slog"
	"time"
)

type ITopicsRepository interface {
	CreateUniqueTopic(attrs *models.CreateTopicAttrs) (*models.Topic, error)
	AssignTopicSupport(topic *models.Topic) (*models.Topic, error)
	LeaveFeedback(feedback *models.CreateTopicFeedbackAttrs) (*models.Topic, error)
	CloseTopic(id int64) error
	FindById(id int64) (*models.Topic, error)
	FindFirstAuthorActiveTopic(id int64) (*models.Topic, error)
	GetTopicRoomBySenderId(ctx context.Context, senderId int64) (*models.TopicRoom, error)
	SaveTopicRoom(ctx context.Context, topic *models.Topic) error
	RemoveTopicRoom(ctx context.Context, topic *models.Topic) error
}

// TopicRepository represent access to user topics
type TopicRepository struct {
	redis           *redis.RedisConnection
	postgres        *postgres.PostgresConnection
	publisher       publisher.IEventPublisher
	publisherConfig *configs.PublisherConfig
}

func NewTopicsRepository(
	redis *redis.RedisConnection,
	postgres *postgres.PostgresConnection,
	publisher publisher.IEventPublisher,
	publisherConfig *configs.PublisherConfig,
) *TopicRepository {
	return &TopicRepository{
		redis,
		postgres,
		publisher,
		publisherConfig,
	}
}

func (r *TopicRepository) FindById(id int64) (*models.Topic, error) {
	query := fmt.Sprintf(`SELECT * FROM %s WHERE id = @id;`, models.TopicsTable)

	rows, err := r.postgres.DB.Query(context.TODO(), query, pgx.NamedArgs{
		"id": id,
	})
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	collectRows, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[models.TopicTable])
	if err != nil {
		return nil, err
	}

	if len(collectRows) == 0 {
		return nil, err
	}

	return mappers.FromTopicsTableToModel(&collectRows[0]), nil
}

func (r *TopicRepository) AssignTopicSupport(topic *models.Topic) (*models.Topic, error) {
	query := fmt.Sprintf(`UPDATE %s SET support_tid = @support_tid, support_tun = @support_tun WHERE id = @id;`, models.TopicsTable)

	result, err := r.postgres.DB.Exec(context.TODO(), query, pgx.NamedArgs{
		"support_tid": topic.Support.TelegramId,
		"support_tun": topic.Support.TelegramUsername,
		"id":          topic.Id,
	})

	if err != nil {
		return nil, err
	}

	if result.RowsAffected() == 0 {
		return nil, pgx.ErrNoRows
	}

	return topic, nil
}

func (r *TopicRepository) SaveTopicRoom(ctx context.Context, topic *models.Topic) error {
	var payloads = []*models.TopicRoom{
		{
			Id:   topic.Id,
			From: topic.Creator.TelegramId,
			To:   topic.Support.TelegramId,
		},
		{
			Id:        topic.Id,
			From:      topic.Support.TelegramId,
			To:        topic.Creator.TelegramId,
			IsSupport: true,
		},
	}

	pipeliner := r.redis.Client.WithContext(ctx).Pipeline()

	for i := 0; i < len(payloads); i++ {
		data, err := json.Marshal(payloads[i])
		if err != nil {
			// TODO: log
			continue
		}

		pipeliner.Set(r.prefixRoomsKey(payloads[i].From), data, 0)
	}

	_, err := pipeliner.Exec()
	if err != nil {
		return err
	}

	return nil
}

func (r *TopicRepository) GetTopicRoomBySenderId(ctx context.Context, senderId int64) (*models.TopicRoom, error) {
	res, err := r.redis.Client.WithContext(ctx).Get(r.prefixRoomsKey(senderId)).Result()
	if err != nil {
		return nil, err
	}

	var result *models.TopicRoom
	if err := json.Unmarshal([]byte(res), &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *TopicRepository) RemoveTopicRoom(ctx context.Context, topic *models.Topic) error {
	pipeliner := r.redis.Client.WithContext(ctx).Pipeline()

	pipeliner.Del(r.prefixRoomsKey(topic.Creator.TelegramId))
	pipeliner.Del(r.prefixRoomsKey(topic.Support.TelegramId))

	_, err := pipeliner.Exec()
	if err != nil {
		return err
	}

	return nil
}

func (r *TopicRepository) FindFirstAuthorActiveTopic(id int64) (*models.Topic, error) {
	query := fmt.Sprintf(`SELECT * FROM %s WHERE author_tid = @author_tid AND status = 1 ORDER BY created_at LIMIT 1`, models.TopicsTable)

	rows, err := r.postgres.DB.Query(context.TODO(), query, pgx.NamedArgs{
		"author_tid": id,
	})
	if err != nil {
		return nil, err
	}

	collectRows, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[models.TopicTable])
	if err != nil {
		return nil, err
	}

	if len(collectRows) == 0 {
		return nil, pgx.ErrNoRows
	}

	return mappers.FromTopicsTableToModel(&collectRows[0]), nil
}

func (r *TopicRepository) CreateUniqueTopic(attrs *models.CreateTopicAttrs) (*models.Topic, error) {
	//query := fmt.Sprintf("INSERT INTO %s (author_tid, author_tun, msg_tid, content, feedback_id) VALUES (@author_tid, @author_tun, @msg_tid, @content, null) ON CONFLICT ON CONSTRAINT unique_active_topic_per_author DO UPDATE SET status = 0 RETURNING id;", models.TopicsTable)
	query := fmt.Sprintf("INSERT INTO %s (author_tid, author_tun, msg_tid, content, feedback_id) VALUES (@author_tid, @author_tun, @msg_tid, @content, null) RETURNING id;", models.TopicsTable)

	entity := &models.Topic{
		Message:   attrs.Message,
		Creator:   attrs.Creator,
		Status:    models.ActiveTopic,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := r.postgres.DB.QueryRow(context.TODO(), query, pgx.NamedArgs{
		"author_tid": entity.Creator.TelegramId,
		"author_tun": entity.Creator.TelegramUsername,
		"msg_tid":    entity.Message.TelegramMessageId,
		"content":    entity.Message.Question,
	}).Scan(&entity.Id)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *TopicRepository) CloseTopic(id int64) error {
	query := fmt.Sprintf("UPDATE %s SET status = @status WHERE id = @id", models.TopicsTable)

	result, err := r.postgres.DB.Exec(context.TODO(), query, pgx.NamedArgs{
		"id":     id,
		"status": models.ClosedTopic,
	})

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *TopicRepository) LeaveFeedback(feedback *models.CreateTopicFeedbackAttrs) (*models.Topic, error) {
	feedback.FeedbackId = uuid.New().String()
	feedback.UpdatedAt = time.Now()

	query := fmt.Sprintf("UPDATE %s SET feedback_id = @feedback_id, updated_at = @updated_at WHERE id = @id", models.TopicsTable)

	result, err := r.postgres.DB.Exec(context.TODO(), query, pgx.NamedArgs{
		"id":          feedback.Topic.Id,
		"feedback_id": feedback.FeedbackId,
		"updated_at":  feedback.UpdatedAt,
	})
	if err != nil {
		return nil, err
	}

	err = r.publisher.Publish(context.TODO(), r.publisherConfig.Topic, &publisher.Event{
		Data: &publisher.TopicFeedbackEvent{
			Pattern: "topic_feedbacks",
			Data: publisher.TopicFeedbackEventData{
				ID:                      feedback.FeedbackId,
				TopicID:                 feedback.Topic.Id,
				SupportTelegramID:       feedback.Topic.Support.TelegramId,
				SupportTelegramUsername: feedback.Topic.Support.TelegramUsername,
				SupportUserID:           "",
				CreatorTelegramID:       feedback.Topic.Creator.TelegramId,
				CreatorTelegramUsername: feedback.Topic.Creator.TelegramUsername,
				Rating:                  feedback.Rating,
			},
		},
	})
	if err != nil {
		slog.Error("", logModels.LogEntryAttr(&logModels.LogEntry{
			Err: err,
		}))
		return feedback.Topic, err
	}

	if result.RowsAffected() == 0 {
		return nil, pgx.ErrNoRows
	}

	return feedback.Topic, nil
}

func (r *TopicRepository) prefixRoomsKey(id int64) string {
	return fmt.Sprintf("rooms:%d", id)
}
