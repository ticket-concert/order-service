package request

type QueueReq struct {
	UserId  string `json:"userId" validate:"required"`
	EventId string `json:"eventId" validate:"required"`
}
