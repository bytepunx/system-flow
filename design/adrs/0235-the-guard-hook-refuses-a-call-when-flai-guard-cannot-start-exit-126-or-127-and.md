---
id: ADR-0235
title: "The guard hook refuses a call when flai guard cannot start, exit 126 or 127, and inside the agent container the guard is the image's flai reached through the project's own hook"
status: accepted
date: 2026-11-01
supersedes: []
superseded_by: []
refines: [ADR-0060, ADR-0203]
topics: [cli, conventions, template, harness]
---

# ADR-0235 The guard hook refuses a call when flai guard cannot start, exit 126 or 127, and inside the agent container the guard is the image's flai reached through the project's own hook

## Context

ADR-0060 made `flai guard` a Claude Code `PreToolUse` hook that fails open, "so a missing or older `flai` costs the guarantee, not the session". The guard exits only 0, to allow, or 2, to refuse; a guard that cannot start at all is the hook's case alone, outside the binary. Both hooks judge it with the same shell line, `out=$(<guard> 2>&1); [ $? -eq 2 ] || exit 0`, which allows every exit but 2, and so allows a guard that was never found.

I-0216 records what that cost. In S-0505's spike (ADR-0203), S-0513 was worked inside `flai-agent-claude-code:local`. This repository's hook runs `"$CLAUDE_PROJECT_DIR"/scripts/flai.sh guard`, which runs `go build` when `bin/flai` is missing; a fresh worktree has no `bin/` and the image has no `go`, so the build exited 127 and every call of the session and its sub-agents passed unguarded. The template's hook runs `flai guard` from `PATH`, which the image carries, so a project made from the template was guarded there; but a host or image without `flai` on `PATH` is unguarded the same way.

## Decision

The guard hook refuses a call when `flai guard` cannot start, exit 126 or 127, and inside the agent container the guard is the image's flai reached through the project's own hook.

The `PreToolUse` hooks, in the template's `.claude/settings.json` and in this repository's, keep the guard command's exit code and judge it:

| Exit | The hook |
|------|----------|
| 0 | allows the call |
| 2 | refuses it with the guard's output, as before |
| 126 or 127 | refuses it with `flai guard did not run (exit N): <output>`, saying no guarded call passes until flai is installed and on `PATH` |
| any other | allows it, as ADR-0060 has it |

Only "the command is not there or cannot be executed" fails closed. An older flai without `guard`, a failed build of this tree (exit 1 from `go build`), and input the guard cannot read still fail open, so a session is never locked out of `Bash` and `Edit` together for a cause its agent cannot fix.

Inside the agent container the guard is the image's flai, reached through the project's own hook. The template's hook names `flai` on `PATH` already. This repository's `scripts/flai.sh` runs the installed flai, with a warning on standard error, when `go` is absent and `bin/flai` is not built, and exits 127 when neither is there, so the hook refuses. pilot passes no hook of its own.

## Consequences

- A host or an image without `flai` on `PATH` refuses every guarded call with the reason, which the agent reports, in place of running unguarded.
- This repository's container runs the image's flai, built from the tree when the image was built, until the project layer gives it `go` (S-0510), when `scripts/flai.sh` builds again as on the host.
- On a developer's host, where `go` is present, `scripts/flai.sh` never runs an older installed flai in place of the tree's.

## Alternatives considered

- Fail closed on every non-zero exit: a broken build of this tree, or an older flai without `guard`, would lock every story's agent out of `Bash` and `Edit` at once.
- A settings file of flai's passed to the inner harness with `--settings`: Claude Code adds hooks from several settings sources together, so the project's hook would still run beside it, and excluding the project's settings drops its permissions too.
- Wait for the project layer's `go` (S-0510): it leaves the template's fail-open shape as it is, and builds flai on the first hook call of every session.
