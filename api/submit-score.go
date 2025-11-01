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

type ScoreSubmission struct {
	PlayerName      string  `json:"player_name"`
	FinalNetWorth   int64   `json:"final_net_worth"`
	ROI             float64 `json:"roi"`
	SuccessfulExits int     `json:"successful_exits"`
	TurnsPlayed     int     `json:"turns_played"`
	Difficulty      string  `json:"difficulty"`
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
	var submission ScoreSubmission
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

	if submission.Difficulty == "" {
		submission.Difficulty = "Medium"
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
	CREATE TABLE IF NOT EXISTS game_scores (
		id TEXT PRIMARY KEY,
		player_name TEXT NOT NULL,
		final_net_worth BIGINT NOT NULL,
		roi REAL NOT NULL,
		successful_exits INTEGER NOT NULL,
		turns_played INTEGER NOT NULL,
		difficulty TEXT NOT NULL,
		played_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_net_worth ON game_scores(final_net_worth DESC);
	CREATE INDEX IF NOT EXISTS idx_roi ON game_scores(roi DESC);
	CREATE INDEX IF NOT EXISTS idx_player ON game_scores(player_name);
	CREATE INDEX IF NOT EXISTS idx_difficulty ON game_scores(difficulty);
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
		INSERT INTO game_scores (id, player_name, final_net_worth, roi, successful_exits, turns_played, difficulty, played_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err = db.Exec(insertSQL,
		scoreID,
		submission.PlayerName,
		submission.FinalNetWorth,
		submission.ROI,
		submission.SuccessfulExits,
		submission.TurnsPlayed,
		submission.Difficulty,
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
		Message: "Score submitted successfully!",
		ID:      scoreID,
	})
}

