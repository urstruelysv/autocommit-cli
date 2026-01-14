import argparse
import subprocess

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
    detect_changes()


if __name__ == "__main__":
    main()
