// internal/http-server/dto/user.go

package dto

import (
	"time"

	"multibank/backend/internal/domain"
)

type UserResponse struct {
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

func UserResponseFromDomain(d domain.User) UserResponse {
	return UserResponse{
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
