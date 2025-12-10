package repository

import (
	"student-report/app/model"
	"database/sql"
	"fmt"
	"strings"
)

// CreateUser membuat user baru dengan role yang spesifik
func CreateUser(db *sql.DB, req model.CreateUserRequest, hashedPassword string) (model.UserResponse, error) {
	var userResp model.UserResponse
	
	tx, err := db.Begin()
	if err != nil {
		return userResp, err
	}
	defer tx.Rollback()

	// Insert user
	err = tx.QueryRow(`
		INSERT INTO users (username, email, password_hash, role, is_active)
		VALUES ($1, $2, $3, $4, true)
		RETURNING id, username, email, role, is_active, created_at
	`, req.Username, req.Email, hashedPassword, req.Role).Scan(
		&userResp.ID, &userResp.Username, &userResp.Email, 
		&userResp.Role, &userResp.IsActive, &userResp.CreatedAt,
	)
	if err != nil {
		return userResp, fmt.Errorf("error creating user: %w", err)
	}

	userResp.FirstName = req.FirstName
	userResp.LastName = req.LastName

	// Get permissions for this role
	permissions, err := GetPermissionsByRole(tx, req.Role)
	if err != nil {
		return userResp, fmt.Errorf("error getting permissions: %w", err)
	}
	userResp.Permissions = permissions

	if err = tx.Commit(); err != nil {
		return userResp, err
	}

	return userResp, nil
}

// GetAllUsers mendapatkan semua users dengan filter optional
func GetAllUsers(db *sql.DB) ([]model.UserResponse, error) {
	rows, err := db.Query(`
		SELECT id, username, email, role, is_active, created_at 
		FROM users 
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.UserResponse
	for rows.Next() {
		var user model.UserResponse
		err := rows.Scan(
			&user.ID, &user.Username, &user.Email, 
			&user.Role, &user.IsActive, &user.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Get permissions
		permissions, err := GetPermissionsByRoleString(db, user.Role)
		if err == nil {
			user.Permissions = permissions
		}

		users = append(users, user)
	}

	return users, rows.Err()
}

// GetUserByID mendapatkan user berdasarkan ID
func GetUserByID(db *sql.DB, userID int) (model.UserResponse, error) {
	var user model.UserResponse
	err := db.QueryRow(`
		SELECT id, username, email, role, is_active, created_at 
		FROM users 
		WHERE id = $1
	`, userID).Scan(
		&user.ID, &user.Username, &user.Email, 
		&user.Role, &user.IsActive, &user.CreatedAt,
	)
	if err != nil {
		return user, err
	}

	permissions, err := GetPermissionsByRoleString(db, user.Role)
	if err == nil {
		user.Permissions = permissions
	}

	return user, nil
}

// UpdateUser melakukan update pada user
func UpdateUser(db *sql.DB, userID int, req model.UpdateUserRequest) (model.UserResponse, error) {
	var user model.UserResponse

	query := `UPDATE users SET `
	var args []interface{}
	var setClauses []string
	
	argNum := 1
	if req.Email != "" {
		setClauses = append(setClauses, fmt.Sprintf("email = $%d", argNum))
		args = append(args, req.Email)
		argNum++
	}
	
	setClauses = append(setClauses, fmt.Sprintf("is_active = $%d", argNum))
	args = append(args, req.IsActive)
	argNum++

	query += strings.Join(setClauses, ", ") + fmt.Sprintf(" WHERE id = $%d RETURNING id, username, email, role, is_active, created_at", argNum)
	args = append(args, userID)

	err := db.QueryRow(query, args...).Scan(
		&user.ID, &user.Username, &user.Email, 
		&user.Role, &user.IsActive, &user.CreatedAt,
	)
	if err != nil {
		return user, err
	}

	permissions, err := GetPermissionsByRoleString(db, user.Role)
	if err == nil {
		user.Permissions = permissions
	}

	return user, nil
}

// DeleteUser melakukan soft delete user
func DeleteUser(db *sql.DB, userID int) error {
	result, err := db.Exec(`
		UPDATE users SET is_active = false WHERE id = $1
	`, userID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// GetPermissionsByRole mengambil permissions untuk suatu role dari tx
func GetPermissionsByRole(tx *sql.Tx, roleName string) ([]string, error) {
	rows, err := tx.Query(`
		SELECT p.name FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		JOIN roles r ON rp.role_id = r.id
		WHERE r.name = $1
	`, roleName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []string
	for rows.Next() {
		var perm string
		if err := rows.Scan(&perm); err != nil {
			return nil, err
		}
		permissions = append(permissions, perm)
	}

	return permissions, rows.Err()
}

// GetPermissionsByRoleString mengambil permissions untuk suatu role
func GetPermissionsByRoleString(db *sql.DB, roleName string) ([]string, error) {
	rows, err := db.Query(`
		SELECT p.name FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		JOIN roles r ON rp.role_id = r.id
		WHERE r.name = $1
	`, roleName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []string
	for rows.Next() {
		var perm string
		if err := rows.Scan(&perm); err != nil {
			return nil, err
		}
		permissions = append(permissions, perm)
	}

	return permissions, rows.Err()
}

// CreateStudentProfile membuat profil student setelah user dibuat
func CreateStudentProfile(db *sql.DB, userID int, req model.StudentProfileRequest) error {
	_, err := db.Exec(`
		INSERT INTO students (user_id, student_id, program_study, academic_year, advisor_id)
		VALUES ($1, $2, $3, $4, $5)
	`, userID, req.StudentID, req.ProgramStudy, req.AcademicYear, req.AdvisorID)
	return err
}

// CreateLecturerProfile membuat profil lecturer setelah user dibuat
func CreateLecturerProfile(db *sql.DB, userID int, req model.LecturerProfileRequest) error {
	_, err := db.Exec(`
		INSERT INTO lecturers (user_id, lecturer_id, department)
		VALUES ($1, $2, $3)
	`, userID, req.LecturerID, req.Department)
	return err
}

// GetStudentByUserID mengambil data student berdasarkan user_id
func GetStudentByUserID(db *sql.DB, userID int) (model.Students, error) {
	var student model.Students
	err := db.QueryRow(`
		SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at
		FROM students
		WHERE user_id = $1
	`, userID).Scan(
		&student.ID, &student.UserID, &student.StudentID, 
		&student.Prodi, &student.Year, &student.AdvisorID, &student.CreatedAt,
	)
	return student, err
}

// GetLecturerByUserID mengambil data lecturer berdasarkan user_id
func GetLecturerByUserID(db *sql.DB, userID int) (model.Lecturers, error) {
	var lecturer model.Lecturers
	err := db.QueryRow(`
		SELECT id, user_id, lecturer_id, department, created_at
		FROM lecturers
		WHERE user_id = $1
	`, userID).Scan(
		&lecturer.ID, &lecturer.UserID, &lecturer.LecturerID, 
		&lecturer.Department, &lecturer.CreatedAt,
	)
	return lecturer, err
}
