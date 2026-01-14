import argparse
import subprocess
import os
from collections import Counter

def detect_changes():
    """
    Detects staged and unstaged changes in the repository.
    """
    print("Detecting changes...")
    try:
        result = subprocess.run(
            ["git", "status", "--porcelain"],
            capture_output=True,
            text=True,
            check=True
        )
        changes = result.stdout.strip()
        if changes:
            print("Found changes:")
            print(changes)
        else:
            print("No changes found.")
        return changes
    except subprocess.CalledProcessError as e:
        print(f"Error detecting changes: Not a git repository or no git installed?")
        print(e.stderr)
        return None
    except FileNotFoundError:
        print("Error: 'git' command not found. Is Git installed and in your PATH?")
        return None

def generate_commit_message(changes):
    """
    Generates a commit message based on change classification.
    """
    print("\n--- Commit Message Generation ---")
    
    change_types = []
    file_paths = [line.split()[-1] for line in changes.splitlines()]

    for file_path in file_paths:
        if "tests/" in file_path or "test_" in file_path:
            change_types.append("test")
            continue
        
        if file_path.endswith(".md"):
            change_types.append("docs")
            continue

        try:
            diff_result = subprocess.run(
                ["git", "diff", "--", file_path],
                capture_output=True, text=True, check=True
            )
            diff = diff_result.stdout.lower()

            if "fix" in diff or "bug" in diff:
                change_types.append("fix")
            elif "add" in diff or "feature" in diff:
                change_types.append("feat")
            else:
                change_types.append("chore")
        except Exception as e:
            print(f"Could not get diff for {file_path}: {e}")
            change_types.append("chore")

    # Determine the most common change type
    if not change_types:
        commit_type = "chore"
    else:
        most_common_type = Counter(change_types).most_common(1)[0][0]
        commit_type = most_common_type

    # Generate a summary based on the type
    summaries = {
        "feat": "implement new features",
        "fix": "apply automatic fixes",
        "test": "add or update tests",
        "docs": "update documentation",
        "chore": "perform routine maintenance"
    }
    summary = summaries.get(commit_type, "perform routine maintenance")

    message = f"{commit_type}: {summary}"
    print(f"Generated message: {message}")
    return message

def commit_changes(message):
    """
    Stages all changes and commits them with the given message.
    """
    print("\n--- Committing Changes ---")
    try:
        subprocess.run(["git", "add", "."], check=True)
        print("Staged all changes.")
        subprocess.run(["git", "commit", "-m", message], check=True)
        print("Committed changes.")
        return True
    except subprocess.CalledProcessError as e:
        print(f"Error during commit process: {e}")
        print(e.stderr)
        return False
    except FileNotFoundError:
        print("Error: 'git' command not found. Is Git installed and in your PATH?")
        return False

def push_changes():
    """
    Pushes changes to the remote repository after safety checks.
    """
    print("\n--- Pushing Changes ---")
    try:
        # Check for detached HEAD
        branch_result = subprocess.run(
            ["git", "symbolic-ref", "--short", "HEAD"],
            capture_output=True, text=True
        )
        if branch_result.returncode != 0:
            print("Error: Detached HEAD state. Aborting push.")
            return

        branch_name = branch_result.stdout.strip()

        # Check if remote is configured
        remote_result = subprocess.run(
            ["git", "config", f"branch.{branch_name}.remote"],
            capture_output=True, text=True
        )
        if remote_result.returncode != 0:
            print(f"Error: No remote configured for branch '{branch_name}'. Aborting push.")
            return

        print(f"Pushing changes to remote for branch '{branch_name}'...")
        subprocess.run(["git", "push"], check=True)
        print("Push successful.")

    except subprocess.CalledProcessError as e:
        print(f"Error during push: {e}")
        print(e.stderr)
    except FileNotFoundError:
        print("Error: 'git' command not found. Is Git installed and in your PATH?")


def main():
    """
    Main function for the AutoCommit AI CLI.
    """
    parser = argparse.ArgumentParser(
        description="AutoCommit AI: AI-powered Git automation tool."
    )

    parser.add_argument(
        "--review",
        action="store_true",
        help="Enable review mode to inspect commits before they are made."
    )

    parser.add_argument(
        "--y-run",
        action="store_true",
        help="Not implemented yet."
    )

    parser.add_argument(
        "--no-push",
        action="store_true",
        help="Create commits but do not push them to the remote repository."
    )

    parser.add_argument(
        "--ci",
        action="store_true",
        help="Enable CI mode for non-interactive, deterministic execution."
    )

    parser.add_argument(
        "--verbose",
        action="store_true",
        help="Enable verbose output for debugging purposes."
    )

    args = parser.parse_args()

    print("AutoCommit AI running...")
    print("Arguments provided:")
    if args.review:
        print("- Review mode enabled.")
    if args.y_run:
        print("- y-run flag set.")
    if args.no_push:
        print("- No-push mode enabled.")
    if args.ci:
        print("- CI mode enabled.")
    if args.verbose:
        print("- Verbose mode enabled.")

    print("\n--- Change Detection ---")
    changes = detect_changes()

    if changes:
        message = generate_commit_message(changes)
        commit_successful = commit_changes(message)
        if commit_successful and not args.no_push:
            push_changes()
    else:
        print("\nNo changes to commit. Exiting.")


# Application entry point
if __name__ == "__main__":
    main()
