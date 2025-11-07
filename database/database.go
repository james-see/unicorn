package database

import (
	"database/sql"
	"fmt"
	"math"
	"strings"
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

	CREATE TABLE IF NOT EXISTS player_upgrades (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		player_name TEXT NOT NULL,
		upgrade_id TEXT NOT NULL,
		purchased_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(player_name, upgrade_id)
	);

	CREATE TABLE IF NOT EXISTS player_profiles (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		player_name TEXT UNIQUE NOT NULL,
		level INTEGER DEFAULT 1,
		experience_points INTEGER DEFAULT 0,
		total_points_earned INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		last_played DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS player_level_history (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		player_name TEXT NOT NULL,
		level INTEGER NOT NULL,
		reached_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS achievement_progress (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		player_name TEXT NOT NULL,
		achievement_id TEXT NOT NULL,
		current_progress INTEGER DEFAULT 0,
		max_progress INTEGER NOT NULL,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(player_name, achievement_id)
	);

	CREATE TABLE IF NOT EXISTS game_history_detailed (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		player_name TEXT NOT NULL,
		game_mode TEXT NOT NULL,
		difficulty TEXT NOT NULL,
		final_net_worth INTEGER,
		roi REAL,
		successful_exits INTEGER,
		turns_played INTEGER,
		investments_made INTEGER,
		best_investment_roi REAL,
		worst_investment_roi REAL,
		avg_investment_amount INTEGER,
		total_invested INTEGER,
		played_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_net_worth ON game_scores(final_net_worth DESC);
	CREATE INDEX IF NOT EXISTS idx_roi ON game_scores(roi DESC);
	CREATE INDEX IF NOT EXISTS idx_player ON game_scores(player_name);
	CREATE INDEX IF NOT EXISTS idx_difficulty ON game_scores(difficulty);
	CREATE INDEX IF NOT EXISTS idx_player_achievements ON player_achievements(player_name);
	CREATE INDEX IF NOT EXISTS idx_player_upgrades ON player_upgrades(player_name);
	CREATE INDEX IF NOT EXISTS idx_player_profiles ON player_profiles(player_name);
	CREATE INDEX IF NOT EXISTS idx_level_history ON player_level_history(player_name);
	CREATE INDEX IF NOT EXISTS idx_achievement_progress ON achievement_progress(player_name);
	CREATE INDEX IF NOT EXISTS idx_game_history ON game_history_detailed(player_name, played_at DESC);
	`

	_, err = db.Exec(createTablesSQL)
	if err != nil {
		return fmt.Errorf("failed to create tables: %v", err)
	}

	// Add level_up_points column if it doesn't exist (migration)
	_, err = db.Exec(`
		ALTER TABLE player_profiles 
		ADD COLUMN level_up_points INTEGER DEFAULT 0
	`)
	// Ignore error if column already exists
	if err != nil && !strings.Contains(err.Error(), "duplicate column") {
		// Only return error if it's not a "column already exists" error
		return fmt.Errorf("failed to add level_up_points column: %v", err)
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

// GetPlayerStatsByMode returns aggregate statistics for a player filtered by game mode
// gameMode: "vc" for VC mode (Easy/Medium/Hard/Expert) or "founder" for Founder mode
func GetPlayerStatsByMode(playerName string, gameMode string) (*PlayerStats, error) {
	var difficultyFilter string
	if gameMode == "vc" {
		// VC modes: Easy, Medium, Hard, Expert
		difficultyFilter = "difficulty IN ('Easy', 'Medium', 'Hard', 'Expert')"
	} else if gameMode == "founder" {
		// Founder mode
		difficultyFilter = "difficulty = 'Founder'"
	} else {
		return nil, fmt.Errorf("invalid game mode: %s", gameMode)
	}
	
	query := fmt.Sprintf(`
		SELECT 
			COUNT(*) as total_games,
			MAX(final_net_worth) as best_net_worth,
			MAX(roi) as best_roi,
			SUM(successful_exits) as total_exits,
			AVG(final_net_worth) as avg_net_worth,
			SUM(CASE WHEN roi > 0 THEN 1 ELSE 0 END) * 100.0 / COUNT(*) as win_rate
		FROM game_scores
		WHERE player_name = ? AND %s
	`, difficultyFilter)

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
		// No games played yet for this mode
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
	
	// Calculate actual points from achievement definitions
	// Note: This requires importing achievements package, but we'll do it differently
	// For now, return count * 10 as placeholder - will be fixed when we integrate
	return len(achievements) * 10, nil
}

// AchievementLeaderboardEntry represents a player's achievement stats for leaderboard
type AchievementLeaderboardEntry struct {
	PlayerName        string
	AchievementCount  int
	TotalPoints       int
	CareerLevel       int
	CareerTitle       string
	RareCount         int
	EpicCount         int
	LegendaryCount    int
}

// GetAchievementLeaderboard returns top players by achievement count and points
func GetAchievementLeaderboard(limit int) ([]AchievementLeaderboardEntry, error) {
	// Get all unique players with achievements
	query := `
		SELECT DISTINCT player_name
		FROM player_achievements
		ORDER BY player_name
	`
	
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query players: %v", err)
	}
	defer rows.Close()
	
	var players []string
	for rows.Next() {
		var playerName string
		if err := rows.Scan(&playerName); err != nil {
			continue
		}
		players = append(players, playerName)
	}
	
		// Build leaderboard entries
		leaderboard := []AchievementLeaderboardEntry{}
		for _, playerName := range players {
			achievements, err := GetPlayerAchievements(playerName)
			if err != nil {
				continue
			}
			
			entry := AchievementLeaderboardEntry{
				PlayerName:       playerName,
				AchievementCount: len(achievements),
				TotalPoints:      0, // Will be calculated in UI layer
			}
			
			// Get career level from profile
			profile, err := GetPlayerProfile(playerName)
			if err == nil {
				entry.CareerLevel = profile.Level
			}
			
			leaderboard = append(leaderboard, entry)
		}
	
	// Sort by achievement count (descending), then by points
	// Simple bubble sort for small datasets
	for i := 0; i < len(leaderboard); i++ {
		for j := i + 1; j < len(leaderboard); j++ {
			if leaderboard[i].AchievementCount < leaderboard[j].AchievementCount ||
				(leaderboard[i].AchievementCount == leaderboard[j].AchievementCount &&
					leaderboard[i].TotalPoints < leaderboard[j].TotalPoints) {
				leaderboard[i], leaderboard[j] = leaderboard[j], leaderboard[i]
			}
		}
	}
	
	// Limit results
	if limit > 0 && limit < len(leaderboard) {
		leaderboard = leaderboard[:limit]
	}
	
	return leaderboard, nil
}

// Upgrade functions

// PurchaseUpgrade saves a purchased upgrade for a player
func PurchaseUpgrade(playerName, upgradeID string) error {
	query := `
		INSERT OR IGNORE INTO player_upgrades (player_name, upgrade_id, purchased_at)
		VALUES (?, ?, ?)
	`
	
	_, err := db.Exec(query, playerName, upgradeID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to purchase upgrade: %v", err)
	}
	
	return nil
}

// GetPlayerUpgrades returns all purchased upgrades for a player
func GetPlayerUpgrades(playerName string) ([]string, error) {
	query := `
		SELECT upgrade_id
		FROM player_upgrades
		WHERE player_name = ?
		ORDER BY purchased_at ASC
	`
	
	rows, err := db.Query(query, playerName)
	if err != nil {
		return nil, fmt.Errorf("failed to query upgrades: %v", err)
	}
	defer rows.Close()
	
	var upgrades []string
	for rows.Next() {
		var upgradeID string
		if err := rows.Scan(&upgradeID); err != nil {
			return nil, fmt.Errorf("failed to scan upgrade: %v", err)
		}
		upgrades = append(upgrades, upgradeID)
	}
	
	return upgrades, nil
}

// HasUpgrade checks if a player has a specific upgrade
func HasUpgrade(playerName, upgradeID string) (bool, error) {
	query := `
		SELECT COUNT(*)
		FROM player_upgrades
		WHERE player_name = ? AND upgrade_id = ?
	`
	
	var count int
	err := db.QueryRow(query, playerName, upgradeID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check upgrade: %v", err)
	}
	
	return count > 0, nil
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

// PlayerProfile represents a player's progression data
type PlayerProfile struct {
	PlayerName         string
	Level              int
	ExperiencePoints   int
	TotalPointsEarned  int
	LevelUpPoints      int // Points earned from leveling up
	NextLevelXP        int
	ProgressPercent    float64
	CreatedAt          time.Time
	LastPlayed         time.Time
}

// GetPlayerProfile retrieves or creates a player's profile
func GetPlayerProfile(playerName string) (*PlayerProfile, error) {
	// Try to get existing profile
	query := `
		SELECT player_name, level, experience_points, total_points_earned, COALESCE(level_up_points, 0), created_at, last_played
		FROM player_profiles
		WHERE player_name = ?
	`
	
	var profile PlayerProfile
	err := db.QueryRow(query, playerName).Scan(
		&profile.PlayerName,
		&profile.Level,
		&profile.ExperiencePoints,
		&profile.TotalPointsEarned,
		&profile.LevelUpPoints,
		&profile.CreatedAt,
		&profile.LastPlayed,
	)
	
	if err == sql.ErrNoRows {
		// Create new profile
		insertQuery := `
			INSERT INTO player_profiles (player_name, level, experience_points, total_points_earned, level_up_points)
			VALUES (?, 1, 0, 0, 0)
		`
		_, err := db.Exec(insertQuery, playerName)
		if err != nil {
			return nil, fmt.Errorf("failed to create player profile: %v", err)
		}
		
		// Return new profile
		profile = PlayerProfile{
			PlayerName:         playerName,
			Level:              1,
			ExperiencePoints:   0,
			TotalPointsEarned:  0,
			LevelUpPoints:      0,
			CreatedAt:          time.Now(),
			LastPlayed:         time.Now(),
		}
	} else if err != nil {
		return nil, fmt.Errorf("failed to get player profile: %v", err)
	}
	
	// Calculate next level requirements
	profile.NextLevelXP = GetLevelRequirement(profile.Level + 1)
	profile.ProgressPercent = float64(profile.ExperiencePoints) / float64(profile.NextLevelXP) * 100
	
	return &profile, nil
}

// GetLevelRequirement returns the XP needed to reach a specific level
func GetLevelRequirement(level int) int {
	// Level N requires: 200 * (N ^ 1.5) total XP
	// This creates an exponential curve:
	// Level 2: 200 * (2^1.5) ≈ 283 XP
	// Level 3: 200 * (3^1.5) ≈ 519 XP  
	// Level 5: 200 * (5^1.5) ≈ 1118 XP
	// Level 10: 200 * (10^1.5) ≈ 6325 XP
	if level <= 1 {
		return 0
	}
	
	return int(200 * math.Pow(float64(level), 1.5))
}

// AddExperience adds XP to a player and handles level ups
// Returns: leveledUp, newLevel, pointsEarned (from level ups), error
func AddExperience(playerName string, xpAmount int) (leveledUp bool, newLevel int, pointsEarned int, err error) {
	// Get current profile
	profile, err := GetPlayerProfile(playerName)
	if err != nil {
		return false, 0, 0, err
	}
	
	// Add XP
	newXP := profile.ExperiencePoints + xpAmount
	newTotal := profile.TotalPointsEarned + xpAmount
	currentLevel := profile.Level
	totalPointsEarned := 0
	
	// Check for level ups
	for {
		requiredXP := GetLevelRequirement(currentLevel + 1)
		if newXP >= requiredXP {
			currentLevel++
			newXP -= requiredXP
			leveledUp = true
			
			// Award points for leveling up: 10 * newLevel
			levelUpPoints := 10 * currentLevel
			totalPointsEarned += levelUpPoints
			
			// Record level up in history
			historyQuery := `
				INSERT INTO player_level_history (player_name, level, reached_at)
				VALUES (?, ?, CURRENT_TIMESTAMP)
			`
			_, err := db.Exec(historyQuery, playerName, currentLevel)
			if err != nil {
				return false, 0, 0, fmt.Errorf("failed to record level history: %v", err)
			}
		} else {
			break
		}
	}
	
	// Calculate new level_up_points total
	newLevelUpPoints := profile.LevelUpPoints + totalPointsEarned
	
	// Update profile
	updateQuery := `
		UPDATE player_profiles
		SET level = ?, experience_points = ?, total_points_earned = ?, level_up_points = ?, last_played = CURRENT_TIMESTAMP
		WHERE player_name = ?
	`
	_, err = db.Exec(updateQuery, currentLevel, newXP, newTotal, newLevelUpPoints, playerName)
	if err != nil {
		return false, 0, 0, fmt.Errorf("failed to update player profile: %v", err)
	}
	
	return leveledUp, currentLevel, totalPointsEarned, nil
}

// GetTotalPlayerPoints returns the total points a player has (achievement points + level-up points)
func GetTotalPlayerPoints(playerName string) (int, error) {
	// Get level-up points from profile
	profile, err := GetPlayerProfile(playerName)
	if err != nil {
		return 0, err
	}
	
	// Return level-up points (achievement points are calculated in UI)
	return profile.LevelUpPoints, nil
}

// UpdateLastPlayed updates the last played timestamp for a player
func UpdateLastPlayed(playerName string) error {
	query := `
		UPDATE player_profiles
		SET last_played = CURRENT_TIMESTAMP
		WHERE player_name = ?
	`
	_, err := db.Exec(query, playerName)
	if err != nil {
		return fmt.Errorf("failed to update last played: %v", err)
	}
	return nil
}

// UpdateAchievementProgress updates progress toward a progressive achievement
func UpdateAchievementProgress(playerName, achievementID string, currentProgress, maxProgress int) error {
	query := `
		INSERT INTO achievement_progress (player_name, achievement_id, current_progress, max_progress, updated_at)
		VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(player_name, achievement_id) 
		DO UPDATE SET current_progress = ?, updated_at = CURRENT_TIMESTAMP
	`
	_, err := db.Exec(query, playerName, achievementID, currentProgress, maxProgress, currentProgress)
	if err != nil {
		return fmt.Errorf("failed to update achievement progress: %v", err)
	}
	return nil
}

// GetAchievementProgress gets the current progress for an achievement
func GetAchievementProgress(playerName, achievementID string) (current, max int, err error) {
	query := `
		SELECT current_progress, max_progress
		FROM achievement_progress
		WHERE player_name = ? AND achievement_id = ?
	`
	err = db.QueryRow(query, playerName, achievementID).Scan(&current, &max)
	if err == sql.ErrNoRows {
		return 0, 0, nil // No progress yet
	}
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get achievement progress: %v", err)
	}
	return current, max, nil
}

// ProgressInfo represents progress toward an achievement
type ProgressInfo struct {
	AchievementID   string
	CurrentProgress int
	MaxProgress     int
	UpdatedAt       time.Time
}

// GetAllProgress gets all achievement progress for a player
func GetAllProgress(playerName string) (map[string]ProgressInfo, error) {
	query := `
		SELECT achievement_id, current_progress, max_progress, updated_at
		FROM achievement_progress
		WHERE player_name = ?
	`
	
	rows, err := db.Query(query, playerName)
	if err != nil {
		return nil, fmt.Errorf("failed to query achievement progress: %v", err)
	}
	defer rows.Close()
	
	progress := make(map[string]ProgressInfo)
	for rows.Next() {
		var info ProgressInfo
		if err := rows.Scan(&info.AchievementID, &info.CurrentProgress, &info.MaxProgress, &info.UpdatedAt); err != nil {
			return nil, err
		}
		progress[info.AchievementID] = info
	}
	
	return progress, nil
}

// GetTopScoresByPlayer returns top scores for a specific player
func GetTopScoresByPlayer(playerName string, limit int) ([]GameScore, error) {
	query := `
		SELECT id, player_name, final_net_worth, roi, successful_exits, turns_played, difficulty, played_at
		FROM game_scores
		WHERE player_name = ?
		ORDER BY played_at DESC
		LIMIT ?
	`
	
	rows, err := db.Query(query, playerName, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query scores: %v", err)
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
			return nil, err
		}
		scores = append(scores, score)
	}
	
	return scores, nil
}
