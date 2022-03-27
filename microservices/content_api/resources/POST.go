package resources

import (
	"content_api/utils"
	"context"
	"encoding/json"
	"time"

	"shared"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostResponseData struct {
	ID shared.Snowflake `json:"id"`
}

type PostResponseError string

func Post[T any](db *mongo.Database) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		data := new(T)
		err := json.Unmarshal(c.Body(), data)
		if err != nil {
			res := shared.ErrorResponse[PostResponseError]{Error: "Invalid Body"}
			return shared.SendResponse[PostResponseData, PostResponseError](res, c)
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
			return c.SendStatus(500)
		}

		res := shared.SuccessfulResponse[PostResponseData]{
			Data: PostResponseData{ID: resource.ID},
		}
		return shared.SendResponse[PostResponseData, PostResponseError](res, c)
	}
}
