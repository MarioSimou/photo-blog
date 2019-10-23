package utils

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/MarioSimou/photo-blog-in-golang/models"
	"github.com/gbrlsnchs/jwt/v3"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type Response struct {
	Status  int64       `json:"status"`
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type Payload struct {
	jwt.Payload
	Email string `json:"email,omitempty"`
}

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

func (u Utils) GenerateToken(user models.User, s string) ([]byte, bool) {
	hs := jwt.NewHS256([]byte(s))
	now := time.Now()
	pl := Payload{
		Email: user.Email,
		Payload: jwt.Payload{
			ExpirationTime: jwt.NumericDate(now.Add(time.Hour * 1)),
		},
	}
	if token, e := jwt.Sign(pl, hs); e == nil {
		return token, true
	}
	return nil, false
}

func (u Utils) VerifyToken(t []byte, s string) (*Payload, bool) {
	var pl Payload
	hs := jwt.NewHS256([]byte(s))
	now := time.Now()
	fmt.Println("verifyingggg")
	expValidator := jwt.ExpirationTimeValidator(now)
	validatePayload := jwt.ValidatePayload(&pl.Payload, expValidator)

	fmt.Println(pl)
	fmt.Println(validatePayload)
	if hd, e := jwt.Verify(t, hs, &pl, validatePayload); e == nil {
		fmt.Println(hd)
		return &pl, true
	}
	return nil, false
}
