package service

import (
	"database/sql"
	"os"
	"student-report/app/repository"

	"github.com/gofiber/fiber/v2"
)

type StudentService struct {
	repo *repository.StudentRepository
}

func NewStudentService(repo *repository.StudentRepository) *StudentService {
	return &StudentService{repo: repo}
}

func (s *StudentService) GetStudentsService(c *fiber.Ctx) error {
	key := c.Params("key")
	if key != os.Getenv("API_KEY") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "API tidak sesuai",
			"success": false,
		})
	}

	students, err := s.repo.GetStudentsRepository()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal terhubung ke database.",
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    students,
		"message": "Berhasil mendapatkan data Students",
		"success": true,
	})
}

func (s *StudentService) GetStudentByUserID(c *fiber.Ctx) error {
	key := c.Params("key")
	if key != os.Getenv("API_KEY") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "API tidak sesuai",
			"success": false,
		})
	}

	userID := c.Params("userId")
	student, err := s.repo.GetStudentByUserID(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Student tidak ditemukan",
				"success": false,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mendapatkan data student",
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    student,
		"message": "Berhasil mendapatkan data student",
		"success": true,
	})
}

func (s *StudentService) GetStudentsByAdvisorID(c *fiber.Ctx) error {
	key := c.Params("key")
	if key != os.Getenv("API_KEY") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "API tidak sesuai",
			"success": false,
		})
	}

	advisorID := c.Params("advisorId")
	students, err := s.repo.GetStudentsByAdvisorID(advisorID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mendapatkan data students",
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    students,
		"message": "Berhasil mendapatkan data students by advisor",
		"success": true,
	})
}
