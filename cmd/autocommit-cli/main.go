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

/*
Pixel-style ASCII inspired by Press Start 2P.
Terminal-safe. No invalid escapes.
*/
func printWelcomeMessage() {
	brightRed := "\033[91m"
	reset := "\033[0m"

	ascii := `
 █████╗ ██╗   ██╗████████╗ ██████╗  ██████╗ ██████╗ ███╗   ███╗███╗   ███╗██╗████████╗
██╔══██╗██║   ██║╚══██╔══╝██╔═══██╗██╔════╝██╔═══██╗████╗ ████║████╗ ████║██║╚══██╔══╝
███████║██║   ██║   ██║   ██║   ██║██║     ██║   ██║██╔████╔██║██╔████╔██║██║   ██║
██╔══██║██║   ██║   ██║   ██║   ██║██║     ██║   ██║██║╚██╔╝██║██║╚██╔╝██║██║   ██║
██║  ██║╚██████╔╝   ██║   ╚██████╔╝╚██████╗╚██████╔╝██║ ╚═╝ ██║██║ ╚═╝ ██║██║   ██║
╚═╝  ╚═╝ ╚═════╝    ╚═╝    ╚═════╝  ╚═════╝ ╚═════╝ ╚═╝     ╚═╝╚═╝     ╚═╝╚═╝   ╚═╝
`

	fmt.Printf("%s%s%s\n", brightRed, ascii, reset)
	fmt.Println("autocommit-cli — commit smarter, not harder\n")
	fmt.Println("Tips:")
	fmt.Println("  • Press Enter to use AI-Commit (default)")
	fmt.Println("  • Use --ci for non-interactive mode")
	fmt.Println("  • Add GROQ_API_KEY or OPENAI_API_KEY to your .env file\n")
}

type AppMode struct {
	Review   bool
	NoPush   bool
	CI       bool
	Verbose  bool
	AICommit bool
}

func promptForMode() AppMode {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Select a mode (default: AI-Commit):")
	fmt.Println("1. AI-Commit (default)")
	fmt.Println("2. Normal (no AI)")
	fmt.Println("3. Review before commit")
	fmt.Println("4. No-push")
	fmt.Println("5. Verbose")
	fmt.Print("Enter choice (1-5 or Enter): ")

	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				os.Exit(0)
			}
			log.Fatalf("Input error: %v", err)
		}

		switch strings.TrimSpace(input) {
		case "", "1":
			return AppMode{AICommit: true}
		case "2":
			return AppMode{}
		case "3":
			return AppMode{Review: true, AICommit: true}
		case "4":
			return AppMode{NoPush: true, AICommit: true}
		case "5":
			return AppMode{Verbose: true, AICommit: true}
		default:
			fmt.Print("Invalid choice. Enter 1–5: ")
		}
	}
}

