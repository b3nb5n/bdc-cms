package resource

import (
	"content_api/utils"
	"context"
	"shared"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type GetResponseData[T any] shared.Resource[T]

type GetResponseError string

func Get[T any](client *mongo.Client) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			res := shared.ErrorResponse[GetResponseError]{Error: "Invalid id"}
			return shared.SendResponse[GetResponseData[T], GetResponseError](res, c)
		}

		queryCtx, cancelQueryCtx := context.WithTimeout(context.Background(), 6*time.Second)
		defer cancelQueryCtx()
		collection := utils.ResolveCollection(c.Path())
		queryResult := client.Database("content").Collection(collection).FindOne(queryCtx, bson.M{"_id": id})
		resource := new(GetResponseData[T])
		err = queryResult.Decode(resource)
		if err != nil {
			return c.SendStatus(500)
		}

		res := shared.SuccessfulResponse[GetResponseData[T]]{
			Data: *resource,
		}
		return shared.SendResponse[GetResponseData[T], GetResponseError](res, c)
	}
}