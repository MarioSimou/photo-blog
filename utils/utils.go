// Package utils contains code that is used for differnet utilities
package utils

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"

	"projects/users-auth-api/models"

	"github.com/gbrlsnchs/jwt/v3"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// Payload is a custom type used to map a JWT token
type Payload struct {
	jwt.Payload
	Email string              `json:"email,omitempty"`
	Id    *primitive.ObjectID `json:"id,omitempty"`
}

// MongoClient is a custom type used to map a TCP connection to Mongodb
type MongoClient struct {
	URI      string
	Database string
	Client   *mongo.Client
}

// Connect is a method that initiates a tcp connection to mongodb, returning the result of the operation
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

// Utils is a custom type used to represent a utilities object
type Utils struct{}

// LoadDotEnv is used to load the environment variables
func (u Utils) LoadDotEnv(fnames ...string) {
	e := godotenv.Load(fnames...)
	if e != nil {
		log.Fatal("Error loading .env file")
	}
}

func (u Utils) ConnectDatabase(uri string, dbName string) *MongoClient {
	mcli := MongoClient{URI: uri, Database: dbName}
	_, e := mcli.Connect()
	if e != nil {
		log.Fatal(e)
	}
	return &mcli
}

// HashPassword is used to hash a given password
func (u Utils) HashPassword(pwd string) string {
	hpwd, e := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if e != nil {
		errors.New("Unable to hash password")
	}
	return string(hpwd)
}

// GenerateToken is used to generate a JWT token
func (u Utils) GenerateToken(user models.User, s string) ([]byte, bool) {
	hs := jwt.NewHS256([]byte(s))
	now := time.Now()
	pl := Payload{
		Email: user.Email,
		Id:    user.Id,
		Payload: jwt.Payload{
			ExpirationTime: jwt.NumericDate(now.Add(time.Hour * 1)),
		},
	}
	if token, e := jwt.Sign(pl, hs); e == nil {
		return token, true
	}
	return nil, false
}

// VerifyToken is used to verify a JWT token
func (u Utils) VerifyToken(t []byte, s string) (*Payload, bool) {
	var pl Payload
	hs := jwt.NewHS256([]byte(s))
	now := time.Now()
	expValidator := jwt.ExpirationTimeValidator(now)
	validatePayload := jwt.ValidatePayload(&pl.Payload, expValidator)

	if _, e := jwt.Verify(t, hs, &pl, validatePayload); e == nil {
		return &pl, true
	}
	return nil, false
}

// ExtractPayload is used to extract the payload of a JWT token (header.payload.signature)
func (u Utils) ExtractPayload(t string) *Payload {
	var p Payload

	parts := strings.Split(t, ".")
	dec, _ := base64.RawStdEncoding.DecodeString(parts[1])
	json.Unmarshal(dec, &p)

	return &p
}
