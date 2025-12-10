package model

import "time"

// Request models untuk User Management
type CreateUserRequest struct {
	Username    string `json:"username" validate:"required,min=3,max=50"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8"`
	Role        string `json:"role" validate:"required,oneof=admin lecturer student"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
}

type UpdateUserRequest struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	IsActive  bool   `json:"is_active"`
}

type UserResponse struct {
	ID          int       `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Role        string    `json:"role"`
	IsActive    bool      `json:"is_active"`
	Permissions []string  `json:"permissions"`
	CreatedAt   time.Time `json:"created_at"`
}

// Student profile request
type StudentProfileRequest struct {
	StudentID      string `json:"student_id" validate:"required"`
	ProgramStudy   string `json:"program_study" validate:"required"`
	AcademicYear   string `json:"academic_year" validate:"required"`
	AdvisorID      int    `json:"advisor_id"` // lecturer's user_id
}

// Lecturer profile request
type LecturerProfileRequest struct {
	LecturerID string `json:"lecturer_id" validate:"required"`
	Department string `json:"department" validate:"required"`
}

// Role response
type RoleResponse struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Permissions []string  `json:"permissions"`
	CreatedAt   time.Time `json:"created_at"`
}
