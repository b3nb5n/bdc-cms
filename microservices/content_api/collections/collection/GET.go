package collection

import (
	"content_api/collections"
	"context"
	"shared"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type GetResponseData collections.Collection

type GetResponseError string

func Get(client *mongo.Client) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			res := shared.ErrorResponse[GetResponseError]{Error: "Invalid id"}
			return shared.SendResponse[GetResponseData, GetResponseError](res, c)
		}

		queryCtx, cancelQueryCtx := context.WithTimeout(context.Background(), 6*time.Second)
		defer cancelQueryCtx()
		queryResult := client.Database("content").Collection(collections.COLLECTION).FindOne(queryCtx, bson.M{"_id": id})
		data := new(GetResponseData)
		err = queryResult.Decode(&data)
		if err != nil {
			return c.SendStatus(500)
		}

		res := shared.SuccessfulResponse[GetResponseData]{
			Data: *data,
		}
		return shared.SendResponse[GetResponseData, GetResponseError](res, c)
	}
}