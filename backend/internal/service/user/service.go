package user

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"multibank/backend/internal/domain"
	"multibank/backend/internal/logger"
	"multibank/backend/internal/storage"
)

type Service struct {
	log  *slog.Logger
	repo Repository
}

type Repository interface {
	Create(ctx context.Context, u domain.User) (int64, error)
	GetByID(ctx context.Context, id int64) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
}

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrEmailAlreadyUsed = errors.New("email already used")
)

// New created a new instance of User Service
func New(log *slog.Logger, repo Repository) *Service {
	return &Service{log: log, repo: repo}
}

// Create creates new user (password is already hashed)
// Returns ErrEmailAlreadyUsed if the email is already used
func (s *Service) Create(ctx context.Context, u domain.User) (int64, error) {
	const op = "service.user.Create"

	log := s.log.With(
		slog.String("op", op),
		slog.String("email", u.Email),
	)

	log.Info("creating new user")

	id, err := s.repo.Create(ctx, u)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("email already used", logger.Err(err))
			return 0, fmt.Errorf("%s: %w", op, ErrEmailAlreadyUsed)
		}
		log.Error("failed to create user", logger.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user created")

	return id, nil
}

// GetByID gets a User with given ID
// Returns ErrUserNotFound if there is no such User
func (s *Service) GetByID(ctx context.Context, id int64) (domain.User, error) {
	const op = "service.user.GetByID"

	log := s.log.With(
		slog.String("op", op),
		slog.Int64("ID", id),
	)

	log.Info("getting user by ID")

	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found", logger.Err(err))
			return domain.User{}, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user successfully got by ID")

	return u, nil
}

// GetByEmail gets a User with given email
// Returns ErrUserNotFound if there is no such User
func (s *Service) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	const op = "service.user.GetByEmail"

	log := s.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)
	log.Info("getting user by email")

	u, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found", logger.Err(err))
			return domain.User{}, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user successfully got by email")

	return u, nil

}
