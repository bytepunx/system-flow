---
id: S-0368
type: story
nature: improvement
title: flai guard does not refuse a sub-agent's shell commands that write files (sed -i, perl -i, tee, redirection)
status: review
owner: alex
created: 2026-10-12T16:52:41Z
updated: 2026-10-12T17:42:29Z
transitions:
  - to: ready
    at: 2026-10-12T17:11:33Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-12T17:18:30Z
    by: agent-S-0368
  - to: review
    at: 2026-10-12T17:42:29Z
    by: agent-S-0368
tags: [flai, guard]
topics: [cli, conventions]
touches: [design/adrs/README.md, design/system/agent-adapters.md, flai/internal/guard/guard.go, flai/internal/guard/guard_test.go, flai/internal/guard/writes.go, flai/internal/guard/writes_test.go, flai/cmd/guard.go, flai/cmd/guard_test.go, docs/users/flai-reference.md, docs/users/flai.md, design/issues/I-0052-flai-guard-does-not-refuse-a-sub-agent-s-shell-commands-that-write-files-sed-i-perl-i-tee-redirection.md, design/issues/summary.md, design/adrs/0141-flai-guard-refuses-a-read-only-sub-agent-a-shell-command-that-writes-a-file.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1456
  turns:
    - day: 2026-10-12
      ceremony: 1
      hand_edits: 1
      work: 32
  models:
    - model: claude-opus-5-5
      input: 182
      output: 92637
      cache_read: 10056126
      cache_write: 448586
      cost: 6.7709
  strategic:
    - kind: orchestrator
      seconds: 307
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 73
          output: 1289
          cache_read: 20086498
          cache_write: 31242
          cost: 5.0715
cost_of_delay:
  inputs:
    penalty_per_week: 10
    by: orchestrator
    at: 2026-10-12T16:55:07Z
  value: 10
  by: planner-S-0368
  at: 2026-10-12T16:56:42Z
forecast:
  duration: 45m
  delivery: 2026-10-13T02:06:00Z
  basis: "flai forecast's 19m (78 s a unit over 61 large-band improvement stories, size 14) raised to 45m: S-0246, the last guard change to its shell handling, took 31m without an ADR, and this story adds an ADR, a write detector, and a read-only role check; delivery is flai's 01:40Z moved by the extra 26m."
  by: planner-S-0368
  at: 2026-10-12T16:56:42Z
finalized:
  by: orchestrator
  at: 2026-10-12T16:58:50Z
---
# S-0368 flai guard does not refuse a sub-agent's shell commands that write files (sed -i, perl -i, tee, redirection)

## Goal

This story remediates [I-0052](../../../design/issues/I-0052-flai-guard-does-not-refuse-a-sub-agent-s-shell-commands-that-write-files-sed-i-perl-i-tee-redirection.md), "flai guard does not refuse a sub-agent's shell commands that write files (sed -i, perl -i, tee, redirection)". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0052 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0052 is closed with `flai issue close I-0052 --reason` saying what fixed it

## Tasks
- T-1460 An ADR refining ADR-0060 proposes which sub-agents flai guard refuses a shell write and which commands count as one, and agent-adapters.md records it
- T-1461 flai guard refuses a read-only sub-agent's shell command that writes a file, with tests that reproduce I-0052
- T-1462 docs/users/flai.md says which sub-agents flai guard refuses a shell write, what counts as one, and what passes
- T-1463 I-0052 is closed with a reason naming the ADR and the guard's new refusal

## Notes

I-0052's cost "0m" is not a duration longer than zero, so no time_lost_per_cycle was derived from it.

### Planning

The story declared no touches, so `flai touches suggest` ran from `flai/internal/guard/guard.go` and `flai/cmd/guard.go`, where the guard is.

| Touch | Source | Why |
|-------|--------|-----|
| `flai/internal/guard/guard.go` | layout | `Decide` and `refuse` judge a sub-agent's `Bash` call; no rule there looks for a file write |
| `flai/internal/guard/guard_test.go` | co-change, 78% | Holds the verifier's table of Bash cases, where I-0052 is reproduced |
| `flai/internal/guard/writes.go`, `writes_test.go` | layout | New files predicted for the write detector and the read-only role check |
| `flai/cmd/guard.go` | layout | Its help says what the guard refuses |
| `flai/cmd/guard_test.go` | co-change, 44% | The hook-JSON table |
| `docs/users/flai-reference.md` | co-change, 44% | Generated from the help by `make flai-reference` |
| `docs/users/flai.md` | design | `## Sub-agents` says the guard is not a shell and lets file writes pass |
| `design/system/agent-adapters.md` | design | `### Sub-agents and strategic agents` and `### Refusing a sub-agent's write` |
| `design/adrs/README.md` | co-change, 9% | Index of the new ADR refining ADR-0060 |
| `design/issues/I-0052-…md`, `design/issues/summary.md` | goal | Criterion 2 closes the issue |

No folder touch is kept. The new ADR's file cannot be named before `flai adr new` gives it a number; `flai task done` widens T-1460's touches with it.

Left out on purpose:

- `template/root/design/conventions/delegation.md`: it already tells the explorer and the verifier not to edit files, and the guard now enforces that.
- `.claude/settings.json`: the hook already runs before every `Bash` call.
- `flai/go.mod`: no shell parser is added, since the guard stays not a shell (ADR-0060).

Figures:

- **Forecast: 45m, delivery 2026-10-13T02:06:00Z.** `flai forecast` gave 19m: 78 s a unit over 61 large-band improvement stories, times size 14. Raised because S-0246, the last change to the guard's shell handling, took 31m with no ADR. This story adds an ADR, a write detector, and a read-only role check. The delivery is flai's 01:40Z plus the extra 26m.
- **Cost of delay: 10 USD a week**, as `flai cod` works it out from the orchestrator's `penalty_per_week: 10` (TH-0414). Kept as is: I-0052 has one instance and no measured cost.
