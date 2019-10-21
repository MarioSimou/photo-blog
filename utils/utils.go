package utils

import (
	"context"
	"log"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	URI      string
	Database string
	Client   *mongo.Client
}

func (mcli *MongoClient) Connect() (*mongo.Client, error) {
	client, e := mongo.Connect(context.TODO(), options.Client().ApplyURI(mcli.URI))
	if e != nil {
		return nil, e
	}

	e = client.Ping(context.TODO(), nil)
	if e != nil {
		return nil, e
	}

	mcli.Client = client
	return client, nil
}

func LoadDotEnv() {
	e := godotenv.Load()
	if e != nil {
		log.Fatal("Error loading .env file")
	}
}
