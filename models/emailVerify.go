package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type EmailVerify struct {
	ID  primitive.ObjectID `json:"id"`
	OTP string             `json:"otp"`
}
