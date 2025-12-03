package assets

import (
	"embed"
	"fmt"
	"io/fs"
)

//go:embed startups/*.json rounds/*.json founder/*.json
var EmbeddedFiles embed.FS

// ReadStartupFile reads a startup JSON file by ID
func ReadStartupFile(id int) ([]byte, error) {
	return EmbeddedFiles.ReadFile(fmt.Sprintf("startups/%d.json", id))
}

// ReadRoundOptions reads the round-options.json file
func ReadRoundOptions() ([]byte, error) {
	return EmbeddedFiles.ReadFile("rounds/round-options.json")
}

// ReadFounderStartups reads the founder/startups.json file
func ReadFounderStartups() ([]byte, error) {
	return EmbeddedFiles.ReadFile("founder/startups.json")
}

// GetStartupsFS returns a filesystem for the startups directory
func GetStartupsFS() (fs.FS, error) {
	return fs.Sub(EmbeddedFiles, "startups")
}
