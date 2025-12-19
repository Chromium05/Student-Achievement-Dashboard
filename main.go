package main

import (
	"log"
	"os"
	"student-report/config"
	"student-report/database"
	_ "student-report/docs"
	"student-report/middleware"
	"student-report/route"

	"github.com/gofiber/fiber/v2"
)

// @title Student Report API
// @version 1.0
// @description API untuk sistem pelaporan prestasi mahasiswa
// @description
// @description ## Authentication
// @description API ini menggunakan JWT Bearer token untuk autentikasi.
// @description Setelah login, gunakan token yang diterima dengan menambahkan header:
// @description `Authorization: Bearer {token}`
// @description
// @description ## Roles & Permissions
// @description - Student: Dapat membuat, melihat, update, submit dan delete prestasi sendiri
// @description - Lecturer: Dapat melihat dan verify/reject prestasi mahasiswa bimbingan
// @description - Admin: Full access ke semua endpoints dan user management
// @BasePath /api/v1
// @schemes http
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
