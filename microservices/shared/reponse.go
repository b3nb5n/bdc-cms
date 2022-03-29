package shared

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

type Response[D any] struct {
	Data  D      `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}

func (res Response[_]) Send(c *fiber.Ctx) error {
	data, err := json.Marshal(res)
	if err != nil {
		return c.SendStatus(500)
	}

	return c.Send(data)
}
