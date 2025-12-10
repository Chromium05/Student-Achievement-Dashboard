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

	// User Management (Admin only)
	api := app.Group("/api/v1")
	users := api.Group("/users", middleware.AuthRequired(), middleware.RequirePermission("user:read"))
	users.Post("/", middleware.RequirePermission("user:create"), func(c *fiber.Ctx) error {
		return service.CreateUserService(c, db)
	})
	users.Get("/", func(c *fiber.Ctx) error {
		return service.GetAllUsersService(c, db)
	})
	users.Get("/:id", func(c *fiber.Ctx) error {
		return service.GetUserByIDService(c, db)
	})
	users.Put("/:id", middleware.RequirePermission("user:update"), func(c *fiber.Ctx) error {
		return service.UpdateUserService(c, db)
	})
	users.Delete("/:id", middleware.RequirePermission("user:delete"), func(c *fiber.Ctx) error {
		return service.DeleteUserService(c, db)
	})

	// Student Profile Management (Admin only)
	users.Post("/:id/student-profile", middleware.RequirePermission("user:create"), func(c *fiber.Ctx) error {
		return service.CreateStudentProfileService(c, db)
	})

	// Lecturer Profile Management (Admin only)
	users.Post("/:id/lecturer-profile", middleware.RequirePermission("user:create"), func(c *fiber.Ctx) error {
		return service.CreateLecturerProfileService(c, db)
	})

	app.Post("/api/v1/auth/login", func(c *fiber.Ctx) error {
		return service.LoginService(c, db)
	})

	// Logout endpoint (POST /api/v1/auth/logout)
	app.Post("/api/v1/auth/logout", middleware.AuthRequired(), func(c *fiber.Ctx) error {
		return service.LogoutService(c)
	})

	app.Post("/api/v1/auth/refresh", func(c *fiber.Ctx) error {
		return service.RefreshTokenService(c, db)
	})

	// Implementasi Middleware
	protected := app.Group("", middleware.AuthRequired())

	// Untuk routing students 
	students := protected.Group("/:key/v1/students")

	// Menampilkan semua data students (GET 127.0.0.1:3000/api/v1/students)
	students.Get("/", func(c *fiber.Ctx) error {
		return service.GetStudentsService(c, db)
	})
}
