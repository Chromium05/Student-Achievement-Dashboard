package repository

import (
	"student-report/app/model"
	// "context"
	// "go.mongodb.org/mongo-driver/bson"
    // "go.mongodb.org/mongo-driver/mongo"
    // "go.mongodb.org/mongo-driver/mongo/options"
	"database/sql"
)

func GetStudentsRepository(db *sql.DB) ([]model.Students, error) {
	rows, err := db.Query(`
		SELECT id, user_id, student_id, program_study, academic_year, 
		advisor_id, created_at
		FROM students
	`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var StudentsList []model.Students

	for rows.Next() {
		var students model.Students
		err := rows.Scan(
			&students.ID,
			&students.UserID,
			&students.StudentID,
			&students.Prodi,
			&students.Year,
			&students.AdvisorID,
			&students.CreatedAt,
		)

		if err != nil {
			return nil, err
		}
		StudentsList = append(StudentsList, students)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return StudentsList, nil
}