---
id: T-2223
type: task
nature: remediation
title: An ADR refining ADR-0060 and ADR-0203 records that the guard hook refuses a call when flai guard cannot start, and that inside the container the guard is the image's flai through the project's hook
status: backlog
parent: S-0518
owner: alex
created: 2026-10-28T13:32:54Z
updated: 2026-10-28T13:32:54Z
transitions: []
stream: S-0518
tags: [cli, template]
touches: [design/adrs/README.md]
---
# T-2223 An ADR refining ADR-0060 and ADR-0203 records that the guard hook refuses a call when flai guard cannot start, and that inside the container the guard is the image's flai through the project's hook

## Work

Write the ADR with `flai adr new`, refining ADR-0060 and ADR-0203, from I-0216 and S-0505's results. Context: the guard exits only 0 or 2, so a guard that cannot start is the hook's case alone, and the hook's shell line, `[ $? -eq 2 ] || exit 0`, allows a missing binary; inside `flai-agent-claude-code:local` this repository's hook runs `scripts/flai.sh guard`, whose `go build` exits 127, so every call passed. Decision, in two parts. First, the `PreToolUse` hook in `.claude/settings.json`, in the template and in this repository, refuses the call with exit 2 and `flai guard did not run (exit N): <output>` when the guard command exits 126 or 127, not found or not executable, and still allows on exit 0 and on any other exit: ADR-0060's rule that an older flai costs the guarantee, not the session, stands for exit 1, and a failed build of this tree keeps the session, since both hooks run the guard and a fail-closed on every exit would lock the agent out of `Bash` and `Edit` together. Second, inside the container the guard is the image's flai reached through the project's own hook: the template's hook names `flai` on `PATH`, and this repository's `scripts/flai.sh` runs the installed flai, with a warning on standard error, when `go` is absent and `bin/flai` is not built; pilot passes no hook of its own. Alternatives considered: fail closed on every non-zero exit (the lock-out above); a settings file of flai's passed to the inner harness with `--settings` (hooks from several settings sources add up, so the project's hook still runs, and excluding the project's settings drops its permissions); the project layer's `go` (S-0510; a build per hook call, and the template's shape stays open). Consequences: a host or image without `flai` on `PATH` now refuses every guarded call with the reason, which the agent reports, in place of running unguarded; this repository's container runs the image's flai, built from the tree at image build, until S-0510's layer gives it `go`. Give it the topics of ADR-0060, `cli`, `conventions`, and `template`, with `flai adr topics`, and accept it; `flai adr new` writes the index row in `design/adrs/README.md`.

## Done when

The ADR is under `design/adrs/` with status `accepted`, `refines` names ADR-0060 and ADR-0203, its row is in `design/adrs/README.md`, and `flai check --strict` is clean.

## Notes

The ADR's number is the next free one when the task runs: S-0502 (in progress) has taken 0206.
