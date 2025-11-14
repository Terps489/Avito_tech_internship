package main

import (
	"log"

	"github.com/terps489/avito_tech_internship/internal/app"
	httpTransport "github.com/terps489/avito_tech_internship/internal/http"
	"github.com/terps489/avito_tech_internship/internal/repository/postgres"
)

func main() {
	db, err := postgres.NewFromEnv()
	if err != nil {
		log.Fatalf("failed to init postgres: %v", err)
	}
	defer db.Close()

	userRepo := postgres.NewUserRepository(db)
	teamRepo := postgres.NewTeamRepository(db)
	prRepo := postgres.NewPullRequestRepository(db)

	service := app.NewService(userRepo, teamRepo, prRepo)

	server := httpTransport.NewServer(":8080", service)

	if err := server.Run(); err != nil {
		log.Fatalf("server stopped with error: %v", err)
	}
}
