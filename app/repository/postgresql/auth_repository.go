package repository

import (
	"database/sql"
	"go_clean/app/models/postgresql"
)


type AuthRepository struct {
	DB *sql.DB
}

func (r *AuthRepository) GetByUsernameOrEmail(identifier string) (*models.User, string, error) {
	u := models.User{}
	var hash string
	err := r.DB.QueryRow(`
		SELECT id, username, email, password_hash, role
		FROM users
		WHERE username = $1 OR email = $1
	`, identifier).Scan(&u.ID, &u.Username, &u.Email, &hash, &u.Role)
	if err != nil {
		return nil, "", err
	}
	return &u, hash, nil
}

func (r *AuthRepository) ExistsByUsernameOrEmail(username, email string) (bool, error) {
	var exists bool
	err := r.DB.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM users WHERE username = $1 OR email = $2
		)
	`, username, email).Scan(&exists)
	return exists, err
}

func (r *AuthRepository) LoginRepo(identifier string) (*models.User, string, error) {
	var u models.User
	var hash string

	err := r.DB.QueryRow(`
		SELECT id, username, email, password_hash, role
		FROM users
		WHERE username = $1 OR email = $1
	`, identifier).Scan(&u.ID, &u.Username, &u.Email, &hash, &u.Role)

	if err != nil {
		return nil, "", err
	}

	return &u, hash, nil
}

