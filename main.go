package main

import (
	"student-report/config"
	"student-report/database"
	"log"
	"os"
)

func main() {
	config.LoadEnv()
	postgres := database.PostgreConn()
	mongoDB := database.MongoConn()
	_ = mongoDB // Gunakan mongoDB sesuai kebutuhan aplikasi Anda
	app := config.NewApp(postgres)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen(":" + port))
}