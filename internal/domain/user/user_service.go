package user

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/itsLeonB/go-crud"
	"github.com/itsLeonB/ungerr"
	"github.com/reflect-homini/stora/internal/core/mail"
	"github.com/reflect-homini/stora/internal/core/otel"
	"github.com/reflect-homini/stora/internal/domain/mapper"
)

type Service interface {
	// Public
	Me(ctx context.Context, id uuid.UUID) (UserResponse, error)

	// Internal
	GetByID(ctx context.Context, id uuid.UUID) (User, error)
	FindByEmail(ctx context.Context, email string) (User, error)
	CreateNew(ctx context.Context, request NewUserRequest) (User, error)
	Verify(ctx context.Context, id uuid.UUID, email string, name string, avatar string) (User, error)
	GeneratePasswordResetToken(ctx context.Context, userID uuid.UUID) (string, error)
	ResetPassword(ctx context.Context, userID uuid.UUID, email, resetToken, password string) (User, error)
}

type userServiceImpl struct {
	transactor             crud.Transactor
	userRepo               crud.Repository[User]
	passwordResetTokenRepo crud.Repository[PasswordResetToken]
	mailSvc                mail.Service
}

func NewUserService(
	transactor crud.Transactor,
	userRepo crud.Repository[User],
	passwordResetTokenRepo crud.Repository[PasswordResetToken],
	mailSvc mail.Service,
) *userServiceImpl {
	return &userServiceImpl{
		transactor,
		userRepo,
		passwordResetTokenRepo,
		mailSvc,
	}
}

func (us *userServiceImpl) CreateNew(ctx context.Context, request NewUserRequest) (User, error) {
	ctx, span := otel.Tracer.Start(ctx, "UserService.CreateNew")
	defer span.End()

	var response User
	err := us.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		newUser := User{
			Email:        strings.ToLower(request.Email),
			PasswordHash: request.Password,
			Profile: UserProfile{
				Name: request.Name,
				Avatar: sql.NullString{
					String: request.Avatar,
					Valid:  request.Avatar != "",
				},
			},
		}

		if request.VerifyNow {
			newUser.VerifiedAt = sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			}
		}

		user, err := us.userRepo.Insert(ctx, newUser)
		if err != nil {
			return err
		}

		response = user
		return nil
	})

	return response, err
}

func (us *userServiceImpl) FindByEmail(ctx context.Context, email string) (User, error) {
	ctx, span := otel.Tracer.Start(ctx, "UserService.FindByEmail")
	defer span.End()

	userSpec := crud.Specification[User]{}
	userSpec.Model.Email = strings.ToLower(email)
	userSpec.PreloadRelations = []string{"Profile"}
	return us.userRepo.FindFirst(ctx, userSpec)
}

func (us *userServiceImpl) Verify(ctx context.Context, id uuid.UUID, email string, name string, avatar string) (User, error) {
	ctx, span := otel.Tracer.Start(ctx, "UserService.Verify")
	defer span.End()

	user, err := us.GetByID(ctx, id)
	if err != nil {
		return User{}, err
	}
	if user.Email != email {
		return User{}, ungerr.Unknown("email does not match")
	}

	user.VerifiedAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	user.Profile.Name = name
	user.Profile.Avatar = sql.NullString{
		String: avatar,
		Valid:  avatar != "",
	}

	return us.userRepo.Update(ctx, user)
}

func (us *userServiceImpl) GeneratePasswordResetToken(ctx context.Context, userID uuid.UUID) (string, error) {
	ctx, span := otel.Tracer.Start(ctx, "UserService.GeneratePasswordResetToken")
	defer span.End()

	token, err := us.generateRandomToken(255)
	if err != nil {
		return "", err
	}
	resetToken := PasswordResetToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}
	if _, err := us.passwordResetTokenRepo.Insert(ctx, resetToken); err != nil {
		return "", err
	}
	return token, nil
}

func (us *userServiceImpl) ResetPassword(ctx context.Context, userID uuid.UUID, email, resetToken, password string) (User, error) {
	ctx, span := otel.Tracer.Start(ctx, "UserService.ResetPassword")
	defer span.End()

	var resp User
	err := us.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		spec := crud.Specification[User]{}
		spec.Model.ID = userID
		spec.Model.Email = email
		spec.PreloadRelations = []string{"PasswordResetTokens"}
		user, err := us.getBySpec(ctx, spec)
		if err != nil {
			return err
		}

		if !us.validateToken(user.PasswordResetTokens, resetToken) {
			return ungerr.BadRequestError("invalid or expired reset token")
		}

		user.PasswordHash = password
		resp, err = us.userRepo.Update(ctx, user)
		if err != nil {
			return err
		}

		if err = us.passwordResetTokenRepo.DeleteMany(ctx, user.PasswordResetTokens); err != nil {
			return err
		}

		return nil
	})
	return resp, err
}

func (us *userServiceImpl) validateToken(resetTokens []PasswordResetToken, resetToken string) bool {
	if len(resetTokens) < 1 {
		return false
	}
	if len(resetTokens) == 1 {
		return resetTokens[0].IsValid() && resetTokens[0].Token == resetToken
	}
	sort.Slice(resetTokens, func(i, j int) bool {
		return resetTokens[i].CreatedAt.After(resetTokens[j].CreatedAt)
	})
	return resetTokens[0].IsValid() && resetTokens[0].Token == resetToken
}

func (us *userServiceImpl) generateRandomToken(length int) (string, error) {
	tokenBytes := make([]byte, length)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", ungerr.Wrap(err, "error generating random token")
	}
	return base64.URLEncoding.EncodeToString(tokenBytes)[:length], nil
}

func (us *userServiceImpl) GetByID(ctx context.Context, id uuid.UUID) (User, error) {
	ctx, span := otel.Tracer.Start(ctx, "UserService.GetByID")
	defer span.End()

	spec := crud.Specification[User]{}
	spec.Model.ID = id
	spec.PreloadRelations = []string{"Profile"}
	return us.getBySpec(ctx, spec)
}

func (us *userServiceImpl) Me(ctx context.Context, id uuid.UUID) (UserResponse, error) {
	ctx, span := otel.Tracer.Start(ctx, "UserService.Me")
	defer span.End()

	user, err := us.GetByID(ctx, id)
	if err != nil {
		return UserResponse{}, err
	}

	return UserResponse{
		BaseDTO: mapper.BaseToDTO(user.BaseEntity),
		Email:   user.Email,
		Profile: ProfileResponse{
			BaseDTO: mapper.BaseToDTO(user.Profile.BaseEntity),
			UserID:  user.Profile.UserID,
			Name:    user.Profile.Name,
			Avatar:  user.Profile.Avatar.String,
		},
	}, nil
}

func (us *userServiceImpl) getBySpec(ctx context.Context, spec crud.Specification[User]) (User, error) {
	user, err := us.userRepo.FindFirst(ctx, spec)
	if err != nil {
		return User{}, err
	}
	if user.IsZero() {
		return User{}, ungerr.NotFoundError("user is not found")
	}
	return user, nil
}
