package main

import(
	"log"
	"github.com/rizkyswandy/TeamSeekerBackend/api"
	"github.com/rizkyswandy/TeamSeekerBackend/internal/database/postgres"
)

func main() {
	connStr := "postgres://ilb:@localhost:5432/team_seeker?sslmode=disable"

	db, err := postgres.NewPostgresDB(connStr)
	if err != nil {
		log.Fatalf("Failed to initialized database : %v", err)
	}

	server :=  api.NewAPIServer(db)

	log.Println("Server starting on port 8080")
	if err := server.Start(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}