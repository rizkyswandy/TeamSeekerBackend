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

	//Search Operations
	SearchProfiles(filter SearchFilters) ([]StudentProfile, error)
	GetAllProfiles() ([]StudentProfile, error)
}

func NewAPIServer(db Database) *APIServer {
	server := &APIServer{
		router: mux.NewRouter(),
		db:     db,
	}
	server.setupRoutes()
	return server
}

func (s *APIServer) setupRoutes() {

	s.router.Use(middleware.Logger)
    s.router.Use(middleware.CORS)

	s.router.HandleFunc("/api/profiles", s.handleCreateProfile).Methods("POST")
	s.router.HandleFunc("/api/profiles", s.handleGetAllProfiles).Methods("GET")
	s.router.HandleFunc("/api/profiles/search", s.handleSearchProfiles).Methods("GET")
	s.router.HandleFunc("/api/profiles/{id}", s.handleGetProfile).Methods("GET")
	s.router.HandleFunc("/api/profiles/{id}", s.handleUpdateProfile).Methods("PUT")
	s.router.HandleFunc("/api/profiles/{id}", s.handleDeleteProfile).Methods("DELETE")
}

func (s *APIServer) handleCreateProfile(writer http.ResponseWriter, request *http.Request) {
	var newProfile StudentProfile
	if err := json.NewDecoder(request.Body).Decode(&newProfile); err != nil {
		http.Error(writer, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := s.db.CreateProfile(&newProfile); err != nil {
		http.Error(writer, "Failed to create profile", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)

	json.NewEncoder(writer).Encode(newProfile)
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
    vars := mux.Vars(r)
    id := vars["id"]

    if id == "" {
        http.Error(w, "Profile ID is required", http.StatusBadRequest)
        return
    }

    var updatedProfile StudentProfile
    if err := json.NewDecoder(r.Body).Decode(&updatedProfile); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if err := s.db.UpdateProfile(id, &updatedProfile); err != nil {
        if err.Error() == "profile not found" {
            http.Error(w, "Profile not found", http.StatusNotFound)
            return
        }
        http.Error(w, "Failed to update profile", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
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
		http.Error(writer, "Failed to fetch profiels!", http.StatusInternalServerError)
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
