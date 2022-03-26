package shared

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

type SuccessfulResponse[T any] struct {
	Data T `json:"data"`
}

type ErrorResponse[T any] struct {
	Error T `json:"error"`
}

type Response[D, E any] interface {
	SuccessfulResponse[D] | ErrorResponse[E]
}

func SendResponse[D, E any, T Response[D, E]](res T, c *fiber.Ctx) error {
	data, err := json.Marshal(res)
	if err != nil {
		return c.SendStatus(500)
	}

	return c.Send(data)
}