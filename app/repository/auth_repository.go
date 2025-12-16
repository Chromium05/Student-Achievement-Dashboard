package repository

import (
	"student-report/app/model"
	"database/sql"
	"fmt"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) Login(username string, password string) (model.User, error) {
	var user model.User
	
	row := r.db.QueryRow(`
		SELECT u.id, u.username, u.email, u.password_hash, u.full_name, u.role_id, r.name as role_name 
		FROM users u
		JOIN roles r ON u.role_id = r.id
		WHERE u.username = $1 AND u.is_active = true
	`, username)
	
	var hashedPassword string
	err := row.Scan(&user.ID, &user.Username, &user.Email, &hashedPassword, &user.FullName, &user.RoleID, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("user not found")
		}
		return user, err
	}
	
	user.Permissions = make([]string, 0)
	
	permissions, err := r.getPermissionsByRole(user.Role)
	if err != nil {
		// Log the error but don't fail the login
		fmt.Printf("Warning: Error getting permissions for role %s: %v\n", user.Role, err)
	} else {
		user.Permissions = permissions
	}
	
	return user, nil
}

func (r *AuthRepository) getPermissionsByRole(roleName string) ([]string, error) {
	rows, err := r.db.Query(`
		SELECT p.name 
		FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		JOIN roles r ON rp.role_id = r.id
		WHERE r.name = $1
	`, roleName)
	if err != nil {
		return nil, fmt.Errorf("error querying permissions: %w", err)
	}
	defer rows.Close()

	permissions := make([]string, 0)
	for rows.Next() {
		var perm string
		if err := rows.Scan(&perm); err != nil {
			return nil, fmt.Errorf("error scanning permission: %w", err)
		}
		permissions = append(permissions, perm)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating permissions: %w", err)
	}

	return permissions, nil
}

// Legacy function for backward compatibility
func Login(db *sql.DB, username string, password string) (model.User, error) {
	repo := NewAuthRepository(db)
	return repo.Login(username, password)
}
