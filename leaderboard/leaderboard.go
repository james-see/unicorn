package leaderboard

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	// Update this URL after deploying to Vercel
	// Production URL: https://unicorn.vercel.app/api/submit-score
	// Preview URL (temporary): https://unicorn-g5xp1fzsv-james-campbells-projects-98ba50e1.vercel.app/api/submit-score
	DefaultAPIEndpoint = "https://unicorn-g5xp1fzsv-james-campbells-projects-98ba50e1.vercel.app/api/submit-score"
)

// ScoreSubmission represents a score to be submitted to the global leaderboard
type ScoreSubmission struct {
	PlayerName      string  `json:"player_name"`
	FinalNetWorth   int64   `json:"final_net_worth"`
	ROI             float64 `json:"roi"`
	SuccessfulExits int     `json:"successful_exits"`
	TurnsPlayed     int     `json:"turns_played"`
	Difficulty      string  `json:"difficulty"`
}

// APIResponse represents the response from the API
type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	ID      string `json:"id,omitempty"`
}

// SubmitScore submits a score to the global leaderboard
func SubmitScore(submission ScoreSubmission, apiURL string) error {
	if apiURL == "" {
		apiURL = DefaultAPIEndpoint
	}

	// Marshal the submission to JSON
	jsonData, err := json.Marshal(submission)
	if err != nil {
		return fmt.Errorf("failed to encode submission: %v", err)
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Create request
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Unicorn-Game/1.0")

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}

	// Parse response
	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return fmt.Errorf("failed to parse response: %v", err)
	}

	// Check for success
	if !apiResp.Success {
		return fmt.Errorf("API error: %s", apiResp.Message)
	}

	return nil
}

// IsAPIAvailable checks if the leaderboard API is reachable
func IsAPIAvailable(apiURL string) bool {
	if apiURL == "" {
		apiURL = DefaultAPIEndpoint
	}

	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	resp, err := client.Get(apiURL)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return true
}
