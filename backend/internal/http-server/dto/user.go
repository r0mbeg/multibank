// internal/http-server/dto/user.go
package dto

import (
	"time"

	"multibank/backend/internal/domain"
)

type User struct {
	ID         int64     `json:"id"`
	Email      string    `json:"email"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Patronymic string    `json:"patronymic"`
	BirthDate  string    `json:"birthdate"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type RegisterRequest struct {
	Email      string `json:"email"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Patronymic string `json:"patronymic"`
	BirthDate  string `json:"birthdate"`
	Password   string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func UserFromDomain(d domain.User) User {
	return User{
		ID:         d.ID,
		Email:      d.Email,
		FirstName:  d.FirstName,
		LastName:   d.LastName,
		Patronymic: d.Patronymic,
		BirthDate:  d.BirthDate,
		CreatedAt:  d.CreatedAt,
		UpdatedAt:  d.UpdatedAt,
	}
}
