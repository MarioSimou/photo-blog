package utils

import (
	"context"
	"errors"
	"log"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
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

type Response struct {
	Status  int64       `json:"status"`
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type Utils struct{}

func (u Utils) LoadDotEnv(fnames ...string) {
	e := godotenv.Load(fnames...)
	if e != nil {
		log.Fatal("Error loading .env file")
	}
}

func (u Utils) HashPassword(pwd string) string {
	hpwd, e := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if e != nil {
		errors.New("Unable to hash password")
	}
	return string(hpwd)
}
