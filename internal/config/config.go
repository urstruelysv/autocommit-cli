package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath" // Added for filepath.Glob
	"regexp" // Added for regexp.MustCompile
	"strings" // Already present, but ensuring it's there

	"github.com/BurntSushi/toml"
)

// Config holds the application's configuration settings.
type Config struct {
	AutoPush       bool `toml:"auto_push"`
	ReviewMode     bool `toml:"review_mode"`
	LearnFromHistory bool `toml:"learn_from_history"`
	AICommit       bool `toml:"ai_commit"`
	CI             bool `toml:"ci"`
	Verbose        bool `toml:"verbose"`
}

// LoadConfig reads the configuration from a .autocommitrc file.
// It looks for the file in the current directory.
func LoadConfig() (Config, error) {
	var cfg Config
	configPath := ".autocommitrc"

	// Set default values
	cfg.AutoPush = true
	cfg.ReviewMode = false
	cfg.LearnFromHistory = true
	cfg.AICommit = false
	cfg.CI = false
	cfg.Verbose = false

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// If config file doesn't exist, return default config
		return cfg, nil
	}

	// Read and parse the config file
	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		return cfg, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	if _, err := toml.Decode(string(configData), &cfg); err != nil {
		return cfg, fmt.Errorf("failed to decode config file %s: %w", configPath, err)
	}

	return cfg, nil
}

// CommitRules holds the extracted static rules for commit message validation.
type CommitRules struct {
	CommitMessageRegex string
}

// ParseCommitGuides parses the detected commit guide files to extract static rules.
func ParseCommitGuides(guides []string) (CommitRules, error) {
	var rules CommitRules

	for _, guidePath := range guides {
		content, err := ioutil.ReadFile(guidePath)
		if err != nil {
			return rules, fmt.Errorf("failed to read commit guide file %s: %w", guidePath, err)
		}

		// Regex-based extraction for CONTRIBUTING.md
		if guidePath == "CONTRIBUTING.md" {
			// Very basic attempt to find a regex pattern for commit messages
			// This is highly brittle and assumes a specific format.
			re := regexp.MustCompile(`(?i)commit message(?: format)?(?: regex)?:?\s*(\/.*\/)`)
			matches := re.FindStringSubmatch(string(content))
			if len(matches) > 1 {
				rules.CommitMessageRegex = strings.TrimSpace(matches[1])
				fmt.Printf("Extracted commit message regex from CONTRIBUTING.md: %s\n", rules.CommitMessageRegex)
				// Prioritize CONTRIBUTING.md if found
				return rules, nil
			}
		}

		// Regex-based extraction for commitlint.config.* (assuming JSON for simplicity)
		if strings.HasPrefix(guidePath, "commitlint.config.") && strings.HasSuffix(guidePath, ".json") {
			// Attempt to extract a regex from a JSON structure
			// This is also highly brittle and assumes a specific JSON path.
			re := regexp.MustCompile(`"pattern":\s*"(.*?)"`)
			matches := re.FindStringSubmatch(string(content))
			if len(matches) > 1 {
				rules.CommitMessageRegex = strings.TrimSpace(matches[1])
				fmt.Printf("Extracted commit message regex from %s: %s\n", guidePath, rules.CommitMessageRegex)
				return rules, nil
			}
		}
	}

	return rules, nil
}

// DetectCommitGuides checks for the presence of commit guide files.
func DetectCommitGuides() ([]string, error) {
	var guides []string

	// Check for CONTRIBUTING.md
	if _, err := os.Stat("CONTRIBUTING.md"); err == nil {
		guides = append(guides, "CONTRIBUTING.md")
	}

	// Check for commitlint.config.* files
	matches, err := filepath.Glob("commitlint.config.*")
	if err != nil {
		return nil, fmt.Errorf("error globbing for commitlint config files: %w", err)
	}
	guides = append(guides, matches...)

	return guides, nil
}
