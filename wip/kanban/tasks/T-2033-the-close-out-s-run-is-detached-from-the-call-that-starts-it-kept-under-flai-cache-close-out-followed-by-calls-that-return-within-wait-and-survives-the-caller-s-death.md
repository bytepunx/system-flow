---
id: T-2033
type: task
nature: improvement
title: The close-out's run is detached from the call that starts it, kept under .flai-cache/close-out, followed by calls that return within --wait, and survives the caller's death
status: backlog
parent: S-0466
owner: alex
created: 2026-10-22T20:58:30Z
updated: 2026-10-22T21:01:13Z
transitions: []
stream: S-0466
tags: [cli]
touches: [flai/internal/closeout/detach.go, flai/internal/closeout/detach_test.go, flai/cmd/closeout.go, flai/cmd/closeout_test.go]
after: [T-2031]
---
# T-2033 The close-out's run is detached from the call that starts it, kept under .flai-cache/close-out, followed by calls that return within --wait, and survives the caller's death

## Work

Detach the run T-2031 added from the process that starts it, as T-2030's ADR says. `flai close-out <story> ...` starts the rounds in a child flai of its own process group (`execx`'s detach, as `flai checks run` and `flai serve`'s detached writes use it), with the child's stdout and stderr in a log under `.flai-cache/close-out/` and its state in `.flai-cache/close-out/<story>.json`: pid, started, round, step, outcome, ended, and the log path, written as each step starts and ends, in the shape of `serve.ChecksRun`. The caller follows the log, printing it as it comes, and returns with the run's exit status when the run ends, or after `--wait` (default `9m`) with exit 0 and a last line `close-out: <story> still running, round r of 3 at <step>, after <wait>; run flai close-out <story> again to keep following`. A later call while the state says running and the pid is alive attaches and follows, whatever its message; a call after the run ended prints the ended run's last line and, when it carries a message or file, starts a new run. `--status` prints the state and returns; `--cancel` terminates the run's process group, then kills it, as `flai checks cancel` does. A caller killed with SIGKILL, as Claude Code's Bash tool kills a call at its limit, leaves the run alive; SIGINT or SIGTERM to the caller leaves it alive too and prints how to follow or cancel it. Refuse a second concurrent run for one story.

Waits for T-2031, whose command and package this changes.

## Done when

- A test reproduces I-0148's cause: a close-out whose tiers outlast the caller. With a fixture tier that sleeps, the caller is started, killed with SIGKILL while the tier runs, and the run goes on: its state reaches `outcome: passed`, and `flai close-out <story>` called again prints the passing line.
- A test shows `--wait` passing: the caller returns at the wait with the still-running line and exit 0, and a second call attaches to the same run (one pid, one log) rather than starting another.
- Tests cover `--status` on a running and an ended run, `--cancel` ending a run with `outcome: cancelled`, and a new call with a message after an ended run starting a new run.
- `flai test` on the changed paths passes and `flai check --strict` is clean.

## Notes

The third layer, with T-2046 (the scripts); the two share no path. The default wait of 9m is a flag's default, not a setting, so `docs/operators/settings.md` gains no row; the ADR may choose otherwise.
