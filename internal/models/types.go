package models

type MessageKey string

const (
	// HelloMessageKey greeting message
	HelloMessageKey MessageKey = "hello"
	// QueueStartMessageKey user in queue state
	QueueStartMessageKey MessageKey = "queue_start"
	// RoomStartMessageKey user in room state
	RoomStartMessageKey MessageKey = "room_start"
	// GotTicketKey service got user question
	GotTicketKey MessageKey = "got_ticket"
	// NewTicketCreatedMessageKey send message to support group
	NewTicketCreatedMessageKey MessageKey = "new_ticket"
	// TicketAcceptedMessageKey when request accepted
	TicketAcceptedMessageKey MessageKey = "ticket_accepted"
	// TicketDeclinedMessageKey when request declined
	TicketDeclinedMessageKey MessageKey = "ticket_declined"
	// SuccessAcceptTicketMessageKey when help accepted
	SuccessAcceptTicketMessageKey MessageKey = "success_accept_ticket"
	// SuccessDeclineTicketMessageKey when help declined
	SuccessDeclineTicketMessageKey MessageKey = "success_decline_ticket"
	// AcceptTicketMessageKey accept help button text
	AcceptTicketMessageKey MessageKey = "accept_ticket"
	// DeclineTicketMessageKey decline help button text
	DeclineTicketMessageKey MessageKey = "decline_ticket"
	// NewTicketMessageKey new message topic (for user)
	NewTicketMessageKey MessageKey = "new_message"
	// NewTicketCommentMessageKey new comment (for support)
	NewTicketCommentMessageKey MessageKey = "new_comment"
	// CompletePrevTicketMessageKey suggestion for support
	CompletePrevTicketMessageKey MessageKey = "complete_last_ticket"
	// UserNoExpectHelpMessageKey if user already got help
	UserNoExpectHelpMessageKey MessageKey = "user_no_expect"
	// LeaveFeedbackMessageKey leave feedback message
	LeaveFeedbackMessageKey MessageKey = "leave_rating"
	// FeedbackSendMessageKey feedback send success
	FeedbackSendMessageKey MessageKey = "feedback_received"
	RatingOk               MessageKey = "rating_ok"
	RatingNormal           MessageKey = "rating_normal"
	RatingBad              MessageKey = "rating_bad"
	NoRating               MessageKey = "no_rating"
	TicketClosedMessageKey MessageKey = "ticket_end"
	NoRightMessageKey      MessageKey = "no_rights"
)

const (
	// DefaultState if sender don't have active topic or not in queue or just any message
	DefaultState = "default-state"
	// QueueState add user to queue when topic created but wait decision
	QueueState = "queue-state"
	// RoomState when support accept topic and start "chat" with user
	RoomState = "room-state"
	// UnknownState something wrong
	UnknownState = "unknown-state"
)

const (
	// BotWasBlockedError telegram api message
	BotWasBlockedError = "Forbidden: bot was blocked by the user"
)

const (
	StartBotCommand   = "/start"
	BreakTopicCommand = "/break"
)

const (
	AcceptHelp  = "accept"
	DeclineHelp = "decline"
	NoRight     = "no_right"
	Rating      = "rating"
)

const (
	TopicsTable = "public.topics"
)
