package token

import (
	"encoding/json"
	"fmt"
	"os"
	"shared"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Payload struct {
	IssuedAt time.Time `json:"iat"`
	ExpiresAt time.Time `json:"exp"`
	Subject shared.Snowflake `json:"sub"`
}

func (payload *Payload) MarshalJSON() (data []byte, err error) {
	type Alias Payload
	return json.Marshal(struct{
		*Alias
		IssuedAt int64 `json:"iat"`
		ExipresAt int64 `json:"exp"`
	}{
		Alias: (*Alias)(payload),
		IssuedAt: payload.IssuedAt.Unix(),
		ExipresAt: payload.ExpiresAt.Unix(),
	})
}

func (payload *Payload) UnmarshalJSON(data []byte) error {
	type Alias Payload
	aux := &struct{
		*Alias
		IssuedAt int64 `json:"iat"`
		ExipresAt int64 `json:"exp"`
	}{
		Alias: (*Alias)(payload),
	}

	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}

	payload.IssuedAt = time.Unix(aux.IssuedAt, 0)
	payload.ExpiresAt = time.Unix(aux.ExipresAt, 0)
	return err
}

func (payload *Payload) Valid() error {
	return nil
}

func NewPayload(sub shared.Snowflake, ttl time.Duration) *Payload {
	return &Payload{
		IssuedAt: time.Now(),
		ExpiresAt: time.Now().Add(ttl),
		Subject: sub,
	}
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