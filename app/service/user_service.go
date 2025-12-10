package service

import (
	"student-report/app/model"
	"student-report/app/repository"
	// "student-report/utils"
	"database/sql"
	"errors"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// CreateUserService membuat user baru (Admin only)
func CreateUserService(c *fiber.Ctx, db *sql.DB) error {
	var req model.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
			"success": false,
		})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal membuat password",
			"error":   err.Error(),
			"success": false,
		})
	}

	// Create user
	user, err := repository.CreateUser(db, req, string(hashedPassword))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal membuat user",
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data":    user,
		"message": "User berhasil dibuat",
		"success": true,
	})
}

// GetAllUsersService mengambil semua users (Admin only)
func GetAllUsersService(c *fiber.Ctx, db *sql.DB) error {
	users, err := repository.GetAllUsers(db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data users",
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    users,
		"message": "Data users berhasil diambil",
		"success": true,
	})
}

// GetUserByIDService mengambil user berdasarkan ID
func GetUserByIDService(c *fiber.Ctx, db *sql.DB) error {
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "User tidak valid",
			"success": false,
		})
	}

	// Admin dapat melihat user lain, selain itu hanya bisa lihat diri sendiri
	idParam := c.Params("id")
	if idParam != "" {
		adminRole := c.Locals("role")
		if adminRole != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "Tidak memiliki akses",
				"success": false,
			})
		}
	}

	user, err := repository.GetUserByID(db, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "User tidak ditemukan",
				"success": false,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data user",
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    user,
		"message": "Data user berhasil diambil",
		"success": true,
	})
}

// UpdateUserService mengupdate data user
func UpdateUserService(c *fiber.Ctx, db *sql.DB) error {
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "User tidak valid",
			"success": false,
		})
	}

	var req model.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
			"success": false,
		})
	}

	user, err := repository.UpdateUser(db, userID, req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "User tidak ditemukan",
				"success": false,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengupdate user",
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    user,
		"message": "User berhasil diupdate",
		"success": true,
	})
}

// DeleteUserService menghapus user (Admin only)
func DeleteUserService(c *fiber.Ctx, db *sql.DB) error {
	userID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID user tidak valid",
			"success": false,
		})
	}

	err = repository.DeleteUser(db, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "User tidak ditemukan",
				"success": false,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menghapus user",
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User berhasil dihapus",
		"success": true,
	})
}

// CreateStudentProfileService membuat profil student (Admin only)
func CreateStudentProfileService(c *fiber.Ctx, db *sql.DB) error {
	userID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID user tidak valid",
			"success": false,
		})
	}

	var req model.StudentProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
			"success": false,
		})
	}

	// Verify user exists dan role adalah student
	user, err := repository.GetUserByID(db, userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User tidak ditemukan",
			"success": false,
		})
	}

	if user.Role != "student" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "User bukan student",
			"success": false,
		})
	}

	// Create student profile
	err = repository.CreateStudentProfile(db, userID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal membuat profil student",
			"error":   err.Error(),
			"success": false,
		})
	}

	// Get student profile
	student, err := repository.GetStudentByUserID(db, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil profil student",
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data":    student,
		"message": "Profil student berhasil dibuat",
		"success": true,
	})
}

// CreateLecturerProfileService membuat profil lecturer (Admin only)
func CreateLecturerProfileService(c *fiber.Ctx, db *sql.DB) error {
	userID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID user tidak valid",
			"success": false,
		})
	}

	var req model.LecturerProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
			"success": false,
		})
	}

	// Verify user exists dan role adalah lecturer
	user, err := repository.GetUserByID(db, userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User tidak ditemukan",
			"success": false,
		})
	}

	if user.Role != "lecturer" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "User bukan lecturer",
			"success": false,
		})
	}

	// Create lecturer profile
	err = repository.CreateLecturerProfile(db, userID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal membuat profil lecturer",
			"error":   err.Error(),
			"success": false,
		})
	}

	// Get lecturer profile
	lecturer, err := repository.GetLecturerByUserID(db, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil profil lecturer",
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data":    lecturer,
		"message": "Profil lecturer berhasil dibuat",
		"success": true,
	})
}