package response

import "time"

type QueueResp struct {
	QueueId     string    `json:"queueId" bson:"queueId"`
	UserId      string    `json:"userId" bson:"userId"`
	QueueNumber int       `json:"queueNumber" bson:"queueNumber"`
	CountryCode string    `json:"countryCode" bson:"countryCode"`
	CreatedAt   time.Time `json:"createdAt" bson:"createdAt"`
}
