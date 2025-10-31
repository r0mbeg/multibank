package domain

import "time"

type User struct {
	ID           int64     `db:"id"            json:"id"`
	Email        string    `db:"email"         json:"email"`
	FirstName    string    `db:"first_name"    json:"first_name"`
	LastName     string    `db:"last_name"     json:"last_name"`
	Patronymic   string    `db:"patronymic"    json:"patronymic"`
	BirthDate    string    `db:"birthdate"     json:"birthdate"` // YYYY-MM-DD
	PasswordHash string    `db:"password_hash" json:"-"`         // не отдаём наружу
	CreatedAt    time.Time `db:"created_at"    json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"    json:"updated_at"`
}

func NewUser(email, first, last, patr, birthDate, passHash string) User {
	return User{
		Email:        email,
		FirstName:    first,
		LastName:     last,
		Patronymic:   patr,
		BirthDate:    birthDate, // ожидаем "YYYY-MM-DD"
		PasswordHash: passHash,  // уже захешированный
	}
}
