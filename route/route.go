package route

import (
	// "student-report/app/service"
	"student-report/config"
	"student-report/middleware"
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterRoutes(app *fiber.App, db *sql.DB, mongoDB *mongo.Database, services *config.ServiceContainer) {
	// Homepage (GET 127.0.0.1:3000)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Selamat datang di Student Report Dashboard!")
	})

	app.Post("/:key/v1/auth/login", func(c *fiber.Ctx) error {
		return services.AuthService.LoginService(c)
	})

	app.Post("/:key/v1/auth/logout", middleware.AuthRequired(), func(c *fiber.Ctx) error {
		return services.AuthService.LogoutService(c)
	})

	app.Post("/:key/v1/auth/refresh", func(c *fiber.Ctx) error {
		return services.AuthService.RefreshTokenService(c)
	})

	api := app.Group("/:key/v1")

	users := api.Group("/users", middleware.AuthRequired(), middleware.RequirePermission("user:read"))
	
	users.Post("/", middleware.RequirePermission("user:create"), func(c *fiber.Ctx) error {
		return services.UserService.CreateUserService(c)
	})
	users.Get("/", func(c *fiber.Ctx) error {
		return services.UserService.GetAllUsersService(c)
	})
	users.Get("/:id", func(c *fiber.Ctx) error {
		return services.UserService.GetUserByIDService(c)
	})
	users.Put("/:id", middleware.RequirePermission("user:update"), func(c *fiber.Ctx) error {
		return services.UserService.UpdateUserService(c)
	})
	users.Delete("/:id", middleware.RequirePermission("user:delete"), func(c *fiber.Ctx) error {
		return services.UserService.DeleteUserService(c)
	})

	users.Post("/:id/student-profile", middleware.RequirePermission("user:create"), func(c *fiber.Ctx) error {
		return services.UserService.CreateStudentProfileService(c)
	})

	users.Post("/:id/lecturer-profile", middleware.RequirePermission("user:create"), func(c *fiber.Ctx) error {
		return services.UserService.CreateLecturerProfileService(c)
	})

	// Protected routes group
	protected := api.Group("", middleware.AuthRequired())

	students := protected.Group("/students")
	
	students.Get("/", func(c *fiber.Ctx) error {
		return services.StudentService.GetStudentsService(c)
	})

	students.Get("/user/:userId", func(c *fiber.Ctx) error {
		return services.StudentService.GetStudentByUserID(c)
	})

	students.Get("/advisor/:advisorId", func(c *fiber.Ctx) error {
		return services.StudentService.GetStudentsByAdvisorID(c)
	})
}
