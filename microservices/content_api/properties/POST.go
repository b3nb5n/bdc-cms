package properties

import (
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

func Post(client *mongo.Client) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var data PropertyData
		err := json.Unmarshal(c.Body(), &data)
		if err != nil {
			res := shared.ErrorResponse[PostResponseError]{Error: "Invalid Body"}
			return shared.SendResponse[PostResponseData, PostResponseError](res, c)
		}

		property, err := shared.NewResource(data)
		if err != nil {
			return c.SendStatus(500)
		}

		ctx, cancelWriteCtx := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelWriteCtx()
		_, err = client.Database("content").Collection("properties").InsertOne(ctx, property)
		if err != nil {
			return c.SendStatus(500)
		}

		res := shared.SuccessfulResponse[PostResponseData]{
			Data: PostResponseData{ID: property.ID},
		}
		return shared.SendResponse[PostResponseData, PostResponseError](res, c)
	}
}