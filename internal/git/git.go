package git

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func CheckGitStatus(verbose bool) error {
	if verbose {
		fmt.Println("Verbose: Checking git status...")
	}
	// Check for staged but uncommitted changes
	cmdStaged := exec.Command("git", "diff", "--cached", "--quiet")
	if err := cmdStaged.Run(); err != nil {
		return fmt.Errorf("there are staged but uncommitted changes. Please commit or unstage them before running autocommit")
	}

	// Check for detached HEAD
	cmdBranch := exec.Command("git", "symbolic-ref", "--short", "HEAD")
	_, err := cmdBranch.Output()
	if err != nil {
		return fmt.Errorf("detached HEAD state detected. Please checkout a branch before running autocommit")
	}

	// Check if the current branch has an upstream configured
	cmdRevParse := exec.Command("git", "rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{u}")
	if err := cmdRevParse.Run(); err != nil {
		return fmt.Errorf("current branch does not have an upstream branch configured. Please set an upstream branch (e.g., 'git push -u origin <branch_name>')")
	}

	if verbose {
		fmt.Println("Verbose: Git status checks passed.")
	}
	return nil
}

func DetectChanges(verbose bool) (string, error) {
	if verbose {
		fmt.Println("Verbose: Detecting changes...")
	}
	fmt.Println("Detecting changes...")
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Error detecting changes: %v", err)
		return "", err
	}

	changes := strings.TrimSpace(string(output))
	if changes != "" {
		if verbose {
			fmt.Println("Verbose: Found changes.")
		}
		fmt.Println("Found changes:")
		fmt.Println(changes)
	} else {
		if verbose {
			fmt.Println("Verbose: No changes found.")
		}
		fmt.Println("No changes found.")
	}
	return changes, nil
}

func CommitChanges(message string, files []string, verbose bool) error {
	if verbose {
		fmt.Printf("Verbose: Committing group with message: %s\n", message)
	}
	fmt.Printf("\n--- Committing Group: %s ---\n", message)

	addArgs := append([]string{"add"}, files...)
	addCmd := exec.Command("git", addArgs...)
	if output, err := addCmd.CombinedOutput(); err != nil {
		log.Printf("Error staging files %v: %s\n%v", files, string(output), err)
		return err
	}
	if verbose {
		fmt.Printf("Verbose: Staged files: %v\n", files)
	}
	fmt.Printf("Staged files: %v\n", files)

	commitCmd := exec.Command("git", "commit", "-m", message)
	if output, err := commitCmd.CombinedOutput(); err != nil {
		log.Printf("Error committing group: %s\n%v", string(output), err)
		return err
	}
	if verbose {
		fmt.Println("Verbose: Committed group.")
	}
	fmt.Println("Committed group.")
	return nil
}

func PushChanges(verbose bool) error {
	if verbose {
		fmt.Println("Verbose: Pushing changes...")
	}
	fmt.Println("\n--- Pushing Changes ---")
	// Check if remote is configured
	branchCmd := exec.Command("git", "symbolic-ref", "--short", "HEAD")
	branchOutput, err := branchCmd.Output()
	if err != nil {
		return fmt.Errorf("could not determine current branch: %w", err)
	}
	branchName := strings.TrimSpace(string(branchOutput))

	remoteCmd := exec.Command("git", "config", fmt.Sprintf("branch.%s.remote", branchName))
	if err := remoteCmd.Run(); err != nil {
		return fmt.Errorf("no remote configured for branch '%s'", branchName)
	}

	if verbose {
		fmt.Printf("Verbose: Pushing changes to remote for branch '%s'...\n", branchName)
	}
	fmt.Printf("Pushing changes to remote for branch '%s'...\n", branchName)
	pushCmd := exec.Command("git", "push")
	if output, err := pushCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("error during push: %s\n%w", string(output), err)
	}
	if verbose {
		fmt.Println("Verbose: Push successful.")
	}
	fmt.Println("Push successful.")
	return nil
}
