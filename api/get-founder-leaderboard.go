package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

type FounderScore struct {
	ID                string  `json:"id"`
	PlayerName        string  `json:"player_name"`
	FinalValuation    int64   `json:"final_valuation"`
	FounderEquity     float64 `json:"founder_equity"`
	FounderPayout     int64   `json:"founder_payout"`
	ExitType          string  `json:"exit_type"`
	ExitMonth         int     `json:"exit_month"`
	MaxARR            int64   `json:"max_arr"`
	StartupTemplate   string  `json:"startup_template"`
	FundingRaised     int64   `json:"funding_raised"`
	CustomersAcquired int     `json:"customers_acquired"`
	PlayedAt          string  `json:"played_at"`
}

type LeaderboardResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message,omitempty"`
	Count   int            `json:"count"`
	Scores  []FounderScore `json:"scores"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	// Handle preflight
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Only accept GET
	if r.Method != "GET" {
		json.NewEncoder(w).Encode(LeaderboardResponse{
			Success: false,
			Message: "Method not allowed. Use GET.",
		})
		return
	}

	// Parse query parameters
	query := r.URL.Query()
	limit := 10
	if limitStr := query.Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	sortBy := query.Get("sort_by")
	if sortBy == "" {
		sortBy = "payout" // Default: sort by founder payout
	}

	template := query.Get("template")
	exitType := query.Get("exit_type")

	// Connect to Vercel Postgres database
	postgresURL := os.Getenv("POSTGRES_URL")
	if postgresURL == "" {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(LeaderboardResponse{
			Success: false,
			Message: "Database configuration error: POSTGRES_URL not set",
		})
		return
	}

	db, err := sql.Open("postgres", postgresURL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(LeaderboardResponse{
			Success: false,
			Message: fmt.Sprintf("Database connection error: %v", err),
		})
		return
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(LeaderboardResponse{
			Success: false,
			Message: fmt.Sprintf("Database ping failed: %v", err),
		})
		return
	}

	// Build query
	var orderColumn string
	switch sortBy {
	case "valuation":
		orderColumn = "final_valuation"
	case "arr":
		orderColumn = "max_arr"
	case "equity":
		orderColumn = "founder_equity"
	default:
		orderColumn = "founder_payout"
	}

	querySQL := fmt.Sprintf(`
		SELECT 
			id, player_name, final_valuation, founder_equity, founder_payout,
			exit_type, exit_month, max_arr, startup_template, funding_raised,
			customers_acquired, played_at
		FROM founder_scores
		WHERE 1=1
	`)

	args := []interface{}{}
	argCount := 1

	if template != "" {
		querySQL += fmt.Sprintf(" AND startup_template = $%d", argCount)
		args = append(args, template)
		argCount++
	}

	if exitType != "" {
		querySQL += fmt.Sprintf(" AND exit_type = $%d", argCount)
		args = append(args, exitType)
		argCount++
	}

	querySQL += fmt.Sprintf(" ORDER BY %s DESC LIMIT $%d", orderColumn, argCount)
	args = append(args, limit)

	// Execute query
	rows, err := db.Query(querySQL, args...)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(LeaderboardResponse{
			Success: false,
			Message: fmt.Sprintf("Query failed: %v", err),
		})
		return
	}
	defer rows.Close()

	// Parse results
	var scores []FounderScore
	for rows.Next() {
		var score FounderScore
		err := rows.Scan(
			&score.ID,
			&score.PlayerName,
			&score.FinalValuation,
			&score.FounderEquity,
			&score.FounderPayout,
			&score.ExitType,
			&score.ExitMonth,
			&score.MaxARR,
			&score.StartupTemplate,
			&score.FundingRaised,
			&score.CustomersAcquired,
			&score.PlayedAt,
		)
		if err != nil {
			continue
		}
		scores = append(scores, score)
	}

	// Return results
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(LeaderboardResponse{
		Success: true,
		Count:   len(scores),
		Scores:  scores,
	})
}

