---
id: S-0845
type: story
nature: improvement
title: A story in review whose close-out fails after main moved waits for the operator to send it back, while its after chain holds the backlog and the board starves
status: backlog
owner: alex
created: 2026-11-17T08:32:56Z
updated: 2026-11-17T12:39:22Z
transitions: []
tags: [flai, orchestrator]
topics: [acceptance, verify, orchestration]
touches: [design/adrs, design/adrs/README.md, design/system/strategic-agents.md, design/system/workflow.md, flai/cmd/accept_resolve.go, flai/cmd/accept_resolve_test.go, flai/internal/mcpserver/resolve.go, flai/internal/mcpserver/resolve_test.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/folder_test.go, flai/internal/guard/guard.go, flai/internal/guard/guard_rules.go, flai/internal/guard/orchestrate_test.go, design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/orchestrator.md, flai/internal/serve/resolvejob_test.go, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, docs/users/flai.md, docs/users/flai-reference.md, template/CHANGELOG.md, design/issues/I-0386-a-story-in-review-whose-close-out-fails-after-main-moved-waits-for-the-operator.md, design/issues/summary.md]
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
      seconds: 65
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 14
          output: 34
          cache_read: 1708568
          cache_write: 14655
          cost: 0.1177
        - model: claude-opus-5-5
          input: 67
          output: 8638
          cache_read: 5803243
          cache_write: 81180
          cost: 2.1123
    - kind: orchestrator
      seconds: 18
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 4
          output: 48
          cache_read: 320670
          cache_write: 5680
          cost: 0.1022
    - kind: analyzer
      seconds: 99
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 47
          output: 6584
          cache_read: 349360
          cache_write: 34306
          cost: 1.1031
        - model: claude-haiku-5-5
          input: 36
          output: 14843
          cache_read: 2398931
          cache_write: 47318
          cost: 0.1557
cost_of_delay:
  inputs:
    time_lost_per_cycle: 30h
    by: flai
    at: 2026-11-17T08:32:56Z
  value: 4500
  by: planner-S-0845
  at: 2026-11-17T12:34:56Z
forecast:
  duration: 3h
  delivery: 2026-11-18T02:24:00Z
  basis: "Its own forecast of 3h; 22nd in the pull order with an in-progress limit of 3, behind S-0806, S-0834, S-0838, S-0841, S-0811, S-0829, S-0840, S-0809, S-0842, S-0812, S-0808, S-0807, S-0810, S-0241, S-0813, S-0814, S-0815, S-0822, S-0824, S-0835, S-0837, S-0839, S-0843 and S-0844."
  by: flai
  at: 2026-11-17T12:39:22Z
finalized:
  by: orchestrator
  at: 2026-11-17T12:38:55Z
---
# S-0845 A story in review whose close-out fails after main moved waits for the operator to send it back, while its after chain holds the backlog and the board starves

## Goal

