// only for tests!
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync/atomic"
	"time"
)

const baseURL = "http://localhost:8080"

type CreatePRRequest struct {
	ID     string `json:"pull_request_id"`
	Name   string `json:"pull_request_name"`
	Author string `json:"author_id"`
}

type MergePRRequest struct {
	ID string `json:"pull_request_id"`
}

var prCounter uint64

func main() {
	var durationStr string
	var rps int
	flag.StringVar(&durationStr, "duration", "30s", "test duration (e.g. 30s, 1m)")
	flag.IntVar(&rps, "rps", 5, "requests per second (approximate)")
	flag.Parse()

	dur, err := time.ParseDuration(durationStr)
	if err != nil {
		log.Fatalf("invalid duration: %v", err)
	}

	log.Printf("Starting load test: duration=%s, rps=%d", dur, rps)

	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	authors := make([]string, 0, 80)
	for i := 1; i <= 80; i++ {
		authors = append(authors, fmt.Sprintf("u%d", i))
	}

	interval := time.Second / time.Duration(rps)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	stopAt := time.Now().Add(dur)

	var success uint64
	var failed uint64

	for now := range ticker.C {
		if now.After(stopAt) {
			break
		}

		go func() {
			author := authors[rand.Intn(len(authors))]

			id := atomic.AddUint64(&prCounter, 1)
			prID := fmt.Sprintf("lt-pr-%d", id)

			if err := scenarioCreateAndMerge(client, prID, author); err != nil {
				atomic.AddUint64(&failed, 1)
				return
			}
			atomic.AddUint64(&success, 1)
		}()
	}

	time.Sleep(3 * time.Second)

	log.Printf("Load test finished. Success=%d, Failed=%d", success, failed)
}

func scenarioCreateAndMerge(client *http.Client, prID string, author string) error {
	createReq := CreatePRRequest{
		ID:     prID,
		Name:   "Load test PR " + prID,
		Author: author,
	}

	body, err := json.Marshal(createReq)
	if err != nil {
		return fmt.Errorf("marshal create: %w", err)
	}

	resp, err := client.Post(baseURL+"/pullRequest/create", "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusConflict {
		return nil
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("create %s: unexpected status %s", prID, resp.Status)
	}

	mergeReq := MergePRRequest{ID: prID}
	mergeBody, err := json.Marshal(mergeReq)
	if err != nil {
		return fmt.Errorf("marshal merge: %w", err)
	}

	resp, err = client.Post(baseURL+"/pullRequest/merge", "application/json", bytes.NewReader(mergeBody))
	if err != nil {
		return fmt.Errorf("merge request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("merge %s: unexpected status %s", prID, resp.Status)
	}

	return nil
}
