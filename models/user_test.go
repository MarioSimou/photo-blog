package models

import (
	"log"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var user User
var luser LoginUser

func init() {
	userID, e := primitive.ObjectIDFromHex("5db5b5b06507b38887bedc88")
	if e != nil {
		log.Fatal("Error while generating user Id")
	}
	user = User{
		Id:       &userID,
		Username: "paul",
		Email:    "paul@gmail.com",
		Password: "$2a$04$DzlgE3dAEEynd4Ed9z0oY.MafLBCoZl815bXXeOjekaZztjwDLcdm",
		Role:     "BASIC",
	}

	luser = LoginUser{Email: "paul@gmail.com", Password: "12345678"}
}

func TestUserName(t *testing.T) {
	if n := user.Name(); n != "user" {
		t.Errorf("The method Name() of user model should have returned 'user' rather than %v", n)
	}
}

func TestUserValidateUsername(t *testing.T) {
	if b := user.ValidateUsername(); !b {
		t.Errorf("The method ValidateUsername() for the user should have returned 'true' rather than %v", b)
	}
}

func TestNullUserValidateUsername(t *testing.T) {
	u := User{}
	if b := u.ValidateUsername(); b {
		t.Errorf("The method ValidateUsername() for an empty username should have returned 'false' rather than %v", b)
	}
}

func TestUserValidateEmail(t *testing.T) {
	if b := user.ValidateEmail(); !b {
		t.Errorf("The method ValidateEmail() for the user should have returned 'true' rather than %v", b)
	}
}

func TestNullUserValidateEmail(t *testing.T) {
	u := User{}
	if b := u.ValidateEmail(); b {
		t.Errorf("The method ValidateEmail() for an empty email should have returned 'false' rather than %v", b)
	}
}

func TestUserValidatePassword(t *testing.T) {
	if b := user.ValidatePassword(); !b {
		t.Errorf("The method ValidatePassword() for the user should have returned 'true' rather than %v", b)
	}
}

func TestNullUserValidatePassword(t *testing.T) {
	if b := user.ValidatePassword(); !b {
		t.Errorf("The method ValidatePassword() for the user should have returned 'false' rather than %v", b)
	}
}

func TestUserValidateRole(t *testing.T) {
	u := User{Role: "BASIC"}
	if u.ValidateRole(); u.Role != "BASIC" {
		t.Errorf("Should have returned a role of BASIC rather than %v", u.Role)
	}
}

func TestNullUserValidateRole(t *testing.T) {
	u := User{}
	if u.ValidateRole(); u.Role != "BASIC" {
		t.Errorf("Should have returned a role of BASIC rather than %v", u.Role)
	}
}

func TestAdminUserValidateRole(t *testing.T) {
	u := User{Role: "ADMIN"}
	if u.ValidateRole(); u.Role != "ADMIN" {
		t.Errorf("Should have returned a role of ADMIN rather than %v", u.Role)
	}
}

func TestUserComparePassword(t *testing.T) {
	if b := user.ComparePassword("12345678"); !b {
		t.Errorf("Should have returned 'true' for a valid password rather than %v", b)
	}
}

func TestInvalidUserComparePassword(t *testing.T) {
	if b := user.ComparePassword("inalidpassword"); b {
		t.Errorf("Should have returned 'false' for an invalid password rather than %v", b)
	}
}

func TestMapToSecureUser(t *testing.T) {
	su := user.MapToSecureUser()
	if su.Id != user.Id {
		t.Errorf("Should have returned a user id of %v rather than %v", user.Id, su.Id)
	}
	if su.Email != user.Email {
		t.Errorf("Should have returned an email of %v rather than %v", user.Email, su.Email)
	}
	if su.Username != user.Username {
		t.Errorf("Should have returned a username of %v rather than %v", user.Username, su.Username)
	}
	if su.Role != user.Role {
		t.Errorf("Should have returned a role of %v rather than %v", user.Role, su.Role)
	}
}

func TestLoginUserEmail(t *testing.T) {
	if b := luser.ValidateEmail(); !b {
		t.Errorf("Should have returned 'true' rather than %v", b)
	}
}

func TestLoginUserPassword(t *testing.T) {
	if b := luser.ValidatePassword(); !b {
		t.Errorf("Should have returned 'true' rather than %v", b)
	}
}
