package model

import "time"

type Lecturers struct {
	ID 	 		string    `json:"id"`
	UserID 		string    `json:"user_id"`
	LecturerID 	string    `json:"lecturer_id"`
	FullName    string    `json:"full_name"`
	Department 	string    `json:"department"`
	CreatedAt 	time.Time `json:"created_at"`
}