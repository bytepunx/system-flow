---
id: S-0395
type: story
nature: feature
title: Under pull-request acceptance, flai accept and the orchestrator's acceptance refuse a story while a required check has not passed or a required review is outstanding, naming each
status: review
parent: E-0021
owner: alex
created: 2026-10-16T04:37:54Z
updated: 2026-10-17T18:22:20Z
transitions:
  - to: ready
    at: 2026-10-17T06:43:53Z
    by: alex
  - to: in-progress
    at: 2026-10-17T17:33:51Z
    by: agent-S-0395
  - to: review
    at: 2026-10-17T18:22:20Z
    by: agent-S-0395
tags: [cli, dashboard]
topics: [release, orchestration]
touches: [flai/internal/preview/accept.go, flai/internal/preview/accept_forge_test.go, flai/cmd/accept_forge_test.go, design/adrs, design/system/workflow.md, design/system/strategic-agents.md, design/system/flai-cli.md, docs/users/flai.md, flai/internal/forge/forge.go, flai/internal/forge/github.go, flai/internal/forge/github_test.go, flai/internal/forge/state.go, flai/internal/forge/state_test.go, flai/internal/serve/forge_test.go, flai/internal/serve/forgestatus.go, flai/internal/preview/orchestrator.go, flaiover/src/lib/components/Review.svelte, flaiover/src/lib/components/Review.svelte.test.ts, design/conventions/strategic-agents.md, design/system/flaiover-dashboard.md, template/CHANGELOG.md, template/root/design/conventions/strategic-agents.md, design/issues/I-0162-the-components-dev-server-test-picks-a-free-port-and-closes-it-before-the-fake-vite-binds-it-so-a-concurrent-test-run-on-the-host-can-take-the-port-and-fail-the-integration-tier.md, design/issues/summary.md]
after: [S-0394]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 2929
  priced_by: harness
  turns:
    - day: 2026-10-17
      ceremony: 3
      hand_edits: 2
      work: 79
  models:
    - model: claude-opus-5-5
      input: 436
      output: 137112
      cache_read: 24130169
      cache_write: 744033
      cost: 12.016
  strategic:
    - kind: orchestrator
      seconds: 3015
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 156
          output: 2653
          cache_read: 62580401
          cache_write: 45826
          cost: 15.6997
        - model: claude-sonnet-5-5
          input: 26
          output: 135
          cache_read: 530306
          cache_write: 56525
          cost: 0.4803
forecast:
  duration: 28m
  delivery: 2026-10-17T18:04:00Z
  basis: "Its own forecast of 28m; 1st in the pull order with an in-progress limit of 5, with nothing ahead of it."
  by: flai
  at: 2026-10-17T17:33:42Z
finalized:
  by: alex
  at: 2026-10-17T06:42:38Z
---
# S-0395 Under pull-request acceptance, flai accept and the orchestrator's acceptance refuse a story while a required check has not passed or a required review is outstanding, naming each

## Goal

Under pull-request acceptance, CI and the team's reviewers are the gate, whoever accepts. `flai accept` previews and refuses a story whose pull request has a required check that has not passed, a requested change still standing, or fewer approvals than the main branch requires. It names each one, for the operator, the dashboard, and the orchestrator alike. A project that accepts locally sees no change.

## Acceptance criteria

- [x] Under `acceptance.mode: pull-request`, the preview of `flai accept S-nnnn` lists as blockers, each by name: every required check that is failing or still pending, every review requesting changes that still stands, too few approving reviews for the main branch's rule, and a story with no pull request. The real run refuses on them before anything is merged, as it does on any blocker. Under `local`, the preview is as today.
- [x] The gate reads the state S-0394 keeps. A state recorded for a commit other than the story branch's head is itself a blocker, saying that flai serve has not yet read the head, so a stale pass never lets a story through.
- [x] The orchestrator's acceptance, `flai accept --by orchestrator` under ADR-0093, obeys the same blockers. The orchestrator never approves a review itself, as TH-0463 answers question 3.
- [x] The dashboard's acceptance confirmation on the review page lists the same blockers, from the preview's `--json`.
- [x] An ADR, refining ADR-0093, records the gate. `design/system/workflow.md` (review to done), `design/system/strategic-agents.md` (Accepting a story), `design/system/flai-cli.md`, and `docs/users/flai.md` describe it.

## Tasks
- T-1623 The acceptance preview lists each required check not passed, each standing requested change, missing approvals, a missing pull request, and a stale state as blockers under pull-request acceptance
- T-1624 An ADR records that under pull-request acceptance every acceptance, the orchestrator's included, waits for the required checks and reviews
- T-1625 flai accept refuses on the forge blockers before merging, for the operator, the dashboard's acceptance, and the orchestrator's alike
- T-1626 The dashboard's acceptance confirmation lists the forge blockers by name
- T-1627 The workflow, the strategic agents' design, and the CLI docs describe the forge gate on acceptance
- T-1895 The forge contract reports the approvals the base branch requires, and flai serve keeps them with the checks and reviews

## Notes

- Drafted by the planner for E-0021, after S-0394, whose kept state this story reads. It needs no token: the gate reads what flai serve kept, so the orchestrator and an operator's shell judge the same state.
- The merge through the remote is S-0397, after this story. Until it lands, a story that passes this gate is still merged locally.

### Planning

- **Touches.** None were declared. Every touch is the planner's, and each is a file except `design/adrs`:
  - **Layout:**
    - `flai/internal/preview/accept.go`, with a new test, `accept_forge_test.go`: `AcceptWith` and its blockers. Criteria 1 and 2.
    - `flai/cmd/accept.go`, with a new test, `accept_forge_test.go`, and `flai/internal/hostapi/writes.go`, with its test: the refusals, `accept.run`, and `orchestratorAccepts`. Criteria 1 and 3.
    - `flaiover/src/routes/review/[id]/+page.svelte`: the confirmation. Criterion 4.
  - **Design:** `workflow.md`, `strategic-agents.md`, `flai-cli.md`, and `docs/users/flai.md`, for criterion 5.
- **Folder touch kept.** `design/adrs`, a shared path, for an ADR `flai adr new` names.
- **Co-change left out.**
  - `docs/users/flai-reference.md`: no new flag.
  - `docs/operators/*` and `flaiover/src/lib/server/agent.ts`: no setting and no host method.
- **Tag added.** `dashboard`, second, for the review page.
- **Paths it shares with other stories.** E-0022's S-0400 changes `flai/cmd/accept.go` and `flai/internal/preview/accept.go`, and S-0399 changes `writes.go`. They hold each other while either is in progress (TH-0464).
- **Layers.**
  - First: T-1623 and T-1624 together.
  - Second: T-1625 and T-1626, each after T-1623.
  - Third: T-1627, after T-1624, T-1625, and T-1626.
- **Forecast: 28m.** `flai forecast` gave 28m: 96 s per unit over 61 done feature stories on claude-opus-5-5 in the large band, times size 17. Kept: the gate reads state S-0394 keeps, and the dashboard part is one page.
- **Cost of delay: not yet set.** `flai cod` refuses until E-0021 has inputs, asked for on TH-0463.
