package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"shared"

	"content_api/resources"
	"content_api/resources/resource"

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

	db := client.Database("content")
	app := fiber.New()

	app.Post("/properties", resources.Post[resources.PropertyData](db))
	app.Get("/properties", resources.Get[resources.PropertyData](db))
	app.Get("/properties/:id", resource.Get[resources.PropertyData](db))
	app.Delete("/properties/:id", resource.Delete(db))

	app.Post("/collections", resources.Post[resources.CollectionData](db))
	app.Get("/collections", resources.Get[resources.CollectionData](db))
	app.Get("/collections/:id", resource.Get[resources.CollectionData](db))
	app.Delete("/collections/:id", resource.Delete(db))

	app.Post("/items", resources.Post[resources.ItemData](db))
	app.Get("/items", resources.Get[resources.ItemData](db))
	app.Get("/items/:id", resource.Get[resources.ItemData](db))
	app.Delete("/items/:id", resource.Delete(db))

	app.Listen(":3000")
}