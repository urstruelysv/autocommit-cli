package classify

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func ClassifyAndGroupChanges(changes string) map[string][]string {
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
		groups[commitType] = append(groups[commitType], filePath)
	}

	for commitType, files := range groups {
		fmt.Printf("Group '%s': %v\n", commitType, files)
	}

	return groups
}
