---
id: ADR-0280
title: "scripts/flai.sh guard runs an older bin/flai, or else the installed flai, when the build of this tree fails, and exits 127 when neither is there, so this repository's guard hooks never pass a call for a failed build"
status: proposed
date: 2026-11-10
supersedes: []
superseded_by: []
refines: [ADR-0235]
topics: [cli, harness]
---

# ADR-0280 scripts/flai.sh guard runs an older bin/flai, or else the installed flai, when the build of this tree fails, and exits 127 when neither is there, so this repository's guard hooks never pass a call for a failed build

## Context

ADR-0235 made this repository's guard hooks, in `.claude/settings.json`, refuse a call when `scripts/flai.sh guard` exits 126 or 127, and allow every exit but 2, 126, and 127. It kept a failed build of this tree failing open: "a failed build of this tree (exit 1 from `go build`) … still fail open, so a session is never locked out of `Bash` and `Edit` together for a cause its agent cannot fix". Its consequences add that where `go` is present `scripts/flai.sh` never runs an older installed flai in place of the tree's.

I-0275 records what that leaves open. S-0567's live planner ran inside the project's layer image, which has `go`. There `scripts/flai.sh guard` builds `bin/flai` whenever the sources are newer than it. When that `go build` fails, `scripts/flai-build.sh` exits 1 under `set -eu`, `scripts/flai.sh` exits 1 with it, and the hook lets the call through, unguarded. The same holds on a developer's host. A branch that does not build is common while flai is being changed, and every call of every session in that worktree then passes the guard, sub-agents' included.

This repository's OpenCode plugin, `.opencode/plugins/flai-guard.js`, runs `scripts/flai.sh` too. It treats only a spawn error as a missing guard, so it passes an exit 127, and a failed build, as any status but 2.

## Decision

`scripts/flai.sh guard` always reaches a guard. When `go` is present and the build of `bin/flai` fails, it runs, in order:

| What is there | `scripts/flai.sh guard` |
|---------------|-------------------------|
| `bin/flai`, older than the sources | runs it, with a warning on standard error that the build failed |
| no `bin/flai`, a `flai` on `PATH` | runs the installed flai, with the same warning naming it |
| neither | exits 127, with the build's failure on standard error, so the hooks refuse |

`scripts/flai-build.sh` builds to a temporary file and renames it into place, so a failed build leaves the last good `bin/flai` as it was.

Every other command keeps failing on a failed build: `scripts/flai.sh check`, `test`, `task done`, and the rest exit with the build's status, so the agent working on flai still sees its build fail.

The hooks' text in `.claude/settings.json` stays as ADR-0235 has it. The OpenCode plugin judges the guard's exit as the hooks do: 2 refuses with the guard's output, 126 or 127 is a missing guard, which refuses in a session flai started and passes, said once, otherwise, as a spawn error does today; any other status passes.

This reverses, for `guard` alone, ADR-0235's rule that a failed build of this tree fails open, and its consequence that `scripts/flai.sh` never runs an older installed flai where `go` is present.

## Consequences

- A guarded session in this repository never passes a call because the tree does not build. It is guarded by the last flai that built, or by the installed flai, which may lack the guard's newest rules; ADR-0235 already accepts an older flai as the guard where `go` is absent.
- With no `bin/flai` and no installed flai, a failed build refuses every guarded call, with the build's output as the reason. That is a fresh worktree on a branch that does not build, on a host with no flai installed; fixing the build, or installing flai, clears it.
- The template is unchanged: its hooks and plugin run `flai` from `PATH`, which builds nothing.

## Alternatives considered

- Make the hooks refuse every exit but 0 and 2: an older flai without `guard`, or input it cannot read, would lock a session out of `Bash` and `Edit`, which ADR-0235 rejected.
- Exit 127 on every failed build, with no fallback: a branch that does not build would lock its own agent out of the edits that fix it.
- Fall back on a failed build for every command: the agent working on flai would run an older flai without seeing its build fail.
