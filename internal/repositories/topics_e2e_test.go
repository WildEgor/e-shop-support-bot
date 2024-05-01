package repositories_test

import (
	"context"
	"github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTopicRepository_CreateUniqueTopic(t *testing.T) {
	attrs := &models.CreateTopicAttrs{
		Message: models.TopicMessage{
			TelegramMessageId: 1,
			Question:          "test1",
		},
		Creator: models.TopicCreator{
			TelegramId:       2,
			TelegramUsername: "test2",
		},
	}

	topic, err := TopicRepository.CreateUniqueTopic(attrs)

	assert.Nil(t, err)
	assert.Equal(t, topic.Creator.TelegramId, attrs.Creator.TelegramId)
}

func TestTopicRepository_FindById(t *testing.T) {
	attrs := &models.CreateTopicAttrs{
		Message: models.TopicMessage{
			TelegramMessageId: 1,
			Question:          "test1",
		},
		Creator: models.TopicCreator{
			TelegramId:       3,
			TelegramUsername: "test2",
		},
	}

	topic, _ := TopicRepository.CreateUniqueTopic(attrs)

	existed, err := TopicRepository.FindById(topic.Id)

	assert.Nil(t, err)
	assert.NotNil(t, existed)
	assert.Equal(t, topic.Creator.TelegramId, existed.Creator.TelegramId)
}

func TestTopicRepository_AssignTopicSupport(t *testing.T) {
	attrs := &models.CreateTopicAttrs{
		Message: models.TopicMessage{
			TelegramMessageId: 1,
			Question:          "test1",
		},
		Creator: models.TopicCreator{
			TelegramId:       4,
			TelegramUsername: "test2",
		},
	}

	existedTopic, _ := TopicRepository.CreateUniqueTopic(attrs)

	existedTopic.Support = models.TopicSupport{
		TelegramId:       1,
		TelegramUsername: "test",
		TelegramChatId:   777,
	}

	support, err := TopicRepository.AssignTopicSupport(existedTopic)

	assert.Nil(t, err)
	assert.Equal(t, support.Id, existedTopic.Id)
}

func TestTopicRepository_CloseTopic(t *testing.T) {
	attrs := &models.CreateTopicAttrs{
		Message: models.TopicMessage{
			TelegramMessageId: 1,
			Question:          "test1",
		},
		Creator: models.TopicCreator{
			TelegramId:       5,
			TelegramUsername: "test2",
		},
	}

	existedTopic, _ := TopicRepository.CreateUniqueTopic(attrs)

	err := TopicRepository.CloseTopic(existedTopic.Id)
	assert.Nil(t, err)
}

func TestTopicRepository_FindFirstAuthorActiveTopic(t *testing.T) {
	attrs := &models.CreateTopicAttrs{
		Message: models.TopicMessage{
			TelegramMessageId: 1,
			Question:          "test1",
		},
		Creator: models.TopicCreator{
			TelegramId:       6,
			TelegramUsername: "test2",
		},
	}

	existedTopic, _ := TopicRepository.CreateUniqueTopic(attrs)

	topic, err := TopicRepository.FindFirstAuthorActiveTopic(existedTopic.Creator.TelegramId)
	assert.Nil(t, err)
	assert.Equal(t, topic.Id, existedTopic.Id)
}

func TestTopicRepository_LeaveFeedback(t *testing.T) {
	attrs := &models.CreateTopicAttrs{
		Message: models.TopicMessage{
			TelegramMessageId: 1,
			Question:          "test1",
		},
		Creator: models.TopicCreator{
			TelegramId:       7,
			TelegramUsername: "test2",
		},
	}

	existedTopic, _ := TopicRepository.CreateUniqueTopic(attrs)

	topic, err := TopicRepository.LeaveFeedback(&models.CreateTopicFeedbackAttrs{
		Topic: existedTopic,
	})
	assert.Nil(t, err)
	assert.Equal(t, topic.Id, existedTopic.Id)
}

func TestTopicRepository_SaveTopicRoom(t *testing.T) {
	attrs := &models.CreateTopicAttrs{
		Message: models.TopicMessage{
			TelegramMessageId: 1,
			Question:          "test1",
		},
		Creator: models.TopicCreator{
			TelegramId:       8,
			TelegramUsername: "test2",
		},
	}

	existedTopic, _ := TopicRepository.CreateUniqueTopic(attrs)

	err := TopicRepository.SaveTopicRoom(context.TODO(), existedTopic)
	assert.Nil(t, err)
}

func TestTopicRepository_GetTopicRoomBySenderId(t *testing.T) {
	attrs := &models.CreateTopicAttrs{
		Message: models.TopicMessage{
			TelegramMessageId: 1,
			Question:          "test1",
		},
		Creator: models.TopicCreator{
			TelegramId:       9,
			TelegramUsername: "test2",
		},
	}

	existedTopic, _ := TopicRepository.CreateUniqueTopic(attrs)

	err := TopicRepository.SaveTopicRoom(context.TODO(), existedTopic)

	room, err := TopicRepository.GetTopicRoomBySenderId(context.TODO(), attrs.Creator.TelegramId)
	assert.Nil(t, err)
	assert.Equal(t, room.From, attrs.Creator.TelegramId)
}
