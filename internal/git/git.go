package git

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
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

// GetDiffForAI returns a richer diff for AI prompts, including staged/unstaged diffs and untracked files list.
func GetDiffForAI(log logger.Logger) (string, error) {
	log.Debug("Collecting diff for AI...")

	var parts []string

	unstagedCmd := exec.Command("git", "diff", "--no-color")
	unstagedOut, err := unstagedCmd.Output()
	if err != nil {
		return "", fmt.Errorf("error getting unstaged diff: %w", err)
	}
	if len(unstagedOut) > 0 {
		parts = append(parts, "UNSTAGED DIFF:\n"+string(unstagedOut))
	}

	stagedCmd := exec.Command("git", "diff", "--cached", "--no-color")
	stagedOut, err := stagedCmd.Output()
	if err != nil {
		return "", fmt.Errorf("error getting staged diff: %w", err)
	}
	if len(stagedOut) > 0 {
		parts = append(parts, "STAGED DIFF:\n"+string(stagedOut))
	}

	untrackedCmd := exec.Command("git", "ls-files", "--others", "--exclude-standard")
	untrackedOut, err := untrackedCmd.Output()
	if err != nil {
		return "", fmt.Errorf("error getting untracked files: %w", err)
	}
	untracked := strings.TrimSpace(string(untrackedOut))
	if untracked != "" {
		parts = append(parts, "UNTRACKED FILES:\n"+untracked+"\n")
	}

	if len(parts) == 0 {
		return "", nil
	}

	diff := strings.Join(parts, "\n")

	maxChars := 8000
	if v := strings.TrimSpace(os.Getenv("AI_DIFF_MAX_CHARS")); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
			maxChars = parsed
		}
	}

	if maxChars > 0 && len(diff) > maxChars {
		keepHead := maxChars / 2
		keepTail := maxChars - keepHead
		diff = diff[:keepHead] + "\n... [diff truncated] ...\n" + diff[len(diff)-keepTail:]
	}

	return diff, nil
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
