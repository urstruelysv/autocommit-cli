package git

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/urstruelysv/autocommit-cli/internal/logger"
)

func CheckGitStatus(log logger.Logger) error {
	log.Debug("Checking git status...")
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

	log.Debug("Git status checks passed.")
	return nil
}

func DetectChanges(log logger.Logger) (string, error) {
	log.Debug("Detecting changes...")
	log.Info("Detecting changes...")
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		log.Error("Error detecting changes: %v", err)
		return "", err
	}

	changes := strings.TrimSpace(string(output))
	if changes != "" {
		log.Debug("Found changes.")
		log.Info("Found changes:")
		log.Info(changes)
	} else {
		log.Debug("No changes found.")
		log.Info("No changes found.")
	}
	return changes, nil
}

func CommitChanges(log logger.Logger, message string, files []string) error {
	log.Debug("Committing group with message: %s", message)
	log.Info("\n--- Committing Group: %s ---", message)

	addArgs := append([]string{"add"}, files...)
	addCmd := exec.Command("git", addArgs...)
	if output, err := addCmd.CombinedOutput(); err != nil {
		log.Error("Error staging files %v: %s\n%v", files, string(output), err)
		return err
	}
	log.Debug("Staged files: %v", files)
	log.Info("Staged files: %v", files)

	commitCmd := exec.Command("git", "commit", "-m", message)
	if output, err := commitCmd.CombinedOutput(); err != nil {
		log.Error("Error committing group: %s\n%v", string(output), err)
		return err
	}
	log.Debug("Committed group.")
	log.Info("Committed group.")
	return nil
}

func PushChanges(log logger.Logger) error {
	log.Debug("Pushing changes...")
	log.Info("\n--- Pushing Changes ---")
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

	log.Debug("Pushing changes to remote for branch '%s'...", branchName)
	log.Info("Pushing changes to remote for branch '%s'...", branchName)
	pushCmd := exec.Command("git", "push")
	if output, err := pushCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("error during push: %s\n%w", string(output), err)
	}
	log.Debug("Push successful.")
	log.Info("Push successful.")
	return nil
}
