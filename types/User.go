package types

type User struct{
	ID string `json:"id"`
	Email string `json:"email"`
	Password string `json:"-"`
	Role string `json:"role"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}