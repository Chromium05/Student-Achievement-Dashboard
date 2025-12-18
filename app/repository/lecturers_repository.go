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

func (r *LecturerRepository) GetLecturerIDByUserID(userID string) (string, error) {
	query := `SELECT id FROM lecturers WHERE user_id = $1`
	var lecturerID string
	err := r.db.QueryRow(query, userID).Scan(&lecturerID)
	if err != nil {
		return "", err
	}
	return lecturerID, nil
}

// Legacy function for backward compatibility
func GetLecturersRepository(db *sql.DB) ([]model.Lecturers, error) {
	repo := NewLecturerRepository(db)
	return repo.GetLecturersRepository()
}
