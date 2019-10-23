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
	router.GET("/api/users", v.Authorization((v.ValidateRequest(c.GetUsers))))
	router.GET("/api/users/:id", v.Authorization((v.ValidateRequest(c.GetUser))))
	router.POST("/api/users", v.ValidateRequest(v.ValidateCreateUser(c.CreateUser)))
	router.DELETE("/api/users/:id", v.Authorization((v.ValidateRequest(c.DeleteUser))))
	router.PUT("/api/users/:id", v.Authorization((v.ValidateRequest(c.UpdateUser))))
	router.POST("/api/users/login", v.ValidateRequest((c.SignIn)))
	log.Fatal(http.ListenAndServe(":8080", router))
}
