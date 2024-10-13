package models

import "github.com/go-playground/validator/v10"

type User struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=20"`
}

var Validate *validator.Validate

func init() {
	Validate = validator.New()
}
