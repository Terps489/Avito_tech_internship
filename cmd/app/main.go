package main

import (
	"log"

	httpTransport "github.com/terps489/avito_tech_internship/internal/http"
)

func main() {
	server := httpTransport.NewServer(":8080")

	if err := server.Run(); err != nil {
		log.Fatalf("server stopped with error: %v", err)
	}
}
