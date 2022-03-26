package properties

import (
	"context"
	"time"

	"shared"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type GetResponseData []Property

type GetResponseError string

func Get(client *mongo.Client) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		queryCtx, cancelQueryCtx := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelQueryCtx()
		queryResult, err := client.Database("content").Collection("properties").Find(queryCtx, bson.D{})
		if err != nil {
			return c.SendStatus(500)
		}
		defer queryResult.Close(context.Background())

		documents := make([]Property, 0)
		decodeCtx, cancelDecodeCtx := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelDecodeCtx()
		err = queryResult.All(decodeCtx, &documents)
		if err != nil {
			return c.SendStatus(500)
		}

		res := shared.SuccessfulResponse[GetResponseData] {
			Data: documents,
		}
		return shared.SendResponse[GetResponseData, GetResponseError](res, c)
	}
}