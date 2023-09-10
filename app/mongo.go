package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func mongoConnect(username, password, host string, port int) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	connectionString := fmt.Sprintf("mongodb://%s:%s@%s:%d", username, password, host, port)
	return mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
}
