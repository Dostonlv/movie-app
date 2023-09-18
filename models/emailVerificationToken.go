package models

import "time"

type EmailVerificationToken struct {
	OwnerID   interface{} `bson:"_id"`
	Token     string      `bson:"token" json:"token"`
	CreatedAt time.Time   `bson:"created_at" json:"created_at"`
}
