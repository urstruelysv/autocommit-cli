<!-- You are Gemini, acting as a **bounded implementation assistant**, not an autonomous agent.

Your task is to **complete all remaining TODOs** in the AutoCommit AI project safely, deterministically, and without introducing reasoning loops, architectural drift, or hidden behavior.

This document is your **only source of truth**.

---

## ABSOLUTE RULES (DO NOT VIOLATE)

1. **ONE TASK PER RUN**
   - Select exactly **one unchecked TODO**
   - Complete it fully
   - Stop immediately after

2. **NO SELF-REANALYSIS**
   - Once a design decision is made for a task, do not revisit or optimize it.
   - Do not refactor unrelated code.

3. **NO RUNTIME FEEDBACK LOOPS**
   - Do not read logs, outputs, or runtime state to adjust behavior.
   - Logic must be static and deterministic.

4. **NO CROSS-TASK COUPLING**
   - Touch the minimum number of files.
   - No “while here” changes.

5. **ABORT > GUESS**
   - If information is missing, fail loudly and leave a TODO comment.
   - Never invent behavior.

---

## EXECUTION ORDER (MANDATORY)

Follow this order strictly. Do not skip ahead.

### PHASE A — DETERMINISM FIRST

1. Ordered Commit Plans  
   - Convert commit grouping output from `map` → `[]CommitPlan`
   - Enforce deterministic ordering (lexical path sort)
   - No intelligence changes

2. Strict Exit Codes  
   - Centralize exit handling
   - Enforce:
     - `0` success
     - `1` error
     - `2` no changes

3. CI Mode Hard Locks  
   - Enforce:
     - No prompts
     - No learning writes
     - No history mutation
   - Use explicit guards, not scattered conditionals

STOP AFTER PHASE A

---

### PHASE B — SAFEGUARDS & ABORTS

4. Detect Rebase / Merge in Progress  
   - Check for:
     - `.git/rebase-apply`
     - `.git/rebase-merge`
     - `.git/MERGE_HEAD`
   - Abort immediately

5. Dirty Index After Snapshot  
   - Hash index state at snapshot time
   - Verify before commit execution
   - Abort on mismatch

6. Partial Commit Failure Detection  
   - Detect mid-sequence failure
   - Do NOT retry

7. Rollback Mechanism  
   - Only for multi-commit runs
   - Use `git reset --hard <original HEAD>`
   - Never rollback in CI mode

STOP AFTER PHASE B

---

### PHASE C — INTELLIGENCE (BOUNDED)

8. Directory Proximity Grouping  
   - Rule-based only
   - Same parent directory → same group
   - No learning
   - No recursion

9. Dependency Hints (OPTIONAL)  
   - Static heuristics only:
     - import/include paths
   - Max depth = 1
   - Abort if cycles detected

STOP IF COMPLEXITY GROWS

---

### PHASE D — COMMIT GUIDE SUPPORT

10. Commit Guide Detection  
    - Detect:
      - `CONTRIBUTING.md`
      - `commitlint.config.*`

11. Commit Guide Parsing  
    - Regex-based extraction only
    - Compile once into static rules
    - No runtime reinterpretation

12. AI Output Validation  
    - Validate commit messages
    - Reject invalid ones
    - Fallback to safe conventional commit

---

### PHASE E — CI & OBSERVABILITY

13. JSON Logs (CI Mode Only)  
    - Fixed schema
    - No human-readable output

14. Disable Learning in CI  
    - Hard-disable cache writes at IO boundary

---

### PHASE F — PACKAGING & DOCS

15. Static Binary Builds  
16. Homebrew Formula  
17. Scoop / Winget  
18. GitHub Releases Automation  

Docs Tasks (one per run):
- 5-minute Quick Start
- CI Recipes
- Failure Explanations
- Tool Comparisons

---

## HARD NOT-TO-DOS

- ❌ No background processes
- ❌ No auto-reruns
- ❌ No AI memory across runs
- ❌ No re-analysis after execution starts
- ❌ No modifying Git history beyond current HEAD
- ❌ No silent fallbacks

---

## FINAL GUARANTEE

There must be **no code path** where:
- AI is called twice in one run
- A commit plan is regenerated
- Execution restarts automatically

If a rule is violated → ABORT the task.

---

### BEGIN NOW
Pick the **next unchecked TODO** and implement **only that**. -->
