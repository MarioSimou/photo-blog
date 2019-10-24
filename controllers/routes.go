package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/MarioSimou/photo-blog-in-golang/models"
	"github.com/MarioSimou/photo-blog-in-golang/utils"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Controller struct {
	Mongo *mongo.Database
	Utils *utils.Utils
}

func NewController(mcli *utils.MongoClient, utils *utils.Utils) *Controller {
	return &Controller{mcli.Client.Database(mcli.Database), utils}
}

func (c Controller) Ping(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintln(w, "alive")
}

func (c Controller) GetUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var users models.Users
	cur, e := c.Mongo.Collection("users").Find(context.TODO(), bson.M{})
	if e != nil {
		json.NewEncoder(w).Encode(utils.Response{
			Status:  500,
			Success: false,
			Message: "The server was unable to parse the users",
		})
		return
	}

	for cur.Next(context.TODO()) {
		var user models.User
		e := cur.Decode(&user)
		if e != nil {
			log.Fatal(e)
		}
		users = append(users, user)
	}

	if e := cur.Err(); e != nil {
		json.NewEncoder(w).Encode(utils.Response{
			Status:  500,
			Success: false,
			Message: "The server was unable to parse the users",
		})
		return
	}

	cur.Close(context.TODO()) // closes the cursor
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(utils.Response{
		Status:  200,
		Success: true,
		Message: "Successful fetch",
		Data:    users,
	})
}

func (c Controller) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var user models.User
	id := p.ByName("id")
	if id == "" {
		json.NewEncoder(w).Encode(utils.Response{
			Status:  400,
			Success: false,
			Message: "Unable to parse target user id",
		})
		return
	}

	oid, _ := primitive.ObjectIDFromHex(id)
	c.Mongo.Collection("users").FindOne(context.TODO(), bson.M{"_id": oid}).Decode(&user)

	if user.Id == nil {
		json.NewEncoder(w).Encode(utils.Response{
			Status:  404,
			Success: false,
			Message: "The user does not exists",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(utils.Response{
		Status:  200,
		Success: true,
		Message: "Successful fetch",
		Data:    user,
	})
}

func (c Controller) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var body models.User
	var user models.User
	json.NewDecoder(r.Body).Decode(&body)

	result, e := c.Mongo.Collection("users").InsertOne(context.TODO(), body)
	if e != nil {
		json.NewEncoder(w).Encode(utils.Response{
			Status:  400,
			Success: false,
			Message: "Unable to store new user",
		})
		return
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		w.Header().Set("Location", strings.Join([]string{r.URL.Path, oid.Hex()}, "/"))
	}

	c.Mongo.Collection("users").FindOne(context.TODO(), bson.M{"_id": result.InsertedID}).Decode(&user)

	token, ok := c.Utils.GenerateToken(user, os.Getenv("JWT_SECRET"))
	if !ok {
		json.NewEncoder(w).Encode(utils.Response{
			Status:  500,
			Success: true,
			Message: "Unable to generate a token",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(utils.Response{
		Status:  200,
		Success: true,
		Message: "Successful Creation",
		Data:    user,
		Token:   string(token),
	})
}

func (c Controller) DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	if id == "" {
		json.NewEncoder(w).Encode(utils.Response{
			Status:  400,
			Success: false,
			Message: "Unable to parse target user id",
		})
		return
	}

	oid, _ := primitive.ObjectIDFromHex(id)
	result, e := c.Mongo.Collection("users").DeleteOne(context.TODO(), bson.M{"_id": oid})
	if e != nil {
		http.Error(w, e.Error(), 400)
		return
	}

	if result.DeletedCount == 0 {
		json.NewEncoder(w).Encode(utils.Response{
			Status:  404,
			Success: false,
			Message: "User does not exists",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(204)
}

func (c Controller) UpdateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var body interface{}
	var user models.User
	id := p.ByName("id")
	if id == "" {
		json.NewEncoder(w).Encode(utils.Response{
			Status:  400,
			Success: false,
			Message: "Unable to parse target user id",
		})
		return
	}

	json.NewDecoder(r.Body).Decode(&body)
	oid, _ := primitive.ObjectIDFromHex(id)
	result, e := c.Mongo.Collection("users").UpdateOne(context.TODO(), bson.M{"_id": oid}, bson.M{"$set": body})
	if e != nil {
		json.NewEncoder(w).Encode(utils.Response{
			Status:  400,
			Success: false,
			Message: "Unable to update user",
		})
		return
	}

	if result.MatchedCount == 0 {
		json.NewEncoder(w).Encode(utils.Response{
			Status:  404,
			Success: false,
			Message: "User does not exists",
		})
		return
	}

	c.Mongo.Collection("users").FindOne(context.TODO(), bson.M{"_id": oid}).Decode(&user)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(utils.Response{
		Status:  200,
		Success: true,
		Message: "Successful Update",
		Data:    user,
	})
}

func (c Controller) SignIn(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var body models.LoginUser
	var user models.User
	json.NewDecoder(r.Body).Decode(&body)

	fmt.Println(body)
	c.Mongo.Collection("users").FindOne(context.TODO(), bson.M{"email": body.Email}).Decode(&user)

	if user.Id == nil {
		json.NewEncoder(w).Encode(utils.Response{
			Status:  404,
			Success: false,
			Message: "The user does not exists",
		})
		return
	}

	if !user.ComparePassword(body.Password) {
		json.NewEncoder(w).Encode(utils.Response{
			Status:  401,
			Success: false,
			Message: "Invalid password",
		})
		return
	}

	// write logic to sign in a user
	token, ok := c.Utils.GenerateToken(user, os.Getenv("JWT_SECRET"))
	if !ok {
		json.NewEncoder(w).Encode(utils.Response{
			Status:  500,
			Success: false,
			Message: "Unable to generate user token",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(utils.Response{
		Status:  200,
		Success: true,
		Message: "The user has been successfully logged in",
		Data:    user,
		Token:   token,
	})
}
