// internal/service/auth/service.go

package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"multibank/backend/internal/domain"
	"multibank/backend/internal/logger"
	authjwt "multibank/backend/internal/service/auth/jwt"
	usrsvc "multibank/backend/internal/service/user"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// RegisterInput struct for service layer
// repeats structure RegisterRequest from dto,
// but must be independent from http-layer
type RegisterInput struct {
	Email      string
	FirstName  string
	LastName   string
	Patronymic string
	BirthDate  string
	Password   string
}

type Auth struct {
	log *slog.Logger
	u   Service
	jwt *authjwt.Manager
}

type Service interface {
	Create(ctx context.Context, user domain.User) (int64, error)
	GetByID(ctx context.Context, id int64) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
}

func New(log *slog.Logger, userSvc Service, jwt *authjwt.Manager) *Auth {
	return &Auth{log: log, u: userSvc, jwt: jwt}
}

// Register registers a new user
func (a *Auth) Register(ctx context.Context, in RegisterInput) (domain.User, error) {
	const op = "service.auth.Register"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", in.Email),
	)

	log.Info("registering new user")

	email := strings.ToLower(strings.TrimSpace(in.Email))
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", logger.Err(err))
		return domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	u := domain.User{
		Email:        email,
		FirstName:    strings.TrimSpace(in.FirstName),
		LastName:     strings.TrimSpace(in.LastName),
		Patronymic:   strings.TrimSpace(in.Patronymic),
		BirthDate:    strings.TrimSpace(in.BirthDate),
		PasswordHash: string(hash),
	}

	id, err := a.u.Create(ctx, u)
	if err != nil {

		if errors.Is(err, usrsvc.ErrEmailAlreadyUsed) {
			log.Error("email already used", logger.Err(err))
			return domain.User{}, usrsvc.ErrEmailAlreadyUsed
		}

		return domain.User{}, fmt.Errorf("%s: %w", op, err)
	}
	created, err := a.u.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, usrsvc.ErrUserNotFound) {
			log.Error("user not found", logger.Err(err))
			return domain.User{}, ErrInvalidCredentials
		}
		return domain.User{}, fmt.Errorf("%s: %w", op, err)
	}
	return created, nil
}

func (a *Auth) Login(ctx context.Context, email, password string) (string, error) {
	const op = "service.auth.Login"

	email = strings.ToLower(strings.TrimSpace(email))

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("attempting to login user")

	u, err := a.u.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, usrsvc.ErrUserNotFound) {
			log.Warn("user not found", logger.Err(err))
			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		log.Error("failed to get user", logger.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		log.Info("invalid credentials", logger.Err(err))
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	tok, _, err := a.jwt.Issue(u.ID)
	if err != nil {
		log.Error("failed to create token", logger.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("successfully logged in")
	return tok, nil
}
