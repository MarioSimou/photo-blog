package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

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

func (u *User) ValidateUsername() bool {
	if u.Username == "" {
		return false
	}
	return true
}

func (u *User) ValidateEmail() bool {
	if u.Email == "" {
		return false
	}
	return true
}

func (u *User) ValidatePassword() bool {
	if len(u.Password) < 8 {
		return false
	}
	return true
}
func (u *User) ValidateRole() {
	if u.Role == "" {
		u.Role = "BASIC"
	}
}

func (u *User) ComparePassword(s string) bool {
	e := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(s))
	if e != nil {
		return false
	}
	return true
}

type Users []User

func (u Users) Name() string {
	return "users"
}

type LoginUser struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

func (cu LoginUser) ValidateEmail() bool {
	if cu.Email != "" {
		return true
	}
	return false
}

func (cu LoginUser) ValidatePassword() bool {
	if len(cu.Password) >= 8 {
		return true
	}
	return false
}
