package resource

import (
	"context"
	"shared"
	"strconv"
	"time"

	"content_api/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type GetResponseData[T any] shared.Resource[T]

type GetResponseError string

func Get[T any](db *mongo.Database) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		res := new(shared.Response[GetResponseData[T], GetResponseError])

		idParam := c.Params("id")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			return res.Send(c.Status(400))
		}

		queryCtx, cancelQueryCtx := context.WithTimeout(context.Background(), 6*time.Second)
		defer cancelQueryCtx()
		collection := utils.ResolveCollection(c.Path())
		queryResult := db.Collection(collection).FindOne(queryCtx, bson.M{"_id": id})
		err = queryResult.Decode(&res.Data)
		if err != nil {
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