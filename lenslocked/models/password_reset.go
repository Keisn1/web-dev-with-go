package models

import "time"
import "database/sql"
import "fmt"

type PasswordReset struct {
	ID     int
	UserID int
	// Token is only set when PasswordReset is being created
	Token     string
	TokenHash string
	ExpiresAt time.Time
}

type PasswordResetService struct {
	DB            *sql.DB
	Duration      time.Duration
	BytesPerToken int
}

const (
	// Default time that a passwordreset is valid
	DefaultResetDuration = 1 * time.Hour
)

func (prs *PasswordResetService) Create(email string) (*PasswordReset, error) {
	return nil, nil
}

// We are going to consume a token and return the user associated with it, or return an error if not present
func (service *PasswordResetService) Consume(token string) (*User, error) {
	return nil, fmt.Errorf("TODO: Implement PasswordResetService.Consume")
}
