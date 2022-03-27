package resources

import (
	"context"
	"encoding/json"
	"time"

	"shared"

	"content_api/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostResponseData struct {
	ID shared.Snowflake `json:"id"`
}

type PostResponseError string

func Post[T any](client *mongo.Client) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var data T
		err := json.Unmarshal(c.Body(), &data)
		if err != nil {
			res := shared.ErrorResponse[PostResponseError]{Error: "Invalid Body"}
			return shared.SendResponse[PostResponseData, PostResponseError](res, c)
		}

		resource, err := shared.NewResource(data)
		if err != nil {
			return c.SendStatus(500)
		}

		ctx, cancelWriteCtx := context.WithTimeout(context.Background(), 6*time.Second)
		defer cancelWriteCtx()
		collection := utils.ResolveCollection(c.Path())
		_, err = client.Database("content").Collection(collection).InsertOne(ctx, resource)
		if err != nil {
			return c.SendStatus(500)
		}

		res := shared.SuccessfulResponse[PostResponseData]{
			Data: PostResponseData{ID: resource.ID},
		}
		return shared.SendResponse[PostResponseData, PostResponseError](res, c)
	}
}
