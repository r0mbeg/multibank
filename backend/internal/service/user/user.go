package user

import (
	"context"
	"database/sql"
	"errors"
	"multibank/backend/internal/domain"
	"time"
)

type Repository interface {
	GetByID(ctx context.Context, id int64) (domain.User, error)
}

type Service struct {
	repo    Repository
	timeout time.Duration
}

func New(repo Repository, timeout time.Duration) *Service {
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &Service{repo: repo, timeout: timeout}
}

func (s *Service) GetByID(ctx context.Context, id int64) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, ErrUserNotFound
		}
		return domain.User{}, err
	}
	return u, nil
}
