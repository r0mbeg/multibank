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
	IsAdmin    bool      `json:"is_admin"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func UserFromDomain(d domain.User) User {
	return User{
		ID:         d.ID,
		Email:      d.Email,
		FirstName:  d.FirstName,
		LastName:   d.LastName,
		Patronymic: d.Patronymic,
		BirthDate:  d.BirthDate,
		IsAdmin:    d.IsAdmin,
		CreatedAt:  d.CreatedAt,
		UpdatedAt:  d.UpdatedAt,
	}
}
