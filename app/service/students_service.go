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

// GetStudents godoc
// @Summary Get all students
// @Description Get list of all students
// @Tags Students
// @Accept json
// @Produce json
// @Param key path string true "API Key"
// @Security BearerAuth
// @Success 200 {object} object{success=bool,message=string,data=[]object{id=string,user_id=string,student_id=string,full_name=string,major=string,batch=string,gpa=number,advisor_id=string}} "Students retrieved successfully"
// @Failure 401 {object} object{success=bool,message=string} "Unauthorized"
// @Failure 403 {object} object{success=bool,message=string} "Forbidden - insufficient permissions"
// @Failure 500 {object} object{success=bool,message=string,error=string} "Internal server error"
// @Router /{key}/v1/students [get]
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

// GetStudentByUserID godoc
// @Summary Get student by user ID
// @Description Get student profile by user ID
// @Tags Students
// @Accept json
// @Produce json
// @Param key path string true "API Key"
// @Param userId path string true "User ID"
// @Security BearerAuth
// @Success 200 {object} object{success=bool,message=string,data=object{id=string,user_id=string,student_id=string,full_name=string,major=string,batch=string,gpa=number,advisor_id=string}} "Student retrieved successfully"
// @Failure 401 {object} object{success=bool,message=string} "Unauthorized"
// @Failure 404 {object} object{success=bool,message=string} "Student not found"
// @Failure 500 {object} object{success=bool,message=string,error=string} "Internal server error"
// @Router /{key}/v1/students/user/{userId} [get]
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

// GetStudentsByAdvisorID godoc
// @Summary Get students by advisor ID
// @Description Get all students under a specific advisor
// @Tags Students
// @Accept json
// @Produce json
// @Param key path string true "API Key"
// @Param advisorId path string true "Advisor ID"
// @Security BearerAuth
// @Success 200 {object} object{success=bool,message=string,data=[]object} "Students retrieved successfully"
// @Failure 401 {object} object{success=bool,message=string} "Unauthorized"
// @Failure 500 {object} object{success=bool,message=string,error=string} "Internal server error"
// @Router /{key}/v1/students/advisor/{advisorId} [get]
func (s *StudentService) GetStudentsByAdvisorID(c *fiber.Ctx) error {
	key := c.Params("key")
	if key != os.Getenv("API_KEY") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "API tidak sesuai",
			"success": false,
		})
	}

	advisorID := c.Params("advisorId")
	students, err := s.repo.GetStudentByAdvisorID(advisorID)
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
