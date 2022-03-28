package shared

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

type ErrorResponse[T any] struct {
	Global string `json:"global,omitempty"`
	Body T `json:"body,omitempty"`
}

type Response[D, E any] struct {
	Data D `json:"data,omitempty"`
	Error ErrorResponse[E] `json:"error,omitempty"`
}

func (res Response[_, _]) Send(c *fiber.Ctx) error {
	data, err := json.Marshal(res)
	if err != nil {
		return c.SendStatus(500)
	}

	return c.Send(data)
}