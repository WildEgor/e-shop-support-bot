package handlers

import (
	accept_callback_handler "github.com/WildEgor/e-shop-support-bot/internal/handlers/accept_callback"
	break_action_handler "github.com/WildEgor/e-shop-support-bot/internal/handlers/break_action"
	decline_callback_handler "github.com/WildEgor/e-shop-support-bot/internal/handlers/decline_callback"
	edit_message_handler "github.com/WildEgor/e-shop-support-bot/internal/handlers/edit_message"
	eh "github.com/WildEgor/e-shop-support-bot/internal/handlers/errors"
	hch "github.com/WildEgor/e-shop-support-bot/internal/handlers/health_check"
	new_message_handler "github.com/WildEgor/e-shop-support-bot/internal/handlers/new_message"
	no_right_callback_handler "github.com/WildEgor/e-shop-support-bot/internal/handlers/no_right_callback"
	rating_callback_handler "github.com/WildEgor/e-shop-support-bot/internal/handlers/rating_callback"
	rch "github.com/WildEgor/e-shop-support-bot/internal/handlers/ready_check"
	start_action_handler "github.com/WildEgor/e-shop-support-bot/internal/handlers/start_action"
	"github.com/WildEgor/e-shop-support-bot/internal/repositories"
	"github.com/WildEgor/e-shop-support-bot/internal/services"
	"github.com/google/wire"
)

var HandlersSet = wire.NewSet(
	services.ServicesSet,
	repositories.RepositoriesSet,
	eh.NewErrorsHandler,
	hch.NewHealthCheckHandler,
	rch.NewReadyCheckHandler,
	accept_callback_handler.NewAcceptCallbackHandler,
	break_action_handler.NewBreakActionHandler,
	decline_callback_handler.NewDeclineCallbackHandler,
	edit_message_handler.NewEditMessageHandler,
	new_message_handler.NewNewMessageHandler,
	no_right_callback_handler.NewNoRightCallbackHandler,
	rating_callback_handler.NewRatingCallbackHandler,
	start_action_handler.NewStartActionHandler,
)
