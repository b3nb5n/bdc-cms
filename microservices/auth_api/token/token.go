package token

import "github.com/go-playground/validator/v10"

type Credential struct {
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

var validate = validator.New()