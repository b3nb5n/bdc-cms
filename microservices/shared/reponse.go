package shared

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

type Response[D, E any] struct {
	Data D `json:"data,omitempty"`
	Error E `json:"error,omitempty"`
}

func SendResponse[D, E any](res Response[D, E], c *fiber.Ctx) error {
	data, err := json.Marshal(res)
	if err != nil {
		return c.SendStatus(500)
	}

	return c.Send(data)
}