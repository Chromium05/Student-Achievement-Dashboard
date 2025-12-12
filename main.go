package main

import (
	"student-report/config"
	"student-report/database"
	"student-report/route"
	"student-report/middleware"
	"log"
	"os"
	"github.com/gofiber/fiber/v2"
)

func main() {
	config.LoadEnv()
	postgres := database.PostgreConn()
	mongoDB := database.MongoConn()
	
	services := config.InitializeServices(postgres, mongoDB)
	
	app := fiber.New()
	app.Use(middleware.LoggerMiddleware)
	
	route.RegisterRoutes(app, postgres, mongoDB, services)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen(":" + port))
}
