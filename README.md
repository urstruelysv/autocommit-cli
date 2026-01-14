# AutoCommit AI (CLI-first)

## AI-powered Git Automation Tool

---

## Overview

AutoCommit AI is an **AI-powered Git automation tool** designed for **solo developers** and **CI/CD systems**. It automatically detects code changes, intelligently groups them into logical commits, generates high-quality commit messages, optionally allows quick review/editing, and safely pushes them to the correct remote branch.

The product focuses on **zero-friction usage**, **easy installation across platforms**, and **learning from the developer’s past commit history and guidelines** to continuously improve commit quality.

---

## Problem Statement

Git workflows introduce unnecessary friction:

*   Deciding *what belongs in a commit*
*   Writing accurate commit messages
*   Maintaining consistency with commit standards
*   Repeating the same steps across machines, OSes, and CI pipelines

Existing tools either:

*   Only generate commit messages, or
*   Require too much interaction, or
*   Don’t learn historical commit patterns

There is no **simple, intelligent, fully-automated, and install-anywhere** solution.

---

## Goals

*   Fully automate Git commits **smartly and safely**
*   Group changes into **logical commits** (feature / fix / refactor / docs / chore)
*   Learn from **past commits and commit guides**
*   Provide **optional human review** without slowing flow
*   Be trivial to install on **any platform**
*   Work reliably in **local dev and CI/CD**

### Success Metrics

*   First successful commit within 2 minutes of install
*   Zero required configuration for basic usage
*   Consistent commit quality aligned with repo history
*   CI-safe with deterministic behavior

---

## Key Differentiators (vs Existing Tools)

| Feature                    | Existing Tools | AutoCommit AI  |
| :------------------------- | :------------- | :------------- |
| Logical commit grouping    | ❌ Rare         | ✅ Core feature |
| Learns from commit history | ❌              | ✅              |
| Commit guide awareness     | ❌              | ✅              |
| Fully automated push       | ⚠️ Partial     | ✅              |
| Review + instant accept    | ⚠️             | ✅              |
| Multi-platform install     | ⚠️             | ✅              |

---

## User Stories

*   As a solo dev, I want my changes committed without thinking about Git.
*   As a developer, I want commits to match my existing style.
*   As a CI system, I want deterministic, non-interactive commits.
*   As a new contributor, I want commit rules to be enforced automatically.

---

## Core Features (Current & Planned)

### Current (Go MVP)

*   **Change Detection:** Automatically detects staged and unstaged changes in a Git repository.
*   **Logical Commit Grouping:** Groups detected changes into logical categories (e.g., `feat`, `fix`, `test`, `docs`, `chore`) based on file paths and diff content. Each group results in a separate commit.
*   **Basic Commit Message Generation:** Generates conventional commit messages (e.g., `fix: apply automatic fixes`) for each logical group.
*   **Safe Commit & Push:** Stages and commits changes, with safeguards to prevent pushing from a detached HEAD or to a branch without a configured remote. Includes a `--no-push` flag.
*   **History Learning (Initial):** Extracts potential commit scopes from `git log` for future intelligent message generation.

### Planned (from PRD)

*   **Easy Installation:** Homebrew, PowerShell/Scoop/Winget, single binary, npm/pip (via wrapper), VS Code Extension.
*   **Intelligent Change Classification:** Advanced language-aware heuristics and learned patterns.
*   **Learning From History:** Continuously improves commit quality based on past commits and guidelines.
*   **Commit Guide Awareness:** Automatically detects and adheres to project-specific commit guidelines.
*   **Review & Edit Mode:** Optional interactive mode to review and edit proposed commits before finalization.
*   **CI/CD Mode:** `--ci` flag for non-interactive, deterministic execution in CI environments.

---

## Installation

### Prerequisites

*   Go (version 1.18 or higher recommended)
*   Git (installed and configured)

### From Source

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/your-repo/autocommit-ai.git # Replace with actual repo URL
    cd autocommit-ai
    ```
2.  **Build the executable:**
    ```bash
    go build -o autocommit main.go
    ```
3.  **Move to your PATH (optional):**
    ```bash
    sudo mv autocommit /usr/local/bin/
    ```

### Future Installation Methods

*   Homebrew (macOS/Linux)
*   Scoop/Winget (Windows)
*   Pre-built binaries

---

## CLI Usage

```bash
autocommit [flags]
```

### Flags

*   `--review`: Enable review mode to inspect commits before they are made (Planned).
*   `--y-run`: Not implemented yet.
*   `--no-push`: Create commits but do not push them to the remote repository.
*   `--ci`: Enable CI mode for non-interactive, deterministic execution (Planned).
*   `--verbose`: Enable verbose output for debugging purposes (Planned).

### Examples

*   **Automatically commit and push all changes:**
    ```bash
    autocommit
    ```
*   **Commit changes without pushing:**
    ```bash
    autocommit --no-push
    ```
*   **Review proposed commits before committing (Planned):**
    ```bash
    autocommit --review
    ```

---

## Current Status

The project is currently in active development. The core MVP features (change detection, basic commit message generation, logical grouping, safe commit/push) have been implemented in Go. We are now working on enhancing the intelligence and learning capabilities.

---

## Future Scope

*   PR creation
*   GitHub/GitLab App
*   IDE-native integrations
*   Semantic version automation

---
sai vamshi