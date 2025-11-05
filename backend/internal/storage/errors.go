// internal/storage/errors.go

package storage

import (
	"errors"
)

var (
	ErrUserExists    = errors.New("user already exists")
	ErrUserNotFound  = errors.New("user not found")
	ErrBanksNotFound = errors.New("banks not found")
)
