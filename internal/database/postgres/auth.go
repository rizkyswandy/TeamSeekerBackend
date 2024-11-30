package postgres

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/rizkyswandy/TeamSeekerBackend/api"
	"golang.org/x/crypto/bcrypt"
)

func (p *PostgresDB) CreateUser(user *api.User) error{
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return fmt.Errorf("error hashing password: %v", err)
	}

	query := `
		INSERT INTO users (email, password_hash, role)
		VALUES ($1,$2,$3)
		RETURNING id, created_at, updated_at`

	err = p.db.QueryRow(
		query,
		user.Email,
		string(hashedPassword),
		"user",
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err!= nil{
		log.Printf("Database error creating user: %v", err)
	}

	return nil
}

func (p *PostgresDB) GetUserByEmail(email string) (api.User, error){
	var user api.User

	query := `
		SELECT id, email, password_hash, role, created_at, updated_at
		FROM users WHERE email = $1`

	err := p.db.QueryRow(query, email).Scan(
		&user.ID,
        &user.Email,
        &user.Password,
        &user.Role,
        &user.CreatedAt,
        &user.UpdatedAt,
	)

	if err == sql.ErrNoRows{
		return api.User{}, fmt.Errorf("user not found")
	}

	if err != nil{
		return api.User{}, err
	}

	return user, nil
}

func (p *PostgresDB) GetUserByID(id string) (api.User, error){
	var user api.User

	userID, err := uuid.Parse(id)
	if err!= nil {
		return api.User{}, fmt.Errorf("invalid user ID format")
	}

	query := `
		SELECT id, email, password_hash, role, created_at, updated_at
		FROM users WHERE id = $1`

    err = p.db.QueryRow(query, userID).Scan(
        &user.ID,
        &user.Email,
        &user.Password,
        &user.Role,
        &user.CreatedAt,
        &user.UpdatedAt,
    )

	if err == sql.ErrNoRows{
		return api.User{}, fmt.Errorf("user not found")
	}

	if err != nil{
		return api.User{}, err
	}

	return user, nil
}
