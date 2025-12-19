package service

import (
	"student-report/app/model"
	"student-report/app/repository"
	"database/sql"
	"errors"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// CreateUser godoc
// @Summary Create new user
// @Description Create new user (Admin only)
// @Tags Users
// @Accept json
// @Produce json
// @Param key path string true "API Key"
// @Security BearerAuth
// @Param user body object true "User data" example({"username":"johndoe","password":"password123","full_name":"John Doe","email":"john@example.com","role":"student"})
// @Success 201 {object} object{success=bool,message=string,data=object{id=string,username=string,full_name=string,email=string,role=string,created_at=string}} "User created successfully"
// @Failure 400 {object} object{success=bool,message=string,error=string} "Invalid request body"
// @Failure 401 {object} object{success=bool,message=string} "Unauthorized"
// @Failure 403 {object} object{success=bool,message=string} "Forbidden - insufficient permissions"
// @Failure 500 {object} object{success=bool,message=string,error=string} "Internal server error"
// @Router /{key}/v1/users [post]
// CreateUserService membuat user baru (Admin only)
func (s *UserService) CreateUserService(c *fiber.Ctx) error {
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
	user, err := s.repo.CreateUser(req, string(hashedPassword))
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

// GetAllUsers godoc
// @Summary Get all users
// @Description Get all users in system (Admin only)
// @Tags Users
// @Accept json
// @Produce json
// @Param key path string true "API Key"
// @Security BearerAuth
// @Success 200 {object} object{success=bool,message=string,data=[]object{id=string,username=string,full_name=string,email=string,role=string,created_at=string}} "Users retrieved successfully"
// @Failure 401 {object} object{success=bool,message=string} "Unauthorized"
// @Failure 403 {object} object{success=bool,message=string} "Forbidden - insufficient permissions"
// @Failure 500 {object} object{success=bool,message=string,error=string} "Internal server error"
// @Router /{key}/v1/users [get]
// GetAllUsersService mengambil semua users (Admin only)
func (s *UserService) GetAllUsersService(c *fiber.Ctx) error {
	users, err := s.repo.GetAllUsers()
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

// GetUserByID godoc
// @Summary Get user by ID
// @Description Get user details by ID (Admin can see any user, others can only see themselves)
// @Tags Users
// @Accept json
// @Produce json
// @Param key path string true "API Key"
// @Param id path string true "User ID"
// @Security BearerAuth
// @Success 200 {object} object{success=bool,message=string,data=object{id=string,username=string,full_name=string,email=string,role=string,created_at=string}} "User retrieved successfully"
// @Failure 401 {object} object{success=bool,message=string} "Unauthorized"
// @Failure 403 {object} object{success=bool,message=string} "Forbidden - insufficient permissions"
// @Failure 404 {object} object{success=bool,message=string} "User not found"
// @Failure 500 {object} object{success=bool,message=string,error=string} "Internal server error"
// @Router /{key}/v1/users/{id} [get]
// GetUserByIDService mengambil user berdasarkan ID
func (s *UserService) GetUserByIDService(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
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
		userID = idParam
	}

	user, err := s.repo.GetUserByID(userID)
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

// UpdateUser godoc
// @Summary Update user
// @Description Update user information (Admin can update any user, others can only update themselves)
// @Tags Users
// @Accept json
// @Produce json
// @Param key path string true "API Key"
// @Param id path string true "User ID"
// @Security BearerAuth
// @Param user body object true "Update data" example({"full_name":"John Doe Updated","email":"john.new@example.com"})
// @Success 200 {object} object{success=bool,message=string,data=object{id=string,username=string,full_name=string,email=string,role=string}} "User updated successfully"
// @Failure 400 {object} object{success=bool,message=string,error=string} "Invalid request body"
// @Failure 401 {object} object{success=bool,message=string} "Unauthorized"
// @Failure 403 {object} object{success=bool,message=string} "Forbidden - insufficient permissions"
// @Failure 404 {object} object{success=bool,message=string} "User not found"
// @Failure 500 {object} object{success=bool,message=string,error=string} "Internal server error"
// @Router /{key}/v1/users/{id} [put]
// UpdateUserService mengupdate data user
func (s *UserService) UpdateUserService(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
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

	user, err := s.repo.UpdateUser(userID, req)
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

// DeleteUser godoc
// @Summary Delete user
// @Description Delete user (Admin only)
// @Tags Users
// @Accept json
// @Produce json
// @Param key path string true "API Key"
// @Param id path string true "User ID"
// @Security BearerAuth
// @Success 200 {object} object{success=bool,message=string} "User deleted successfully"
// @Failure 401 {object} object{success=bool,message=string} "Unauthorized"
// @Failure 403 {object} object{success=bool,message=string} "Forbidden - insufficient permissions"
// @Failure 404 {object} object{success=bool,message=string} "User not found"
// @Failure 500 {object} object{success=bool,message=string,error=string} "Internal server error"
// @Router /{key}/v1/users/{id} [delete]
// DeleteUserService menghapus user (Admin only)
func (s *UserService) DeleteUserService(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID user tidak valid",
			"success": false,
		})
	}

	err := s.repo.DeleteUser(userID)
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

// CreateStudentProfile godoc
// @Summary Create student profile
// @Description Create student profile for user with role student (Admin only)
// @Tags Users
// @Accept json
// @Produce json
// @Param key path string true "API Key"
// @Param id path string true "User ID"
// @Security BearerAuth
// @Param profile body object true "Student profile data" example({"student_id":"12345678","major":"Teknik Informatika","batch":"2021","gpa":3.75,"advisor_id":"uuid-advisor"})
// @Success 201 {object} object{success=bool,message=string,data=object{id=string,user_id=string,student_id=string,major=string,batch=string,gpa=number,advisor_id=string}} "Student profile created successfully"
// @Failure 400 {object} object{success=bool,message=string,error=string} "Invalid request body or user is not student"
// @Failure 401 {object} object{success=bool,message=string} "Unauthorized"
// @Failure 403 {object} object{success=bool,message=string} "Forbidden - insufficient permissions"
// @Failure 404 {object} object{success=bool,message=string} "User not found"
// @Failure 500 {object} object{success=bool,message=string,error=string} "Internal server error"
// @Router /{key}/v1/users/{id}/student-profile [post]
// CreateStudentProfileService membuat profil student (Admin only)
func (s *UserService) CreateStudentProfileService(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
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
	user, err := s.repo.GetUserByID(userID)
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
	err = s.repo.CreateStudentProfile(userID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal membuat profil student",
			"error":   err.Error(),
			"success": false,
		})
	}

	// Get student profile
	student, err := s.repo.GetStudentByUserID(userID)
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

// CreateLecturerProfile godoc
// @Summary Create lecturer profile
// @Description Create lecturer profile for user with role lecturer (Admin only)
// @Tags Users
// @Accept json
// @Produce json
// @Param key path string true "API Key"
// @Param id path string true "User ID"
// @Security BearerAuth
// @Param profile body object true "Lecturer profile data" example({"lecturer_id":"198765432","department":"Teknik Informatika","position":"Lektor"})
// @Success 201 {object} object{success=bool,message=string,data=object{id=string,user_id=string,lecturer_id=string,department=string,position=string}} "Lecturer profile created successfully"
// @Failure 400 {object} object{success=bool,message=string,error=string} "Invalid request body or user is not lecturer"
// @Failure 401 {object} object{success=bool,message=string} "Unauthorized"
// @Failure 403 {object} object{success=bool,message=string} "Forbidden - insufficient permissions"
// @Failure 404 {object} object{success=bool,message=string} "User not found"
// @Failure 500 {object} object{success=bool,message=string,error=string} "Internal server error"
// @Router /{key}/v1/users/{id}/lecturer-profile [post]
// CreateLecturerProfileService membuat profil lecturer (Admin only)
func (s *UserService) CreateLecturerProfileService(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
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
	user, err := s.repo.GetUserByID(userID)
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
	err = s.repo.CreateLecturerProfile(userID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal membuat profil lecturer",
			"error":   err.Error(),
			"success": false,
		})
	}

	// Get lecturer profile
	lecturer, err := s.repo.GetLecturerByUserID(userID)
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
