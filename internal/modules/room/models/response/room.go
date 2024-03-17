package response

type QueueResp struct {
	UserId      string `json:"userId" bson:"userId"`
	QueueNumber int    `json:"queueNumber" bson:"queueNumber"`
	CountryCode string `json:"countryCode" bson:"countryCode"`
}
