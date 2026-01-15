# AutoCommit AI — Final Product & Technical Specification (Loop-Safe)

> **Purpose**  
AutoCommit AI is a deterministic, loop-safe CLI tool for solo developers and CI/CD systems that automatically creates **logical, reviewable Git commits** and optionally pushes them—without human micromanagement and without AI runaway loops.

---

## 1. Core Principles (Non-Negotiable)

1. **Single-Pass Execution**
   - Each run performs **exactly one analysis → one plan → one execution → exit**.
   - No re-analysis after execution begins.
   - No self-triggering reruns.

2. **No AI Feedback Loops**
   - AI never reads its own output.
   - AI never re-evaluates commits it generated.
   - AI suggestions are immutable once confirmed or auto-accepted.

3. **Determinism First**
   - Same inputs → same outputs (especially in CI).
   - All randomness disabled by default.

4. **Human Override Ends AI Authority**
   - If the user edits a commit message, AI stops contributing immediately for that commit.

5. **Fail Closed, Not Open**
   - On ambiguity, conflict, or unsafe Git state → abort with explanation.

---

## 2. Target Users

- **Primary**
  - Solo developers
  - Indie hackers
  - OSS maintainers
- **Secondary**
  - CI/CD pipelines
  - Automation bots
  - Monorepo maintainers

---

## 3. Supported Environments

- CLI-first (mandatory)
- CI/CD (non-interactive mode)
- Optional:
  - VS Code extension (wrapper only)
- Platforms:
  - macOS, Linux, Windows

---

## 4. Feature Set (Final)

### 4.1 Change Detection
- Reads **unstaged + staged** diffs once.
- Snapshot-based: diff is frozen at start.
- Respects `.gitignore`.

### 4.2 Logical Commit Grouping
- Groups files by:
  - Change intent (`feat`, `fix`, `refactor`, `docs`, `test`, `chore`)
  - Directory proximity
  - Dependency hints (imports, references)
- Output:
  - Ordered list of **commit plans**
  - Each plan = files + commit message

> **Important:**  
Grouping is computed once and cached for the run.  
No regrouping is allowed later.

### 4.3 Commit Message Generation
- Conventional Commits compliant.
- Inputs:
  - Diff snapshot
  - Repo history (read-only)
  - Optional commit guide (compiled once)
- Output:
  - One message per commit group
- Fallback:
  - `chore: automated commit`

### 4.4 Review Mode (`--review`)
- Shows:
  - Commit groups
  - Messages
- Allowed actions:
  - Accept all
  - Edit messages
  - Abort
- **Editing a message disables AI for that commit permanently.**

### 4.5 Execution Engine
- Sequential commit application:
  1. `git add <group files>`
  2. `git commit -m "<message>"`
- Push behavior:
  - Default: enabled
  - `--no-push`: commits only
- Push happens **once after all commits succeed**.

---

## 5. CI/CD Mode (`--ci`)

- Non-interactive
- No prompts
- JSON logs
- Strict exit codes:
  - `0` → success
  - `1` → error
  - `2` → no changes
- No learning
- No history mutation

---

## 6. Learning System (Bounded & Safe)

### What It Learns
- Common scopes
- Message phrasing patterns

### What It NEVER Does
- Rewrites history
- Influences current-run decisions
- Triggers re-analysis

### Storage
- Repo-local cache (`.autocommit_cache`)
- Write-once per run
- Read-only at startup

---

## 7. Commit Guide Support

- Auto-detect:
  - `CONTRIBUTING.md`
  - `commitlint.config.*`
- Optional manual link in config
- Guides are:
  - Parsed once
  - Compiled into static rules
  - Never reinterpreted mid-run

---

## 8. Safeguards & Abort Conditions

AutoCommit AI **aborts immediately** if:

- Detached HEAD
- Rebase or merge in progress
- Dirty index after snapshot
- Push conflicts
- Partial commit failure
- AI output violates commit rules

Rollback:
- Enabled by default for multi-commit failures.

---

## 9. Loop Prevention Rules (Critical)

| Risk | Mitigation |
|----|----|
AI re-analyzing commits | Snapshot diff, immutable plan |
Unstaged loop | Freeze diff at start |
`--no-push` reruns | Execution flag locks flow |
Multi-group commit loop | Commit plan index cursor (monotonic) |
Learning recursion | No same-run reads |
Human edit loop | AI disabled instantly |

**There is no code path where AI is called twice for the same run.**

---

## 10. Installation (Easy by Design)

- Homebrew (`brew install autocommit-ai`)
- Scoop / Winget (Windows)
- `pipx install autocommit-ai`
- Single static binaries
- GitHub Releases
- VS Code Marketplace (wrapper)

---

## 11. Tech Stack (Final)

### Core
- **Rust**
  - Speed
  - Safety
  - Single binary
- **libgit2**
  - Git operations
- **Tree-sitter**
  - Code structure hints

### AI Layer
- Pluggable LLM interface
- Strict input/output contracts
- Token-limited, no memory

### CLI
- `clap` (Rust)
- Structured logging

---

## 12. Data Flow (One Run)
Start
↓
Git State Validation
↓
Diff Snapshot (Frozen)
↓
Commit Guide Compilation
↓
Grouping + Message Generation (ONCE)
↓
[Optional Review]
↓
Sequential Commit Execution
↓
Optional Push
↓
Learning Write
↓
Exit

---

## 13. Documentation as a Feature

- 5-minute quick start
- CI recipes
- Failure explanations
- Comparison with:
  - Git hooks
  - Husky
  - Commitizen
  - AI copilots

---

## 14. Explicit Non-Goals

- No background daemon
- No auto-reruns
- No rewriting Git history
- No hidden AI decisions
- No cloud dependency required

---

## 15. Final Guarantee

> **AutoCommit AI will never get stuck in a loop, never commit twice unintentionally, and never surprise the user.**

That is a product guarantee, not an aspiration.
