package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type FounderScoreSubmission struct {
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
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	ID      string `json:"id,omitempty"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	// Handle preflight
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Only accept POST
	if r.Method != "POST" {
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Method not allowed. Use POST.",
		})
		return
	}

	// Parse request
	var submission FounderScoreSubmission
	if err := json.NewDecoder(r.Body).Decode(&submission); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Invalid JSON: " + err.Error(),
		})
		return
	}

	// Validate input
	if submission.PlayerName == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Player name is required",
		})
		return
	}

	if submission.StartupTemplate == "" {
		submission.StartupTemplate = "Unknown"
	}

	if submission.ExitType == "" {
		submission.ExitType = "time_limit"
	}

	// Connect to Vercel Postgres database
	postgresURL := os.Getenv("POSTGRES_URL")
	if postgresURL == "" {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Database configuration error: POSTGRES_URL not set",
		})
		return
	}

	db, err := sql.Open("postgres", postgresURL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: fmt.Sprintf("Database connection error: %v", err),
		})
		return
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: fmt.Sprintf("Database ping failed: %v", err),
		})
		return
	}

	// Create table if not exists (Postgres syntax)
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS founder_scores (
		id TEXT PRIMARY KEY,
		player_name TEXT NOT NULL,
		final_valuation BIGINT NOT NULL,
		founder_equity REAL NOT NULL,
		founder_payout BIGINT NOT NULL,
		exit_type TEXT NOT NULL,
		exit_month INTEGER NOT NULL,
		max_arr BIGINT NOT NULL,
		startup_template TEXT NOT NULL,
		funding_raised BIGINT NOT NULL,
		customers_acquired INTEGER NOT NULL,
		played_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_founder_valuation ON founder_scores(final_valuation DESC);
	CREATE INDEX IF NOT EXISTS idx_founder_payout ON founder_scores(founder_payout DESC);
	CREATE INDEX IF NOT EXISTS idx_founder_arr ON founder_scores(max_arr DESC);
	CREATE INDEX IF NOT EXISTS idx_founder_player ON founder_scores(player_name);
	CREATE INDEX IF NOT EXISTS idx_founder_exit ON founder_scores(exit_type);
	CREATE INDEX IF NOT EXISTS idx_founder_template ON founder_scores(startup_template);
	`
	
	_, err = db.Exec(createTableSQL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Failed to initialize database",
		})
		return
	}

	// Generate UUID for score
	scoreID := uuid.New().String()

	// Insert score (Postgres uses $1, $2, etc. instead of ?)
	insertSQL := `
		INSERT INTO founder_scores (
			id, player_name, final_valuation, founder_equity, founder_payout,
			exit_type, exit_month, max_arr, startup_template, funding_raised,
			customers_acquired, played_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err = db.Exec(insertSQL,
		scoreID,
		submission.PlayerName,
		submission.FinalValuation,
		submission.FounderEquity,
		submission.FounderPayout,
		submission.ExitType,
		submission.ExitMonth,
		submission.MaxARR,
		submission.StartupTemplate,
		submission.FundingRaised,
		submission.CustomersAcquired,
		time.Now().UTC(),
	)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: fmt.Sprintf("Failed to save score: %v", err),
		})
		return
	}

	// Success response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: "Founder score submitted successfully!",
		ID:      scoreID,
	})
}

