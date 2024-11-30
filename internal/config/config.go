package config

import (
    "log"
    "os"
)

type Config struct {
    DBConnString string
    JWTSecret    []byte
    ServerPort   string
}

func LoadConfig() *Config {
    jwtSecret := os.Getenv("JWT_SECRET")
    if jwtSecret == "" {
        log.Fatal("JWT_SECRET environment variable is required")
    }

    dbConnString := os.Getenv("DB_CONN_STRING")
    if dbConnString == "" {
        dbConnString = "postgres://ilb:@localhost:5432/team_seeker?sslmode=disable"
        log.Println("Warning: Using default database connection. Set DB_CONN_STRING environment variable in production.")
    }

    serverPort := os.Getenv("SERVER_PORT")
    if serverPort == "" {
        serverPort = "8080"
        log.Println("Using default port 8080")
    }

    return &Config{
        DBConnString: dbConnString,
        JWTSecret:    []byte(jwtSecret),
        ServerPort:   serverPort,
    }
}