package repository

import (
	"database/sql"
	"student-report/app/model"
)

type StudentRepository struct {
	db *sql.DB
}

func NewStudentRepository(db *sql.DB) *StudentRepository {
	return &StudentRepository{db: db}
}

func (r *StudentRepository) GetStudentsRepository() ([]model.Students, error) {
	rows, err := r.db.Query(`
		SELECT s.id, s.user_id, s.student_id, u.full_name, s.program_study, s.academic_year, 
		s.advisor_id, s.created_at
		FROM students AS s
		JOIN users AS u ON s.user_id = u.id
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
			&students.FullName,
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

func (r *StudentRepository) GetStudentByUserID(userID string) (*model.Students, error) {
	query := `
		SELECT s.id, s.user_id, s.student_id, u.full_name, s.program_study, s.academic_year, 
		s.advisor_id, s.created_at
		FROM students AS s
		JOIN users AS u ON s.user_id = u.id
		WHERE s.user_id = $1
	`
	var student model.Students
	err := r.db.QueryRow(query, userID).Scan(
		&student.ID,
		&student.UserID,
		&student.StudentID,
		&student.FullName,
		&student.Prodi,
		&student.Year,
		&student.AdvisorID,
		&student.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &student, nil
}

func (r *StudentRepository) GetStudentByAdvisorID(advisorID string) ([]model.Students, error) {
	query := `
		SELECT s.id, s.user_id, s.student_id, u.full_name, s.program_study, s.academic_year, 
		s.advisor_id, s.created_at
		FROM students s
		JOIN users u ON s.user_id = u.id
		WHERE advisor_id = $1
	`
	rows, err := r.db.Query(query, advisorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []model.Students
	for rows.Next() {
		var student model.Students
		err := rows.Scan(
			&student.ID,
			&student.UserID,
			&student.StudentID,
			&student.FullName,
			&student.Prodi,
			&student.Year,
			&student.AdvisorID,
			&student.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		students = append(students, student)
	}

	return students, nil
}

// Legacy function for backward compatibility
func GetStudentsRepository(db *sql.DB) ([]model.Students, error) {
	repo := NewStudentRepository(db)
	return repo.GetStudentsRepository()
}
