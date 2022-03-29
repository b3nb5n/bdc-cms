package resources

import (
	"content_api/utils"
	"context"
	"time"

	"shared"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostResponseData struct {
	ID shared.Snowflake `json:"id,omitempty"`
}

type PostResponseError string

var validate = validator.New()

func Post[T any](db *mongo.Database) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		res := new(shared.Response[PostResponseData])

		data := new(T)
		c.BodyParser(&data)
		err := validate.Struct(data)
		if err != nil {
			return res.Send(c.Status(400))	
		}

		resource, err := shared.NewResource(*data)
		if err != nil {
			return res.Send(c.Status(500))
		}

		ctx, cancelWriteCtx := context.WithTimeout(context.Background(), 6*time.Second)
		defer cancelWriteCtx()
		collection := utils.ResolveCollection(c.Path())
		_, err = db.Collection(collection).InsertOne(ctx, resource)
		if err != nil {
			return res.Send(c.Status(500))
		}

		res.Data.ID = resource.ID
		return res.Send(c.Status(201))
	}
}
