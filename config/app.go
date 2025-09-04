package config

import (
    "database/sql"
    "github.com/gofiber/fiber/v2"
    "tugasminggu3/middleware"
    "tugasminggu3/route"
)

func NewApp(db *sql.DB) *fiber.App {
    app := fiber.New()
    app.Use(middleware.LoggerMiddleware)
    route.RegisterRoutes(app, db)
    return app
}