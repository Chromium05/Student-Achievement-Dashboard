package main

import (
	"crud-alumni-5/config"
	"crud-alumni-5/database"
	"log"
	"os"
)

func main() {
	config.LoadEnv()
	db := database.ConnectDB()
	app := config.NewApp(db)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen(":" + port))
}
