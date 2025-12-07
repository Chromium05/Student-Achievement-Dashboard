package service

import (
	// "student-report/app/model"
	"student-report/app/repository"
	"database/sql"
	"os"

	"github.com/gofiber/fiber/v2"
)

func GetStudentsService(c *fiber.Ctx, db *sql.DB) error {
	key := c.Params("key")
	if key != os.Getenv("API_KEY") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "API tidak sesuai",
			"success": false,
		})
	}

	students, err := repository.GetStudentsRepository(db)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal terhubung ke database.",
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    students,
		"message": "Berhasil mendapatkan data Students alumni",
		"success": true,
	})
}