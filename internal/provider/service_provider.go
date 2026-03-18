package provider

import (
	"net/http"

	"github.com/itsLeonB/sekure"
	"github.com/reflect-homini/stora/internal/core/config"
	"github.com/reflect-homini/stora/internal/domain/auth"
	"github.com/reflect-homini/stora/internal/domain/entry"
	"github.com/reflect-homini/stora/internal/domain/project"
	"github.com/reflect-homini/stora/internal/domain/projectdetails"
	"github.com/reflect-homini/stora/internal/domain/summary"
	"github.com/reflect-homini/stora/internal/domain/user"
)

type Services struct {
	// Auth
	Auth    auth.Service
	OAuth   auth.OAuthService
	Session auth.SessionService

	// Users
	User user.Service

	// Projects
	Project        project.Service
	Entry          entry.Service
	ProjectDetails projectdetails.Service

	// Summaries
	ProjectSummary summary.ProjectSummaryService
}

func ProvideServices(
	repos *Repositories,
	coreSvc *CoreServices,
) *Services {
	authConfig := config.Global.Auth
	appConfig := config.Global.App

	jwt := sekure.NewJwtService(authConfig.Issuer, authConfig.SecretKey, authConfig.TokenDuration)
	user := user.NewUserService(repos.Transactor, repos.User, repos.PasswordResetToken, coreSvc.Mail)
	session := auth.NewSessionService(jwt, user, repos.Transactor, repos.Session, repos.RefreshToken)

	entrySvc := entry.NewService(repos.Entry)
	projectSvc := project.NewService(repos.Project, entrySvc)

	entrySummarizer := summary.NewEntrySummarizerService(coreSvc.LLM)

	return &Services{
		Auth:    auth.NewAuthService(jwt, repos.Transactor, user, coreSvc.Mail, appConfig.RegisterVerificationUrl, appConfig.ResetPasswordUrl, authConfig.HashCost, session),
		OAuth:   auth.NewOAuthService(repos.Transactor, repos.OAuthAccount, coreSvc.State, user, &http.Client{Timeout: appConfig.Timeout}, session),
		Session: session,

		User: user,

		Project:        projectSvc,
		Entry:          entrySvc,
		ProjectDetails: projectdetails.NewService(repos.ProjectSummary, repos.Entry, projectSvc),

		ProjectSummary: summary.NewProjectSummaryService(repos.ProjectSummary, repos.Project, repos.Entry, entrySummarizer),
	}
}
