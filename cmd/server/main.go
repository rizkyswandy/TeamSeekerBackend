package main

import (
    "log"
    "github.com/rizkyswandy/TeamSeekerBackend/api"
    "github.com/rizkyswandy/TeamSeekerBackend/internal/config"
    "github.com/rizkyswandy/TeamSeekerBackend/internal/database/postgres"
    "github.com/joho/godotenv"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }
    
    cfg := config.LoadConfig()

    db, err := postgres.NewPostgresDB(cfg.DBConnString)
    if err != nil {
        log.Fatal(err)
    }

    server := api.NewAPIServer(db, cfg.JWTSecret)

    //Note: Keep for development
    // cfg.ServerPort = "3001"

    log.Printf("Server starting on port %s", cfg.ServerPort)
    if err := server.Start("0.0.0.0:" + cfg.ServerPort); err != nil {
        log.Fatal(err)
    }
}