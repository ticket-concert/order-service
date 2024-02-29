package response

import (
	"order-service/internal/pkg/constants"
	"time"
)

type OrderResp struct {
	QueueId      string    `json:"queueId"`
	UserId       string    `json:"userId"`
	EventId      string    `json:"eventId"`
	TicketType   string    `json:"ticketType"`
	CountryCode  string    `json:"countryCode"`
	OrderTime    time.Time `json:"orderTime"`
	Price        int       `json:"price"`
	TicketNumber string    `json:"ticketNumber"`
}

type OrderList struct {
	FullName     string    `json:"fullName"`
	TicketType   string    `json:"ticketType"`
	TicketNumber string    `json:"ticketNumber"`
	TicketPrice  int       `json:"ticketPrice"`
	SeatNumber   int       `json:"seatNumber"`
	EventName    string    `json:"eventName"`
	EventTime    time.Time `json:"eventTime"`
	EventPlace   string    `json:"eventPlace"`
	EventId      string    `json:"eventId"`
	TicketId     string    `json:"ticketId"`
	OrderTime    time.Time `json:"orderTime"`
}

type OrderListResp struct {
	CollectionData []OrderList
	MetaData       constants.MetaData
}

type PreOrderList struct {
	UserId       string    `json:"userId"`
	TicketType   string    `json:"ticketType"`
	TicketNumber string    `json:"ticketNumber"`
	TicketPrice  int       `json:"ticketPrice"`
	OrderTime    time.Time `json:"orderTime"`
	EventId      string    `json:"eventId"`
	MaxWaitTime  string    `json:"maxWaitTime"`
}

type PreOrderListResp struct {
	CollectionData []PreOrderList
	MetaData       constants.MetaData
}
