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

type DeleteResponseData struct{}

type DeleteResponseError struct{}

func Delete(db *mongo.Database) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		res := shared.Response[DeleteResponseData, DeleteResponseError]{}

		idParam := c.Params("id")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			res.Error.Global = "Invalid resource id"
			res.Send(c.Status(400))
		}

		queryCtx, cancelQueryCtx := context.WithTimeout(context.Background(), 6*time.Second)
		defer cancelQueryCtx()
		collection := utils.ResolveCollection(c.Path())
		queryResult := db.Collection(collection).FindOneAndDelete(queryCtx, bson.M{"_id": id})
		if err = queryResult.Err(); err != nil {
			switch err {
			case mongo.ErrNoDocuments:
				return res.Send(c.Status(404))
			default:
				return res.Send(c.Status(500))
			}
		}

		return res.Send(c)
	}
}
