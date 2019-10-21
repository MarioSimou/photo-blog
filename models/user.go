package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id       *primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Username string              `json:"username,omitempty" bson:"username"`
	Email    string              `json:"email,omitempty" bson:"email"`
	Password string              `json:"password,omitempty" bson:"password"`
	Role     string              `json:"role,omitempty" bson:"role"`
}

func (u User) Name() string {
	return "user"
}

type Users []User

func (u Users) Name() string {
	return "users"
}
