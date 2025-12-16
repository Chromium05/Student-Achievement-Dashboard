package model

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	ID          string    `json:"id"` // UUID instead of int
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	FullName    string    `json:"full_name"` // full_name field
	Role        string    `json:"role"`
	RoleID      string    `json:"role_id"` // UUID for role
	Permissions []string  `json:"permissions"`
	CreatedAt   time.Time `json:"created_at"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User         User   `json:"user"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type JWTClaims struct {
	UserID      string   `json:"user_id"` // UUID instead of int
	Username    string   `json:"username"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
	TokenType   string   `json:"token_type"`
	jwt.RegisteredClaims
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}
