package service

import (
	"crud-alumni-5/app/repository"
	"database/sql"
	"os"

	"github.com/gofiber/fiber/v2"
)

func GetAllEmployedAlumniService(c *fiber.Ctx, db *sql.DB) error {
	key := c.Params("key")
	if key != os.Getenv("API_KEY") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "API tidak sesuai",
			"success": false,
		})
	}

	pekerjaan, err := repository.GetAllEmployedAlumni(db)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal terhubung ke database.",
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"count":   len(pekerjaan),
		"data":    pekerjaan,
		"message": "Berhasil mendapatkan data pekerjaan alumni",
		"success": true,
	})
}

func GetEmployedAlumniLessThreeYearsService(c *fiber.Ctx, db *sql.DB) error {
	key := c.Params("key")
	if key != os.Getenv("API_KEY") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "API tidak sesuai",
			"success": false,
		})
	}

	pekerjaan, err := repository.GetEmployedAlumniLessThreeYears(db)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal terhubung ke database.",
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"count":   len(pekerjaan),
		"data":    pekerjaan,
		"message": "Berhasil mendapatkan data pekerjaan alumni",
		"success": true,
	})
}
