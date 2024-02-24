package models

import "database/sql"

type Session struct {
	Id     int
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

// Create will create a new session for the user provided. The session token
// will be returned as the Token field on the Session type, but only the hashed
// session token is stored in the database
func (ss *SessionService) Create(userID int) (*Session, error) {
	var session Session
	return &session, nil
}

func (ss *SessionService) User(token string) (*User, error) {
	return nil, nil
}
