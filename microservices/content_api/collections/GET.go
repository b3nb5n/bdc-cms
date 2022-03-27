package collections

import (
	"context"
	"time"

	"shared"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type GetResponseData []Collection

type GetResponseError string

func Get(client *mongo.Client) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		queryCtx, cancelQueryCtx := context.WithTimeout(context.Background(), 6*time.Second)
		defer cancelQueryCtx()
		queryResult, err := client.Database("content").Collection(COLLECTION).Find(queryCtx, bson.D{})
		if err != nil {
			return c.SendStatus(500)
		}
		defer queryResult.Close(context.Background())

		decodeCtx, cancelDecodeCtx := context.WithTimeout(context.Background(), time.Second)
		defer cancelDecodeCtx()
		data := new(GetResponseData)
		err = queryResult.All(decodeCtx, &data)
		if err != nil {
			return c.SendStatus(500)
		}

		res := shared.SuccessfulResponse[GetResponseData] {
			Data: *data,
		}
		return shared.SendResponse[GetResponseData, GetResponseError](res, c)
	}
}