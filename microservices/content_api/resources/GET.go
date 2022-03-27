package resources

import (
	"content_api/utils"
	"context"
	"time"

	"shared"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type GetResponseData[T any] []shared.Resource[T]

type GetResponseError string

func Get[T any](client *mongo.Client) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		queryCtx, cancelQueryCtx := context.WithTimeout(context.Background(), 6*time.Second)
		defer cancelQueryCtx()
		collection := utils.ResolveCollection(c.Path())
		queryResult, err := client.Database("content").Collection(collection).Find(queryCtx, bson.D{})
		if err != nil {
			return c.SendStatus(500)
		}
		defer queryResult.Close(context.Background())

		resources := make(GetResponseData[T], 0)
		decodeCtx, cancelDecodeCtx := context.WithTimeout(context.Background(), time.Second)
		defer cancelDecodeCtx()
		err = queryResult.All(decodeCtx, &resources)
		if err != nil {
			return c.SendStatus(500)
		}

		res := shared.SuccessfulResponse[GetResponseData[T]] {
			Data: resources,
		}
		return shared.SendResponse[GetResponseData[T], GetResponseError](res, c)
	}
}