package main

import (
	"encoding/json"
	"shared"
	"strings"
)

type UserData struct {
	FirstName string `bson:"firstName" json:"firstName" valdiate:"required"`
	LastName string `bson:"lastName" json:"lastName" validate:"required"`
	Email string `bson:"email" json:"email" validate:"required,email"`
	Password string `bson:"password" json:"password" validate:"required"`
}

type User shared.Resource[UserData]

func (user *UserData) UnmarshalJSON(data []byte) error {
	type Alias UserData
	aux := &struct {
		*Alias
		Email string `json:"email"`
	} {
		Alias: (*Alias)(user),
	}

	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}

	user.Email = strings.ToLower(aux.Email)
	return nil
}