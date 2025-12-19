package route

import (
	// "student-report/app/service"
	"database/sql"
	"student-report/app/service"
	"student-report/config"
	"student-report/middleware"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/gofiber/swagger"
)

func RegisterRoutes(app *fiber.App, db *sql.DB, mongoDB *mongo.Database, services *config.ServiceContainer) {
	// Homepage (GET 127.0.0.1:3000)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Selamat datang di Student Report Dashboard!")
	})

	// Swagger Documentation
	app.Get("/swagger/*", swagger.HandlerDefault)

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

	profile := api.Group("/auth/profile", middleware.AuthRequired())

	profile.Get("/", func(c *fiber.Ctx) error {
		return services.AuthService.GetProfileService(c)
	})

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
	
	students.Get("/", middleware.RequirePermission("student:read"), func(c *fiber.Ctx) error {
		return services.StudentService.GetStudentsService(c)
	})

	students.Get("/user/:userId", func(c *fiber.Ctx) error {
		return services.StudentService.GetStudentByUserID(c)
	})

	students.Get("/advisor/:advisorId", func(c *fiber.Ctx) error {
		return services.StudentService.GetStudentsByAdvisorID(c)
	})

	lecturers := protected.Group("/lecturers")

	lecturers.Get("/", func(c *fiber.Ctx) error {
		return services.LecturerService.GetLecturersService(c)
	})

	achievements := protected.Group("/achievements")

	// FR-003: Create Achievement (Mahasiswa)
	achievements.Post("/", middleware.RequirePermission("achievement:create"), func(c *fiber.Ctx) error {
		return service.CreateAchievementService(c, db, mongoDB)
	})

	// Get My Achievements (Mahasiswa)
	achievements.Get("/my", middleware.RequirePermission("achievement:read"), func(c *fiber.Ctx) error {
		return service.GetMyAchievementsService(c, db, mongoDB)
	})

	// FR-006: Get Advisees Achievements (Dosen Wali)
	achievements.Get("/advisees", middleware.RequirePermission("achievement:verify"), func(c *fiber.Ctx) error {
		return service.GetAdviseesAchievementsService(c, db, mongoDB)
	})

	// FR-010: Get All Achievements (Admin)
	achievements.Get("/", middleware.RequirePermission("report:view"), func(c *fiber.Ctx) error {
		return service.GetAllAchievementsWithFilterService(c, db, mongoDB)
	})

	// Get Achievement by ID
	achievements.Get("/:id", middleware.RequirePermission("achievement:read"), func(c *fiber.Ctx) error {
		return service.GetAchievementByIDService(c, db, mongoDB)
	})

	// Update Achievement (Mahasiswa - draft only)
	achievements.Put("/:id", middleware.RequirePermission("achievement:update"), func(c *fiber.Ctx) error {
		return service.UpdateAchievementService(c, db, mongoDB)
	})

	// FR-004: Submit for Verification (Mahasiswa)
	achievements.Post("/:id/submit", middleware.RequirePermission("achievement:submit"), func(c *fiber.Ctx) error {
		return service.SubmitAchievementService(c, db, mongoDB)
	})

	// FR-007, FR-008: Verify/Reject Achievement (Dosen Wali)
	achievements.Post("/:id/verify", middleware.RequirePermission("achievement:verify"), func(c *fiber.Ctx) error {
		return service.VerifyAchievementService(c, db, mongoDB)
	})

	// FR-005: Delete Achievement (Mahasiswa - draft only)
	achievements.Delete("/:id", middleware.RequirePermission("achievement:delete"), func(c *fiber.Ctx) error {
		return service.DeleteAchievementService(c, db, mongoDB)
	})

	// Upload Attachment
	achievements.Post("/:id/attachments", middleware.RequirePermission("achievement:create"), func(c *fiber.Ctx) error {
		return service.UploadAttachmentService(c, db, mongoDB)
	})

	reports := protected.Group("/reports")
	
	// FR-011: Get Statistics (role-based)
	reports.Get("/statistics", func(c *fiber.Ctx) error {
		return service.GetStatisticsService(c, db, mongoDB)
	})
	
	// Get Student Report
	reports.Get("/student/:id", middleware.RequirePermission("report:view"), func(c *fiber.Ctx) error {
		return service.GetStudentReportService(c, db, mongoDB)
	})
}