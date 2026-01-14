package history

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
)

func LearnScopesFromHistory() {
	fmt.Println("\n--- Learning from Commit History ---")
	logCmd := exec.Command("git", "log", "--pretty=format:%s")
	logOutput, err := logCmd.Output()
	if err != nil {
		log.Printf("Could not get git log: %v", err)
		return
	}

	// Regex to find text in parentheses, like (scope)
	re := regexp.MustCompile(`\((.*?)\)`)
	matches := re.FindAllStringSubmatch(string(logOutput), -1)

	if len(matches) > 0 {
		fmt.Println("Found potential scopes:")
		for _, match := range matches {
			if len(match) > 1 {
				fmt.Printf("- %s\n", match[1])
			}
		}
	} else {
		fmt.Println("No conventional commit scopes found in history.")
	}
}

