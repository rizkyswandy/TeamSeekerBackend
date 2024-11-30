package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (s *APIServer) handleRegister(w http.ResponseWriter, r *http.Request) {
    var req RegisterRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

	_, err := s.db.GetUserByEmail(req.Email)
    if err == nil {
        http.Error(w, "Email already registered", http.StatusBadRequest)
        return
    }

    user := &User{
        Email:    req.Email,
        Password: req.Password,
    }

    if err := s.db.CreateUser(user); err != nil {
        http.Error(w, "Failed to create user", http.StatusInternalServerError)
        return
    }

    token, err := s.generateJWT(user)  
    if err != nil {
        http.Error(w, "Failed to generate token", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(AuthResponse{
        Token: token,
        User:  *user,
    })
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) {
    var req LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
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
    json.NewEncoder(w).Encode(AuthResponse{
        Token: token,
        User:  user,
    })
}

func (s *APIServer) generateJWT(user *User) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID,
        "email":   user.Email,
        "role":    user.Role,
        "exp":     time.Now().Add(time.Hour * 24).Unix(),
    })

    return token.SignedString(s.jwtSecret)
}