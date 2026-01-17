# AutoCommit AI - TODO List

This document outlines the remaining tasks and areas for improvement for the AutoCommit AI project, based on the `ACTION_PLAN.md` and current implementation status.

## 1. Feature Set Enhancements

### 1.1 Logical Commit Grouping (ACTION_PLAN 4.2)
- [x] **Ordered Commit Plans:** Converted commit grouping output from `map` to `[]CommitPlan` and enforced deterministic ordering (lexical sort by GroupKey).
- [x] **Improve "Directory Proximity" Grouping:** Enhanced `classify.ClassifyAndGroupChanges` to use the immediate parent directory as a scope when no learned scope is available, ensuring rule-based grouping without recursion.
- [ ] **Implement "Dependency Hints" for Grouping:** (Skipped due to complexity growth as per execution rules) Explore parsing code to identify dependencies (imports, references) and use this information to create more logical commit groups.

### 1.2 CI/CD Mode (ACTION_PLAN 5)
- [x] **Prevent Interactive Prompts:** Implemented logic to bypass interactive review mode when in CI mode.
- [x] **Strict Exit Codes:** Centralized exit handling to enforce `0` for success, `1` for errors, and `2` for no changes.
- [x] **CI Mode Hard Locks:** Implemented explicit guards to prevent learning data writes and history mutation in CI mode.
- [ ] **JSON Logs:** Implement structured JSON logging for CI mode to facilitate machine readability.

### 1.3 Commit Guide Support (ACTION_PLAN 7)
- [x] **Auto-detect Commit Guides:** Implemented logic to automatically detect `CONTRIBUTING.md` and `commitlint.config.*` files.
- [x] **Parse and Apply Rules:** Developed functionality to parse commit guides using regex-based extraction and compile them into static rules (specifically, a commit message regex).

### 1.4 Other Features
- [x] **Review & Edit Mode:** Optional interactive mode to review and edit proposed commits before finalization.
- [x] **Verbose Mode:** Enable verbose output for debugging purposes.

## 2. Safeguards & Abort Conditions (ACTION_PLAN 8)



- [x] **Detect Detached HEAD:** (Partially implemented via `git.CheckGitStatus()`) Ensure comprehensive checks for a detached HEAD state.

- [x] **Detect Rebase / Merge in Progress:** Implemented checks for `.git/rebase-apply`, `.git/rebase-merge`, and `.git/MERGE_HEAD` to abort if a rebase or merge is in progress.

- [x] **Detect Dirty Index After Snapshot:** Hashed index state at snapshot time and verified it before each commit execution, aborting on mismatch.

- [ ] **Handle Push Conflicts:** Implement robust error handling and abort conditions for Git push conflicts.

- [x] **Handle Partial Commit Failures:** Implicitly handled by existing error handling; the application aborts immediately on any commit failure and does not retry or continue with subsequent commits.

- [x] **AI Output Validation:** Implemented validation of AI-generated commit messages against extracted regex rules, with a fallback to a safe conventional commit if validation fails.

- [x] **Rollback Mechanism:** Implemented a deferred rollback to `git reset --hard <original HEAD>` for multi-commit runs, ensuring it does not activate in CI mode.

## 3. Installation & Distribution (ACTION_PLAN 10)

- [ ] **Package for Homebrew:** Create and maintain a Homebrew formula for macOS and Linux users.
- [ ] **Package for Scoop/Winget:** Develop packaging for Windows users via Scoop or Winget.
- [ ] **Package for pipx:** Create a pipx package for Python users (if applicable for a Go application, or clarify alternative).
- [x] **Single Static Binaries:** Provided a `build.sh` script to build and install the application.
- [ ] **GitHub Releases Automation:** Automate the release process on GitHub.

## 4. Documentation (ACTION_PLAN 13)

- [x] **5-minute Quick Start Guide:** Updated the `README.md` with installation and usage instructions.
- [ ] **CI Recipes:** Provide examples and configurations for integrating AutoCommit AI into common CI/CD pipelines.
- [ ] **Failure Explanations:** Document common failure modes and their explanations/resolutions.
- [ ] **Comparison with Alternatives:** Document how AutoCommit AI compares to other tools like Git hooks, Husky, Commitizen, and AI copilots.

---

## Completed Features (Not originally in TODO.md)

### Improved Rule-based Commit Message Generation (ACTION_PLAN 4.3)
- [x] Enhanced the rule-based commit message generation by leveraging learned subject lines from commit history, making messages more context-aware and aligned with historical patterns.
- [x] Added interactive mode selection to make the application easier to use.
- [x] Added support for `.env` files to manage environment variables.
- [x] Implemented a retry mechanism with exponential backoff for API calls.
- [x] Handled `EOF` error in the interactive prompt.
- [x] **Default to AI:** The application now defaults to using the AI to generate commit messages.
