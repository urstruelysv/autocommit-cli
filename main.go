package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func detectChanges() (string, error) {
	fmt.Println("Detecting changes...")
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Error detecting changes: %v\n", err)
		return "", err
	}

	changes := strings.TrimSpace(string(output))
	if changes != "" {
		fmt.Println("Found changes:")
		fmt.Println(changes)
	} else {
		fmt.Println("No changes found.")
	}
	return changes, nil
}

func classifyAndGroupChanges(changes string) map[string][]string {
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

func commitChanges(message string, files []string) error {
	fmt.Printf("\n--- Committing Group: %s ---\n", message)

	addArgs := append([]string{"add"}, files...)
	addCmd := exec.Command("git", addArgs...)
	if output, err := addCmd.CombinedOutput(); err != nil {
		log.Printf("Error staging files %v: %s\n%v", files, string(output), err)
		return err
	}
	fmt.Printf("Staged files: %v\n", files)

	commitCmd := exec.Command("git", "commit", "-m", message)
	if output, err := commitCmd.CombinedOutput(); err != nil {
		log.Printf("Error committing group: %s\n%v", string(output), err)
		return err
	}
	fmt.Println("Committed group.")
	return nil
}

func pushChanges() {
	fmt.Println("\n--- Pushing Changes ---")
	// Check for detached HEAD
	branchCmd := exec.Command("git", "symbolic-ref", "--short", "HEAD")
	branchOutput, err := branchCmd.Output()
	if err != nil {
		log.Println("Error: Detached HEAD state or other error getting branch. Aborting push.")
		return
	}
	branchName := strings.TrimSpace(string(branchOutput))

	// Check if remote is configured
	remoteCmd := exec.Command("git", "config", fmt.Sprintf("branch.%s.remote", branchName))
	if err := remoteCmd.Run(); err != nil {
		log.Printf("Error: No remote configured for branch '%s'. Aborting push.\n", branchName)
		return
	}

	fmt.Printf("Pushing changes to remote for branch '%s'...\n", branchName)
	pushCmd := exec.Command("git", "push")
	if output, err := pushCmd.CombinedOutput(); err != nil {
		log.Printf("Error during push: %s\n%v", string(output), err)
		return
	}
	fmt.Println("Push successful.")
}

// Application entry point
func main() {
	// Define flags
	review := flag.Bool("review", false, "Enable review mode to inspect commits before they are made.")
	yRun := flag.Bool("y-run", false, "Not implemented yet.")
	noPush := flag.Bool("no-push", false, "Create commits but do not push them to the remote repository.")
	ci := flag.Bool("ci", false, "Enable CI mode for non-interactive, deterministic execution.")
	verbose := flag.Bool("verbose", false, "Enable verbose output for debugging purposes.")

	// Parse the flags
	flag.Parse()

	fmt.Println("AutoCommit AI (Go version) running...")
	fmt.Println("Arguments provided:")

	if *review {
		fmt.Println("- Review mode enabled.")
	}
	if *yRun {
		fmt.Println("- y-run flag set.")
	}
	if *noPush {
		fmt.Println("- No-push mode enabled.")
	}
	if *ci {
		fmt.Println("- CI mode enabled.")
	}
	if *verbose {
		fmt.Println("- Verbose mode enabled.")
	}

	fmt.Println("\n--- Change Detection ---")
	changes, err := detectChanges()
	if err != nil {
		log.Fatalf("Failed to detect changes: %v", err)
	}

	if changes != "" {
		groups := classifyAndGroupChanges(changes)
		
		summaries := map[string]string{
			"feat":   "implement new features",
			"fix":    "apply automatic fixes",
			"test":   "add or update tests",
			"docs":   "update documentation",
			"chore":  "perform routine maintenance",
		}

		commitCount := 0
		for commitType, files := range groups {
			summary := summaries[commitType]
			message := fmt.Sprintf("%s: %s", commitType, summary)
			
			if err := commitChanges(message, files); err != nil {
				log.Printf("Failed to commit group '%s'. Aborting.", commitType)
				// Optionally, could try to reset changes here
				return
			}
			commitCount++
		}

		if commitCount > 0 && !*noPush {
			pushChanges()
		}
	} else {
		fmt.Println("\nNo changes to commit. Exiting.")
	}
}