func main() {
	printWelcomeMessage()

	ciFlag := flag.Bool("ci", false, "Run in CI mode")
	flag.Parse()

	// Load environment variables from .env when present.
	_ = godotenv.Load()

	var appMode AppMode
	var logg logger.Logger

	if *ciFlag {
		appMode = AppMode{CI: true, AICommit: true}
		logg = logger.NewJSONLogger(false)
	} else {
		appMode = promptForMode()
		// Safety net: ensure Enter defaults to AI-Commit.
		if !appMode.AICommit && !appMode.Review && !appMode.NoPush && !appMode.Verbose {
			appMode.AICommit = true
		}
		logg = logger.NewHumanReadableLogger(appMode.Verbose)
	}

	logg.Info("autocommit-cli started")

	if err := git.CheckGitStatus(logg); err != nil {
		logg.Fatal(1, "Git status check failed: %v", err)
	}

	learnedData, err := history.LoadLearnedData(logg)
	if err != nil {
		learnedData = history.LearnFromHistory(logg)
		_ = history.SaveLearnedData(logg, learnedData)
	}

	changes, err := git.DetectChanges(logg)
	if err != nil {
		logg.Fatal(1, "Change detection failed: %v", err)
	}

	if changes == "" {
		logg.Info("No changes detected. Clean working tree.")
		return
	}

	if appMode.AICommit {
		// Build a richer diff for AI (staged + unstaged + untracked list).
		aiInput, diffErr := git.GetDiffForAI(logg)
		if diffErr != nil {
			logg.Fatal(1, "Diff collection failed: %v", diffErr)
		}
		if strings.TrimSpace(aiInput) == "" {
			logg.Info("No diff content available for AI.")
			return
		}

		provider := strings.ToLower(strings.TrimSpace(os.Getenv("AI_PROVIDER")))
		if provider == "" {
			provider = "groq"
		}
		logg.Info("AI provider: %s", provider)

		var (
			message string
			err     error
		)

		switch provider {
		case "groq":
			if os.Getenv("GROQ_API_KEY") == "" {
				logg.Fatal(1, "GROQ_API_KEY not set")
			}
			message, err = ai.GenerateGroqCommitMessage(logg, aiInput)
			if err == nil {
				logg.Info("AI commit message: %s", message)
			}
			if err != nil {
				fallback := strings.ToLower(strings.TrimSpace(os.Getenv("AI_FALLBACK")))
				if fallback == "" {
					fallback = "basic"
				}

				switch fallback {
				case "basic":
					logg.Info("Groq failed; falling back to basic commit message")
					message = "chore: update changes"
					err = nil
				case "none":
					// Keep original error handling below.
				default:
					logg.Fatal(1, "Unknown AI_FALLBACK: %s (use basic or none)", fallback)
				}
			}
		case "openai", "chatgpt":
			if os.Getenv("OPENAI_API_KEY") == "" {
				logg.Fatal(1, "OPENAI_API_KEY not set")
			}
			message, err = ai.GenerateOpenAICommitMessage(logg, aiInput)
			if err == nil {
				logg.Info("AI commit message: %s", message)
			}
			if err != nil {
				fallback := strings.ToLower(strings.TrimSpace(os.Getenv("AI_FALLBACK")))
				if fallback == "" {
					fallback = "basic"
				}

				switch fallback {
				case "basic":
					logg.Info("OpenAI failed; falling back to basic commit message")
					message = "chore: update changes"
					err = nil
				case "none":
					// Keep original error handling below.
				default:
					logg.Fatal(1, "Unknown AI_FALLBACK: %s (use basic or none)", fallback)
				}
			}
		default:
			logg.Fatal(1, "Unknown AI provider: %s (use groq or openai)", provider)
		}

		if err != nil {
			errMsg := err.Error()
			switch {
			case strings.Contains(errMsg, "exceeded max retries"):
				logg.Fatal(1, "AI commit failed after multiple retries (rate limit). Try again later or set AI_FALLBACK=basic to proceed. Details: %s", errMsg)
			case strings.Contains(errMsg, "invalid_api_key") || strings.Contains(errMsg, "Incorrect API key"):
				logg.Fatal(1, "AI commit failed due to invalid API key. Check GROQ_API_KEY/OPENAI_API_KEY. Details: %s", errMsg)
			default:
				logg.Fatal(1, "AI commit failed: %v", err)
			}
		}

		if err := git.CommitChanges(logg, message, []string{"--all"}); err != nil {
			logg.Fatal(1, "Commit failed: %v", err)
		}

		if !appMode.NoPush {
			_ = git.PushChanges(logg)
		}
		return
	}

	// Non-AI path
	groups := classify.ClassifyAndGroupChanges(logg, changes, learnedData)

	summaries := map[string]string{
		"feat":     "add new functionality",
		"fix":      "fix bugs",
		"docs":     "update documentation",
		"chore":    "maintenance",
		"refactor": "refactor code",
		"test":     "update tests",
	}

	for groupKey, files := range groups {
		commitType := groupKey
		if i := strings.Index(groupKey, "("); i != -1 {
			commitType = groupKey[:i]
		}

		message := fmt.Sprintf("%s: %s", groupKey, summaries[commitType])

		if err := git.CommitChanges(logg, message, files); err != nil {
			logg.Fatal(1, "Commit failed: %v", err)
		}
	}

	if !appMode.NoPush {
		_ = git.PushChanges(logg)
	}
}
