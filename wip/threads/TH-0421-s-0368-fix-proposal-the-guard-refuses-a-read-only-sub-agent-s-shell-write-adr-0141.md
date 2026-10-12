---
id: TH-0421
title: "S-0368 fix proposal: the guard refuses a read-only sub-agent's shell write (ADR-0141)"
anchor:
  path: wip/kanban/stories/S-0368-flai-guard-does-not-refuse-a-sub-agent-s-shell-commands-that-write-files-sed-i-perl-i-tee-redirection.md
  item: S-0368
status: resolved
participants: [agent-S-0368, orchestrator, alex]
created: 2026-10-12T17:21:40Z
updated: 2026-10-12T18:20:02Z
---

# TH-0421 S-0368 fix proposal: the guard refuses a read-only sub-agent's shell write (ADR-0141)

On wip/kanban/stories/S-0368-flai-guard-does-not-refuse-a-sub-agent-s-shell-commands-that-write-files-sed-i-perl-i-tee-redirection.md.

## Entries

### 2026-10-12T17:21:40Z agent-S-0368
**Recommendation: accept ADR-0141 as proposed.** I am building it now. If you answer differently, I change the ADR and the guard to match.

## The proposal

- **Who is refused.** A read-only sub-agent:
  - one whose `.claude/agents/<agent_type>.md` lists tools with neither `Edit` nor `Write`, which covers the explorer and the verifier
  - Claude Code's built-in `Explore` and `Plan`

  Task sub-agents and `general-purpose` keep writing through the shell. A definition the guard cannot read fails open, as ADR-0060 says.
- **What counts as a write:**
  - a redirection to a file: `>`, `>>`, `&>`, or `2>f`
  - `tee <file>`, `sed -i`, or `perl -i`
  - `cp`, `mv`, `rm`, `touch`, `mkdir`, `rmdir`, `truncate`, `ln`, `install`, `patch`, or `dd of=`
  - a formatter's write flag, such as `gofmt -w`, `prettier --write`, or `eslint --fix`
- **What passes:**
  - a target under the system temporary folder (`/tmp`)
  - `/dev/null` and any other `/dev/` path
  - a descriptor duplication such as `2>&1`

  A relative target is refused.
- **Limits.** The guard is still not a shell, so it cannot see a write made by an interpreter or a script (`python -c`, `make`). A file written that way still shows in `git status` before review.

## Why

- The verifier needs `Bash` for `flai verify --last` and `git diff`, so removing `Bash` is not an option.
- Task sub-agents use `sed -i` and `gofmt -w` as part of their work, so refusing every sub-agent would break them.

Source: design/adrs/0141-flai-guard-refuses-a-read-only-sub-agent-a-shell-command-that-writes-a-file.md

### 2026-10-12T17:58:53Z orchestrator (recommendation)
**Recommendation:** accept ADR-0141 as proposed.

- **Who is refused** matches the plan I approved on TH-0415. It enforces ADR-0059's rule that the explorer and the verifier only read, and it leaves task sub-agents their shell writes (ADR-0072).
- **What counts as a write, and what passes,** keeps the verifier able to work. `/tmp`, `/dev/`, and `2>&1` pass for test runs and captured output.
- **The stated limit stands as the backstop:** interpreter or script writes are not seen, and `git status` before review still catches them.

Accepting an ADR is a recorded decision, so it is yours to confirm. S-0368 is in review but lacks main; see TH-0426.

This is a recommendation, not an answer: confirm it, or answer otherwise.

Source: design/adrs/0059-a-story-s-agent-hands-search-test-runs-and-verification-to-an-explorer-and-a.md

### 2026-10-12T18:20:02Z alex
Resolved: S-0368 was accepted
