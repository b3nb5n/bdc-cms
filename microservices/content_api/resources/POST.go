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
	ID shared.Snowflake `json:"id"`
}

type PostResponseError string

var validate = validator.New()

func Post[T any](db *mongo.Database) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		data := new(T)
		c.BodyParser(&data)
		err := validate.Struct(data)
		if err != nil {
			return c.SendStatus(400)
		}

		resource, err := shared.NewResource(*data)
		if err != nil {
			return c.SendStatus(500)
		}

		ctx, cancelWriteCtx := context.WithTimeout(context.Background(), 6*time.Second)
		defer cancelWriteCtx()
		collection := utils.ResolveCollection(c.Path())
		_, err = db.Collection(collection).InsertOne(ctx, resource)
		if err != nil {
			switch err {
			case mongo.ErrNoDocuments:
				return c.SendStatus(404)
			default:
				return c.SendStatus(500)
			}
		}

		res := shared.Response[PostResponseData, any]{
			Data: PostResponseData{ID: resource.ID},
		}
		return shared.SendResponse(res, c)
	}
}
