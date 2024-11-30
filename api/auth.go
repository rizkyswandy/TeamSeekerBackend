package api

type User struct{
	ID string `json:"id"`
	Email string `json:"email"`
	Password string `json:"-"`
	Role string `json:"role"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type LoginRequest struct{
	Email string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type AuthResponse struct{
	Token string `json:"token"`
	User User `json:"user"`
}