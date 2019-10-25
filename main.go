package main

import (
	"log"
	"net/http"
	"os"

	"github.com/MarioSimou/photo-blog-in-golang/controllers"
	"github.com/MarioSimou/photo-blog-in-golang/utils"
	"github.com/MarioSimou/photo-blog-in-golang/utils/middlewares"
	"github.com/julienschmidt/httprouter"
)

func main() {
	u := utils.Utils{}
	m := middlewares.Middleware{&u}

	u.LoadDotEnv()
	mcli := utils.MongoClient{URI: os.Getenv("MONGO_URI"), Database: os.Getenv("DB_NAME")}
	_, e := mcli.Connect()
	if e != nil {
		log.Fatal(e)
	}

	c := controllers.NewController(&mcli, &u)

	// routes wrapped within middlewares that check the requests
	getUsers := middlewares.Handler(m.ValidateRequest(m.Authorization((c.GetUsers))))
	getUser := middlewares.Handler(m.ValidateRequest(m.Authorization((c.GetUser))))
	createUser := middlewares.Handler(m.ValidateRequest(m.ValidateCreateUser(c.CreateUser)))
	deleteUser := middlewares.Handler(m.ValidateRequest(m.Authorization((c.DeleteUser))))
	updateUser := middlewares.Handler(m.ValidateRequest(m.Authorization(c.UpdateUser)))
	signin := middlewares.Handler(m.ValidateRequest(m.ValidateSignIn(c.SignIn)))

	router := httprouter.New()
	router.GET("/ping", c.Ping)
	router.GET("/api/v1/users", getUsers)
	router.GET("/api/v1/users/:id", getUser)
	router.POST("/api/v1/users", createUser)
	router.DELETE("/api/v1/users/:id", deleteUser)
	router.PUT("/api/v1/users/:id", updateUser)
	router.POST("/api/v1/users/signin", signin)
	log.Fatal(http.ListenAndServe(":8080", router))
}
