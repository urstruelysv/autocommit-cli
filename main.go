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

func generateCommitMessage() string {
	fmt.Println("\n--- Commit Message Generation ---")
	message := "chore: automatic commit of all changes"
	fmt.Printf("Generated message: %s\n", message)
	return message
}

func commitChanges(message string) error {
	fmt.Println("\n--- Committing Changes ---")
	addCmd := exec.Command("git", "add", ".")
	if output, err := addCmd.CombinedOutput(); err != nil {
		log.Printf("Error staging changes: %s\n%v", string(output), err)
		return err
	}
	fmt.Println("Staged all changes.")

	commitCmd := exec.Command("git", "commit", "-m", message)
	if output, err := commitCmd.CombinedOutput(); err != nil {
		log.Printf("Error committing changes: %s\n%v", string(output), err)
		return err
	}
	fmt.Println("Committed changes.")
	return nil
}

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
		message := generateCommitMessage()
		if err := commitChanges(message); err != nil {
			log.Fatalf("Failed to commit changes: %v", err)
		}
	} else {
		fmt.Println("\nNo changes to commit. Exiting.")
	}
}
