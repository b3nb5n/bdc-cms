package token

import (
	"encoding/json"
	"fmt"
	"os"
	"shared"
	"strconv"

	"github.com/golang-jwt/jwt/v4"
)

type Payload struct {
	UID shared.Snowflake `validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

func (payload *Payload) MarshalJSON() (data []byte, err error) {
	type Alias Payload
	return json.Marshal(struct{
		*Alias
		UID string `json:"uid"`
	}{
		UID: strconv.FormatInt(int64(payload.UID), 10),
		Alias: (*Alias)(payload),
	})
}

func (payload *Payload) UnmarshalJSON(data []byte) error {
	type Alias Payload
	aux := &struct{
		*Alias
		UID string `json:"uid"`
	}{
		Alias: (*Alias)(payload),
	}

	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}

	uid, err := strconv.ParseInt(aux.UID, 10, 64)
	payload.UID = shared.Snowflake(uid)
	return err
}

func (payload *Payload) Valid() error {
	return nil
}

func ParsePayload(jwtString string) (payload Payload, err error) {
	token, err := jwt.ParseWithClaims(jwtString, &payload, func(token *jwt.Token) (interface{}, error) {
    return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if _, ok := token.Claims.(*Payload); !ok || !token.Valid {
		return payload, err
	}

	err = validate.Struct(payload)
	if err != nil {
		return payload, fmt.Errorf("Invalid claims")
	}

	return payload, err
}