package provider

import (
	"github.com/itsLeonB/go-crud"
	"github.com/reflect-homini/stora/internal/domain/auth"
	"github.com/reflect-homini/stora/internal/domain/project"
	"github.com/reflect-homini/stora/internal/domain/user"
	"gorm.io/gorm"
)

type Repositories struct {
	Transactor crud.Transactor

	// Users
	User               crud.Repository[user.User]
	Profile            crud.Repository[user.UserProfile]
	PasswordResetToken crud.Repository[user.PasswordResetToken]
	OAuthAccount       crud.Repository[auth.OAuthAccount]
	Session            crud.Repository[auth.Session]
	RefreshToken       crud.Repository[auth.RefreshToken]

	// Projects
	Project crud.Repository[project.Project]
	Entry   project.EntryRepository

	// Summaries
	ProjectSummary project.ProjectSummaryRepository
}

func ProvideRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		Transactor: crud.NewTransactor(db),

		User:               crud.NewRepository[user.User](db),
		Profile:            crud.NewRepository[user.UserProfile](db),
		PasswordResetToken: crud.NewRepository[user.PasswordResetToken](db),
		OAuthAccount:       crud.NewRepository[auth.OAuthAccount](db),
		Session:            crud.NewRepository[auth.Session](db),
		RefreshToken:       crud.NewRepository[auth.RefreshToken](db),

		Project:        crud.NewRepository[project.Project](db),
		Entry:          project.NewEntryRepository(db),
		ProjectSummary: project.NewProjectSummaryRepository(db),
	}
}
