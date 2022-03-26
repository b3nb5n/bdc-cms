package property

import (
	"context"
	"shared"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DeleteResponseData struct {}

type DeleteResponseError string

func Delete(client *mongo.Client) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			res := shared.ErrorResponse[DeleteResponseError]{Error: "Invalid id"}
			return shared.SendResponse[DeleteResponseData, DeleteResponseError](res, c)
		}

		queryCtx, cancelQueryCtx := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelQueryCtx()
		queryResult := client.Database("content").Collection("properties").FindOneAndDelete(queryCtx, bson.M{"_id": id})
		if err = queryResult.Err(); err != nil {
			res := shared.ErrorResponse[DeleteResponseError]{Error: "An unknown error ocurred"}
			return shared.SendResponse[DeleteResponseData, DeleteResponseError](res, c)
		}

		res := shared.SuccessfulResponse[DeleteResponseData] {
			Data: struct {} {},
		}
		return shared.SendResponse[DeleteResponseData, DeleteResponseError](res, c)
	}
}