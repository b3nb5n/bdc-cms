package shared

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func NewDBClient(ctx context.Context, endpoint string) (client *mongo.Client, err error) {
	client, err = mongo.NewClient(options.Client().ApplyURI(endpoint))
	if err != nil {
		return client, fmt.Errorf("Error connecting to database: %v", err)
	}
	
	err = client.Connect(ctx)
	if err != nil {
		return client, fmt.Errorf("Error connecting to database: %v", err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return client, fmt.Errorf("Error pinging db: %v", err)
	}

	return client, err
}