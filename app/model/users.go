package model

import "time"

type Users struct {
	ID 	 		string    `json:"id"`
	Username 	string    `json:"username"`
	Email 		string    `json:"email"`
	Password	string    `json:"password"`
	FullName	string    `json:"full_name"`
	RoleID 		string    `json:"role"`
	IsActive	bool      `json:"is_active"`
	CreatedAt 	time.Time `json:"created_at"`
	UpdatedAt 	time.Time `json:"updated_at"`
}