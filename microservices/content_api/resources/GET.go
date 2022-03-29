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

func Get[T any](db *mongo.Database) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		res := new(shared.Response[GetResponseData[T]])

		queryCtx, cancelQueryCtx := context.WithTimeout(context.Background(), 6*time.Second)
		defer cancelQueryCtx()
		collection := utils.ResolveCollection(c.Path())
		queryResult, err := db.Collection(collection).Find(queryCtx, bson.M{})
		if err != nil {
			switch err {
			case mongo.ErrNoDocuments:
				return res.Send(c.Status(404))
			default:
				return res.Send(c.Status(500))
			}
		}
		defer queryResult.Close(context.Background())

		decodeCtx, cancelDecodeCtx := context.WithTimeout(context.Background(), time.Second)
		defer cancelDecodeCtx()
		err = queryResult.All(decodeCtx, &res.Data)
		if err != nil {
			return res.Send(c.Status(500))
		}

		return res.Send(c)
	}
}