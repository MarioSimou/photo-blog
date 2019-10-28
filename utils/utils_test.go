package utils

import (
	"fmt"
	"os"
	"projects/users-auth-api/models"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var u Utils

func init() {
	u = Utils{}
}

func TestUtilsLoadDotEnv(t *testing.T) {
	u.LoadDotEnv("../.test.env")
	if env := os.Getenv("GO_ENV"); env != "test" {
		t.Errorf("Should have returned an environment value of %v rather than %v", "test", env)
	}
}

func TestUtilsHashPassword(t *testing.T) {
	if hash := u.HashPassword("12345678"); hash == "" {
		t.Errorf("Should have returned a hahed password that is not empty")
	}
}

func TestGenerateToken(t *testing.T) {
	userId, _ := primitive.ObjectIDFromHex("5db5b5b06507b38887bedc87")
	user := models.User{Id: &userId, Email: "paul@gmail.com"}
	if token, ok := u.GenerateToken(user, "secret", time.Hour); !ok {
		t.Errorf("Should have generated a token rather than %v", token)
	}
}

func TestGenerateTokenWithInvalidUser(t *testing.T) {
	user := models.User{}
	if token, ok := u.GenerateToken(user, "secret", time.Hour); token != nil && ok {
		t.Errorf("Should have returned nil and false for token and boolean flag respectively")
	}
}

func TestVerifyToken(t *testing.T) {
	userId, _ := primitive.ObjectIDFromHex("5db5b5b06507b38887bedc87")
	user := models.User{Id: &userId, Email: "paul@gmail.com"}
	token, _ := u.GenerateToken(user, "secret", time.Hour)
	payload, ok := u.VerifyToken([]byte(token), "secret")
	if !ok {
		t.Errorf("Should have returned a true boolean flag for the token")
	}
	if payload.Email != "paul@gmail.com" {
		t.Errorf("The payload should have included the email %v rather than %v", user.Email, payload.Email)
	}
	if payload.Id.Hex() != user.Id.Hex() {
		t.Errorf("The payload should have included an Id of %v rather than %v", user.Id, payload.Id)
	}
	if !payload.ExpirationTime.Before(time.Now().Add(time.Hour)) {
		t.Errorf("The token should have expired an hour later")
	}
}

func TestVerifyTokenWithInvalidSecret(t *testing.T) {
	userId, _ := primitive.ObjectIDFromHex("5db5b5b06507b38887bedc87")
	user := models.User{Id: &userId, Email: "paul@gmail.com"}
	token, _ := u.GenerateToken(user, "secret", time.Hour)
	if payload, ok := u.VerifyToken([]byte(token), "wrongsecret"); payload != nil || ok {
		fmt.Println(token)
		fmt.Println(ok)
		t.Errorf("Should returned a 'false' flag for an invalid secret")
	}
}

func TestExtractPayload(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NzIzMDYzNDMsImVtYWlsIjoicGF1bEBnbWFpbC5jb20iLCJpZCI6IjVkYjViNWIwNjUwN2IzODg4N2JlZGM4NyJ9.Fhd4bVk2ZRGCb8whtPAbOtFVSXDA-SQCQ9qTE8jvhhU"
	userId, _ := primitive.ObjectIDFromHex("5db5b5b06507b38887bedc87")
	payload := u.ExtractPayload(token)
	if payload.Email != "paul@gmail.com" {
		t.Errorf("The payload should have included the email %v rather than %v", "paul@gmail.com", payload.Email)
	}
	if payload.Id.Hex() != userId.Hex() {
		t.Errorf("The payload should have included an Id of %v rather than %v", userId, payload.Id)
	}
}
