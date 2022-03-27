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

type DeleteResponseData struct {}

type DeleteResponseError string

func Delete(db *mongo.Database) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			res := shared.ErrorResponse[DeleteResponseError]{Error: "Invalid id"}
			return shared.SendResponse[DeleteResponseData, DeleteResponseError](res, c)
		}

		queryCtx, cancelQueryCtx := context.WithTimeout(context.Background(), 6*time.Second)
		defer cancelQueryCtx()
		collection := utils.ResolveCollection(c.Path())
		queryResult := db.Collection(collection).FindOneAndDelete(queryCtx, bson.M{"_id": id})
		if err = queryResult.Err(); err != nil {
			res := shared.ErrorResponse[DeleteResponseError]{Error: "An unknown error ocurred"}
			return shared.SendResponse[DeleteResponseData, DeleteResponseError](res, c)
		}

		res := shared.SuccessfulResponse[DeleteResponseData] {
			Data: DeleteResponseData {},
		}
		return shared.SendResponse[DeleteResponseData, DeleteResponseError](res, c)
	}
}