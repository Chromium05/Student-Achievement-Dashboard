package service

import (
	"student-report/app/model"
	"student-report/app/repository"
	"database/sql"
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func GetAllPekerjaanAlumniService(c *fiber.Ctx, db *sql.DB) error {
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

	pekerjaan, err := repository.GetAllPekerjaanAlumni(db, search, sortBy, order, limit, offset)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal terhubung ke database.",
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    pekerjaan,
		"message": "Berhasil mendapatkan data pekerjaan alumni",
		"success": true,
	})
}

func GetPekerjaanByIDService(c *fiber.Ctx, db *sql.DB) error {
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
			"message": "ID harus berupa angka.",
			"error":   err.Error(),
			"success": false,
		})
	}

	pekerjaan, err := repository.GetPekerjaanByID(db, id)

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
		"message": "Berhasil mendapatkan data pekerjaan alumni",
		"success": true,
	})
}

func GetPekerjaanByAlumniIDService(c *fiber.Ctx, db *sql.DB) error {
	key := c.Params("key")
	if key != os.Getenv("API_KEY") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "API tidak sesuai",
			"success": false,
		})
	}

	id, err := c.ParamsInt("alumni_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID harus berupa angka.",
			"error":   err.Error(),
			"success": false,
		})
	}

	pekerjaan, err := repository.GetPekerjaanByAlumniID(db, id)

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
		"message": "Berhasil mendapatkan data pekerjaan alumni",
		"success": true,
	})
}

func PostNewPekerjaanAlumniService(c *fiber.Ctx, db *sql.DB) error {
	key := c.Params("key")
	if key != os.Getenv("API_KEY") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "API tidak sesuai",
			"success": false,
		})
	}

	var input model.CreatePekerjaan
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Input tidak valid.",
			"error":   err.Error(),
			"success": false,
		})
	}

	pekerjaan, err := repository.PostNewPekerjaanAlumni(db, &input)

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

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data":    pekerjaan,
		"message": "Berhasil menambahkan data pekerjaan alumni",
		"success": true,
	})
}

func UpdatePekerjaanAlumniService(c *fiber.Ctx, db *sql.DB) error {
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

	var input model.CreatePekerjaan
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Input tidak valid.",
			"error":   err.Error(),
			"success": false,
		})
	}

	pekerjaan, err := repository.UpdatePekerjaanAlumni(db, &input, id)

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
		"message": "Berhasil memperbarui data pekerjaan alumni",
		"success": true,
	})
}

func DeletePekerjaanAlumniService(c *fiber.Ctx, db *sql.DB) error {
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

	pekerjaan, err := repository.DeletePekerjaanAlumni(db, id)
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

func SoftDeletePekerjaanAlumniService(c *fiber.Ctx, db *sql.DB) error {
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

	pekerjaan, err := repository.SoftDeletePekerjaanAlumni(db, id)
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
		"message": "Pekerjaan alumni berhasil 'dihapus'",
		"data":    pekerjaan,
		"success": true,
	})
}
