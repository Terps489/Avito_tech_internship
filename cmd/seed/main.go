// only for tests!
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const baseURL = "http://localhost:8080"

type TeamMember struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type TeamPayload struct {
	TeamName string       `json:"team_name"`
	Members  []TeamMember `json:"members"`
}

func main() {
	log.Println("Seeding teams and users...")

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	userCounter := 1

	// 15* teams with 5-6 members each
	for teamIdx := 1; teamIdx <= 15; teamIdx++ {
		teamName := fmt.Sprintf("team-%02d", teamIdx)

		var membersCount int
		if teamIdx%2 == 0 {
			membersCount = 5
		} else {
			membersCount = 6
		}

		members := make([]TeamMember, 0, membersCount)
		for i := 0; i < membersCount; i++ {
			userID := fmt.Sprintf("u%d", userCounter)
			username := fmt.Sprintf("User%d", userCounter)

			isActive := true
			//inactivate every 7th user
			if userCounter%7 == 0 {
				isActive = false
			}

			members = append(members, TeamMember{
				UserID:   userID,
				Username: username,
				IsActive: isActive,
			})

			userCounter++
		}

		payload := TeamPayload{
			TeamName: teamName,
			Members:  members,
		}

		body, err := json.Marshal(payload)
		if err != nil {
			log.Fatalf("marshal team %s: %v", teamName, err)
		}

		resp, err := client.Post(baseURL+"/team/add", "application/json", bytes.NewReader(body))
		if err != nil {
			log.Fatalf("request team %s: %v", teamName, err)
		}
		_ = resp.Body.Close()

		if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
			log.Fatalf("team %s: unexpected status %s", teamName, resp.Status)
		}

		log.Printf("Seeded team %s with %d members", teamName, len(members))
	}

	log.Printf("Done. Total users created: %d", userCounter-1)
}
