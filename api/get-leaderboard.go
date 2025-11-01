package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func GetLeaderboardHandler(w http.ResponseWriter, r *http.Request) {
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
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Method not allowed. Use GET.",
		})
		return
	}

	// Connect to Vercel Postgres database
	postgresURL := os.Getenv("POSTGRES_URL")
	if postgresURL == "" {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Database configuration error: POSTGRES_URL not set",
		})
		return
	}

	db, err := sql.Open("postgres", postgresURL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Database connection error: %v", err),
		})
		return
	}
	defer db.Close()

	// Parse query parameters
	size := 10
	if s := r.URL.Query().Get("_size"); s != "" {
		if parsed, err := strconv.Atoi(s); err == nil && parsed > 0 {
			size = parsed
		}
	}

	sortBy := "final_net_worth"
	if s := r.URL.Query().Get("_sort_desc"); s != "" {
		sortBy = s
	} else if s := r.URL.Query().Get("_sort"); s != "" {
		sortBy = s
	}

	difficulty := r.URL.Query().Get("difficulty")

	// Build query
	query := "SELECT id, player_name, final_net_worth, roi, successful_exits, turns_played, difficulty, played_at FROM game_scores"
	whereClauses := []string{}
	args := []interface{}{}
	argIndex := 1

	if difficulty != "" && difficulty != "all" {
		whereClauses = append(whereClauses, fmt.Sprintf("difficulty = $%d", argIndex))
		args = append(args, difficulty)
		argIndex++
	}

	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Get total count
	countQuery := "SELECT COUNT(*) FROM game_scores"
	if len(whereClauses) > 0 {
		countQuery += " WHERE " + strings.Join(whereClauses, " AND ")
	}
	var totalCount int
	err = db.QueryRow(countQuery, args...).Scan(&totalCount)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Failed to get count: %v", err),
		})
		return
	}

	// Add sorting and limit
	query += fmt.Sprintf(" ORDER BY %s DESC LIMIT $%d", sortBy, argIndex)
	args = append(args, size)

	// Execute query
	rows, err := db.Query(query, args...)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Failed to query scores: %v", err),
		})
		return
	}
	defer rows.Close()

	// Fetch results
	var resultRows [][]interface{}
	for rows.Next() {
		var id, playerName, difficulty string
		var finalNetWorth int64
		var roi float64
		var successfulExits, turnsPlayed int
		var playedAt string

		err := rows.Scan(&id, &playerName, &finalNetWorth, &roi, &successfulExits, &turnsPlayed, &difficulty, &playedAt)
		if err != nil {
			continue
		}

		resultRows = append(resultRows, []interface{}{
			id,
			playerName,
			finalNetWorth,
			roi,
			successfulExits,
			turnsPlayed,
			difficulty,
			playedAt,
		})
	}

	// Return Datasette-compatible format
	response := map[string]interface{}{
		"rows":                      resultRows,
		"filtered_table_rows_count": totalCount,
		"database":                  "leaderboard",
		"table":                     "game_scores",
		"columns":                   []string{"id", "player_name", "final_net_worth", "roi", "successful_exits", "turns_played", "difficulty", "played_at"},
	}

	json.NewEncoder(w).Encode(response)
}

// Handler is the entry point for Vercel serverless function
func Handler(w http.ResponseWriter, r *http.Request) {
	GetLeaderboardHandler(w, r)
}
