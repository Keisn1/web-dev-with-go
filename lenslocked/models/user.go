package models

import (
	"database/sql"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int
	Email        string
	PasswordHash string
}

type NewUser struct {
	Email    string
	Password string
}

type UserService struct {
	DB *sql.DB
}

func (us *UserService) Authenticate(nu NewUser) (*User, error) {
	email := strings.ToLower(nu.Email)
	row := us.DB.QueryRow(
		`SELECT id, password_hash FROM users WHERE email = $1`,
		email,
	)
	var user User
	err := row.Scan(&user.ID, &user.PasswordHash)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("User not in DB: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(nu.Password)); err != nil {
		return nil, fmt.Errorf("Wrong password: %w", err)
	}

	return &user, nil
}

func (us *UserService) Create(nu NewUser) (*User, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("Generating hash, UserService.Create: %w", err)
	}
	user := User{
		Email:        strings.ToLower(nu.Email),
		PasswordHash: string(hashedBytes),
	}
	row := us.DB.QueryRow(`
INSERT INTO users  (email, password_hash)
  VALUES ($1, $2) RETURNING id;`,
		user.Email,
		user.PasswordHash,
	)
	err = row.Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("Inserting User, UserService.Create: %w", err)
	}
	return &user, nil
}
