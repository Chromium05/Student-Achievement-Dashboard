package repository

import (
	"database/sql"
	"student-report/app/model"
)

type LecturerRepository struct {
	db *sql.DB
}

func NewLecturerRepository(db *sql.DB) *LecturerRepository {
	return &LecturerRepository{db: db}
}

func (r *LecturerRepository) GetLecturersRepository() ([]model.Lecturers, error) {
	rows, err := r.db.Query(`
		SELECT l.id, l.user_id, l.lecturer_id, u.full_name, l.department, l.created_at
		FROM lecturers AS l
		JOIN users AS u ON l.user_id = u.id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var LecturersList []model.Lecturers

	for rows.Next() {
		var lecturers model.Lecturers
		err := rows.Scan(
			&lecturers.ID,
			&lecturers.UserID,
			&lecturers.LecturerID,
			&lecturers.FullName,
			&lecturers.Department,
			&lecturers.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		LecturersList = append(LecturersList, lecturers)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return LecturersList, nil
}

// func (r *StudentRepository) GetStudentByUserID(userID string) (*model.Students, error) {
// 	query := `
// 		SELECT s.id, s.user_id, s.student_id, u.full_name, s.program_study, s.academic_year, 
// 		s.advisor_id, s.created_at
// 		FROM students AS s
// 		WHERE user_id = $1
// 	`
// 	var student model.Students
// 	err := r.db.QueryRow(query, userID).Scan(
// 		&student.ID,
// 		&student.UserID,
// 		&student.StudentID,
// 		&student.Prodi,
// 		&student.Year,
// 		&student.AdvisorID,
// 		&student.CreatedAt,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &student, nil
// }

// func (r *StudentRepository) GetStudentsByAdvisorID(advisorID string) ([]model.Students, error) {
// 	query := `
// 		SELECT s.id, s.user_id, s.student_id, u.full_name, s.program_study, s.academic_year, 
// 		s.advisor_id, s.created_at
// 		FROM students s
// 		JOIN users u ON s.user_id = u.id
// 		WHERE advisor_id = $1
// 	`
// 	rows, err := r.db.Query(query, advisorID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var students []model.Students
// 	for rows.Next() {
// 		var student model.Students
// 		err := rows.Scan(
// 			&student.ID,
// 			&student.UserID,
// 			&student.StudentID,
// 			&student.FullName,
// 			&student.Prodi,
// 			&student.Year,
// 			&student.AdvisorID,
// 			&student.CreatedAt,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}
// 		students = append(students, student)
// 	}

// 	return students, nil
// }

// Legacy function for backward compatibility
func GetLecturersRepository(db *sql.DB) ([]model.Lecturers, error) {
	repo := NewLecturerRepository(db)
	return repo.GetLecturersRepository()
}
