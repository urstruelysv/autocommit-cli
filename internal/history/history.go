package history

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"regexp"
	"strings"
	"encoding/json"
)

// LearnData holds the learned scopes and types
type LearnData struct {
	Scopes map[string]int
	Types  map[string]int
}

func SaveLearnedData(data LearnData) error {
	cacheFilePath := ".autocommit_cache"
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal learned data: %w", err)
	}

	err = ioutil.WriteFile(cacheFilePath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write learned data to %s: %w", cacheFilePath, err)
	}
	fmt.Printf("Learned data saved to %s\n", cacheFilePath)
	return nil
}

func LoadLearnedData() (LearnData, error) {
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
	fmt.Printf("Learned data loaded from %s\n", cacheFilePath)
	return data, nil
}

func LearnFromHistory() LearnData {
	fmt.Println("\n--- Learning from Commit History ---")
	logCmd := exec.Command("git", "log", "--pretty=format:%s")
	logOutput, err := logCmd.Output()
	if err != nil {
		log.Printf("Could not get git log: %v", err)
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
		fmt.Println("Found potential scopes:")
		for scope, count := range scopes {
			fmt.Printf("- %s (%d)\n", scope, count)
		}
	} else {
		fmt.Println("No conventional commit scopes found in history.")
	}

	if len(types) > 0 {
		fmt.Println("Found potential types:")
		for t, count := range types {
			fmt.Printf("- %s (%d)\n", t, count)
		}
	} else {
		fmt.Println("No conventional commit types found in history.")
	}

	return LearnData{Scopes: scopes, Types: types}
}

