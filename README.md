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

*   **Easy Installation:** Homebrew, PowerShell/Scoop/Winget, single binary, npm/pip (via wrapper), VS Code Extension.
*   **Project Structure:** Refactored into `cmd/autocommit-cli` and `internal/` packages (`git`, `classify`, `history`, `ai`).
*   **Change Detection:** Automatically detects staged and unstaged changes in a Git repository.
*   **Logical Commit Grouping:** Groups detected changes into logical categories (e.g., `feat`, `fix`, `test`, `docs`, `chore`) based on file paths, diff content, and **folder/module structure (e.g., `fix(git):`)**. Each group results in a separate commit. (Note: This is not used when AI-mode is enabled).
*   **Basic Commit Message Generation:** Generates conventional commit messages (e.g., `fix: apply automatic fixes`) for each logical group, now incorporating module scopes.
*   **AI-Assisted Commit Message Generation:** (Default) Uses the Groq or OpenAI API to generate a single commit message for all changes. This will create a single commit for all the changes and does not perform logical commit grouping.
*   **Safe Commit & Push:** Stages and commits changes, with safeguards to prevent pushing from a detached HEAD or to a branch without a configured remote. Includes a `--no-push` flag.
*   **History Learning (Initial):** Extracts potential commit scopes and types from `git log` for future intelligent message generation.
*   **Interactive Mode Selection:** Prompts the user to select a mode of operation.
*   **CI/CD Mode:** `--ci` flag for non-interactive, deterministic execution in CI environments.
*   **Review & Edit Mode:** Optional interactive mode to review and edit proposed commits before finalization.
*   **Verbose Mode:** Enable verbose output for debugging purposes.

### Planned (from PRD)
*   **Intelligent Change Classification:** Advanced language-aware heuristics and learned patterns.
*   **Learning From History:** Continuously improves commit quality based on past commits and guidelines.
*   **Commit Guide Awareness:** Automatically detects and adheres to project-specific commit guidelines.

---

## Installation

Get started with autocommit-cli in minutes. Follow these simple steps to integrate AI-powered commit messages into your workflow.

### Easy Installation

Homebrew (macOS & Linux)

Install the binary using Homebrew:

```bash
brew tap urstruelysv/autocommit-cli
brew install autocommit-cli
```

npm (Cross-platform)

Install a wrapper via npm (requires Node.js):

```bash
npm install -g @urstruelysv/autocommit-cli
```

### Prerequisites

*   Go (version 1.18 or higher recommended)
*   Git (installed and configured)
*   Groq API Key: Obtain a free API key from GroqCloud. Set it as an environment variable `GROQ_API_KEY` or in a `.env` file.

### Installation from Source (Development)

```bash
# Clone the repository
git clone https://github.com/urstruelysv/autocommit-cli.git
cd autocommit-cli

# Set up API Key (Option 1: Export as environment variable)
export GROQ_API_KEY="YOUR_GROQ_API_KEY"

# OR (Option 2: Create a .env file)
echo 'GROQ_API_KEY="YOUR_GROQ_API_KEY"' > .env
# Then load it (e.g., using a tool like `direnv` or manually `source .env`)

# Run the application
go run cmd/autocommit-cli/main.go
```

### Building and Installing the Executable

```bash
# Clone the repository
git clone https://github.com/urstruelysv/autocommit-cli.git
cd autocommit-cli

# Run the build script
chmod +x build.sh
./build.sh
```

This will build the application and install it in your `/usr/local/bin` directory.

---

## CLI Usage

When you run the application, you will be prompted to select a mode of operation:

```
Select a mode to run autocommit-cli (default: AI-Commit):
1. AI-Commit (default) - Use AI to generate commit messages.
2. Normal - Create commits without AI.
3. Review - Inspect commits before they are made.
4. No-push - Create commits but do not push them to the remote repository.
5. CI - Non-interactive, deterministic execution for CI environments.
6. Verbose - Enable verbose output for debugging purposes.
Enter your choice (1-6, or press Enter for default):
```

### CI Mode

For non-interactive environments like CI/CD pipelines, you can use the `--ci` flag:

```bash
autocommit-cli --ci
```

This will run the application in CI mode, which is non-interactive and deterministic.

### AI Providers (Groq or OpenAI)

Environment variables:
* `AI_PROVIDER` — optional; `groq` (default if `GROQ_API_KEY` is set) or `openai`.
* `GROQ_API_KEY` — required when using Groq.
* `GROQ_MODEL` — optional; defaults to `llama-3.1-8b-instant`.
* `GROQ_MAX_TOKENS` — optional; defaults to `80`.
* `GROQ_TEMPERATURE` — optional; defaults to `0.2`.
* `OPENAI_API_KEY` — required when using OpenAI.
* `OPENAI_MODEL` — optional; defaults to `gpt-5-mini`.
* `AI_FALLBACK` — optional; what to do if OpenAI fails (`basic` or `none`). Defaults to `basic`.
* `AI_DIFF_MAX_CHARS` — optional; max chars of diff sent to OpenAI (default: `8000`).
* `OPENAI_MAX_OUTPUT_TOKENS` — optional; max tokens for response (default: `120`).

Example `.env` (Groq default):
```bash
GROQ_API_KEY="YOUR_GROQ_API_KEY"
GROQ_MODEL="llama-3.1-8b-instant"
GROQ_MAX_TOKENS="80"
GROQ_TEMPERATURE="0.2"
AI_FALLBACK="basic"
AI_DIFF_MAX_CHARS="8000"
```

Example `.env` (OpenAI):
```bash
OPENAI_API_KEY="YOUR_OPENAI_API_KEY"
OPENAI_MODEL="gpt-5-mini"
AI_FALLBACK="basic"
AI_DIFF_MAX_CHARS="8000"
OPENAI_MAX_OUTPUT_TOKENS="120"
```

### Examples

*   **Automatically commit and push all changes with AI:**
    Run the application and press Enter to select the default "AI-Commit" mode.
*   **Commit changes without AI:**
    Run the application and select "Normal" mode.
*   **Review commits before they are made:**
    Run the application and select "Review" mode.
*   **Commit changes without pushing:**
    Run the application and select "No-push" mode.
*   **Run in CI mode:**
    ```bash
    autocommit-cli --ci
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
