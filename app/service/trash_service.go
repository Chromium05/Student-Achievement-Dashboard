package service

import (
	// "crud-alumni-5/app/model"
	"crud-alumni-5/app/repository"
	"database/sql"
	"os"
	"github.com/gofiber/fiber/v2"
)

func GetAllTrashService(c *fiber.Ctx, db *sql.DB) error {
	key := c.Params("key")
	if key != os.Getenv("API_KEY") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "API tidak sesuai",
			"success": false,
		})
	}

	pekerjaan, err := repository.GetAllTrash(db)

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

func RestoreDataService(c *fiber.Ctx, db *sql.DB) error {
	key := c.Params("key")
	if key != os.Getenv("API_KEY") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "API tidak sesuai",
			"success": false,
		})
	}

	id, err := c.ParamsInt("id")

	// Validasi ID
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Mohon masukkan ID yang sesuai.",
			"success": false,
		})
	}

	pekerjaan, err := repository.RestoreData(db, id)

	if err != nil {
		if err == repository.ErrInvalidInput {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Input tidak boleh kosong.",
				"error":   err.Error(),
				"success": false,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal terhubung ke database.",
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    pekerjaan,
		"message": "Berhasil me-restore data pekerjaan alumni.",
		"success": true,
	})
}

func PermanentDeleteService(c *fiber.Ctx, db *sql.DB) error {
	key := c.Params("key")
	if key != os.Getenv("API_KEY") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "API tidak sesuai",
			"success": false,
		})
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Mohon masukkan ID yang sesuai.",
			"success": false,
		})
	}

	pekerjaan, err := repository.PermanentDelete(db, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal terhubung ke database.",
			"error":   err.Error(),
			"success": false,
		})
	}

	if pekerjaan == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Data pekerjaan alumni tidak ditemukan.",
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    pekerjaan,
		"message": "Berhasil menghapus data pekerjaan alumni",
		"success": true,
	})
}