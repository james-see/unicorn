package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// ScoreData represents score data for JSON export
type ScoreData struct {
	PlayerName      string  `json:"player_name"`
	FinalNetWorth   int64   `json:"final_net_worth"`
	ROI             float64 `json:"roi"`
	SuccessfulExits int     `json:"successful_exits"`
	TurnsPlayed     int     `json:"turns_played"`
	Difficulty      string  `json:"difficulty"`
	PlayedAt        string  `json:"played_at"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run export_db.go <database.db> [output.json]")
		os.Exit(1)
	}

	dbPath := os.Args[1]
	outputPath := "leaderboard.json"
	if len(os.Args) > 2 {
		outputPath = os.Args[2]
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Printf("Error opening database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	query := `
		SELECT player_name, final_net_worth, roi, successful_exits, turns_played, difficulty, played_at
		FROM game_scores
		ORDER BY final_net_worth DESC
	`

	rows, err := db.Query(query)
	if err != nil {
		fmt.Printf("Error querying database: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	var scores []ScoreData
	for rows.Next() {
		var score ScoreData
		var playedAt sql.NullString
		err := rows.Scan(
			&score.PlayerName,
			&score.FinalNetWorth,
			&score.ROI,
			&score.SuccessfulExits,
			&score.TurnsPlayed,
			&score.Difficulty,
			&playedAt,
		)
		if err != nil {
			fmt.Printf("Error scanning row: %v\n", err)
			continue
		}
		if playedAt.Valid {
			score.PlayedAt = playedAt.String
		}
		scores = append(scores, score)
	}

	jsonData, err := json.MarshalIndent(scores, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}

	err = ioutil.WriteFile(outputPath, jsonData, 0644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Exported %d scores to %s\n", len(scores), outputPath)
}
