package user

import (
	"context"
	"shared"
	"strconv"
	"time"

	"auth_api/users"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Get(db * mongo.Database) func (*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		res := new(shared.Response[users.User])

		id, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			return c.SendStatus(404)
		}

		queryCtx, cancelQueryCtx := context.WithTimeout(context.Background(), 6*time.Second)
		defer cancelQueryCtx()
		queryRes := db.Collection("users").FindOne(queryCtx, bson.M{"_id": id})
		if err := queryRes.Err(); err != nil {
			if err == mongo.ErrNoDocuments {
				return c.SendStatus(404)
			}

			return c.SendStatus(500)
		}

		err = queryRes.Decode(&res.Data)
		if err != nil {
			return c.SendStatus(500)
		}

		return res.Send(c)
	}
}