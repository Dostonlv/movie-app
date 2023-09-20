package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type EmailVerify struct {
	ID        primitive.ObjectID `json:"id"`
	OTP       string             `json:"otp"`
	CreatedAt time.Time          `json:"created_at"`
}
