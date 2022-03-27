package property

import (
	"content_api/properties"
	"context"
	"shared"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type GetResponseData properties.Property

type GetResponseError string

func Get(client *mongo.Client) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			res := shared.ErrorResponse[GetResponseError]{Error: "Invalid id"}
			return shared.SendResponse[GetResponseData, GetResponseError](res, c)
		}

		queryCtx, cancelQueryCtx := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelQueryCtx()
<<<<<<< HEAD:microservices/content_api/resources/resource/GET.go
		collection := utils.ResolveCollection(c.Path())
		queryResult := client.Database("content").Collection(collection).FindOne(queryCtx, bson.M{"_id": id})
		resource := new(GetResponseData[T])
		err = queryResult.Decode(resource)
=======
		queryResult := client.Database("content").Collection("properties").FindOne(queryCtx, bson.M{"_id": id})
		var property properties.Property
		err = queryResult.Decode(&property)
>>>>>>> parent of 950d857 (generic handlers):microservices/content_api/properties/property/GET.go
		if err != nil {
			return c.SendStatus(500)
		}

		res := shared.SuccessfulResponse[GetResponseData]{
			Data: GetResponseData(property),
		}
		return shared.SendResponse[GetResponseData, GetResponseError](res, c)
	}
}