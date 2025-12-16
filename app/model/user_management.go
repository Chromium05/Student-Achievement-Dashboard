package model

import "time"

// Request models untuk User Management
type CreateUserRequest struct {
	Username    string `json:"username" validate:"required,min=3,max=50"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8"`
	FullName    string `json:"full_name" validate:"required"`
	RoleID      string `json:"role_id" validate:"required"` // UUID for role
}

type UpdateUserRequest struct {
	Email     string `json:"email"`
	FullName  string `json:"full_name"` // combined full_name instead of first/last
	IsActive  bool   `json:"is_active"`
}

type UserResponse struct {
	ID          string    `json:"id"` // UUID instead of int
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	FullName    string    `json:"full_name"` // single full_name field
	Role        string    `json:"role"`
	RoleID      string    `json:"role_id"` // UUID for role
	IsActive    bool      `json:"is_active"`
	Permissions []string  `json:"permissions"`
	CreatedAt   time.Time `json:"created_at"`
}

// Student profile request
type StudentProfileRequest struct {
	StudentID      string `json:"student_id" validate:"required"`
	ProgramStudy   string `json:"program_study" validate:"required"`
	AcademicYear   string `json:"academic_year" validate:"required"`
	AdvisorID      string `json:"advisor_id"` // lecturer's user_id as UUID
}

// Lecturer profile request
type LecturerProfileRequest struct {
	LecturerID string `json:"lecturer_id" validate:"required"`
	Department string `json:"department" validate:"required"`
}

// Permission model with resource and action fields
type Permission struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Resource    string    `json:"resource"`
	Action      string    `json:"action"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// Role response
type RoleResponse struct {
	ID          string       `json:"id"` // UUID instead of int
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Permissions []Permission `json:"permissions"`
	CreatedAt   time.Time    `json:"created_at"`
}
