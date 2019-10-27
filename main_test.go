package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"os"
	"projects/users-auth-api/controllers"
	"projects/users-auth-api/models"
	"projects/users-auth-api/utils"
	"strings"
	"testing"
)

var c *controllers.Controller

func mockData(mcli *utils.MongoClient) {
	usersCollection := mcli.Client.Database(mcli.Database).Collection("users")
	e := usersCollection.Drop(context.TODO())
	if e != nil {
		log.Fatal("Unable to drop Users collection")
	}

	john := models.User{Username: "paul", Email: "paul@gmail.com", Password: "$2a$04$DzlgE3dAEEynd4Ed9z0oY.MafLBCoZl815bXXeOjekaZztjwDLcdm", Role: "BASIC"}
	paul := models.User{Username: "john", Email: "john@gmail.com", Password: "$2a$04$DzlgE3dAEEynd4Ed9z0oY.MafLBCoZl815bXXeOjekaZztjwDLcdm", Role: "BASIC"}
	users := []interface{}{john, paul}

	result, _ := usersCollection.InsertMany(context.TODO(), users)
	if result.InsertedIDs == nil {
		log.Fatal("Data has not being loaded to the database")
	}
}

func init() {
	u := utils.Utils{}
	u.LoadDotEnv("./.test.env")
	mcli := u.ConnectDatabase(os.Getenv("MONGO_URI"), os.Getenv("DB_NAME"))
	c = controllers.NewController(mcli, &u)
	mockData(mcli)
}

func TestAddition(t *testing.T) {
	if 1+1 != 2 {
		t.Error("eRRORRRR")
	}
}

func TestGetUsers(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/v1/users", nil)
	c.GetUsers(w, r, nil)

	res := w.Result()
	if code := res.StatusCode; code != 200 {
		t.Errorf("Should return a status code of %d rather than %d", 200, code)
	}

	body, _ := ioutil.ReadAll(res.Body)
	expected := `{"status":200,"success":true,"message":"Successful fetch","data":[{"username":"paul","email":"paul@gmail.com","role":"BASIC"},{"username":"john","email":"john@gmail.com","role":"BASIC"}]}`
	if b := strings.TrimRight(string(body), "\n"); b != expected {
		t.Errorf("Should return a body of %v rather than %v", expected, b)
	}
}
