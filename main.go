package main

import (
	"log"
	"net/http"
	"os"

	"./controllers"
	"./utils"
	"github.com/julienschmidt/httprouter"
)

func main() {
	utils.LoadDotEnv()
	mcli := utils.MongoClient{URI: os.Getenv("MONGO_URI"), Database: os.Getenv("DB_NAME")}
	_, e := mcli.Connect()
	if e != nil {
		log.Fatal(e)
	}

	controller := controllers.NewController(&mcli)

	router := httprouter.New()
	router.GET("/ping", controller.Ping)
	router.GET("/api/users", controller.GetUsers)
	router.GET("/api/users/:id", controller.GetUser)
	router.POST("/api/users", controller.CreateUser)
	router.DELETE("/api/users/:id", controller.DeleteUser)
	router.PUT("/api/users/:id", controller.UpdateUser)
	log.Fatal(http.ListenAndServe(":8080", router))
}
