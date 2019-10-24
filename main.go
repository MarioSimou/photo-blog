package main

import (
	"log"
	"net/http"
	"os"

	"github.com/MarioSimou/photo-blog-in-golang/controllers"
	"github.com/MarioSimou/photo-blog-in-golang/utils"
	"github.com/MarioSimou/photo-blog-in-golang/utils/validator"
	"github.com/julienschmidt/httprouter"
)

func main() {
	u := utils.Utils{}
	v := validator.Middleware{&u}

	u.LoadDotEnv()
	mcli := utils.MongoClient{URI: os.Getenv("MONGO_URI"), Database: os.Getenv("DB_NAME")}
	_, e := mcli.Connect()
	if e != nil {
		log.Fatal(e)
	}

	c := controllers.NewController(&mcli, &u)

	router := httprouter.New()
	router.GET("/ping", c.Ping)
	router.GET("/api/v1/users", v.Authorization((v.ValidateRequest(c.GetUsers))))
	router.GET("/api/v1/users/:id", v.Authorization((v.ValidateRequest(c.GetUser))))
	router.POST("/api/v1/users", v.ValidateRequest(v.ValidateCreateUser(c.CreateUser)))
	router.DELETE("/api/v1/users/:id", v.Authorization((v.ValidateRequest(c.DeleteUser))))
	router.PUT("/api/v1/users/:id", v.Authorization((v.ValidateRequest(c.UpdateUser))))
	router.POST("/api/v1/users/signin", v.ValidateRequest(v.ValidateSignIn(c.SignIn)))
	log.Fatal(http.ListenAndServe(":8080", router))
}
