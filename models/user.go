package models

import (
	_ "github.com/go-playground/validator/v10"
)

type User struct {
	Name     string `validate:"required" bson:"name" json:"name"`
	Email    string `validate:"required,email" bson:"email" json:"email"`
	Password string `validate:"required,min=8,max=20" bson:"password" json:"password"`
	IsVerify bool   `bson:"is_verify" json:"is_verify" default:"false"`
}

type SignIn struct {
	Email    string `validate:"required,email" bson:"email" json:"email"`
	Password string `validate:"required,min=8,max=20" bson:"password" json:"password"`
}
