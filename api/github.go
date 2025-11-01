package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// GitHubIssue represents a GitHub issue for score submission
type GitHubIssue struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Labels []string `json:"labels"`
}

// ScoreSubmission represents a game score to submit
type ScoreSubmission struct {
	PlayerName      string
	FinalNetWorth   int64
	ROI             float64
	SuccessfulExits int
	TurnsPlayed     int
	Difficulty      string
	PlayedAt        time.Time
}

// SubmitScoreToGitHub submits a score to GitHub via Issue API
// This creates a GitHub issue that can be processed by GitHub Actions
func SubmitScoreToGitHub(repoOwner, repoName, token string, score ScoreSubmission) error {
	if token == "" {
		return fmt.Errorf("GitHub token not configured. Set GITHUB_TOKEN environment variable")
	}

	issueTitle := fmt.Sprintf("Score Submission: %s - $%s (%.1f%% ROI)", 
		score.PlayerName, 
		formatMoney(score.FinalNetWorth), 
		score.ROI)

	issueBody := fmt.Sprintf(`# Score Submission

**Player:** %s
**Final Net Worth:** $%s
**ROI:** %.2f%%
**Successful Exits:** %d
**Turns Played:** %d
**Difficulty:** %s
**Played At:** %s

<!-- Score Data (JSON) -->
`+"```json\n"+`{
  "player_name": "%s",
  "final_net_worth": %d,
  "roi": %.2f,
  "successful_exits": %d,
  "turns_played": %d,
  "difficulty": "%s",
  "played_at": "%s"
}
`+"```\n", 
		score.PlayerName,
		formatMoney(score.FinalNetWorth),
		score.ROI,
		score.SuccessfulExits,
		score.TurnsPlayed,
		score.Difficulty,
		score.PlayedAt.Format(time.RFC3339),
		score.PlayerName,
		score.FinalNetWorth,
		score.ROI,
		score.SuccessfulExits,
		score.TurnsPlayed,
		score.Difficulty,
		score.PlayedAt.Format(time.RFC3339),
	)

	issue := GitHubIssue{
		Title: issueTitle,
		Body:  issueBody,
		Labels: []string{"score-submission", "leaderboard"},
	}

	jsonData, err := json.Marshal(issue)
	if err != nil {
		return fmt.Errorf("failed to marshal issue: %v", err)
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues", repoOwner, repoName)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to submit score: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("GitHub API error: %d - %s", resp.StatusCode, string(body))
	}

	return nil
}

// formatMoney formats money with commas
func formatMoney(amount int64) string {
	abs := amount
	if abs < 0 {
		abs = -abs
	}

	s := fmt.Sprintf("%d", abs)
	if len(s) <= 3 {
		if amount < 0 {
			return "-" + s
		}
		return s
	}

	result := ""
	for i, digit := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result += ","
		}
		result += string(digit)
	}

	if amount < 0 {
		return "-" + result
	}
	return result
}
