package config

import (
	"student-report/middleware"
	"student-report/route"
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

func NewApp(db *sql.DB) *fiber.App {
	app := fiber.New()
	app.Use(middleware.LoggerMiddleware)
	route.RegisterRoutes(app, db)
	return app
}
