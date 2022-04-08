package resource

import (
	"shared"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type PatchResponseData struct{}

var validate = validator.New()

func Patch[T any](db *mongo.Database) func (*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var res shared.Response[PatchResponseData]

		var data map[string]any
		c.BodyParser(&data)

		
	}
}