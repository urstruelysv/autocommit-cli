# AutoCommit AI: Action Plan

This document outlines the development roadmap for AutoCommit AI, breaking down the project into phased milestones based on the Product Requirements Document (PRD).

---

## Phase 1: Minimum Viable Product (MVP) - Core Automation

**Goal:** Create a functional CLI tool that can detect all changes, generate a single commit message, and safely push to the remote branch. This phase prioritizes core functionality and safety over intelligence.

1.  **Project Setup & Scaffolding:**
    *   Initialize Git repository.
    *   Choose primary language (e.g., Python or Go for cross-platform CLI capabilities).
    *   Set up project structure (`/src`, `/tests`, `/scripts`).
    *   Implement basic command-line argument parsing.

2.  **Change Detection:**
    *   Implement logic to detect staged and unstaged changes using Git commands.
    *   Ensure `.gitignore` is respected.
    *   Create a function to gather a unified diff of all changes.

3.  **Basic Commit Message Generation:**
    *   Implement a rule-based commit message generator based on Conventional Commits.
    *   Start with a simple implementation: `chore: automatic commit of all changes`.
    *   This will be the fallback mechanism in later phases.

4.  **Core `autocommit` Command:**
    *   Create the main entry point for the `autocommit` command.
    *   The command will:
        1.  Detect changes.
        2.  Generate a commit message.
        3.  Stage all changes (`git add .`).
        4.  Commit the changes (`git commit -m "..."`).

5.  **Safe Push Logic:**
    *   Implement safeguards:
        *   Check for a clean working directory before starting.
        *   Verify the current branch and remote tracking information.
        *   Abort if on a detached `HEAD`.
    *   Implement the push functionality (`git push`).
    *   Add the `--no-push` flag to disable automatic pushing.

6.  **Basic Installation:**
    *   Create a `pip` package (`setup.py` or `pyproject.toml`) or equivalent for the chosen language.
    *   Write a simple installation script (`install.sh`).

---

## Phase 2: Intelligence & Learning

**Goal:** Introduce the "smart" features that provide the core value proposition: logical grouping and history-aware commit messages.

1.  **Intelligent Change Classification:**
    *   Develop a module to analyze diffs and classify them (`feat`, `fix`, `refactor`, `docs`, `test`, `chore`).
    *   Use file extensions, keywords (e.g., "fix," "add," "update"), and path conventions (e.g., `/tests/`) for classification.

2.  **Logical Commit Grouping (Core Innovation):**
    *   Implement the initial algorithm for grouping related files into separate commit candidates.
    *   Start with a heuristic-based approach: group files by change type and directory proximity.
    *   The output should be a list of proposed commit groups, each with its own set of files.

3.  **History-Aware Message Generation:**
    *   Implement a module to parse the repository's `git log`.
    *   Analyze past commit messages to infer common scopes `(<scope>)` and wording patterns.
    *   Enhance the commit message generator to use this historical context. For each commit group, generate a tailored message (e.g., `feat(auth): add user login endpoint`).

4.  **Learning Module:**
    *   Create a mechanism to store learned patterns (e.g., in a local, repo-specific cache file like `.autocommit_cache`).
    *   Ensure this data is never shared externally.

---

## Phase 3: User Experience & Platform Support

**Goal:** Refine the user workflow, add optional review, and broaden platform accessibility.

1.  **Review & Edit Mode (`--review`):**
    *   Implement the interactive review flow.
    *   Display the proposed commit groups and their generated messages in a clean, readable format.
    *   Allow the user to quickly edit messages inline.
    *   Implement a single-key confirmation to accept all proposed commits and proceed with the commit/push sequence.

2.  **Configuration File:**
    *   Implement support for a `.autocommitrc` file (YAML or TOML).
    *   Add initial configuration options: `auto_push`, `review_mode`, `learn_from_history`.

3.  **Expanded Installation Support:**
    *   Create a Homebrew formula for macOS/Linux installation.
    *   Create installation packages for Windows (e.g., Scoop, Winget, or a simple PowerShell script).
    *   Provide single binary downloads for all major platforms.

4.  **VS Code Extension (Wrapper):**
    *   Develop a basic VS Code extension that wraps the CLI tool.
    *   Provide a command palette option and a status bar icon to trigger `autocommit`.

---

## Phase 4: CI/CD & Advanced Automation

**Goal:** Make AutoCommit AI a reliable tool for automated environments and handle more complex repository setups.

1.  **CI/CD Mode (`--ci`):**
    *   Implement the `--ci` flag to enable non-interactive mode.
    *   Ensure deterministic output and disable all prompts.
    *   Provide structured, machine-readable logging (e.g., JSON).
    *   Implement strict exit codes to signal success, failure, or no changes.

2.  **Commit Guide Awareness:**
    *   Implement a feature to automatically detect and parse commit guidelines (e.g., `CONTRIBUTING.md`).
    *   Use simple regex to extract rules or examples.
    *   Allow users to specify a guide via URL in `.autocommitrc`.
    *   Adjust commit generation strategy based on the detected guidelines.

3.  **Advanced Safeguards:**
    *   Implement rollback on push failure (optional, configurable).
    *   Handle large/binary files by automatically creating a separate `chore(assets): add binary files` commit.

---

## Phase 5: Documentation & Release Readiness

**Goal:** Produce high-quality documentation and prepare for a public launch.

1.  **Comprehensive Documentation:**
    *   Write a "5-Minute Quick Start" guide.
    *   Create detailed installation guides for every supported platform.
    *   Develop a "Recipes" section for common use cases (e.g., CI/CD integration with GitHub Actions, GitLab CI).
    *   Document all CLI flags, configuration options, and safeguards.
    *   Create a comparison page: "AutoCommit AI vs. X".

2.  **Website/Landing Page:**
    *   Create a simple, clear landing page to market the tool and host the documentation.

3.  **Testing & Refinement:**
    *   Conduct extensive end-to-end testing across different repository types and edge cases.
    *   Refine the intelligence algorithms based on real-world usage.

4.  **Future Scope Planning:**
    *   Create issues/tickets for future features like PR creation and IDE-native integrations to guide future development.
