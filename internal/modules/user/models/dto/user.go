package dto

import "time"

type UserResp struct {
	UserId    string    `json:"userId" bson:"userId"`
	FullName  string    `json:"fullName" bson:"fullName"`
	Email     string    `json:"email" bson:"email"`
	Role      string    `json:"role" bson:"role"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type UserData struct {
	Data UserResp `json:"data" bson:"data"`
}
