package auth

import (
	"github.com/google/uuid"
	"github.com/itsLeonB/go-crud"
	"github.com/reflect-homini/stora/internal/domain/user"
)

type OAuthAccount struct {
	crud.BaseEntity
	UserID     uuid.UUID
	Provider   string
	ProviderID string
	Email      string

	// Relations
	User user.User
}

func (OAuthAccount) TableName() string {
	return "oauth_accounts"
}
