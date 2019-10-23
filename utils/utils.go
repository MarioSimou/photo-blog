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

type Response struct {
	Status  int64       `json:"status"`
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// type Payload struct {
// 	jwt.Payload
// 	Email string `json:"email,omitempty"`
// }

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

// func (u Utils) GenerateToken(u *models.User, s string) ([]byte, bool) {
// 	hs := jwt.NewHS256([]byte(s))
// 	pl := Payload{
// 		Email: u.Email,
// 		Payload: jwt.Payload{
// 			ExpirationTime: jwt.NumericDate(now.Add())
// 		}
// 	}
// 	if token, e := jwt.Sign(pl, hs); !e {
// 		return token, true
// 	}
// 	return nil, false
// }

// func (u Utils) VerifyToken(t []byte, s string) (*Payload, bool){
// 	var pl Payload
// 	hs := jwt.NewHS256([]byte(s))

// 	if hd, e := jwt.Verify(t,hs, &pl); !e {
// 		fmt.Println(hd)
// 		return &pl, true
// 	}
// 	return nil, false
// }
