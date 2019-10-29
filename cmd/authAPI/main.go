package main

import (
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"

	"github.com/MarioSimou/authAPI/internal/controllers"
	"github.com/MarioSimou/authAPI/internal/utils"
	"github.com/MarioSimou/authAPI/internal/utils/middlewares"
)

type App struct {
	Controller  *controllers.Controller
	Utils       *utils.Utils
	Middlewares *middlewares.Middleware
}

func (a *App) Run() {
	m := a.Middlewares
	c := a.Controller

	// routes wrapped within middlewares that check the requests
	getUsers := middlewares.Handler(m.ValidateRequest(m.Authorization(c.GetUsers)))
	getUser := middlewares.Handler(m.ValidateRequest(m.Authorization(c.GetUser)))
	createUser := middlewares.Handler(m.ValidateRequest(m.ValidateCreateUser(c.CreateUser)))
	deleteUser := middlewares.Handler(m.ValidateRequest(m.Authorization(c.DeleteUser)))
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

func main() {
	var app App
	u := utils.Utils{}
	m := middlewares.Middleware{&u}

	u.LoadDotEnv("./configs/.env")
	mcli := u.ConnectDatabase(os.Getenv("MONGO_URI"), os.Getenv("DB_NAME"))
	c := controllers.NewController(mcli, &u)

	app = App{Controller: c, Utils: &u, Middlewares: &m}
	app.Run()
}
