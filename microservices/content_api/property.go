package main

import (
	"context"
	"encoding/json"
	"shared"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PropertyData struct {
	Hosts []string `bson:"hosts" json:"hosts"`
}

type Property shared.Resource[PropertyData]

type PropertyPostResponseData struct {
	ID shared.Snowflake `json:"id"`
}

type PropertyPostResponseError string

func PropertyPOST(client *mongo.Client) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var data PropertyData
		err := json.Unmarshal(c.Body(), &data)
		if err != nil {
			res := shared.ErrorResponse[PropertyPostResponseError]{Error: "Invalid Body"}
			return shared.SendResponse[PropertyPostResponseData, PropertyPostResponseError](res, c)
		}

		property, err := shared.NewResource(data)
		if err != nil {
			return c.SendStatus(500)
		}

		ctx, cancelWriteCtx := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelWriteCtx()
		_, err = client.Database("content").Collection("properties").InsertOne(ctx, property)
		if err != nil {
			return c.SendStatus(500)
		}

		res := shared.SuccessfulResponse[PropertyPostResponseData]{
			Data: PropertyPostResponseData{ID: property.ID},
		}
		return shared.SendResponse[PropertyPostResponseData, PropertyPostResponseError](res, c)
	}
}

type PropertyGetResponseData Property

type PropertyGetResponseError string

func PropertyGET(client *mongo.Client) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			res := shared.ErrorResponse[PropertyGetResponseError]{Error: "Invalid id"}
			return shared.SendResponse[PropertyGetResponseData, PropertyGetResponseError](res, c)
		}

		queryCtx, cancelQueryCtx := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelQueryCtx()
		queryResult := client.Database("content").Collection("properties").FindOne(queryCtx, bson.M{"_id": id})
		var property Property
		err = queryResult.Decode(&property)
		if err != nil {
			return c.SendStatus(500)
		}

		res := shared.SuccessfulResponse[PropertyGetResponseData]{
			Data: PropertyGetResponseData(property),
		}
		return shared.SendResponse[PropertyGetResponseData, PropertyGetResponseError](res, c)
	}
}

type PropertyDeleteResponseData struct {}

type PropertyDeleteResponseError string

func PropertyDELETE(client *mongo.Client) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			res := shared.ErrorResponse[PropertyDeleteResponseError]{Error: "Invalid id"}
			return shared.SendResponse[PropertyDeleteResponseData, PropertyDeleteResponseError](res, c)
		}

		queryCtx, cancelQueryCtx := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelQueryCtx()
		queryResult := client.Database("content").Collection("properties").FindOneAndDelete(queryCtx, bson.M{"_id": id})
		if err = queryResult.Err(); err != nil {
			res := shared.ErrorResponse[PropertyDeleteResponseError]{Error: "An unknown error ocurred"}
			return shared.SendResponse[PropertyDeleteResponseData, PropertyDeleteResponseError](res, c)
		}

		res := shared.SuccessfulResponse[PropertyDeleteResponseData] {
			Data: struct {} {},
		}
		return shared.SendResponse[PropertyDeleteResponseData, PropertyDeleteResponseError](res, c)
	}
}

type PropertiesGetResponseData []Property

type PropertiesGetResponseError string

func PropertiesGET(client *mongo.Client) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		queryCtx, cancelQueryCtx := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelQueryCtx()
		queryResult, err := client.Database("content").Collection("properties").Find(queryCtx, bson.D{})
		if err != nil {
			return c.SendStatus(500)
		}
		defer queryResult.Close(context.Background())

		documents := make([]Property, 0)
		decodeCtx, cancelDecodeCtx := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelDecodeCtx()
		err = queryResult.All(decodeCtx, &documents)
		if err != nil {
			return c.SendStatus(500)
		}

		res := shared.SuccessfulResponse[PropertiesGetResponseData] {
			Data: documents,
		}
		return shared.SendResponse[PropertiesGetResponseData, PropertiesGetResponseError](res, c)
	}
}
