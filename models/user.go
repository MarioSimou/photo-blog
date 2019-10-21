package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id       primitive.ObjectID `json:"id",omitempty bson:"_id,omitempty"`
	Username string             `json:"username,omitempty" bson:"username,omitempty"`
	Email    string             `json:"email,omitempty" bson:"email,omitempty"`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
	Role     string             `json:"role,omitempty" bson:"role,omitempty"`
}

func (u User) Name() string {
	return "user"
}

type Users []User

func (u Users) Name() string {
	return "users"
}

type Collection interface {
	Name() string
}
