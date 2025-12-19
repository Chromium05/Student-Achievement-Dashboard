package service

import (
	// "database/sql"
	"os"
	"student-report/app/repository"
	"github.com/gofiber/fiber/v2"
)

type LecturerService struct {
	repo *repository.LecturerRepository
}

func NewLecturerService(repo *repository.LecturerRepository) *LecturerService {
	return &LecturerService{repo: repo}
}

// GetLecturers godoc
// @Summary Get all lecturers
// @Description Get list of all lecturers
// @Tags Lecturers
// @Accept json
// @Produce json
// @Param key path string true "API Key"
// @Security BearerAuth
// @Success 200 {object} object{success=bool,message=string,data=[]object{id=string,user_id=string,lecturer_id=string,full_name=string,department=string,position=string}} "Lecturers retrieved successfully"
// @Failure 401 {object} object{success=bool,message=string} "Unauthorized"
// @Failure 500 {object} object{success=bool,message=string,error=string} "Internal server error"
// @Router /{key}/v1/lecturers [get]
func (s *LecturerService) GetLecturersService(c *fiber.Ctx) error {
	key := c.Params("key")
	if key != os.Getenv("API_KEY") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "API tidak sesuai",
			"success": false,
		})
	}

	lecturer, err := s.repo.GetLecturersRepository()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal terhubung ke database.",
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    lecturer,
		"message": "Berhasil mendapatkan data Students",
		"success": true,
	})
}