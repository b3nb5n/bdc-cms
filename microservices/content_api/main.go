package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	endpoint := "mongodb://localhost:27017"

	ctx, cancelDBContext := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancelDBContext()
	
	client, err := mongo.NewClient(options.Client().ApplyURI(endpoint))
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer client.Disconnect(ctx)

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalf("Error pinging db: %v", err)
	}

	fmt.Println("connected")
}