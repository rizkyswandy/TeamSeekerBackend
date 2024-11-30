package api

import (
	"encoding/json"
	"net/http"

	"github.com/rizkyswandy/TeamSeekerBackend/middleware"
	"github.com/gorilla/mux"
)

type APIServer struct {
	router *mux.Router
	db     Database
	jwtSecret []byte
}

type StudentProfile struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Email        string   `json:"email"`
	Faculty      string   `json:"faculty"`
	FieldOfStudy string   `json:"field_of_study"`
	Semester     int      `json:"semester"`
	Skills       []string `json:"skills"`
	Focus        []string `json:"focus"`
	IsAvailable  bool     `json:"is_available"`
	CreatedAt    string   `json:"created_at"`
	UpdatedAt    string   `json:"updated_at"`
}

type SearchFilters struct {
	Faculty      string   `json:"faculty"`
	Skills       []string `json:"skills"`
	Focus        []string `json:"focus"`
	Availability bool     `json:"availability"`
}

type Database interface {
	CreateProfile(profile *StudentProfile) error
	GetProfile(id string) (StudentProfile, error)
	UpdateProfile(id string, profile *StudentProfile) error
	DeleteProfile(id string) error

	// MARK: Search Operations
	SearchProfiles(filter SearchFilters) ([]StudentProfile, error)
	GetAllProfiles() ([]StudentProfile, error)

	// MARK: Auth methods goes here
	CreateUser(user *User) error
	GetUserByEmail(email string) (User, error)
	GetUserByID(id string) (User, error)
}

func NewAPIServer(db Database, jwtSecret []byte) *APIServer {
    server := &APIServer{
        router:    mux.NewRouter(),
        db:        db,
        jwtSecret: jwtSecret,
    }
    server.setupRoutes()
    return server
}

// In setupRoutes()
func (s *APIServer) setupRoutes() {
    s.router.Use(middleware.Logger)
    s.router.Use(middleware.CORS)

    // Public routes (authentication)
    s.router.HandleFunc("/api/auth/register", s.handleRegister).Methods("POST")
    s.router.HandleFunc("/api/auth/login", s.handleLogin).Methods("POST")

    // Protected routes
    api := s.router.PathPrefix("/api").Subrouter()
    api.Use(middleware.Auth(s.jwtSecret))

    // Profile routes (all require authentication)
    api.HandleFunc("/profiles", s.handleCreateProfile).Methods("POST")
    api.HandleFunc("/profiles", s.handleGetAllProfiles).Methods("GET")
    api.HandleFunc("/profiles/search", s.handleSearchProfiles).Methods("GET")
    api.HandleFunc("/profiles/{id}", s.handleGetProfile).Methods("GET")
    api.HandleFunc("/profiles/{id}", s.handleUpdateProfile).Methods("PUT")
    api.HandleFunc("/profiles/{id}", s.handleDeleteProfile).Methods("DELETE")
}

func (s *APIServer) handleCreateProfile(w http.ResponseWriter, r *http.Request) {
    claims, err := middleware.GetUserFromContext(r.Context())
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    var newProfile StudentProfile
    if err := json.NewDecoder(r.Body).Decode(&newProfile); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    newProfile.Email = claims.Email

    if err := s.db.CreateProfile(&newProfile); err != nil {
        http.Error(w, "Failed to create profile", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(newProfile)
}

func (s *APIServer) handleGetProfile(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    if id == "" {
        http.Error(w, "Profile ID is required", http.StatusBadRequest)
        return
    }

    profile, err := s.db.GetProfile(id)
    if err != nil {
        if err.Error() == "profile not found" {
            http.Error(w, "Profile not found", http.StatusNotFound)
            return
        }
        http.Error(w, "Failed to get profile", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(profile)
}

func (s *APIServer) handleUpdateProfile(w http.ResponseWriter, r *http.Request) {
    claims, err := middleware.GetUserFromContext(r.Context())
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    vars := mux.Vars(r)
    id := vars["id"]

    // Get existing profile
    existingProfile, err := s.db.GetProfile(id)
    if err != nil {
        if err.Error() == "profile not found" {
            http.Error(w, "Profile not found", http.StatusNotFound)
            return
        }
        http.Error(w, "Failed to get profile", http.StatusInternalServerError)
        return
    }

    // Verify ownership
    if existingProfile.Email != claims.Email {
        http.Error(w, "Unauthorized to modify this profile", http.StatusForbidden)
        return
    }

    var updatedProfile StudentProfile
    if err := json.NewDecoder(r.Body).Decode(&updatedProfile); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Ensure email cannot be changed
    updatedProfile.Email = claims.Email

    if err := s.db.UpdateProfile(id, &updatedProfile); err != nil {
        http.Error(w, "Failed to update profile", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(updatedProfile)
}

func (s *APIServer) handleDeleteProfile(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["id"]

	if id == "" {
		http.Error(writer, "Profile ID is required", http.StatusBadRequest)
		return
	}

	if err := s.db.DeleteProfile(id); err != nil {
		http.Error(writer, "Failed to delete profile", http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusNoContent)
}

func (s *APIServer) handleGetAllProfiles(writer http.ResponseWriter, request *http.Request) {
	profiles, err := s.db.GetAllProfiles()
	if err != nil {
		http.Error(writer, "Failed to fetch profiles!", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(profiles)
}

func (s *APIServer) handleSearchProfiles(writer http.ResponseWriter, request *http.Request) {
	var filters SearchFilters
	if err := json.NewDecoder(request.Body).Decode(&filters); err != nil {
		http.Error(writer, "Invalid request body", http.StatusBadRequest)
		return
	}

	profiles, err := s.db.SearchProfiles(filters)
	if err != nil {
		http.Error(writer, "Failed to search profiles", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(profiles)
}

func (s *APIServer) Start(addr string) error {
	return http.ListenAndServe(addr, s.router)
}
