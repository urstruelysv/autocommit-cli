package history

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strings"

	"github.com/urstruelysv/autocommit-cli/internal/logger"
)

// LearnData holds the learned scopes and types
type LearnData struct {
	Scopes map[string]int
	Types  map[string]int
}

func SaveLearnedData(log logger.Logger, data LearnData) error {
	log.Debug("Saving learned data...")
	cacheFilePath := ".autocommit_cache"
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal learned data: %w", err)
	}

	err = ioutil.WriteFile(cacheFilePath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write learned data to %s: %w", cacheFilePath, err)
	}
	log.Debug("Learned data saved to %s", cacheFilePath)
	log.Info("Learned data saved to %s", cacheFilePath)
	return nil
}

func LoadLearnedData(log logger.Logger) (LearnData, error) {
	log.Debug("Loading learned data...")
	cacheFilePath := ".autocommit_cache"
	jsonData, err := ioutil.ReadFile(cacheFilePath)
	if err != nil {
		return LearnData{}, fmt.Errorf("failed to read learned data from %s: %w", cacheFilePath, err)
	}

	var data LearnData
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return LearnData{}, fmt.Errorf("failed to unmarshal learned data from %s: %w", cacheFilePath, err)
	}
	log.Debug("Learned data loaded from %s", cacheFilePath)
	log.Info("Learned data loaded from %s", cacheFilePath)
	return data, nil
}

func LearnFromHistory(log logger.Logger) LearnData {
	log.Debug("Learning from commit history...")
	log.Info("\n--- Learning from Commit History ---")
	logCmd := exec.Command("git", "log", "--pretty=format:%s")
	logOutput, err := logCmd.Output()
	if err != nil {
		log.Error("Could not get git log: %v", err)
		return LearnData{}
	}

	scopes := make(map[string]int)
	types := make(map[string]int)

	// Regex to find text in parentheses, like (scope)
	scopeRe := regexp.MustCompile(`\((.*?)\)`)
	// Regex to find commit type, e.g., "feat", "fix"
	typeRe := regexp.MustCompile(`^([a-z]+)(?:\(.*\))?:`)

	commitSubjects := strings.Split(strings.TrimSpace(string(logOutput)), "\n")
	for _, subject := range commitSubjects {
		// Extract scope
		scopeMatches := scopeRe.FindStringSubmatch(subject)
		if len(scopeMatches) > 1 {
			scopes[scopeMatches[1]]++
		}

		// Extract type
		typeMatches := typeRe.FindStringSubmatch(subject)
		if len(typeMatches) > 1 {
			types[typeMatches[1]]++
		}
	}

	if len(scopes) > 0 {
		log.Debug("Found potential scopes.")
		log.Info("Found potential scopes:")
		for scope, count := range scopes {
			log.Info("- %s (%d)", scope, count)
		}
	} else {
		log.Debug("No conventional commit scopes found in history.")
		log.Info("No conventional commit scopes found in history.")
	}

	if len(types) > 0 {
		log.Debug("Found potential types.")
		log.Info("Found potential types:")
		for t, count := range types {
			log.Info("- %s (%d)", t, count)
		}
	} else {
		log.Debug("No conventional commit types found in history.")
		log.Info("No conventional commit types found in history.")
	}

	return LearnData{Scopes: scopes, Types: types}
}
