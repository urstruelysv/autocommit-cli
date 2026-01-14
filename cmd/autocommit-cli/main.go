package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"github.com/urstruelysv/autocommit-cli/internal/git"
	"github.com/urstruelysv/autocommit-cli/internal/classify"
	"github.com/urstruelysv/autocommit-cli/internal/history"
	"github.com/urstruelysv/autocommit-cli/internal/ai"
)

// Application entry point
func main() {
	// Define flags
	review := flag.Bool("review", false, "Enable review mode to inspect commits before they are made.")
	yRun := flag.Bool("y-run", false, "Not implemented yet.")
	noPush := flag.Bool("no-push", false, "Create commits but do not push them to the remote repository.")
	ci := flag.Bool("ci", false, "Enable CI mode for non-interactive, deterministic execution.")
	verbose := flag.Bool("verbose", false, "Enable verbose output for debugging purposes.")
	aiCommit := flag.Bool("ai-commit", false, "Use AI to generate commit messages.") // Added this flag

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
	if *aiCommit {
		fmt.Println("- AI Commit mode enabled.")
	}

	learnedData := history.LearnFromHistory()
	if len(learnedData.Scopes) > 0 || len(learnedData.Types) > 0 {
		fmt.Println("History learning complete.")
	} else {
		fmt.Println("No history data learned.")
	}

	fmt.Println("\n--- Change Detection ---")
	changes, err := git.DetectChanges()
	if err != nil {
		log.Fatalf("Failed to detect changes: %v", err)
	}

	if changes != "" {
		if *aiCommit {
			// AI-driven commit message for all changes
			fmt.Println("\n--- AI Commit Message Generation ---")
			message, err := ai.GenerateAICommitMessage(changes)
			if err != nil {
				log.Fatalf("Failed to generate AI commit message: %v", err)
			}
			if err := git.CommitChanges(message, []string{".", "--all"}); err != nil { // Commit all changes
				log.Fatalf("Failed to commit changes with AI message: %v", err)
			}
			if !*noPush {
				git.PushChanges()
			}
		} else {
			// Rule-based and history-aware commit message generation (existing logic)
			groups := classify.ClassifyAndGroupChanges(changes, learnedData)
			
			summaries := map[string]string{
				"feat":   "implement new features",
				"fix":    "apply automatic fixes",
				"test":   "add or update tests",
				"docs":   "update documentation",
				"chore":  "perform routine maintenance",
			}

			commitCount := 0
			for groupKey, files := range groups {
				var commitType string
				if idx := strings.Index(groupKey, "("); idx != -1 {
					commitType = groupKey[:idx]
				} else {
					commitType = groupKey
				}

				summary := summaries[commitType]
				message := fmt.Sprintf("%s: %s", groupKey, summary)
				
				if err := git.CommitChanges(message, files); err != nil {
					log.Printf("Failed to commit group '%s'. Aborting.", groupKey)
					return
				}
				commitCount++
			}

			if commitCount > 0 && !*noPush {
				git.PushChanges()
			}
		}
	} else {
		fmt.Println("\nNo changes to commit. Exiting.")
	}
}