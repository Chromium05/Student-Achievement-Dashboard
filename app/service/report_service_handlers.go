package service

import (
	"context"
	"database/sql"
	"strconv"

	"student-report/app/model"
	"student-report/app/repository"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

// FR-011: Get Achievement Statistics
func GetStatisticsService(c *fiber.Ctx, db *sql.DB, mongoDB *mongo.Database) error {
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(string)
	
	achievementRepo := repository.NewAchievementRepository(mongoDB, db)
	achievementService := NewAchievementService(achievementRepo)
	
	ctx := context.Background()
	var studentIDs []string
	
	// Filter based on role
	if role == "student" {
		// Get own statistics
		studentRepo := repository.NewStudentRepository(db)
		student, err := studentRepo.GetStudentByUserID(userID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Student profile tidak ditemukan",
				"success": false,
			})
		}
		studentIDs = []string{student.ID}
	} else if role == "lecturer" {
		// Get advisees statistics
		studentRepo := repository.NewStudentRepository(db)
		advisees, err := studentRepo.GetStudentsByAdvisorID(userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Gagal mengambil data mahasiswa bimbingan",
				"error":   err.Error(),
				"success": false,
			})
		}
		
		for _, advisee := range advisees {
			studentIDs = append(studentIDs, advisee.ID)
		}
	}
	// else admin: leave studentIDs empty to get all statistics
	
	stats, err := achievementService.GetStatistics(ctx, studentIDs)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil statistik",
			"error":   err.Error(),
			"success": false,
		})
	}
	
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    stats,
		"success": true,
	})
}

// Get Student Report
func GetStudentReportService(c *fiber.Ctx, db *sql.DB, mongoDB *mongo.Database) error {
	studentID := c.Params("id")
	if studentID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Student ID diperlukan",
			"success": false,
		})
	}
	
	achievementRepo := repository.NewAchievementRepository(mongoDB, db)
	achievementService := NewAchievementService(achievementRepo)
	
	ctx := context.Background()
	report, err := achievementService.GetStudentReport(ctx, studentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil report mahasiswa",
			"error":   err.Error(),
			"success": false,
		})
	}
	
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    report,
		"success": true,
	})
}

// Get All Achievements with Filter (Admin)
func GetAllAchievementsWithFilterService(c *fiber.Ctx, db *sql.DB, mongoDB *mongo.Database) error {
	// Parse query parameters
	filter := model.AchievementFilter{
		Page:  1,
		Limit: 10,
	}
	
	if status := c.Query("status"); status != "" {
		filter.Status = &status
	}
	
	if achievementType := c.Query("achievementType"); achievementType != "" {
		filter.AchievementType = &achievementType
	}
	
	if studentID := c.Query("studentId"); studentID != "" {
		filter.StudentID = &studentID
	}
	
	if dateFrom := c.Query("dateFrom"); dateFrom != "" {
		filter.DateFrom = &dateFrom
	}
	
	if dateTo := c.Query("dateTo"); dateTo != "" {
		filter.DateTo = &dateTo
	}
	
	if sortBy := c.Query("sortBy"); sortBy != "" {
		filter.SortBy = &sortBy
	}
	
	if sortOrder := c.Query("sortOrder"); sortOrder != "" {
		filter.SortOrder = &sortOrder
	}
	
	if page := c.Query("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			filter.Page = p
		}
	}
	
	if limit := c.Query("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 {
			filter.Limit = l
		}
	}
	
	achievementRepo := repository.NewAchievementRepository(mongoDB, db)
	achievementService := NewAchievementService(achievementRepo)
	
	ctx := context.Background()
	achievements, total, err := achievementService.GetAllAchievementsWithFilter(ctx, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data achievements",
			"error":   err.Error(),
			"success": false,
		})
	}
	
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":       achievements,
		"total":      total,
		"page":       filter.Page,
		"limit":      filter.Limit,
		"totalPages": (int(total) + filter.Limit - 1) / filter.Limit,
		"success":    true,
	})
}
