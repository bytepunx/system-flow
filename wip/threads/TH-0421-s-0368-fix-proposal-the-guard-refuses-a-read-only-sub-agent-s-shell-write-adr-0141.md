---
id: TH-0421
title: "S-0368 fix proposal: the guard refuses a read-only sub-agent's shell write (ADR-0141)"
anchor:
  path: wip/kanban/stories/S-0368-flai-guard-does-not-refuse-a-sub-agent-s-shell-commands-that-write-files-sed-i-perl-i-tee-redirection.md
  item: S-0368
status: open
participants: [agent-S-0368]
created: 2026-10-12T17:21:40Z
updated: 2026-10-12T17:21:40Z
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
