package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type ClearResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Deleted int    `json:"deleted"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Admin-Key")
	w.Header().Set("Content-Type", "application/json")

	// Handle preflight
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Only accept DELETE
	if r.Method != "DELETE" {
		json.NewEncoder(w).Encode(ClearResponse{
			Success: false,
			Message: "Method not allowed. Use DELETE.",
		})
		return
	}

	// Simple admin key check (optional - remove if you want it open)
	adminKey := r.Header.Get("X-Admin-Key")
	expectedKey := os.Getenv("ADMIN_KEY")
	
	// If ADMIN_KEY is set, require it; otherwise allow without key
	if expectedKey != "" && adminKey != expectedKey {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ClearResponse{
			Success: false,
			Message: "Unauthorized. Admin key required.",
		})
		return
	}

	// Connect to Vercel Postgres database
	postgresURL := os.Getenv("POSTGRES_URL")
	if postgresURL == "" {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ClearResponse{
			Success: false,
			Message: "Database configuration error: POSTGRES_URL not set",
		})
		return
	}

	db, err := sql.Open("postgres", postgresURL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ClearResponse{
			Success: false,
			Message: fmt.Sprintf("Database connection error: %v", err),
		})
		return
	}
	defer db.Close()

	// Get count before deletion
	var countBefore int
	err = db.QueryRow("SELECT COUNT(*) FROM game_scores").Scan(&countBefore)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ClearResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to count records: %v", err),
		})
		return
	}

	// Delete all test data
	result, err := db.Exec("DELETE FROM game_scores")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ClearResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to delete data: %v", err),
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()

	// Success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ClearResponse{
		Success: true,
		Message: fmt.Sprintf("Successfully deleted all test data (%d records)", rowsAffected),
		Deleted: int(rowsAffected),
	})
}

