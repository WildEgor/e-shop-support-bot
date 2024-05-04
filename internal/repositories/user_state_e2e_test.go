package repositories_test

import (
	"context"
	"github.com/WildEgor/e-shop-support-bot/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserStateRepository_AddUserToQueue(t *testing.T) {

	err := UserStateRepository.AddUserToQueue(context.TODO(), 1)

	errDel := UserStateRepository.DeleteUserFromQueue(context.TODO(), 1)

	assert.Nil(t, err)
	assert.Nil(t, errDel)
}

func TestUserStateRepository_AddUserMessageToBuffer(t *testing.T) {

	msg := &models.UserMessage{
		TelegramMessageId: 1,
		TelegramUserId:    1,
		ChatId:            1,
		Content:           "test",
	}

	err := UserStateRepository.AddUserMessageToBuffer(context.TODO(), msg)
	msgs := UserStateRepository.GetUserMessagesFromBuffer(context.TODO(), 1)

	assert.Nil(t, err)
	assert.NotNil(t, msgs)
	assert.Equal(t, 1, len(msgs))

	errClean := UserStateRepository.CleanUserMessagesFromBuffer(context.TODO(), 1)
	oldMessages := UserStateRepository.GetUserMessagesFromBuffer(context.TODO(), 1)

	assert.Nil(t, errClean)
	assert.NotNil(t, oldMessages)
	assert.Equal(t, 0, len(oldMessages))
}

func TestUserStateRepository_CheckUserOpts(t *testing.T) {
	uopts := &models.UserOptions{
		TelegramId: 1,
		ChatId:     1,
		Lang:       "ru",
	}

	err := UserStateRepository.SaveUserOptions(context.TODO(), uopts)
	assert.Nil(t, err)

	options, err := UserStateRepository.GetUserOptions(context.TODO(), uopts.TelegramId)
	assert.Nil(t, err)
	assert.NotNil(t, options)

	assert.Equal(t, uopts.TelegramId, options.TelegramId)
	assert.Equal(t, uopts.ChatId, options.ChatId)
	assert.Equal(t, uopts.Lang, options.Lang)
}

func TestUserStateRepository_CheckUserState(t *testing.T) {
	state := UserStateRepository.CheckUserState(context.TODO(), 1)

	assert.Equal(t, models.DefaultState, state)
}
