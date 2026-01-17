package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/urstruelysv/autocommit-cli/internal/ai"
	"github.com/urstruelysv/autocommit-cli/internal/classify"
	"github.com/urstruelysv/autocommit-cli/internal/git"
	"github.com/urstruelysv/autocommit-cli/internal/history"
)

// Application entry point
func main() {
	// Define flags
	ciFlag := flag.Bool("ci", false, "Enable CI mode for non-interactive, deterministic execution.")
	flag.Parse()

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, proceeding without it.")
	}

	var mode string
	if *ciFlag {
		mode = "ci"
	} else {
		mode = promptForMode()
	}

	// Set flags based on mode
	review := mode == "review"
	yRun := mode == "y-run"
	noPush := mode == "no-push"
	ci := mode == "ci"
	verbose := mode == "verbose"
	aiCommit := mode == "ai-commit"

	fmt.Println("AutoCommit AI (Go version) running...")

	// Perform initial git status checks
	if err := git.CheckGitStatus(); err != nil {
		log.Fatalf("Git status check failed: %v", err)
	}

	fmt.Println("Arguments provided:")

	if review {
		fmt.Println("- Review mode enabled.")
	}
	if yRun {
		fmt.Println("- y-run flag set.")
	}
	if noPush {
		fmt.Println("- No-push mode enabled.")
	}
	if ci {
		fmt.Println("- CI mode enabled.")
	}
	if verbose {
		fmt.Println("- Verbose mode enabled.")
	}
	if aiCommit {
		fmt.Println("- AI Commit mode enabled.")
	}

	var learnedData history.LearnData

	// Try to load learned data from cache
	learnedData, err = history.LoadLearnedData()
	if err != nil {
		fmt.Println("No cached history data found or error loading. Learning from history...")
		learnedData = history.LearnFromHistory()
		if len(learnedData.Scopes) > 0 || len(learnedData.Types) > 0 {
			fmt.Println("History learning complete.")
			// Save newly learned data to cache
			if err := history.SaveLearnedData(learnedData); err != nil {
				log.Printf("Failed to save learned data: %v", err)
			}
		} else {
			fmt.Println("No history data learned.")
		}
	} else {
		fmt.Println("History data loaded from cache.")
	}

	fmt.Println("\n--- Change Detection ---")
	changes, err := git.DetectChanges()
	if err != nil {
		log.Fatalf("Failed to detect changes: %v", err)
	}

	if changes != "" {
		if aiCommit {
			// Check for API key
			if os.Getenv("GEMINI_API_KEY") == "" {
				log.Fatalf("GEMINI_API_KEY environment variable not set. Please create a .env file and add your API key.")
			}
			// AI-driven commit message for all changes
			fmt.Println("\n--- AI Commit Message Generation ---")
			message, err := ai.GenerateAICommitMessage(changes)
			if err != nil {
				log.Fatalf("Failed to generate AI commit message: %v", err)
			}
			if err := git.CommitChanges(message, []string{".", "--all"}); err != nil { // Commit all changes
				log.Fatalf("Failed to commit changes with AI message: %v", err)
			}
			if !noPush {
				git.PushChanges()
			}
		} else {
			// Rule-based and history-aware commit message generation (existing logic)
			groups := classify.ClassifyAndGroupChanges(changes, learnedData)

			summaries := map[string]string{
				"feat":     "implement new features",
				"fix":      "apply automatic fixes",
				"test":     "add or update tests",
				"docs":     "update documentation",
				"chore":    "perform routine maintenance",
				"refactor": "refactor code", // Added refactor summary
			}

			commitCount := 0
			for groupKey, files := range groups {
				var commitType string

				// Extract type and scope from groupKey
				if idx := strings.Index(groupKey, "("); idx != -1 {
					commitType = groupKey[:idx]
					// scope = groupKey[idx+1 : len(groupKey)-1] // Scope is not directly used here
				} else {
					commitType = groupKey
				}

				summary := summaries[commitType]
				message := fmt.Sprintf("%s: %s", groupKey, summary)

				// Try to make the message more history-aware
				if len(learnedData.Types) > 0 {
					// Find the most common subject for this type
					// This is a very basic implementation, can be improved
					for t, count := range learnedData.Types {
						if t == commitType && count > 1 { // Only if type is common
							// This part needs more sophisticated logic to extract common subject lines
							// For now, we'll stick to the summary
						}
					}
				}

				if err := git.CommitChanges(message, files); err != nil {
					log.Printf("Failed to commit group '%s'. Aborting.", groupKey)
					return
				}
				commitCount++
			}

			if commitCount > 0 && !noPush {
				git.PushChanges()
			}
		}
	} else {
		fmt.Println("\nNo changes to commit. Exiting.")
	}
}

// promptForMode prompts the user to select a mode and returns their choice.
func promptForMode() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Select a mode to run autocommit-cli:")
	fmt.Println("1. Normal - Create commits and push them to the remote repository.")
	fmt.Println("2. Review - Inspect commits before they are made.")
	fmt.Println("3. No-push - Create commits but do not push them to the remote repository.")
	fmt.Println("4. CI - Non-interactive, deterministic execution for CI environments.")
	fmt.Println("5. Verbose - Enable verbose output for debugging purposes.")
	fmt.Println("6. AI-Commit - Use AI to generate commit messages.")
	fmt.Print("Enter your choice (1-6): ")

	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("\nNo input received, exiting.")
				os.Exit(0)
			}
			log.Fatalf("Failed to read input: %v", err)
		}
		input = strings.TrimSpace(input)
		switch input {
		case "1":
			return "normal"
		case "2":
			return "review"
		case "3":
			return "no-push"
		case "4":
			return "ci"
		case "5":
			return "verbose"
		case "6":
			return "ai-commit"
		default:
			fmt.Print("Invalid choice. Please enter a number between 1 and 6: ")
		}
	}
}