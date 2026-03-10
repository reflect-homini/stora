package user

import (
	"database/sql"

	"github.com/itsLeonB/go-crud"
)

type User struct {
	crud.BaseEntity
	Email        string
	PasswordHash string
	VerifiedAt   sql.NullTime

	// Relations
	Profile             UserProfile
	PasswordResetTokens []PasswordResetToken
}

func (u User) IsVerified() bool {
	return u.VerifiedAt.Valid && !u.VerifiedAt.Time.IsZero()
}
