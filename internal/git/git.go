package git

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func CheckGitStatus() error {
	// Check for clean working directory (excluding untracked files)
	cmdStatus := exec.Command("git", "status", "--porcelain", "--untracked-files=no")
	outputStatus, err := cmdStatus.Output()
	if err != nil {
		return fmt.Errorf("failed to get git status: %w", err)
	}
	if strings.TrimSpace(string(outputStatus)) != "" {
		return fmt.Errorf("working directory has uncommitted changes. Please commit or stash your changes before running autocommit")
	}

	// Check for detached HEAD
	cmdBranch := exec.Command("git", "symbolic-ref", "--short", "HEAD")
	_, err = cmdBranch.Output()
	if err != nil {
		return fmt.Errorf("detached HEAD state detected. Please checkout a branch before running autocommit")
	}

	// Check if the current branch has an upstream configured
	cmdRevParse := exec.Command("git", "rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{u}")
	if err := cmdRevParse.Run(); err != nil {
		return fmt.Errorf("current branch does not have an upstream branch configured. Please set an upstream branch (e.g., 'git push -u origin <branch_name>')")
	}

	return nil
}

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
	// Check if remote is configured
	branchCmd := exec.Command("git", "symbolic-ref", "--short", "HEAD")
	branchOutput, err := branchCmd.Output()
	if err != nil {
		log.Println("Error: Could not determine current branch. Aborting push.")
		return
	}
	branchName := strings.TrimSpace(string(branchOutput))

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
