package main

import (
	"context"
	"encoding/json"
	"fmt"
	"shared"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type PropertyData struct {
	Hosts []string `bson:"hosts"`
}

type Property shared.Resource[PropertyData]

func PropertyPOST(client *mongo.Client) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var data PropertyData
		err := json.Unmarshal(c.Body(), &data)
		if err != nil {
			return fmt.Errorf("Error unmarshaling body: %v", err)
		}

		property, err := shared.NewResource(data)
		if err != nil {
			return fmt.Errorf("Error constructing resource: %v", err)
		}

		ctx, cancelWriteCtx := context.WithTimeout(context.Background(), 10 * time.Second)
		defer cancelWriteCtx()
		_, err = client.Database("content").Collection("properties").InsertOne(ctx, property)
		if err != nil {
			return fmt.Errorf("Error writing document: %v", err)
		}

		return err
	}
}