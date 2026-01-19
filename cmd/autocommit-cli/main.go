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
	"github.com/urstruelysv/autocommit-cli/internal/logger"
)

// printWelcomeMessage prints a welcome message with ASCII art and tips.
func printWelcomeMessage() {
	brightRed := "\033[91m" // Bright Red
	reset := "\033[0m"

	asciiArtLines := []string{
		".------..------..------..------..------..------..------..------..------..------.",
		"|A.--. ||U.--. ||T.--. ||O.--. ||C.--. ||O.--. ||M.--. ||M.--. ||I.--. ||T.--. |",
		"| (\\/) || (\\/) || (\\/) || (\\/) || (\\/) || (\\/) || (\\/) || (\\/) || (\\/) || (\\/) |",
		"| :\\/: || :\\/: || :\\/: || :\\/: || :\\/: || :\\/: || :\\/: || :\\/: || :\\/: || :\\/: |",
		"| '--'A|| '--'U|| '--'T|| '--'O|| '--'C|| '--'O|| '--'M|| '--'M|| '--'I|| '--'T|",
		"`------``------``------``------``------``------``------``------``------``------`",
		".------..------..------.",
		"|C.--. ||L.--. ||I.--. |",
		"| (\\/) || (\\/) || (\\/) |",
		"| :\\/: || :\\/: || :\\/: |",
		"| '--'C|| '--'L|| '--'I|",
		"`------``------``------`",
	}

	asciiArt := strings.Join(asciiArtLines, "\n")

	fmt.Printf("%s%s%s\n", brightRed, asciiArt, reset)
	fmt.Println("Welcome to autocommit-cli!")
	fmt.Println("\nTips to get started:")
	fmt.Println("  - Run 'autocommit-cli --help' to see available commands.")
	fmt.Println("  - Configure your preferences in '.autocommitrc'.")
	fmt.Println("  - Make some changes to your git repository and run 'autocommit-cli commit'.")
	fmt.Println("")
}

type AppMode struct {
	Review   bool
	NoPush   bool
	CI       bool
	Verbose  bool
	AICommit bool
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

// Application entry point
func main() {
	printWelcomeMessage()
	// Define flags
	ciFlag := flag.Bool("ci", false, "Enable CI mode for non-interactive, deterministic execution.")
	flag.Parse()

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		// In CI mode, we don't want to see this message
		if !*ciFlag {
			fmt.Println("No .env file found, proceeding without it.")
		}
	}

	var appMode AppMode
	var log logger.Logger

	if *ciFlag {
		appMode = AppMode{CI: true, AICommit: true}
		log = logger.NewJSONLogger()
	} else {
		appMode = promptForMode()
		log = logger.NewHumanReadableLogger()
	}

	log.Info("AutoCommit AI (Go version) running...")

	// Perform initial git status checks
	if err := git.CheckGitStatus(log); err != nil {
		log.Fatal(1, "Git status check failed: %v", err)
	}

	log.Info("Arguments provided:")

	if appMode.Review {
		log.Info("- Review mode enabled.")
	}
	if appMode.NoPush {
		log.Info("- No-push mode enabled.")
	}
	if appMode.CI {
		log.Info("- CI mode enabled.")
	}
	if appMode.Verbose {
		log.Info("- Verbose mode enabled.")
	}
	if appMode.AICommit {
		log.Info("- AI Commit mode enabled.")
	}

	var learnedData history.LearnData

	// Try to load learned data from cache
	learnedData, err = history.LoadLearnedData(log)
	if err != nil {
		log.Info("No cached history data found or error loading. Learning from history...")
		learnedData = history.LearnFromHistory(log)
		if len(learnedData.Scopes) > 0 || len(learnedData.Types) > 0 {
			log.Info("History learning complete.")
			// Save newly learned data to cache
			if err := history.SaveLearnedData(log, learnedData); err != nil {
				log.Error("Failed to save learned data: %v", err)
			}
		} else {
			log.Info("No history data learned.")
		}
	} else {
		log.Info("History data loaded from cache.")
	}

	log.Info("\n--- Change Detection ---")
	changes, err := git.DetectChanges(log)
	if err != nil {
		log.Fatal(1, "Failed to detect changes: %v", err)
	}

	if changes != "" {
		if appMode.AICommit {
			// Check for API key
			if os.Getenv("GEMINI_API_KEY") == "" {
				log.Fatal(1, "GEMINI_API_KEY environment variable not set. Please create a .env file and add your API key.")
			}
			// AI-driven commit message for all changes
			log.Info("\n--- AI Commit Message Generation ---")
			message, err := ai.GenerateAICommitMessage(log, changes)
			if err != nil {
				log.Fatal(1, "Failed to generate AI commit message: %v", err)
			}

			if appMode.Review {
				log.Info("\n--- Review Commit ---")
				log.Info("Commit message:\n%s", message)
				log.Info("\nFiles to be committed:")
				log.Info(changes)
				fmt.Print("Do you want to proceed with this commit? (y/n): ")
				reader := bufio.NewReader(os.Stdin)
				input, _ := reader.ReadString('\n')
				if strings.TrimSpace(input) != "y" {
					log.Info("Commit aborted.")
					return
				}
			}

			if err := git.CommitChanges(log, message, []string{".", "--all"}); err != nil { // Commit all changes
				log.Fatal(1, "Failed to commit changes with AI message: %v", err)
			}
			if !appMode.NoPush {
				if err := git.PushChanges(log); err != nil {
					log.Fatal(1, "Failed to push changes: %v", err)
				}
			}
		} else {
			// Rule-based and history-aware commit message generation (existing logic)
			groups := classify.ClassifyAndGroupChanges(log, changes, learnedData)

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
					log.Info("\n--- Review Commit ---")
					log.Info("Commit message:\n%s", message)
					log.Info("\nFiles to be committed:")
					for _, file := range files {
						log.Info(file)
					}
					fmt.Print("Do you want to proceed with this commit? (y/n): ")
					reader := bufio.NewReader(os.Stdin)
					input, _ := reader.ReadString('\n')
					if strings.TrimSpace(input) != "y" {
						log.Info("Commit aborted.")
						continue
					}
				}

				if err := git.CommitChanges(log, message, files); err != nil {
					log.Error("Failed to commit group '%s'. Aborting.", groupKey)
					return
				}
				commitCount++
			}

			if commitCount > 0 && !appMode.NoPush {
				if err := git.PushChanges(log); err != nil {
					log.Fatal(1, "Failed to push changes: %v", err)
				}
			}
		}
	} else {
		log.Info("\nNo changes to commit. Exiting.")
	}
}