This story remediates [I-0386](../../../design/issues/I-0386-a-story-in-review-whose-close-out-fails-after-main-moved-waits-for-the-operator.md), "A story in review whose close-out fails after main moved waits for the operator to send it back, while its after chain holds the backlog and the board starves", which the analysis in [2026-11-16-all.md](../../../design/analysis/2026-11-16-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0386 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0386 is closed with `flai issue close I-0386 --reason` saying what fixed it

## Tasks
- T-3587 Propose in an ADR who asks for rebase and resolve when a review story's close-out stops after main moved, and write it into the design
- T-3588 flai accept --resolve takes the orchestrator's request under accept_reviews, and an MCP tool resolve asks for it on the host
- T-3589 flai guard and the orchestrator's role file let it ask for rebase and resolve under accept_reviews when a close-out did not pass after main moved
- T-3590 A test reproduces I-0386: a review story whose close-out stopped after main moved is sent back by the resolve job the orchestrator asked for
- T-3591 The strategic-agents convention, the users guide, the reference, and the changelog say the orchestrator asks for rebase and resolve
- T-3592 Close I-0386 with flai issue close, naming the ADR and the change that fixed it

## Notes

The analyzer found I-0386: the analysis the goal links gives its evidence.

Cost of delay inputs set by flai from I-0386. time_lost_per_cycle 30h carried over from I-0386's Impact section. Its Impact time_lost_per_cycle was taken rather than the 30h derived from its cost and count.

Strategic usage carried over from I-0386: analyzer 1m39s, 1.2588 USD.

### Planning

The fix the plan assumes, for the story's agent to propose in an ADR before code: ADR-0313's rebase and resolve job already rebases a review story's branch, closes it out at the rebased head, and sends the story back through ADR-0248's send-back when the rebase stops or a step fails. What I-0386 lacks is a way for anyone but the operator to ask for it: `flai accept --resolve` refuses the orchestrator (`flai/cmd/accept_resolve.go`, `resolveFlagsRefused`), no MCP tool asks for it, and a plain close-out job never sends back (`flai/internal/serve/verifyjob.go`, `startJob` sets `settle` only for `KindResolve`). So the orchestrator, under `accept_reviews`, asks for rebase and resolve, through a new MCP tool, when a story's close-out did not pass after main moved, in place of a thread awaiting the operator. The story's agent may choose otherwise in its ADR, such as a close-out job of a review story that sends back on a stop; the touches then shift from `flai/cmd` and `flai/internal/mcpserver` to `flai/internal/serve/verifyjob.go`.

Touches, and where each came from:

- `design/adrs` (folder) and `design/adrs/README.md`: design. The new ADR's file name comes from its title, which `flai adr new` slugs and no task can name yet, so the folder stays on the story; the index is generated from the ADR files (ADR-0228).
- `design/system/strategic-agents.md`, `design/system/workflow.md`: design. § Accepting a story (S-0221) and § Transitions and who makes them, § What waits on an acceptance (S-0296) say who sends a review story back.
- `flai/cmd/accept_resolve.go`, `flai/cmd/accept_resolve_test.go`: layout; the refusal of the orchestrator's `--resolve`.
- `flai/internal/mcpserver/resolve.go`, `resolve_test.go`, `folder.go`, `folder_test.go`: layout; a new tool beside `accept.go` and `verify.go`, registered in `folder.go`.
- `flai/internal/guard/guard.go`, `guard_rules.go`, `orchestrate_test.go`: layout; the orchestrator's `accept_reviews` calls (`queue_add`, `accept`).
- `design/workflows/roles/orchestrator.md` and its template copy: layout; its tools list and the `accept_reviews` permission's tools, and the paragraph that leaves a story whose verify did not pass in review with a thread.
- `flai/internal/serve/resolvejob_test.go`: layout; the test that reproduces I-0386, a resolve job the orchestrator asked for that stops at a step and sends the story back.
- `design/conventions/strategic-agents.md` and its template copy, `docs/users/flai.md`, `docs/users/flai-reference.md`, `template/CHANGELOG.md`: design and the documentation convention.
- `design/issues/I-0386-...md`, `design/issues/summary.md`: the second criterion, written by `flai issue close`.
- `flai touches suggest S-0845` listed nothing: the story declared no touches.

No `after`: the fix reads ADR-0313's job, which S-0804 delivered and is accepted. S-0824 (backlog) changes `resolvejob.go` and `verifyjob.go`; this story only adds a test to `resolvejob_test.go`, and a trial merge judges any overlap.

Figures:

- Forecast 3h, delivery 2026-11-18T00:30:00Z, adjusted. `flai forecast S-0845` gave 46m with the touches (2m before, with none): 108 s per unit over 186 large improvement stories, times size 25. That median undercounts this story: its first task waits for the operator to accept an ADR before any code, its six tasks run in five layers, and each close-out runs about ten minutes (S-0828's took 9m43s today). The delivery moves from flai's 22:15 by the same 2h15m.
- Cost of delay 4500 USD a week, as `flai cod S-0845` gave it: the operator's input, 30h lost per 168h cycle, at 150 USD an hour. It stands: 30h is the longest of the three review waits I-0386 measured, and the stalled review lane also starved ready on 48 `ready-empty` events.
