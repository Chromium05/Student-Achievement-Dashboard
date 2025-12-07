package model

import "time"

type Students struct {
	ID 	  		string   `json:"id"`
	UserID 		string   `json:"user_id"`
	StudentID 	string   `json:"student_id"`
	Prodi 		string   `json:"prodi"`
	Year 		string   `json:"year"`
	AdvisorID 	string   `json:"advisor_id"`
	CreatedAt   time.Time   `json:"created_at"`
}

