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

type AppMode struct {
	Review   bool
	NoPush   bool
	CI       bool
	Verbose  bool
	AICommit bool
}

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

	var appMode AppMode
	if *ciFlag {
		appMode = AppMode{CI: true, AICommit: true}
	} else {
		appMode = promptForMode()
	}

	fmt.Println("AutoCommit AI (Go version) running...")

	// Perform initial git status checks
	if err := git.CheckGitStatus(appMode.Verbose); err != nil {
		log.Fatalf("Git status check failed: %v", err)
	}

	fmt.Println("Arguments provided:")

	if appMode.Review {
		fmt.Println("- Review mode enabled.")
	}
	if appMode.NoPush {
		fmt.Println("- No-push mode enabled.")
	}
	if appMode.CI {
		fmt.Println("- CI mode enabled.")
	}
	if appMode.Verbose {
		fmt.Println("- Verbose mode enabled.")
	}
	if appMode.AICommit {
		fmt.Println("- AI Commit mode enabled.")
	}

	var learnedData history.LearnData

	// Try to load learned data from cache
	learnedData, err = history.LoadLearnedData(appMode.Verbose)
	if err != nil {
		fmt.Println("No cached history data found or error loading. Learning from history...")
		learnedData = history.LearnFromHistory(appMode.Verbose)
		if len(learnedData.Scopes) > 0 || len(learnedData.Types) > 0 {
			fmt.Println("History learning complete.")
			// Save newly learned data to cache
			if err := history.SaveLearnedData(learnedData, appMode.Verbose); err != nil {
				log.Printf("Failed to save learned data: %v", err)
			}
		} else {
			fmt.Println("No history data learned.")
		}
	} else {
		fmt.Println("History data loaded from cache.")
	}

	fmt.Println("\n--- Change Detection ---")
	changes, err := git.DetectChanges(appMode.Verbose)
	if err != nil {
		log.Fatalf("Failed to detect changes: %v", err)
	}

	if changes != "" {
		if appMode.AICommit {
			// Check for API key
			if os.Getenv("GEMINI_API_KEY") == "" {
				log.Fatalf("GEMINI_API_KEY environment variable not set. Please create a .env file and add your API key.")
			}
			// AI-driven commit message for all changes
			fmt.Println("\n--- AI Commit Message Generation ---")
			message, err := ai.GenerateAICommitMessage(changes, appMode.Verbose)
			if err != nil {
				log.Fatalf("Failed to generate AI commit message: %v", err)
			}

			if appMode.Review {
				fmt.Println("\n--- Review Commit ---")
				fmt.Printf("Commit message:\n%s\n", message)
				fmt.Println("\nFiles to be committed:")
				fmt.Println(changes)
				fmt.Print("Do you want to proceed with this commit? (y/n): ")
				reader := bufio.NewReader(os.Stdin)
				input, _ := reader.ReadString('\n')
				if strings.TrimSpace(input) != "y" {
					fmt.Println("Commit aborted.")
					return
				}
			}

			if err := git.CommitChanges(message, []string{".", "--all"}, appMode.Verbose); err != nil { // Commit all changes
				log.Fatalf("Failed to commit changes with AI message: %v", err)
			}
			if !appMode.NoPush {
				if err := git.PushChanges(appMode.Verbose); err != nil {
					log.Fatalf("Failed to push changes: %v", err)
				}
			}
		} else {
			// Rule-based and history-aware commit message generation (existing logic)
			groups := classify.ClassifyAndGroupChanges(changes, learnedData, appMode.Verbose)

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
							// This part needs more-sophisticated logic to extract common subject lines
							// For now, we'll stick to the summary
						}
					}
				}

				if appMode.Review {
					fmt.Println("\n--- Review Commit ---")
					fmt.Printf("Commit message:\n%s\n", message)
					fmt.Println("\nFiles to be committed:")
					for _, file := range files {
						fmt.Println(file)
					}
					fmt.Print("Do you want to proceed with this commit? (y/n): ")
					reader := bufio.NewReader(os.Stdin)
					input, _ := reader.ReadString('\n')
					if strings.TrimSpace(input) != "y" {
						fmt.Println("Commit aborted.")
						continue
					}
				}

				if err := git.CommitChanges(message, files, appMode.Verbose); err != nil {
					log.Printf("Failed to commit group '%s'. Aborting.", groupKey)
					return
				}
				commitCount++
			}

			if commitCount > 0 && !appMode.NoPush {
				if err := git.PushChanges(appMode.Verbose); err != nil {
					log.Fatalf("Failed to push changes: %v", err)
				}
			}
		}
	} else {
		fmt.Println("\nNo changes to commit. Exiting.")
	}
}

// promptForMode prompts the user to select a mode and returns their choice.
func promptForMode() AppMode {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Select a mode to run autocommit-cli (default: AI-Commit):")
	fmt.Println("1. AI-Commit (default) - Use AI to generate commit messages.")
	fmt.Println("2. Normal - Create commits without AI.")
	fmt.Println("3. Review - Inspect commits before they are made.")
	fmt.Println("4. No-push - Create commits but do not push them to the remote repository.")
	fmt.Println("5. Verbose - Enable verbose output for debugging purposes.")
	fmt.Print("Enter your choice (1-5, or press Enter for default): ")

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
		case "":
			return AppMode{AICommit: true}
		case "1":
			return AppMode{AICommit: true}
		case "2":
			return AppMode{AICommit: false}
		case "3":
			return AppMode{Review: true, AICommit: true}
		case "4":
			return AppMode{NoPush: true, AICommit: true}
		case "5":
			return AppMode{Verbose: true, AICommit: true}
		default:
			fmt.Print("Invalid choice. Please enter a number between 1 and 5: ")
		}
	}
}