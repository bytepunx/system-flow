---
id: S-0537
type: story
nature: improvement
title: flai prime --role orchestrate takes the run's reasons and primes for them, the role's conventions once and per reason what it needs, instead of the general pack
status: ready
parent: E-0030
owner: alex
created: 2026-10-28T21:02:39Z
updated: 2026-10-29T04:22:23Z
transitions:
  - to: ready
    at: 2026-10-28T22:21:29Z
    by: alex
tags: [flai, prime, orchestrator]
topics: [orchestration, cli, conventions, template, mcp]
touches: [flai/internal/context, flai/internal/mcpserver, flai/cmd, design/system/flai-cli.md, design/system/strategic-agents.md, design/adrs, docs, flai/internal/context/reasons.go, flai/internal/context/reasons_test.go, flai/internal/context/reasons_figures.go, flai/internal/context/reasons_figures_test.go, flai/internal/context/strategic.go, flai/internal/context/strategic_test.go, flai/internal/context/context.go, flai/cmd/prime.go, flai/cmd/prime_test.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/prime_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/orchestrator.md, ".claude/agents/orchestrator.md", template/root/.claude/agents/orchestrator.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, design/system/agent-context.md, docs/users/flai.md, docs/users/flai-reference.md, template/CHANGELOG.md]
after: [S-0536]
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
      seconds: 1252
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 62
          output: 150
          cache_read: 3426229
          cache_write: 163428
          cost: 0.2354
        - model: claude-opus-5-5
          input: 66
          output: 4098
          cache_read: 3567001
          cache_write: 151451
          cost: 1.1061
    - kind: orchestrator
      seconds: 1254
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 45
          output: 726
          cache_read: 7744210
          cache_write: 20084
          cost: 1.7085
cost_of_delay:
  inputs:
    time_lost_per_cycle: 1h
    by: orchestrator
    at: 2026-10-28T21:49:26Z
  value: 150
  by: planner-S-0537
  at: 2026-10-28T22:08:55Z
forecast:
  duration: 4h
  delivery: 2026-10-29T09:46:00Z
  basis: "Its own forecast of 4h; 1st in the pull order with an in-progress limit of 5, behind S-0417 and S-0539."
  by: flai
  at: 2026-10-29T04:22:23Z
---
# S-0537 flai prime --role orchestrate takes the run's reasons and primes for them, the role's conventions once and per reason what it needs, instead of the general pack

## Goal

ADR-0209 §3. The orchestrator's general pack is three to four parts of some 35k characters each, loaded at every start and again after every compaction, whatever the run is for. With a run started for named reasons, the pack can be primed for them. `flai prime --role orchestrate --reason <reason>...` and MCP `prime` with `reasons` load the role's conventions once and, per reason, what it needs: for a thread, the thread, its item, and the conventions on threads; for a story in review, the story, its last `verify` result, and the acceptance rules; for a plan, the planner's thread and the items it wrote; for a release, `release_evaluate`'s figures; for a draft, the draft and the finalizing rules; for the ready column, the board and `promote_candidates`. The general pack is not loaded. flai serve's start prompt names the reasons in the form `prime` takes, and the digest's answers ranked against the reasons' threads (ADR-0194) stay within the budget's eighth.

## Acceptance criteria
- [ ] `flai prime --role orchestrate --reason` and MCP `prime` with `reasons` accept each reason the trigger names, load what the goal lists for it, and refuse an unknown reason naming the known ones, with a test per reason
- [ ] A pack primed for one thread is under a quarter of the general pack's size, measured in the test
- [ ] The orchestrator's start prompt names its reasons in the form `prime` takes, and the role's prompt tells it to prime with them
- [ ] `flai-cli.md`, `strategic-agents.md`, and the agent docs describe the reason-primed pack
- [ ] The orchestrator in this repository runs a day on reason-primed packs, and the narrative records the context tokens per action from `flai stats` beside the previous story's figure

## Tasks
- T-2308 The context package primes the orchestrator for named reasons: it parses them, refuses an unknown one naming the known ones, and loads the role's conventions once and, for a thread, a draft, and a plan, their items and threads
- T-2309 A pack primed for a review, a release, or the ready column loads the story and its last verify result, release_evaluate's figures, or the board and promote_candidates
- T-2310 The strategic-agents and work-management conventions, here and in the template, mark the sections each orchestrator reason needs with section topics, so a reason pack takes out the rest and every other pack keeps what it kept
- T-2311 flai prime --role orchestrate takes --reason, repeatable, and prints the reason-primed pack, refusing it with another role or with --epic or --story
- T-2312 The MCP prime tool takes reasons with role orchestrate and returns the reason-primed pack in parts, refusing an unknown reason naming the known ones
- T-2313 The orchestrator's start prompt names its reasons in the form prime takes, and its role prompt tells it to prime with them
- T-2314 flai-cli.md, strategic-agents.md, agent-context.md, the strategic-agents convention and its template copy, and the user docs describe the reason-primed pack
- T-2315 The orchestrator in this repository runs a day on reason-primed packs, and the narrative records the context tokens per action from flai stats beside S-0536's figure

