package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// GameScore represents a completed game
type GameScore struct {
	ID              int
	PlayerName      string
	FinalNetWorth   int64
	ROI             float64
	SuccessfulExits int
	TurnsPlayed     int
	Difficulty      string
	PlayedAt        time.Time
}

// PlayerStats represents aggregate statistics for a player
type PlayerStats struct {
	PlayerName      string
	TotalGames      int
	BestNetWorth    int64
	BestROI         float64
	TotalExits      int
	AverageNetWorth float64
	WinRate         float64 // % of games with positive ROI
}

var db *sql.DB

// InitDB initializes the database connection and creates tables
func InitDB(dbPath string) error {
	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	// Create tables
	createTablesSQL := `
	CREATE TABLE IF NOT EXISTS game_scores (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		player_name TEXT NOT NULL,
		final_net_worth INTEGER NOT NULL,
		roi REAL NOT NULL,
		successful_exits INTEGER NOT NULL,
		turns_played INTEGER NOT NULL,
		difficulty TEXT NOT NULL,
		played_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS player_achievements (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		player_name TEXT NOT NULL,
		achievement_id TEXT NOT NULL,
		unlocked_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(player_name, achievement_id)
	);

	CREATE INDEX IF NOT EXISTS idx_net_worth ON game_scores(final_net_worth DESC);
	CREATE INDEX IF NOT EXISTS idx_roi ON game_scores(roi DESC);
	CREATE INDEX IF NOT EXISTS idx_player ON game_scores(player_name);
	CREATE INDEX IF NOT EXISTS idx_difficulty ON game_scores(difficulty);
	CREATE INDEX IF NOT EXISTS idx_player_achievements ON player_achievements(player_name);
	`

	_, err = db.Exec(createTablesSQL)
	if err != nil {
		return fmt.Errorf("failed to create tables: %v", err)
	}

	return nil
}

// CloseDB closes the database connection
func CloseDB() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

