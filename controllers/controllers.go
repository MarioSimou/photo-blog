// Package controllers represents a composite object with functions used to handle HTTP requests.
package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"projects/users-auth-api/models"
	"projects/users-auth-api/utils"
	"projects/users-auth-api/utils/httpcodes"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Controller custom type
type Controller struct {
	Mongo *mongo.Database
	Utils *utils.Utils
}

// NewController is a function used return an instance of Controller type
func NewController(mcli *utils.MongoClient, utils *utils.Utils) *Controller {
	return &Controller{mcli.Client.Database(mcli.Database), utils}
}

// Ping checks the connection of the API
func (c Controller) Ping(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(200)
	fmt.Fprintln(w, "alive")
}

// GetUsers returns the whole collection of users within the database
func (c Controller) GetUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params, other ...interface{}) {
	var users []models.SecureUser

	cur, e := c.Mongo.Collection("users").Find(context.TODO(), bson.M{})
	if e != nil {
		httpcodes.ResponseError(w, httpcodes.Representation{Message: "The server was unable to parse the users"}.InternalServerError())
		return
	}

	for cur.Next(context.TODO()) {
		var user models.SecureUser

		e := cur.Decode(&user)
		if e != nil {
			log.Fatal(e)
		}
		users = append(users, user)
	}

	if e := cur.Err(); e != nil {
		httpcodes.ResponseError(w, httpcodes.Representation{Message: "The server was unable to parse the users"}.InternalServerError())
		return
	}

	cur.Close(context.TODO()) // closes the cursor
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(httpcodes.Representation{Message: "Successful fetch", Data: users}.Ok())
}

// GetUser is used to return a single user, who is identified based on his/her id
func (c Controller) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params, other ...interface{}) {
	var user models.User
	payload := other[0].(*utils.Payload)
	id := p.ByName("id")

	if id == "" {
		httpcodes.ResponseError(w, httpcodes.Representation{Message: "Invalid target resource"}.BadRequest())
		return
	}

	oid, _ := primitive.ObjectIDFromHex(id)
	c.Mongo.Collection("users").FindOne(context.TODO(), bson.M{"_id": oid}).Decode(&user)

	if user.Id == nil {
		httpcodes.ResponseError(w, httpcodes.Representation{Message: "User does not exists"}.NotFound())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	if payload.Id.Hex() == id {
		json.NewEncoder(w).Encode(httpcodes.Representation{Message: "Successful fetch", Data: user}.Ok())
	} else {
		json.NewEncoder(w).Encode(httpcodes.Representation{Message: "Successful fetch", Data: user.MapToSecureUser()}.Ok())
	}
}

// CreateUser is used to store a user within the database
func (c Controller) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params, other ...interface{}) {
	var body models.User
	var user models.User
	json.NewDecoder(r.Body).Decode(&body)

	result, e := c.Mongo.Collection("users").InsertOne(context.TODO(), body)
	if e != nil {
		httpcodes.ResponseError(w, httpcodes.Representation{Message: "The db was unable to store the user"}.InternalServerError())
		return
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		w.Header().Set("Location", strings.Join([]string{r.URL.Path, oid.Hex()}, "/"))
	}

	c.Mongo.Collection("users").FindOne(context.TODO(), bson.M{"_id": result.InsertedID}).Decode(&user)

	token, ok := c.Utils.GenerateToken(user, os.Getenv("JWT_SECRET"))
	if !ok {
		httpcodes.ResponseError(w, httpcodes.Representation{Message: "Unable to generate a token"}.InternalServerError())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(httpcodes.Representation{Message: "Successful creation", Data: user, Token: string(token)}.Created())
}

// DeleteUser is used to delete a user from the database
func (c Controller) DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params, other ...interface{}) {
	payload := other[0].(*utils.Payload)
	id := p.ByName("id")
	if id == "" {
		httpcodes.ResponseError(w, httpcodes.Representation{Message: "Invalid target resource"}.BadRequest())
		return
	}

	if payload.Id.Hex() != id {
		httpcodes.ResponseError(w, httpcodes.Representation{Message: "Invalid operation for the existing user"}.Forbidden())
		return
	}

	oid, _ := primitive.ObjectIDFromHex(id)
	_, e := c.Mongo.Collection("users").DeleteOne(context.TODO(), bson.M{"_id": oid})
	if e != nil {
		httpcodes.ResponseError(w, httpcodes.Representation{Message: "The db was unable to delete the user"}.InternalServerError())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(204)
}

// UpdateUser is used update a document within the users collection
func (c Controller) UpdateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params, other ...interface{}) {
	var body interface{}
	var user models.User
	id := p.ByName("id")
	payload := other[0].(*utils.Payload)

	if id == "" {
		httpcodes.ResponseError(w, httpcodes.Representation{Message: "Invalid target resource"}.BadRequest())
		return
	}
	if payload.Id.Hex() != id {
		httpcodes.ResponseError(w, httpcodes.Representation{Message: "Invalid operation for the existing user"}.Forbidden())
		return
	}

	json.NewDecoder(r.Body).Decode(&body)
	oid, _ := primitive.ObjectIDFromHex(id)
	_, e := c.Mongo.Collection("users").UpdateOne(context.TODO(), bson.M{"_id": oid}, bson.M{"$set": body})
	if e != nil {
		httpcodes.ResponseError(w, httpcodes.Representation{Message: "The db was unable to update the user"}.InternalServerError())
		return
	}

	c.Mongo.Collection("users").FindOne(context.TODO(), bson.M{"_id": oid}).Decode(&user)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(httpcodes.Representation{Message: "Successful update", Data: user}.Ok())
}

// SignIn is used to login a user in the service
func (c Controller) SignIn(w http.ResponseWriter, r *http.Request, _ httprouter.Params, other ...interface{}) {
	var body models.LoginUser
	var user models.User
	json.NewDecoder(r.Body).Decode(&body)

	c.Mongo.Collection("users").FindOne(context.TODO(), bson.M{"email": body.Email}).Decode(&user)

	if user.Id == nil {
		httpcodes.ResponseError(w, httpcodes.Representation{Message: "The user does not exists"}.NotFound())
		return
	}

	if !user.ComparePassword(body.Password) {
		httpcodes.ResponseError(w, httpcodes.Representation{Message: "Invalid Password"}.Unauthorized())
		return
	}

	// write logic to sign in a user
	token, ok := c.Utils.GenerateToken(user, os.Getenv("JWT_SECRET"))
	fmt.Println(string(token))
	if !ok {
		httpcodes.ResponseError(w, httpcodes.Representation{Message: "Unable to generate user token"}.InternalServerError())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(httpcodes.Representation{Message: "Successful login", Data: user, Token: string(token)}.Ok())
}
