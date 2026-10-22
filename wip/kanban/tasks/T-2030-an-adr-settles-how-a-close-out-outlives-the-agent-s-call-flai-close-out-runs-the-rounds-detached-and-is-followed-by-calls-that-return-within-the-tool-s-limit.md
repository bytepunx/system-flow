---
id: T-2030
type: task
nature: improvement
title: "An ADR settles how a close-out outlives the agent's call: flai close-out runs the rounds detached and is followed by calls that return within the tool's limit"
status: backlog
parent: S-0466
owner: alex
created: 2026-10-22T20:57:35Z
updated: 2026-10-22T20:57:35Z
transitions: []
stream: S-0466
tags: [cli, template]
touches: [design/adrs, design/adrs/README.md]
---
# T-2030 An ADR settles how a close-out outlives the agent's call: flai close-out runs the rounds detached and is followed by calls that return within the tool's limit

## Work

Propose the fix for I-0148 from its instance, as the story's goal asks, and record it as an ADR with `flai adr new`, then `flai adr accept` once the operator agrees on the plan thread or the ADR's thread. Read I-0148, [ADR-0110](../../../design/adrs/0110-a-story-s-agent-runs-the-close-out-which-runs-flai-verify-itself-before-review.md), [ADR-0135](../../../design/adrs/0135-the-sync-step-of-flai-verify-and-the-close-out-s-last-check-pass-over-commits.md), [ADR-0156](../../../design/adrs/0156-a-close-out-that-main-outran-syncs-and-verifies-again-at-most-three-rounds-and.md), and [ADR-0174](../../../design/adrs/0174-while-the-merge-queue-can-batch-flai-verify-and-the-close-out-leave-the-test.md) first: ADR-0174 expected I-0148 to close once the merge queue ran the `queue` tiers, and the operator kept the issue open, so the fix has to hold with the queue off.

The planner's recommendation, to confirm or amend in the ADR:

1. The close-out moves into flai as `flai close-out <story> [-m <message> | -F <file>] [--wait <duration>] [--status] [--cancel] [--json]`, running the rounds `scripts/close-out.sh` runs today (`flai verify --record-issues`, the commit, `flai verify --sync-only`, `flai stream sync`, at most three rounds, the clean worktree check), with the same last line.
2. The run is detached from the call that starts it, as `flai checks run` detaches its steps and as the host write `verify.run` is detached: the first call starts the run in a process group of its own, with its state (pid, started, round, step, outcome, log path) under `.flai-cache/close-out/<story>.json` and its output in a log there, follows the log, and returns when the run ends or when `--wait` passes, default 9m, under Claude Code's 10-minute Bash limit, with a last line saying the run is still going and to call `flai close-out <story>` again. A later call while the run is alive attaches and follows it instead of starting another; a call after it ended prints its last line, and a new call with a message starts a new run. Killing the caller, which the Bash tool does at its limit, leaves the run alive, so a close-out that main outruns finishes its second round on its own.
3. `scripts/close-out.sh`, here and in the template, calls `flai close-out`, and the conventions tell the agent to call it again until its last line names the outcome, instead of one 600000 ms call.

Alternatives to weigh in the ADR: the agent runs the close-out with the Bash tool's `run_in_background` and a 2-hour timeout (harness-specific, and a headless session ends ten minutes after its turn, so the run may still die with it); the close-out stops after round 1 and asks the agent to run it again (the second round still needs its own 12 minutes); the close-out passes over main's commits (rejected by ADR-0156, decision 3). Say which the ADR refines (ADR-0110, ADR-0156) and what the template ships.

## Done when

- The ADR is under `design/adrs/`, listed in `design/adrs/README.md`, with `status: accepted` once the operator has agreed on its thread or the plan thread, and `flai check --strict` is clean.
- The ADR names the command, its flags, where the run's state and log live, what the caller does when `--wait` passes, what a second call does while a run is alive, and what the scripts and conventions change.

## Notes

The first layer of the plan: the code tasks follow its decision. The `design/adrs/` folder touch stands for the new ADR, whose number is not known until `flai adr new` allocates it.
