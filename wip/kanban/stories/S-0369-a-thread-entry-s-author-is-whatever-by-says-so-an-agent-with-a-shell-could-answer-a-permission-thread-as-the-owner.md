---
id: S-0369
type: story
nature: improvement
title: A thread entry's author is whatever --by says, so an agent with a shell could answer a permission thread as the owner
status: in-progress
owner: alex
created: 2026-10-12T16:52:42Z
updated: 2026-10-12T18:18:00Z
transitions:
  - to: ready
    at: 2026-10-12T18:00:46Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-12T18:13:51Z
    by: agent-S-0369
tags: [flai, guard, threads]
topics: [cli, conventions]
touches: [design/adrs/README.md, design/system/agent-adapters.md, design/system/strategic-agents.md, flai/internal/guard/guard.go, flai/internal/guard/guard_test.go, flai/cmd/guard.go, flai/cmd/guard_test.go, flai/cmd/thread.go, flai/cmd/thread_test.go, docs/users/flai-reference.md, docs/users/flai.md, design/issues/I-0091-a-thread-entry-s-author-is-whatever-by-says-so-an-agent-with-a-shell-could-answer-a-permission-thread-as-the-owner.md, design/issues/summary.md, flai/internal/guard/authors.go, flai/internal/guard/authors_test.go]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: planner
      seconds: 496
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 112
          output: 387
          cache_read: 6633400
          cache_write: 179052
          cost: 0.3698
        - model: claude-opus-5-5
          input: 40
          output: 414
          cache_read: 1835220
          cache_write: 109856
          cost: 0.8791
    - kind: orchestrator
      seconds: 533
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 47
          output: 782
          cache_read: 13418258
          cache_write: 20896
          cost: 3.388
cost_of_delay:
  inputs:
    penalty_per_week: 10
    by: orchestrator
    at: 2026-10-12T16:59:25Z
  value: 10
  by: planner-S-0369
  at: 2026-10-12T17:06:49Z
forecast:
  duration: 45m
  delivery: 2026-10-12T18:51:00Z
  basis: "Its own forecast of 45m; 5th in the pull order with an in-progress limit of 5, behind S-0363, S-0370, S-0373 and S-0351."
  by: flai
  at: 2026-10-12T18:00:51Z
finalized:
  by: orchestrator
  at: 2026-10-12T18:00:29Z
---
# S-0369 A thread entry's author is whatever --by says, so an agent with a shell could answer a permission thread as the owner

## Goal

This story remediates [I-0091](../../../design/issues/I-0091-a-thread-entry-s-author-is-whatever-by-says-so-an-agent-with-a-shell-could-answer-a-permission-thread-as-the-owner.md), "A thread entry's author is whatever --by says, so an agent with a shell could answer a permission thread as the owner". The issue recommends this solution:

Directions to weigh: flai refuses a `--by` that is not the session's own name when `FLAI_AGENT` is set, so that an agent's entries always carry its name; and `flai guard` refuses a story's agent `flai thread reply`, `resolve`, and `confirm` with a `--by` at all. Either leaves the operator's own session, where `FLAI_AGENT` is unset, free to name itself. The dashboard already writes the operator's name itself.

## Acceptance criteria
- [ ] The cause I-0091 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0091 is closed with `flai issue close I-0091 --reason` saying what fixed it

## Tasks
- T-1465 An ADR refining ADR-0060 and ADR-0097 decides that flai guard refuses a served agent a flai command that names who wrote it, and agent-adapters.md and strategic-agents.md record it
- T-1466 flai thread commands take the author from the config's author when neither --by nor FLAI_AGENT is given, as the docs say
- T-1467 flai guard refuses a story's agent, the planner, and the analyzer a flai command that names its author with --by or FLAI_AGENT, with tests that reproduce I-0091
- T-1468 docs/users/flai.md says who a thread entry's author is and that flai guard refuses a served agent --by and FLAI_AGENT on flai
- T-1469 I-0091 is closed with a reason naming the ADR and the guard's refusal of --by and FLAI_AGENT to a served agent

## Notes

### Planning

The story declared no touches. `flai touches suggest` ran from `flai/internal/guard/guard.go`, `flai/cmd/thread.go`, and `flai/cmd/guard.go`, where I-0091 names the cause and where the guard is.

| Touch | Source | Why |
|-------|--------|-----|
| `flai/internal/guard/guard.go` | layout | `Decide` lets every `Bash` call of a story's own agent pass; the planner's rules let `thread new` and `reply` pass with any `--by` |
| `flai/internal/guard/guard_test.go` | co-change, 71% | `TestTheStorysAgentIsNeverRefused` changes to the rule's exceptions |
| `flai/internal/guard/thread_test.go` | layout | The guard's thread cases, where I-0091 is reproduced |
| `flai/cmd/guard.go` | layout | Its help says what the guard refuses; it builds the guard from `FLAI_ROLE` and `FLAI_STORY` |
| `flai/cmd/guard_test.go` | co-change, 40% | The hook-JSON table |
| `docs/users/flai-reference.md` | co-change, 43% | Generated from the help by `make flai-reference` |
| `flai/cmd/thread.go`, `flai/cmd/thread_test.go` | layout | `threadAuthor` takes `--by` as given, and never reaches the config's author |
| `docs/users/flai.md` | co-change, 14%; design | `### Threads` says the config's author is used |
| `design/system/agent-adapters.md` | design | `### Permissions` and `### Asking the owner before a protected write` |
| `design/system/strategic-agents.md` | design | The planner's `### The guard` and `### The analyzer's guard` |
| `design/adrs/README.md` | co-change, 9% | Index of the new ADR |
| `design/issues/I-0091-…md`, `design/issues/summary.md` | goal | Criterion 2 closes the issue |

No folder touch is kept. The new ADR's file cannot be named before `flai adr new` numbers it; `flai task done` widens T-1465's touches with it.

Left out on purpose:

- `flai/internal/threads/threads.go` and `flai/internal/mcpserver/permission.go`: they read an entry's author and stay as they are once no served agent can forge one.
- `flai/internal/hostapi/writes.go`: the dashboard runs flai with `FLAI_AGENT=flaiover` and `--by=<owner>`, outside any guard, and keeps doing so.
- `.claude/settings.json`: the hook already runs before every `Bash` call.
- `docs/operators/settings.md` (co-change, 14%): changed only if the guard comes to read `FLAI_AGENT`; T-1468 says so.

S-0368 touches `guard.go`, `guard_test.go`, `flai/cmd/guard.go`, `guard_test.go`, `flai-reference.md`, `flai.md`, `agent-adapters.md`, and `design/adrs/README.md` too. Whichever starts second is held until the other reaches review.

Figures:

- **Forecast: 45m, delivery 2026-10-13T02:21:00Z.** `flai forecast` gave 21m: 78 s a unit over 61 large-band improvement stories, times size 16. Raised as S-0368's was, because this story adds an ADR and the guard's first rule on a story's own agent's `Bash` calls. It also changes the planner's and analyzer's rules and fixes `threadAuthor`. The delivery is flai's 01:57Z plus the extra 24m.
- **Cost of delay: 10 USD a week**, as `flai cod` works it out from the orchestrator's `penalty_per_week: 10` (TH-0416). Kept as is: I-0091 has one instance, found in review, and no measured cost.
