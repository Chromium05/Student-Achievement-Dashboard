package service

import (
	"crud-alumni-5/app/model"
	"crud-alumni-5/app/repository"
	"database/sql"
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func GetAllAlumniService(c *fiber.Ctx, db *sql.DB) error {
	key := c.Params("key")
	if key != os.Getenv("API_KEY") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "API tidak sesuai",
			"success": false,
		})
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sortBy", "id")
	order := c.Query("order", "asc")
	search := c.Query("search", "")
	offset := (page - 1) * limit

	// Validasi input
	sortByWhitelist := map[string]bool{"id": true, "nama": true, "email": true, "created_at": true}

	if !sortByWhitelist[sortBy] {
		sortBy = "id"
	}
	if strings.ToLower(order) != "desc" {
		order = "asc"
	}

	alumni, err := repository.GetAllAlumni(db, search, sortBy, order, limit, offset)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal terhubung ke database.",
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    alumni,
		"message": "Berhasil mendapatkan data alumni",
		"success": true,
	})
}

func GetAlumniByIDService(c *fiber.Ctx, db *sql.DB) error {
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
			"error":   "ID tidak valid",
			"success": false,
		})
	}

	alumni, err := repository.GetAlumniByID(db, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal terhubung ke database karena: ",
			"error":   err.Error(),
			"success": false,
		})
	}

	if alumni == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Alumni not found",
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    alumni,
		"message": "Berhasil mendapatkan data alumni",
		"success": true,
	})
}

func PostNewAlumniService(c *fiber.Ctx, db *sql.DB) error {
	key := c.Params("key")
	if key != os.Getenv("API_KEY") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "API tidak sesuai",
			"success": false,
		})
	}

	var input model.CreateAlumni
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Format JSON tidak valid!",
			"success": false,
		})
	}

	alumni, err := repository.PostNewAlumni(db, &input)
	if err != nil {
		// Handle specific errors dengan response yang sesuai
		switch {
		case errors.Is(err, repository.ErrDuplicateNIM):
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error":   "NIM " + input.NIM + " sudah terdaftar! Silakan gunakan NIM yang berbeda.",
				"success": false,
			})
		case errors.Is(err, repository.ErrDuplicateEmail):
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error":   "Email " + input.Email + " sudah terdaftar! Silakan gunakan email yang berbeda.",
				"success": false,
			})
		case errors.Is(err, repository.ErrInvalidInput):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Semua field harus diisi dengan benar!",
				"success": false,
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Gagal menambahkan data: ",
				"error":   err.Error(),
				"success": false,
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data":    alumni,
		"message": "Data alumni berhasil ditambahkan!",
		"success": true,
	})
}

func UpdateAlumniService(c *fiber.Ctx, db *sql.DB) error {
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

	var input model.CreateAlumni
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Format JSON tidak valid!",
			"success": false,
		})
	}

	alumni, err := repository.UpdateAlumni(db, &input, id)
	if err != nil {
		// Error handling input
		switch {
		case errors.Is(err, repository.ErrDuplicateNIM):
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error":   "NIM " + input.NIM + " sudah terdaftar! Silakan gunakan NIM yang berbeda.",
				"success": false,
			})
		case errors.Is(err, repository.ErrDuplicateEmail):
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error":   "Email " + input.Email + " sudah terdaftar! Silakan gunakan email yang berbeda.",
				"success": false,
			})
		case errors.Is(err, repository.ErrInvalidInput):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Semua field harus diisi dengan benar!",
				"success": false,
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Gagal menambahkan data: ",
				"error":   err.Error(),
				"success": false,
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data":    alumni,
		"message": "Data alumni" + alumni.Nama + " berhasil di-update!",
		"success": true,
	})
}

func DeleteAlumniService(c *fiber.Ctx, db *sql.DB) error {
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

	alumni, err := repository.DeleteAlumni(db, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menambahkan data.",
			"error":   err.Error(),
			"success": false,
		})
	}

	if alumni == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Alumni tidak ditemukan.",
			"success": false,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Alumni " + alumni.Nama + " berhasil dihapus!",
		"success": true,
		"data":    alumni,
	})
}
