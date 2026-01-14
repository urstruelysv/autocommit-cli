package classify

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"autocommit/internal/history" // Added this import
)

func ClassifyAndGroupChanges(changes string, learnedData history.LearnData) map[string][]string {
	fmt.Println("\n--- Classifying and Grouping Changes ---")
	groups := make(map[string][]string)

	lines := strings.Split(changes, "\n")
	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		filePath := parts[len(parts)-1]
		commitType := "chore" // Default type
		scope := ""

		// Determine scope from file path
		pathParts := strings.Split(filePath, "/")
		if len(pathParts) >= 2 {
			if pathParts[0] == "cmd" && len(pathParts) >= 2 {
				scope = pathParts[1] // e.g., cmd/autocommit -> autocommit
			} else if pathParts[0] == "internal" && len(pathParts) >= 2 {
				scope = pathParts[1] // e.g., internal/git -> git
			}
		}

		if strings.Contains(filePath, "tests/") || strings.HasPrefix(filePath, "test_") {
			commitType = "test"
		} else if strings.HasSuffix(filePath, ".md") {
			commitType = "docs"
		} else {
			diffCmd := exec.Command("git", "diff", "--", filePath)
			diffOutput, err := diffCmd.Output()
			if err != nil {
				log.Printf("Could not get diff for %s: %v", filePath, err)
			} else {
				diff := strings.ToLower(string(diffOutput))
				if strings.Contains(diff, "fix") || strings.Contains(diff, "bug") {
					commitType = "fix"
				} else if strings.Contains(diff, "add") || strings.Contains(diff, "feature") {
					commitType = "feat"
				}
			}
		}
		
		groupKey := commitType
		if scope != "" {
			groupKey = fmt.Sprintf("%s(%s)", commitType, scope)
		}
		groups[groupKey] = append(groups[groupKey], filePath)
	}

	for groupKey, files := range groups {
		fmt.Printf("Group '%s': %v\n", groupKey, files)
	}

	return groups
}
