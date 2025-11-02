package dto

type RegisterRequest struct {
	Email      string `json:"email"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Patronymic string `json:"patronymic"`
	BirthDate  string `json:"birthdate"` // YYYY-MM-DD
	Password   string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
