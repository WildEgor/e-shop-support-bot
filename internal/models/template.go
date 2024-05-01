package models

// HINT: using "template" key decouple our templates from struct

type FeedbackTemplatePayload struct {
	TicketId int64 `template:"TicketID"`
}

type TicketDecisionPayload struct {
	TicketId      int64  `template:"TicketID"`
	FromUserId    int64  `template:"FromUserID"`
	FromUsername  string `template:"FromUsername"`
	FromFirstName string `template:"FromFirstName"`
	Text          string `template:"Text"`
}

type TicketPayload struct {
	TicketId      int64  `template:"TicketID"`
	FromUsername  string `template:"FromUsername"`
	FromFirstName string `template:"FromFirstName"`
	Text          string `template:"Text"`
}

type UserInfoPayload struct {
	Username string `template:"Username"`
}
