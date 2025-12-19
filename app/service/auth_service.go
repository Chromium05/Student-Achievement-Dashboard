package service

import (
	"student-report/app/model"
	"student-report/app/repository"
	"student-report/utils"
	"database/sql"
	"errors"

	"github.com/gofiber/fiber/v2"
)

type AuthService struct {
	repo *repository.AuthRepository
}

func NewAuthService(repo *repository.AuthRepository) *AuthService {
	return &AuthService{repo: repo}
}


// Login godoc
// @Summary Login user
// @Description Authenticate user dengan username dan password, mengembalikan access token dan refresh token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param key path string true "API Key"
// @Param login body object true "Login credentials" example({"username":"johndoe","password":"password123"})
// @Success 200 {object} object{success=bool,message=string,token=string,refresh_token=string,data=object{id=string,username=string,full_name=string,email=string,role=string}} "Login successful"
// @Failure 400 {object} object{success=bool,message=string,error=string} "Invalid request body"
// @Failure 401 {object} object{success=bool,message=string,error=string} "Invalid credentials"
// @Failure 500 {object} object{success=bool,message=string,error=string} "Internal server error"
// @Router /{key}/v1/auth/login [post]
func (s *AuthService) LoginService(c *fiber.Ctx) error {
	var loginData model.LoginRequest
	if err := c.BodyParser(&loginData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"success": false,
		})
	}

	if loginData.Username == "" || loginData.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Harap masukkan username dan password",
			"success": false,
		})
	}

	user, err := s.repo.Login(loginData.Username, loginData.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Username atau password salah",
				"error":   err.Error(),
				"success": false,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal terhubung ke database",
			"error":   err.Error(),
			"success": false,
		})
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal membuat token",
			"error":   err.Error(),
			"success": false,
		})
	}

	refreshToken, err := utils.GenerateRefreshToken(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal membuat refresh token",
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token":          token,
		"refresh_token":  refreshToken,
		"data":           user,
		"message":        "Login berhasil",
		"success":        true,
	})
}

// Logout godoc
// @Summary Logout user
// @Description Logout user dan blacklist token JWT
// @Tags Authentication
// @Accept json
// @Produce json
// @Param key path string true "API Key"
// @Security BearerAuth
// @Success 200 {object} object{success=bool,message=string,user_id=string,username=string} "Logout successful"
// @Failure 401 {object} object{success=bool,message=string} "Unauthorized - token invalid or missing"
// @Router /{key}/v1/auth/logout [post]
func (s *AuthService) LogoutService(c *fiber.Ctx) error {
	// Ambil token dari header Authorization
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Token tidak ditemukan",
			"success": false,
		})
	}

	// Validasi token format
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Format token tidak valid",
			"success": false,
		})
	}

	// Extract token
	tokenString := authHeader[7:]

	// Validasi token dan ambil claims untuk mendapat expiry time
	claims, err := utils.ValidateToken(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Token tidak valid",
			"success": false,
		})
	}

	// Ambil user info dari context
	userID := c.Locals("user_id")
	username := c.Locals("username")

	if claims.ExpiresAt != nil {
		utils.AddToBlacklist(tokenString, claims.ExpiresAt.Time)
	}

	// Return success message
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Logout berhasil",
		"user_id": userID,
		"username": username,
		"success": true,
	})
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Generate new access token menggunakan refresh token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param key path string true "API Key"
// @Param refresh body object true "Refresh token" example({"refresh_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."})
// @Success 200 {object} object{success=bool,message=string,token=string,refresh_token=string,expires_in=int} "Token refreshed successfully"
// @Failure 400 {object} object{success=bool,message=string} "Invalid request body"
// @Failure 401 {object} object{success=bool,message=string,error=string} "Invalid or expired refresh token"
// @Failure 500 {object} object{success=bool,message=string,error=string} "Internal server error"
// @Router /{key}/v1/auth/refresh [post]
func (s *AuthService) RefreshTokenService(c *fiber.Ctx) error {
	var refreshReq model.RefreshTokenRequest
	if err := c.BodyParser(&refreshReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"success": false,
		})
	}

	if refreshReq.RefreshToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Refresh token diperlukan",
			"success": false,
		})
	}

	// Validate refresh token
	claims, err := utils.ValidateRefreshToken(refreshReq.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Refresh token tidak valid atau expired",
			"error":   err.Error(),
			"success": false,
		})
	}

	// Ambil user dari database untuk memastikan user masih valid
	user, err := s.repo.Login(claims.Username, "")
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "User tidak ditemukan",
			"error":   err.Error(),
			"success": false,
		})
	}

	// Generate new access token
	newAccessToken, err := utils.GenerateToken(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal membuat token baru",
			"error":   err.Error(),
			"success": false,
		})
	}

	// Generate new refresh token (optional: rotate refresh token)
	newRefreshToken, err := utils.GenerateRefreshToken(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal membuat refresh token baru",
			"error":   err.Error(),
			"success": false,
		})
	}

	response := model.TokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    1200, // 20 minutes in seconds
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token":           response.AccessToken,
		"refresh_token":   response.RefreshToken,
		"expires_in":      response.ExpiresIn,
		"message":         "Token berhasil di-refresh",
		"success":         true,
	})
}

func (s *AuthService) GetProfileService(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	response, err := s.repo.GetUserProfile(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mendapatkan profile user",
			"error":   err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    response,
		"message": "Berhasil mendapatkan profile user",
		"success": true,
	})
}