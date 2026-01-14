package git

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func DetectChanges() (string, error) {
	fmt.Println("Detecting changes...")
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Error detecting changes: %v", err)
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

func CommitChanges(message string, files []string) error {
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

func PushChanges() {
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
		log.Printf("Error: No remote configured for branch '%s'. Aborting push.", branchName)
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