// SaveGameScore saves a completed game to the database
func SaveGameScore(score GameScore) error {
	query := `
		INSERT INTO game_scores (player_name, final_net_worth, roi, successful_exits, turns_played, difficulty, played_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := db.Exec(query,
		score.PlayerName,
		score.FinalNetWorth,
		score.ROI,
		score.SuccessfulExits,
		score.TurnsPlayed,
		score.Difficulty,
		score.PlayedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to save game score: %v", err)
	}

	return nil
}

// GetTopScoresByNetWorth returns the top N scores by net worth
func GetTopScoresByNetWorth(limit int, difficulty string) ([]GameScore, error) {
	var query string
	var args []interface{}

	if difficulty == "" || difficulty == "all" {
		query = `
			SELECT id, player_name, final_net_worth, roi, successful_exits, turns_played, difficulty, played_at
			FROM game_scores
			ORDER BY final_net_worth DESC
			LIMIT ?
		`
		args = []interface{}{limit}
	} else {
		query = `
			SELECT id, player_name, final_net_worth, roi, successful_exits, turns_played, difficulty, played_at
			FROM game_scores
			WHERE difficulty = ?
			ORDER BY final_net_worth DESC
			LIMIT ?
		`
		args = []interface{}{difficulty, limit}
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query top scores: %v", err)
	}
	defer rows.Close()

	var scores []GameScore
	for rows.Next() {
		var score GameScore
		err := rows.Scan(
			&score.ID,
			&score.PlayerName,
			&score.FinalNetWorth,
			&score.ROI,
			&score.SuccessfulExits,
			&score.TurnsPlayed,
			&score.Difficulty,
			&score.PlayedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		scores = append(scores, score)
	}

	return scores, nil
}

// GetTopScoresByROI returns the top N scores by ROI
func GetTopScoresByROI(limit int, difficulty string) ([]GameScore, error) {
	var query string
	var args []interface{}

	if difficulty == "" || difficulty == "all" {
		query = `
			SELECT id, player_name, final_net_worth, roi, successful_exits, turns_played, difficulty, played_at
			FROM game_scores
			ORDER BY roi DESC
			LIMIT ?
		`
		args = []interface{}{limit}
	} else {
		query = `
			SELECT id, player_name, final_net_worth, roi, successful_exits, turns_played, difficulty, played_at
			FROM game_scores
			WHERE difficulty = ?
			ORDER BY roi DESC
			LIMIT ?
		`
		args = []interface{}{difficulty, limit}
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query top scores: %v", err)
	}
	defer rows.Close()

	var scores []GameScore
	for rows.Next() {
		var score GameScore
		err := rows.Scan(
			&score.ID,
			&score.PlayerName,
			&score.FinalNetWorth,
			&score.ROI,
			&score.SuccessfulExits,
			&score.TurnsPlayed,
			&score.Difficulty,
			&score.PlayedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		scores = append(scores, score)
	}

	return scores, nil
}

// GetPlayerStats returns aggregate statistics for a player
func GetPlayerStats(playerName string) (*PlayerStats, error) {
	query := `
		SELECT 
			COUNT(*) as total_games,
			MAX(final_net_worth) as best_net_worth,
			MAX(roi) as best_roi,
			SUM(successful_exits) as total_exits,
			AVG(final_net_worth) as avg_net_worth,
			SUM(CASE WHEN roi > 0 THEN 1 ELSE 0 END) * 100.0 / COUNT(*) as win_rate
		FROM game_scores
		WHERE player_name = ?
	`

	var stats PlayerStats
	stats.PlayerName = playerName

	err := db.QueryRow(query, playerName).Scan(
		&stats.TotalGames,
		&stats.BestNetWorth,
		&stats.BestROI,
		&stats.TotalExits,
		&stats.AverageNetWorth,
		&stats.WinRate,
	)

	if err == sql.ErrNoRows {
		// No games played yet
		return &stats, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get player stats: %v", err)
	}

	return &stats, nil
}

// GetRecentGames returns the most recent N games
func GetRecentGames(limit int) ([]GameScore, error) {
	query := `
		SELECT id, player_name, final_net_worth, roi, successful_exits, turns_played, difficulty, played_at
		FROM game_scores
		ORDER BY played_at DESC
		LIMIT ?
	`

	rows, err := db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query recent games: %v", err)
	}
	defer rows.Close()

	var scores []GameScore
	for rows.Next() {
		var score GameScore
		err := rows.Scan(
			&score.ID,
			&score.PlayerName,
			&score.FinalNetWorth,
			&score.ROI,
			&score.SuccessfulExits,
			&score.TurnsPlayed,
			&score.Difficulty,
			&score.PlayedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		scores = append(scores, score)
	}

	return scores, nil
}

// GetTotalGamesPlayed returns the total number of games in the database
func GetTotalGamesPlayed() (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM game_scores").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get total games: %v", err)
	}
	return count, nil
}

// Achievement tracking

// UnlockAchievement saves an unlocked achievement for a player
func UnlockAchievement(playerName, achievementID string) error {
	query := `
		INSERT OR IGNORE INTO player_achievements (player_name, achievement_id, unlocked_at)
		VALUES (?, ?, ?)
	`
	
	_, err := db.Exec(query, playerName, achievementID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to unlock achievement: %v", err)
	}
	
	return nil
}

// GetPlayerAchievements returns all unlocked achievements for a player
func GetPlayerAchievements(playerName string) ([]string, error) {
	query := `
		SELECT achievement_id
		FROM player_achievements
		WHERE player_name = ?
		ORDER BY unlocked_at ASC
	`
	
	rows, err := db.Query(query, playerName)
	if err != nil {
		return nil, fmt.Errorf("failed to query achievements: %v", err)
	}
	defer rows.Close()
	
	var achievements []string
	for rows.Next() {
		var achievementID string
		if err := rows.Scan(&achievementID); err != nil {
			return nil, fmt.Errorf("failed to scan achievement: %v", err)
		}
		achievements = append(achievements, achievementID)
	}
	
	return achievements, nil
}

// GetPlayerAchievementCount returns total achievements unlocked by player
func GetPlayerAchievementCount(playerName string) (int, error) {
	var count int
	query := `
		SELECT COUNT(*)
		FROM player_achievements
		WHERE player_name = ?
	`
	
	err := db.QueryRow(query, playerName).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get achievement count: %v", err)
	}
	
	return count, nil
}

// GetPlayerAchievementPoints returns total points from unlocked achievements
func GetPlayerAchievementPoints(playerName string) (int, error) {
	achievements, err := GetPlayerAchievements(playerName)
	if err != nil {
		return 0, err
	}
	
	// Note: Points would need to be calculated by looking up achievement definitions
	// For now, return count * 10 as placeholder
	return len(achievements) * 10, nil
}

// GetWinStreak calculates current win streak for a player
func GetWinStreak(playerName string) (int, error) {
	query := `
		SELECT roi
		FROM game_scores
		WHERE player_name = ?
		ORDER BY played_at DESC
		LIMIT 10
	`
	
	rows, err := db.Query(query, playerName)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	
	streak := 0
	for rows.Next() {
		var roi float64
		if err := rows.Scan(&roi); err != nil {
			return 0, err
		}
		
		if roi > 0 {
			streak++
		} else {
			break
		}
	}
	
	return streak, nil
}
