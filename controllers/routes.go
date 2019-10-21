package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"../models"
	"../utils"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Controller struct {
	Mongo *mongo.Database
}

func NewController(mcli *utils.MongoClient) Controller {
	return Controller{mcli.Client.Database(mcli.Database)}
}

func (c Controller) Ping(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintln(w, "alive")
}

func (c Controller) GetUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var users models.Users
	cur, e := c.Mongo.Collection("users").Find(context.TODO(), bson.M{})
	if e != nil {
		http.Error(w, "The server was unable to parse the users collection", 500)
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
		http.Error(w, "The server was unable to parse the users collection", 500)
		return
	}

	cur.Close(context.TODO()) // closes the cursor
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(users)
}

func (c Controller) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var user models.User
	id := p.ByName("id")
	if id == "" {
		http.Error(w, "A user id has not been provider", 400)
		return
	}

	oid, _ := primitive.ObjectIDFromHex(id)
	c.Mongo.Collection("users").FindOne(context.TODO(), bson.M{"_id": oid}).Decode(&user)

	if user.Id == nil {
		http.Error(w, "User not found", 404)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(user)
}

func (c Controller) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var body models.User
	var user models.User
	json.NewDecoder(r.Body).Decode(&body)

	result, e := c.Mongo.Collection("users").InsertOne(context.TODO(), body)
	if e != nil {
		http.Error(w, e.Error(), 400)
		return
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		w.Header().Set("Location", strings.Join([]string{r.URL.Path, oid.Hex()}, "/"))
	}

	c.Mongo.Collection("users").FindOne(context.TODO(), bson.M{"_id": result.InsertedID}).Decode(&user)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(user)
}

func (c Controller) DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	if id == "" {
		http.Error(w, "A user id has not been provider", 400)
		return
	}

	oid, _ := primitive.ObjectIDFromHex(id)
	result, e := c.Mongo.Collection("users").DeleteOne(context.TODO(), bson.M{"_id": oid})
	if e != nil {
		http.Error(w, e.Error(), 400)
		return
	}

	if result.DeletedCount == 0 {
		http.Error(w, "User not found", 404)
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
		http.Error(w, "A user id has not been provider", 400)
		return
	}

	json.NewDecoder(r.Body).Decode(&body)
	oid, _ := primitive.ObjectIDFromHex(id)
	result, e := c.Mongo.Collection("users").UpdateOne(context.TODO(), bson.M{"_id": oid}, bson.M{"$set": body})
	if e != nil {
		http.Error(w, e.Error(), 400)
		return
	}

	if result.MatchedCount == 0 {
		http.Error(w, "User not found", 404)
		return
	}

	c.Mongo.Collection("users").FindOne(context.TODO(), bson.M{"_id": oid}).Decode(&user)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(user)
}
