package service

import (
	"context"
	"database/sql"

	"student-report/app/model"
	"student-report/app/repository"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

// Create Achievement (FR-003) - Mahasiswa
func CreateAchievementService(c *fiber.Ctx, db *sql.DB, mongoDB *mongo.Database) error {
	// Get student_id from context (set by auth middleware)
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "User ID tidak ditemukan",
			"success": false,
		})
	}

	// Get student record
	studentRepo := repository.NewStudentRepository(db)
	student, err := studentRepo.GetStudentByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Student profile tidak ditemukan",
			"success": false,
		})
	}

	var req model.CreateAchievementRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
			"success": false,
		})
	}

	if req.AchievementType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Achievement type wajib diisi",
			"success": false,
		})
	}
	
	if req.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Title wajib diisi",
			"success": false,
		})
	}

	achievementRepo := repository.NewAchievementRepository(mongoDB, db)
	achievementService := NewAchievementService(achievementRepo)

	ctx := context.Background()
	achievement, err := achievementService.CreateAchievement(ctx, student.ID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal membuat achievement",
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data":    achievement,
		"message": "Achievement berhasil dibuat dengan status draft",
		"success": true,
	})
}

// Get Achievement by ID
func GetAchievementByIDService(c *fiber.Ctx, db *sql.DB, mongoDB *mongo.Database) error {
	achievementID := c.Params("id")
	if achievementID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Achievement ID diperlukan",
			"success": false,
		})
	}

	achievementRepo := repository.NewAchievementRepository(mongoDB, db)
	achievementService := NewAchievementService(achievementRepo)

	ctx := context.Background()
	achievement, err := achievementService.GetAchievementByID(ctx, achievementID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Achievement tidak ditemukan",
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    achievement,
		"success": true,
	})
}

// Update Achievement - Mahasiswa (draft only)
func UpdateAchievementService(c *fiber.Ctx, db *sql.DB, mongoDB *mongo.Database) error {
	achievementID := c.Params("id")
	userID := c.Locals("user_id").(string)

	// Get student record
	studentRepo := repository.NewStudentRepository(db)
	student, err := studentRepo.GetStudentByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Student profile tidak ditemukan",
			"success": false,
		})
	}

	var req model.UpdateAchievementRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"success": false,
		})
	}

	achievementRepo := repository.NewAchievementRepository(mongoDB, db)
	achievementService := NewAchievementService(achievementRepo)

	ctx := context.Background()
	achievement, err := achievementService.UpdateAchievement(ctx, achievementID, student.ID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Gagal update achievement",
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    achievement,
		"message": "Achievement berhasil diupdate",
		"success": true,
	})
}

// FR-004: Submit Achievement for Verification
func SubmitAchievementService(c *fiber.Ctx, db *sql.DB, mongoDB *mongo.Database) error {
	achievementID := c.Params("id")
	userID := c.Locals("user_id").(string)

	// Get student record
	studentRepo := repository.NewStudentRepository(db)
	student, err := studentRepo.GetStudentByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Student profile tidak ditemukan",
			"success": false,
		})
	}

	achievementRepo := repository.NewAchievementRepository(mongoDB, db)
	achievementService := NewAchievementService(achievementRepo)

	ctx := context.Background()
	achievement, err := achievementService.SubmitForVerification(ctx, achievementID, student.ID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Gagal submit achievement",
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    achievement,
		"message": "Achievement berhasil disubmit untuk verifikasi",
		"success": true,
	})
}

// FR-005: Delete Achievement (draft only)
func DeleteAchievementService(c *fiber.Ctx, db *sql.DB, mongoDB *mongo.Database) error {
	achievementID := c.Params("id")
	userID := c.Locals("user_id").(string)

	// Get student record
	studentRepo := repository.NewStudentRepository(db)
	student, err := studentRepo.GetStudentByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Student profile tidak ditemukan",
			"success": false,
		})
	}

	achievementRepo := repository.NewAchievementRepository(mongoDB, db)
	achievementService := NewAchievementService(achievementRepo)

	ctx := context.Background()
	err = achievementService.DeleteAchievement(ctx, achievementID, student.ID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Gagal hapus achievement",
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Achievement berhasil dihapus",
		"success": true,
	})
}