## Notes

A reason the trigger gains later needs its pack rule here; the refusal naming the known reasons makes a missing one visible at once.

### Planning

Planned by planner-S-0537 on 2026-10-28. The tasks, their layers, and the assumptions are on the plan's thread.

The general orchestrate pack is built by `ForStrategic` in `flai/internal/context/strategic.go`, its digest step by `orchestratorDigest` and `answersFor` there; `flai prime` is `flai/cmd/prime.go`; the MCP tool's input is `PrimeIn` in `flai/internal/mcpserver/server.go`, registered in `folder.go`; the start prompt is `orchestratePrompt` in `flai/internal/harness/harness.go`. On 2026-10-28 `flai prime --role orchestrate --json` here was 170 KB, of which the role's seven conventions were about 64 KB: over a quarter of the pack on their own, so criterion 2 needs the conventions cut by reason (T-2310), not only the briefs left out.

Touches, and where each came from:

- Folders kept, all declared: `flai/internal/context`, `flai/internal/mcpserver`, `flai/cmd`, and `docs`. The tasks name every file below them the plan expects, so the story's claim narrows each folder to those files (ADR-0096); they stay for a file a reason loader may need that no task can name yet, such as a function moved out of `context` to break an import cycle with `verify` or `release` (T-2309).
- `design/adrs`: declared, and kept. T-2310 may need an ADR refining ADR-0049 and ADR-0075 if cutting convention sections by reason changes how a pack selects sections; its number is not known yet.
- `design/system/flai-cli.md`, `design/system/strategic-agents.md`: declared; criterion 4.
- `flai/internal/context/reasons.go`, `reasons_test.go`, `reasons_figures.go`, `reasons_figures_test.go`, `strategic.go`, `strategic_test.go`, `context.go`: layout; criteria 1 and 2.
- `flai/cmd/prime.go`, `prime_test.go`, `flai/internal/mcpserver/server.go`, `folder.go`, `prime_test.go`: layout; criterion 1.
- `flai/internal/harness/harness.go`, `harness_test.go`, `design/workflows/roles/orchestrator.md` and its template copy, `.claude/agents/orchestrator.md` and its template copy: layout; criterion 3, the role prompt's step 1 and the definitions generated from it.
- `design/conventions/strategic-agents.md`, `work-management.md`, and their template copies, `template/CHANGELOG.md`: co-change (work-management) and layout; criterion 2 on this repository and criterion 4's agent docs.
- `design/system/agent-context.md`, `docs/users/flai.md`, `docs/users/flai-reference.md`: co-change and design; § Strategic packs and § Prime a session describe the pack.
- Left out of the suggestion: `flaiover-dashboard.md`, `workflow.md`, `project-manifest.md`, `metrics.md`, `agent-narrative.md`, `design/issues/summary.md`, the host API, `check.go`, `guard.go`, `manifest.go`, `serve/agents.go`, and `flaiover`. They co-change with past orchestrator stories, but nothing in them changes here: no manifest key, guard rule, host write, or page is added, and the reasons reach the prompt through S-0536's queue.

Forecast: `flai forecast` gave 18m (median 88 s per unit of size over 21 done improvement stories on claude-opus-5-5 in the medium band, size 12: 5 criteria, 7 touches). Raised to 4h: eight tasks in five layers, with six reason loaders, a convention split here and in the template that must leave every other pack's rules in place, the CLI flag and the MCP field, two prompts, and five documents, beside S-0536's 5h for eight. flai's delivery, 2026-10-29T22:42Z, does not count S-0536, which this story waits for (`after`); delivery is S-0536's 2026-10-30T21:27Z plus the 4h and the day criterion 5 runs: 2026-11-01T01:27Z.

Cost of delay: `flai cod` gave 150.00 USD a week from the input the orchestrator set on TH-0718 (1h of time lost per 168h cycle at 150 USD an hour, 1.00 cycles a week). The value stands as computed; the input is not changed.
