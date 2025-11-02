package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	authjwt "multibank/backend/internal/auth/jwt"
	"multibank/backend/internal/domain"
	"multibank/backend/internal/service/user"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailAlreadyUsed   = errors.New("email already used")
)

type Service struct {
	log *slog.Logger
	u   *user.Service
	jwt *authjwt.Manager
}

func New(log *slog.Logger, userSvc *user.Service, jwt *authjwt.Manager) *Service {
	return &Service{log: log, u: userSvc, jwt: jwt}
}

// RegisterInput struct for service layer
// repeats structure RegisterRequest from dto,
// but must be independent from http-layer
type RegisterInput struct {
	Email, FirstName, LastName, Patronymic, BirthDate, Password string
}

// Register registers a new user
func (a *Service) Register(ctx context.Context, in RegisterInput) (domain.User, error) {
	const op = "service.auth.Register"

	email := strings.ToLower(strings.TrimSpace(in.Email))
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
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
		return domain.User{}, fmt.Errorf("%s: %w", op, err)
	}
	created, err := a.u.GetByID(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("%s: %w", op, err)
	}
	return created, nil
}

func (a *Service) Login(ctx context.Context, email, password string) (string, error) {
	const op = "service.auth.Login"

	email = strings.ToLower(strings.TrimSpace(email))

	u, err := a.u.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	tok, _, err := a.jwt.Issue(u.ID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return tok, nil
}
