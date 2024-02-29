package request

import "time"

type UpdateBankTicketReq struct {
	CountryCode   string    `json:"countryCode"`
	TicketType    string    `json:"ticketType"`
	Price         int       `json:"price"`
	UserId        string    `json:"userId"`
	QueueId       string    `json:"queueId"`
	TicketId      string    `json:"ticketId"`
	EventId       string    `json:"eventId"`
	PaymentStatus string    `json:"paymentStatus"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type OrderReq struct {
	UserId     string `json:"userId" validate:"required"`
	TicketType string `json:"ticketType" validate:"required"`
	EventId    string `json:"eventId" validate:"required"`
}

type GetOrderReq struct {
	TicketNumber string `json:"ticketNumber" validate:"required"`
}

type OrderList struct {
	Page   int64  `query:"page" validate:"required"`
	Size   int64  `query:"size" validate:"required"`
	UserId string `query:"userId"`
}

type PreOrderList struct {
	Page   int64  `query:"page" validate:"required"`
	Size   int64  `query:"size" validate:"required"`
	UserId string `query:"userId"`
}
