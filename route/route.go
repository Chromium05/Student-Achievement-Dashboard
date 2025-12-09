package route

import (
	"student-report/app/service"
	"student-report/middleware"
	"database/sql"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, db *sql.DB) {
	// Homepage (GET 127.0.0.1:3000)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Halo trainer!")
	})

	app.Post("/:key/v1/auth/login", func(c *fiber.Ctx) error {
		return service.LoginService(c, db)
	})

	// Logout endpoint (POST /:key/v1/auth/logout)
	app.Post("/:key/v1/auth/logout", middleware.AuthRequired(), func(c *fiber.Ctx) error {
		return service.LogoutService(c)
	})

	app.Post("/:key/v1/auth/refresh", func(c *fiber.Ctx) error {
		return service.RefreshTokenService(c, db)
	})

	// Implementasi Middleware
	protected := app.Group("", middleware.AuthRequired())

	// Untuk routing students 
	students := protected.Group("/:key/v1/students")

	// Menampilkan semua data students (GET 127.0.0.1:3000/:key/v1/students)
	students.Get("/", func(c *fiber.Ctx) error {
		return service.GetStudentsService(c, db)
	})
}
