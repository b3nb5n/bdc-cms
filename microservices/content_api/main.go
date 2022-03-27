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
	ctx, cancelDBContext := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelDBContext()

	client, err := shared.NewDBClient(ctx, endpoint)
	if err != nil {
		log.Fatalf("Error creating database client: %v", err)
	}
	defer client.Disconnect(ctx)

	app := fiber.New()

	app.Post("/properties", resources.Post[resources.PropertyData](client))
	app.Get("/properties", resources.Get[resources.PropertyData](client))
	app.Get("/properties/:id", resource.Get[resources.PropertyData](client))
	app.Delete("/properties/:id", resource.Delete(client))

	app.Listen(":3000")
}
