package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"shared"

	"content_api/collections"
	"content_api/collections/collection"
	"content_api/properties"
	"content_api/properties/property"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading environment: %v", err)
	}

	endpoint := fmt.Sprintf("mongodb://%s:%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"))
	ctx, cancelDBContext := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancelDBContext()

	client, err := shared.NewDBClient(ctx, endpoint)
	if err != nil {
		log.Fatalf("Error creating database client: %v", err)
	}
	defer client.Disconnect(ctx)

	app := fiber.New()

	app.Post("/properties", properties.Post(client))
	app.Get("/properties", properties.Get(client))
	app.Get("/properties/:id", property.Get(client))
	app.Delete("/properties/:id", property.Delete(client))

	app.Post("/properties", collections.Post(client))
	app.Get("/properties", collections.Get(client))
	app.Get("/properties/:id", collection.Get(client))
	app.Delete("/properties/:id", collection.Delete(client))
	
	app.Listen(":3000")
}