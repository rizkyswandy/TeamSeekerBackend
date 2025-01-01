package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
    "github.com/rizkyswandy/TeamSeekerBackend/types"
)

func (s *APIServer) handleRegister(w http.ResponseWriter, r *http.Request) {
    var req types.RegisterRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if req.Email == "" || req.Password == "" {
        http.Error(w, "Email and password are required", http.StatusBadRequest)
        return
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "Failed to process password", http.StatusInternalServerError)
        return
    }

    user := &types.User{
        Email:    req.Email,
        Password: string(hashedPassword),
    }

    if err := s.db.CreateUser(user); err != nil {
        if err.Error() == "email already exists" {
            http.Error(w, "Email already registered", http.StatusBadRequest)
            return
        }
        http.Error(w, "Failed to create user", http.StatusInternalServerError)
        return
    }

    token, err := s.generateJWT(user)
    if err != nil {
        http.Error(w, "Failed to generate token", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(types.AuthResponse{
        Token: token,
        User:  *user,
    })
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) {
    var req types.LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if req.Email == "" || req.Password == "" {
        http.Error(w, "Email and password are required", http.StatusBadRequest)
        return
    }

    user, err := s.db.GetUserByEmail(req.Email)
    if err != nil {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }

    token, err := s.generateJWT(&user)
    if err != nil {
        http.Error(w, "Failed to generate token", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(types.AuthResponse{
        Token: token,
        User:  user,
    })
}

func (s *APIServer) generateJWT(user *types.User) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID,
        "email":   user.Email,
        "role":    user.Role,
        "exp":     time.Now().Add(time.Hour * 24).Unix(),
    })

    return token.SignedString(s.jwtSecret)
}