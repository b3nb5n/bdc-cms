package me

import (
	"auth_api/token"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func Get(db *mongo.Database) func (*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		headers := c.GetReqHeaders()
		jwtString := headers["Authorization"][7:]
		payload, err := token.ParsePayload(jwtString)
		if err != nil {
			return c.SendStatus(400)
		}

		return c.Redirect(fmt.Sprintf("/user/%v", payload.UID), fiber.StatusSeeOther)
	}
}