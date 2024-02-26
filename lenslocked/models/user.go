package models

import (
	"database/sql"
	"fmt"
	"strings"

	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
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

var (
	// A common pattern is to add the package as a prefix to the error for
	// context.
	ErrEmailTaken = errors.New("models: email address is already in use")
)

func (us *UserService) UpdatePassword(userID int, password string) error {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("update password: %w", err)
	}
	passwordHash := string(hashedBytes)
	_, err = us.DB.Exec(`
UPDATE users
SET password_hash = $2
WHERE id = $1;`, userID, passwordHash)
	if err != nil {
		return fmt.Errorf("update password: %w", err)
	}
	return nil
}

func (us *UserService) Authenticate(nu NewUser) (*User, error) {
	var user User
	user.Email = strings.ToLower(nu.Email)
	row := us.DB.QueryRow(
		`SELECT id, password_hash FROM users WHERE email = $1`,
		user.Email,
	)
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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				fmt.Println("here")
				return nil, ErrEmailTaken
			}
		}
		return nil, fmt.Errorf("Inserting User, UserService.Create: %w", err)
	}
	return &user, nil
}
