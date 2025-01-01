package postgres

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
    "github.com/lib/pq"
	"github.com/rizkyswandy/TeamSeekerBackend/types"
)

func (p *PostgresDB) CreateUser(user *types.User) error {
    userID := uuid.New()
    
    query := `
        INSERT INTO users (id, email, password_hash, role, created_at, updated_at)
        VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        RETURNING id, role, created_at, updated_at`

    err := p.db.QueryRow(
        query,
        userID,
        user.Email,
        user.Password, 
        "user",
    ).Scan(&user.ID, &user.Role, &user.CreatedAt, &user.UpdatedAt)

    if err != nil {
        if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
            return fmt.Errorf("email already exists")
        }
        log.Printf("Database error creating user: %v", err)
        return fmt.Errorf("error creating user: %v", err)
    }

    return nil
}

func (p *PostgresDB) GetUserByEmail(email string) (types.User, error) {
    var user types.User

    query := `
        SELECT id, email, password_hash, role, created_at, updated_at
        FROM users 
        WHERE email = $1`

    err := p.db.QueryRow(query, email).Scan(
        &user.ID,
        &user.Email,
        &user.Password, 
        &user.Role,
        &user.CreatedAt,
        &user.UpdatedAt,
    )

    if err == sql.ErrNoRows {
        return types.User{}, fmt.Errorf("user not found")
    }

    if err != nil {
        log.Printf("Database error getting user by email: %v", err)
        return types.User{}, fmt.Errorf("error retrieving user")
    }

    return user, nil
}

func (p *PostgresDB) GetUserByID(id string) (types.User, error) {
    var user types.User

    userID, err := uuid.Parse(id)
    if err != nil {
        return types.User{}, fmt.Errorf("invalid user ID format")
    }

    query := `
        SELECT id, email, password_hash, role, created_at, updated_at
        FROM users 
        WHERE id = $1`

    err = p.db.QueryRow(query, userID).Scan(
        &user.ID,
        &user.Email,
        &user.Password, 
        &user.Role,
        &user.CreatedAt,
        &user.UpdatedAt,
    )

    if err == sql.ErrNoRows {
        return types.User{}, fmt.Errorf("user not found")
    }

    if err != nil {
        log.Printf("Database error getting user by ID: %v", err)
        return types.User{}, fmt.Errorf("error retrieving user")
    }

    return user, nil
}