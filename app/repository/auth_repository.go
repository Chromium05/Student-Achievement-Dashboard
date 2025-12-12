package repository

import (
	"student-report/app/model"
	"database/sql"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) Login(username string, password string) (model.User, error) {
	var user model.User
	row := r.db.QueryRow(`SELECT id, username, email, password_hash, role FROM users WHERE username = $1 AND is_active = true`, username)
	var hashedPassword string
	err := row.Scan(&user.ID, &user.Username, &user.Email, &hashedPassword, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, err
		}
		return user, err
	}

	// Bandingkan password yang diberikan dengan hash yang disimpan
	if hashedPassword != password {
		return user, err
	}

	permissions, err := GetPermissionsByRoleString(r.db, user.Role)
	if err == nil {
		user.Permissions = permissions
	}

	return user, nil
}

// Legacy function for backward compatibility
func Login(db *sql.DB, username string, password string) (model.User, error) {
	repo := NewAuthRepository(db)
	return repo.Login(username, password)
}