// FR-007, FR-008: Verify/Reject Achievement (Dosen Wali)
func VerifyAchievementService(c *fiber.Ctx, db *sql.DB, mongoDB *mongo.Database) error {
	achievementID := c.Params("id")
	userID := c.Locals("user_id").(string)

	var req model.VerifyAchievementRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"success": false,
		})
	}

	if req.Action != "verify" && req.Action != "reject" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Action harus 'verify' atau 'reject'",
			"success": false,
		})
	}

	achievementRepo := repository.NewAchievementRepository(mongoDB, db)
	achievementService := NewAchievementService(achievementRepo)

	ctx := context.Background()
	achievement, err := achievementService.VerifyAchievement(ctx, achievementID, userID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Gagal memproses verifikasi",
			"error":   err.Error(),
			"success": false,
		})
	}

	message := "Achievement berhasil diverifikasi"
	if req.Action == "reject" {
		message = "Achievement berhasil ditolak"
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    achievement,
		"message": message,
		"success": true,
	})
}

// Get My Achievements (Mahasiswa)
func GetMyAchievementsService(c *fiber.Ctx, db *sql.DB, mongoDB *mongo.Database) error {
	userID := c.Locals("user_id").(string)

	// Get student record
	studentRepo := repository.NewStudentRepository(db)
	student, err := studentRepo.GetStudentByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Student profile tidak ditemukan",
			"success": false,
		})
	}

	achievementRepo := repository.NewAchievementRepository(mongoDB, db)
	achievementService := NewAchievementService(achievementRepo)

	ctx := context.Background()
	achievements, err := achievementService.GetAchievementsByStudentID(ctx, student.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data achievements",
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    achievements,
		"total":   len(achievements),
		"success": true,
	})
}

// FR-006: Get Achievements of Advisees (Dosen Wali)
func GetAdviseesAchievementsService(c *fiber.Ctx, db *sql.DB, mongoDB *mongo.Database) error {
	userID := c.Locals("user_id").(string)

	// Get advisees (students under this lecturer)
	studentRepo := repository.NewStudentRepository(db)
	advisees, err := studentRepo.GetStudentsByAdvisorID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data mahasiswa bimbingan",
			"error":   err.Error(),
			"success": false,
		})
	}

	if len(advisees) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"data":    []interface{}{},
			"total":   0,
			"message": "Tidak ada mahasiswa bimbingan",
			"success": true,
		})
	}

	// Get student IDs
	var studentIDs []string
	for _, advisee := range advisees {
		studentIDs = append(studentIDs, advisee.ID)
	}

	achievementRepo := repository.NewAchievementRepository(mongoDB, db)
	achievementService := NewAchievementService(achievementRepo)

	ctx := context.Background()
	achievements, err := achievementService.GetAchievementsForAdvisees(ctx, studentIDs)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data achievements",
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    achievements,
		"total":   len(achievements),
		"success": true,
	})
}

// FR-010: Get All Achievements (Admin)
func GetAllAchievementsService(c *fiber.Ctx, db *sql.DB, mongoDB *mongo.Database) error {
	achievementRepo := repository.NewAchievementRepository(mongoDB, db)
	achievementService := NewAchievementService(achievementRepo)

	ctx := context.Background()
	achievements, err := achievementService.GetAllAchievements(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data achievements",
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    achievements,
		"total":   len(achievements),
		"success": true,
	})
}

// Upload Attachment
func UploadAttachmentService(c *fiber.Ctx, db *sql.DB, mongoDB *mongo.Database) error {
	achievementID := c.Params("id")
	userID := c.Locals("user_id").(string)

	// Get student record
	studentRepo := repository.NewStudentRepository(db)
	student, err := studentRepo.GetStudentByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Student profile tidak ditemukan",
			"success": false,
		})
	}

	var req model.UploadAttachmentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"success": false,
		})
	}

	if req.FileName == "" || req.FileURL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "File name dan URL wajib diisi",
			"success": false,
		})
	}

	achievementRepo := repository.NewAchievementRepository(mongoDB, db)
	achievementService := NewAchievementService(achievementRepo)

	ctx := context.Background()
	err = achievementService.UploadAttachment(ctx, achievementID, student.ID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Gagal upload attachment",
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Attachment berhasil diupload",
		"success": true,
	})
}
