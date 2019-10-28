// Package models contains the collections stored within the database
package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// SecureUser is a custom type used to display none-private information of a user
type SecureUser struct {
	Id       *primitive.ObjectID `json:"id,omitempty", bson:"_id,omitempty"`
	Username string              `json:"username,omitempty" bson:"username"`
	Email    string              `json:"email,omitempty" bson:"email"`
	Role     string              `json:"role,omitempty" bson:"role"`
}

// Name returns the name of the document
func (su SecureUser) Name() string {
	return "secureUser"
}

// User is a custom type used to represent a document in the users collection
type User struct {
	Id       *primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Username string              `json:"username,omitempty" bson:"username"`
	Email    string              `json:"email,omitempty" bson:"email"`
	Password string              `json:"password,omitempty" bson:"password"`
	Role     string              `json:"role,omitempty" bson:"role"`
}

// Name returns the name of the document
func (u User) Name() string {
	return "user"
}

// ValidateUsername is a method used to validate the username of the document
func (u *User) ValidateUsername() bool {
	if u.Username == "" {
		return false
	}
	return true
}

// ValidateEmail is a method used to validate the email of the document
func (u *User) ValidateEmail() bool {
	if u.Email == "" {
		return false
	}
	return true
}

// ValidatePassword is a method used to validate the password of a document
func (u *User) ValidatePassword() bool {
	if len(u.Password) < 8 {
		return false
	}
	return true
}

// ValidateRole is a method used to validate the role of a document
func (u *User) ValidateRole() {
	switch role := u.Role; role {
	case "ADMIN":
	case "BASIC":
	default:
		u.Role = "BASIC"
	}
}

// ComparePassword is a method used to validate a given password with the hashed password of a user
func (u *User) ComparePassword(s string) bool {
	if e := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(s)); e != nil {
		return false
	}
	return true
}

// MapToSecureUser is a method used to convert a User type to SecureUser
func (u *User) MapToSecureUser() *SecureUser {
	return &SecureUser{Id: u.Id, Username: u.Username, Email: u.Email, Role: u.Role}
}

// Users is a custom type used to represent a collection of users (collection)
type Users []User

// Name is a method user to return the name of the collection
func (u Users) Name() string {
	return "users"
}

// LoginUser is custom type used to map the credentials of a user when he/she logs in
type LoginUser struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

// ValidateEmail is a method used to validate the email of a LoginUser
func (cu LoginUser) ValidateEmail() bool {
	tu := User{Email: cu.Email}
	return tu.ValidateEmail()
}

// ValidatePassword is a method used to validate the password of a LoginUser
func (cu LoginUser) ValidatePassword() bool {
	tu := User{Password: cu.Password}
	return tu.ValidatePassword()
}
