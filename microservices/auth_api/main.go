package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"shared"

	"auth_api/token"
	"auth_api/users"
	"auth_api/users/me"
	"auth_api/users/user"

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

	db := client.Database("auth")
	app := fiber.New()

	app.Post("/user", users.Post(db))
	app.Get("/user/me", me.Get(db))
	app.Get("/user/:id", user.Get(db))

	app.Get("/token", token.Get(db))

	app.Listen(":3000")
}