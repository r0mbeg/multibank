package domain

import "time"

type User struct {
	ID           int64     `db:"id"`            // INTEGER PRIMARY KEY AUTOINCREMENT
	Email        string    `db:"email"`         // TEXT NOT NULL UNIQUE
	FirstName    string    `db:"first_name"`    // TEXT NOT NULL
	LastName     string    `db:"last_name"`     // TEXT NOT NULL
	Patronymic   string    `db:"patronymic"`    // TEXT NOT NULL
	BirthDate    string    `db:"birthdate"`     // TEXT (YYYY-MM-DD)
	PasswordHash string    `db:"password_hash"` // TEXT NOT NULL
	CreatedAt    time.Time `db:"created_at"`    // TEXT (datetime('now'))
	UpdatedAt    time.Time `db:"updated_at"`    // TEXT (datetime('now'))
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
