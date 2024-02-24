package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"

	"fmt"

	"github.com/keisn1/lenslocked/rand"
)

type Session struct {
	ID     int
	UserID int
	// Token is only set when creating a new session. When looking up a Session
	// this will be left empty, as we only store the hash of a session Token
	// in our database and we cannot reverse it into a raw token.
	Token     string
	TokenHash string
}

type SessionService struct {
	DB *sql.DB
}

const (
	// the minimum number of bytes to be used for each session token.
	MinBytesPerToken = 32
)

type TokenManager struct {
	BytesPerToken int
}

func (tm TokenManager) New() (token, tokenHash string, err error) {
	bytesPerToken := tm.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}
	token, err = rand.String(bytesPerToken)
	if err != nil {
		return "", "", fmt.Errorf("new: %w", err)
	}
	tokenHash = tm.hash(token)
	return token, tokenHash, nil
}

func (tm TokenManager) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])

}

// Create will create a new session for the user provided. The session token
// will be returned as the Token field on the Session type, but only the hashed
// session token is stored in the database
func (ss *SessionService) Create(userID int) (*Session, error) {
	tm := TokenManager{}
	token, tokenHash, err := tm.New()

	session := Session{UserID: userID, Token: token, TokenHash: tokenHash}
	row := ss.DB.QueryRow(`
		UPDATE sessions SET token_hash=$2
		WHERE user_id=$1 RETURNING id;`,
		session.UserID,
		session.TokenHash,
	)
	err = row.Scan(&session.ID)
	if err == sql.ErrNoRows {
		// if no sesion exists, we will get ErrNoRows. This means we need to
		// create a session object for that user
		row = ss.DB.QueryRow(`
	    INSERT INTO sessions (user_id, token_hash)
		VALUES ($1, $2)
		RETURNING id;`,
			session.UserID,
			session.TokenHash,
		)
		err = row.Scan(&session.ID)
	}

	// if err was not sql.ErrNoRows, need to check to see if it was any
	// other error
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	return &session, nil
}

func (ss *SessionService) User(token string) (*User, error) {
	tm := TokenManager{}
	tokenHash := tm.hash(token)
	var user User
	row := ss.DB.QueryRow(`
		SELECT user_id
		FROM sessions
		WHERE token_hash=$1;`,
		tokenHash,
	)
	err := row.Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("user: %w", err)
	}

	row = ss.DB.QueryRow(
		`SELECT email, password_hash FROM users WHERE id=$1;`,
		user.ID,
	)
	err = row.Scan(&user.Email, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("user: %w", err)
	}
	return &user, nil
}

func (ss *SessionService) Delete(token string) error {
	tm := TokenManager{}
	tokenHash := tm.hash(token)
	_, err := ss.DB.Exec(`
		DELETE FROM sessions
		WHERE token_hash = $1`,
		tokenHash,
	)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}
