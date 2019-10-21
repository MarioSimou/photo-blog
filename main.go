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

	router := httprouter.New()
	router.GET("/", controllers.Index)
	router.GET("/api/users", controllers.GetUsers)
	router.GET("/api/users/:id", controllers.GetUser)
	router.POST("/api/users", controllers.CreateUser)

	log.Fatal(http.ListenAndServe(":8080", router))
}
