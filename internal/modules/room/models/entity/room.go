package entity

import "time"

type QueueRoom struct {
	QueueId     string    `json:"queueId" bson:"queueId"`
	UserId      string    `json:"userId" bson:"userId"`
	EventId     string    `json:"eventId" bson:"eventId"`
	QueueNumber int       `json:"queueNumber" bson:"queueNumber"`
	CountryCode string    `json:"countryCode" bson:"countryCode"`
	CreatedAt   time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt" bson:"updatedAt"`
}
