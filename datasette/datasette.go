package datasette

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	db "github.com/jamesacampbell/unicorn/database"
	"gopkg.in/yaml.v2"
)

// ErrDisabled is returned when Datasette integration is disabled in configuration.
var ErrDisabled = errors.New("datasette integration disabled")

// Config represents Datasette integration configuration
type Config struct {
	Enabled        bool   `yaml:"enabled"`
	BaseURL        string `yaml:"base_url"`
	Database       string `yaml:"database"`
	Table          string `yaml:"table"`
	APIToken       string `yaml:"api_token"`
	TimeoutSeconds int    `yaml:"timeout_seconds"`
}

// LoadConfig loads Datasette configuration from the provided path.
// If the file does not exist, a disabled configuration is returned.
func LoadConfig(path string) (*Config, error) {
	cfg := &Config{}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			cfg.Enabled = false
			return cfg, nil
		}
		return nil, fmt.Errorf("read datasette config: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parse datasette config: %w", err)
	}

	if token := os.Getenv("UNICORN_DATASETTE_TOKEN"); token != "" {
		cfg.APIToken = token
	}

	return cfg, nil
}

// Validate ensures the configuration is usable when enabled.
func (c *Config) Validate() error {
	if !c.Enabled {
		return nil
	}

	missing := make([]string, 0, 3)
	if strings.TrimSpace(c.BaseURL) == "" {
		missing = append(missing, "base_url")
	}
	if strings.TrimSpace(c.Database) == "" {
		missing = append(missing, "database")
	}
	if strings.TrimSpace(c.Table) == "" {
		missing = append(missing, "table")
	}

	if len(missing) > 0 {
		return fmt.Errorf("invalid datasette config: missing %s", strings.Join(missing, ", "))
	}

	if c.TimeoutSeconds <= 0 {
		c.TimeoutSeconds = 10
	}

	return nil
}

// Client provides helper methods for interacting with Datasette APIs.
type Client struct {
	cfg        Config
	httpClient *http.Client
}

// NewClient creates a new Datasette client from the provided configuration.
func NewClient(cfg Config) *Client {
	timeout := time.Duration(cfg.TimeoutSeconds) * time.Second
	if timeout <= 0 {
		timeout = 10 * time.Second
	}

	return &Client{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// SubmitScore sends a score to the configured Datasette instance using the datasette-write API.
func (c *Client) SubmitScore(ctx context.Context, score db.GameScore) error {
	if !c.cfg.Enabled {
		return ErrDisabled
	}

	endpoint := strings.TrimRight(c.cfg.BaseURL, "/") + "/-/insert"

	payload := map[string]any{
		"database": c.cfg.Database,
		"table":    c.cfg.Table,
		"rows": []map[string]any{
			{
				"player_name":      score.PlayerName,
				"final_net_worth":  score.FinalNetWorth,
				"roi":              score.ROI,
				"successful_exits": score.SuccessfulExits,
				"turns_played":     score.TurnsPlayed,
				"difficulty":       score.Difficulty,
				"played_at":        score.PlayedAt.UTC().Format(time.RFC3339),
				"submitted_at":     time.Now().UTC().Format(time.RFC3339),
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal datasette payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create datasette request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.cfg.APIToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.cfg.APIToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send datasette request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		io.Copy(io.Discard, resp.Body)
		return nil
	}

	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 4<<10))
	return fmt.Errorf("datasette insert failed: status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(respBody)))
}
