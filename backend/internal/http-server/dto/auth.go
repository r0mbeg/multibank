// internal/http-server/dto/auth.go

package dto

type RegisterRequest struct {
	Email      string `json:"email" example:"user@example.com"`
	FirstName  string `json:"first_name" example:"Ivan"`
	LastName   string `json:"last_name" example:"Petrov"`
	Patronymic string `json:"patronymic" example:"Ivanovich"`
	BirthDate  string `json:"birthdate" example:"1990-01-15"` // YYYY-MM-DD
	Password   string `json:"password" example:"P@ssw0rd123"`
}

type LoginRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"P@ssw0rd123"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	ExpiresIn   int64  `json:"expires_in"  example:"3600"`
}
