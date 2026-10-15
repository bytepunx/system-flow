---
id: S-0384
type: story
nature: feature
title: An agent's proposal for the design system is a story whose branch changes it, and the Design page previews that branch beside the main branch's before the operator accepts it
status: done
parent: E-0020
owner: alex
created: 2026-10-13T17:45:03Z
updated: 2026-10-15T04:19:34Z
transitions:
  - to: ready
    at: 2026-10-15T01:10:42Z
    by: alex
  - to: in-progress
    at: 2026-10-15T03:15:42Z
    by: agent-S-0384
  - to: review
    at: 2026-10-15T04:18:29Z
    by: system-flow
  - to: done
    at: 2026-10-15T04:19:34Z
    by: orchestrator
tags: [dashboard, cli]
topics: [design]
touches: [flai/internal/designsystem/proposals.go, flai/internal/designsystem/proposals_test.go, flai/cmd/design.go, flai/cmd/design_test.go, flai/internal/hostapi/design.go, flai/internal/hostapi/design_test.go, flaiover/src/lib/server/repo.ts, flaiover/src/lib/server/agent.ts, flaiover/src/routes/api/design/+server.ts, flaiover/src/routes/api/design/design.test.ts, flaiover/src/routes/api/design/proposals/+server.ts, flaiover/src/routes/api/design/proposals/proposals.test.ts, flaiover/src/lib/design/diff.ts, flaiover/src/lib/design/diff.test.ts, flaiover/src/lib/components/design/ProposalCompare.svelte, flaiover/src/lib/components/design/ProposalCompare.svelte.test.ts, flaiover/src/lib/components/design/ProposeChange.svelte, flaiover/src/lib/components/design/ProposeChange.svelte.test.ts, flaiover/src/routes/design/+page.svelte, flaiover/src/routes/design/design.svelte.test.ts, design/system/flai-cli.md, design/system/flaiover-dashboard.md, docs/users/flai.md, docs/users/flaiover.md, docs/users/flai-reference.md, flai/internal/designsystem/css.go, docs/operators/settings.md]
after: [S-0383]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 3791
  priced_by: harness
  turns:
    - day: 2026-10-15
      ceremony: 3
      hand_edits: 2
      work: 75
  models:
    - model: claude-opus-5-5
      input: 604
      output: 217580
      cache_read: 31223276
      cache_write: 947536
      cost: 16.0571
  strategic:
    - kind: orchestrator
      seconds: 3091
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 103
          output: 1780
          cache_read: 30568828
          cache_write: 35245
          cost: 7.6764
        - model: claude-sonnet-5-5
          input: 12
          output: 54
          cache_read: 191237
          cache_write: 63748
          cost: 0.2087
cost_of_delay:
  value: 250
  by: planner-E-0020
  at: 2026-10-13T17:52:06Z
forecast:
  duration: 1h5m
  delivery: 2026-10-15T04:14:00Z
  basis: "Its own forecast of 1h5m; 3rd in the pull order with an in-progress limit of 5, behind S-0385, S-0386 and S-0390."
  by: flai
  at: 2026-10-15T03:00:44Z
finalized:
  by: alex
  at: 2026-10-15T00:36:07Z
---
# S-0384 An agent's proposal for the design system is a story whose branch changes it, and the Design page previews that branch beside the main branch's before the operator accepts it

## Goal

The operator works with agents on the design system through the board's own workflow: they ask for a change from the Design page, which becomes a story; the story's agent changes the design system on its branch; and the Design page previews that branch's tokens on the sample components beside the main branch's, so the operator sees the proposal respond in the browser before accepting it, and accepting it adopts it.

## Acceptance criteria

- [x] `flai design show --story S-nnnn` and the host method `design.get` with a story read the design system from that story's worktree, and refuse a story that is not open or has no worktree.
- [x] `flai design proposals --json` and the host method `design.proposals` list the open stories whose branches change the design system document, each with its ID, title, and state.
- [x] The Design page lists those stories, and choosing one previews its tokens on the gallery beside the main branch's, with each token that differs marked and its two values shown.
- [x] From the Design page the operator writes a request for a change to the design system, which becomes a story under an epic they choose, with topic `design`, its goal from their words, and the design system document in its touches, for a story's agent to work like any other.
- [x] Once such a story is accepted, the Design page shows its tokens as the project's, with no step of its own.
- [x] `design/system/flai-cli.md`, `design/system/flaiover-dashboard.md`, `docs/users/flai.md`, `docs/users/flaiover.md`, and the generated reference say how.

## Tasks
- T-1548 flai/internal/designsystem reads the design system from an open story's worktree and lists the open stories whose branches change it
- T-1549 A comparison shows a proposal's tokens on the gallery beside the main branch's, with each token that differs marked and both values shown
- T-1550 A form on the Design page turns the operator's request for a design change into a story with topic design and the design system document in its touches
- T-1551 flai design show takes --story, and flai design proposals lists the open stories whose branches change the design system
- T-1552 The host method design.get takes a story, and design.proposals lists the open stories whose branches change the design system
- T-1553 The dashboard's /api/design takes a story, and /api/design/proposals answers design.proposals
- T-1554 The Design page lists the open proposals, compares the one chosen with the main branch's tokens, and offers the request form
- T-1555 The design and user docs and the reference describe proposals for the design system, their comparison, and the request form

## Notes

- Drafted by the planner for E-0020. It assumes S-0380 recommends proposals on story branches over a second store of drafts, so that review, acceptance, and history stay the board's. S-0378, in progress, gives `flai doc show` a `--story`; this story reads the design system the same way and should reuse what S-0378 builds.

### Planning

- **Touches.** None declared. All are the planner's, every one a file:
  - `flai/internal/designsystem/proposals.go` with its test: new, in S-0381's package. Criteria 1 and 2.
  - `flai/cmd/design.go` and `flai/internal/hostapi/design.go`, with their tests: S-0381's command group and host methods, extended. Criteria 1 and 2.
  - `flaiover/src/lib/server/repo.ts`, `agent.ts`, `flaiover/src/routes/api/design/+server.ts`, and the new `api/design/proposals/+server.ts`, with their tests: the reads. Criteria 1 to 3.
  - `flaiover/src/lib/design/diff.ts`, `ProposalCompare.svelte`, and `ProposeChange.svelte`, with tests: new. Criteria 3 and 4.
  - `flaiover/src/routes/design/+page.svelte` with its page test: the Design page. Criteria 3 to 5.
  - The five docs: design. Criterion 6.
- **Not touched.** The story request goes through the dashboard's existing story creation (`POST /api/items`, as `/new` does), so `flaiover/src/routes/api/items/+server.ts` is called, not changed. Criterion 5 needs no code: an accepted branch is merged into the main branch, which the page already shows.
- **Co-change.** `flai touches suggest` lists hub docs, `flai/internal/hostapi/writes.go`, and `flai/cmd/serve_actions.go`. None is added: the story adds reads only.
- **No folder touch.**
- **Overlaps.** S-0378, in progress, claims `repo.ts` and `agent.ts` and adds a story's worktree to `flai doc show`; this story waits behind S-0383, so it starts after S-0378 is likely accepted, and T-1548 reuses what S-0378 built.
- **Layers.** T-1548, T-1549, and T-1550 together; then T-1551 and T-1552 together; then T-1553; then T-1554; then T-1555.
- **Forecast: 1h5m, delivery 2026-10-14T00:24:00Z.** `flai forecast` gave 50m: 95 s per unit of size over 50 done feature stories on claude-opus-5-5 in the large band, times size 31. Raised about a quarter to 1h5m, as S-0378 was, because it spans flai and flaiover, runs both test tiers, and adds three Svelte components with tests, across eight tasks in five layers. The delivery adds the 15m to flai's 2026-10-14T00:09Z.
- **Cost of delay: 250.00 USD a week.** `flai cod`: 1h5m of the 5h12m forecast over E-0020's seven open stories, of the epic's 1200 USD a week. It stands.

### Accepted by the orchestrator

- Verified: c9396f086540c49edda8ec6a62a67f80f5a7f825
- At: 2026-10-15T04:19:34Z

Verdict: meets all six criteria; flai verify passed every step at c9396f08; criteria unchanged, with the two refinements confirmed on TH-0457 in the code; every changed path is within the story's touches, none under .claude/.

- 1: flai/internal/designsystem/proposals.go, flai/internal/designsystem/proposals_test.go, flai/cmd/design.go, flai/cmd/design_test.go, flai/internal/hostapi/design.go, flai/internal/hostapi/design_test.go, flai/internal/designsystem/css.go, flaiover/src/lib/server/repo.ts, flaiover/src/routes/api/design/+server.ts, flaiover/src/routes/api/design/design.test.ts
- 2: flai/internal/designsystem/proposals.go, flai/internal/designsystem/proposals_test.go, flai/cmd/design.go, flai/cmd/design_test.go, flai/internal/hostapi/design.go, flai/internal/hostapi/design_test.go, flaiover/src/lib/server/agent.ts, flaiover/src/lib/server/repo.ts, flaiover/src/routes/api/design/proposals/+server.ts, flaiover/src/routes/api/design/proposals/proposals.test.ts
- 3: flaiover/src/lib/design/diff.ts, flaiover/src/lib/design/diff.test.ts, flaiover/src/lib/components/design/ProposalCompare.svelte, flaiover/src/lib/components/design/ProposalCompare.svelte.test.ts, flaiover/src/routes/design/+page.svelte, flaiover/src/routes/design/design.svelte.test.ts
- 4: flaiover/src/lib/components/design/ProposeChange.svelte, flaiover/src/lib/components/design/ProposeChange.svelte.test.ts, flaiover/src/routes/design/+page.svelte, flaiover/src/routes/design/design.svelte.test.ts
- 5: flaiover/src/routes/design/+page.svelte, flaiover/src/routes/design/design.svelte.test.ts
- 6: design/system/flai-cli.md, design/system/flaiover-dashboard.md, docs/users/flai.md, docs/users/flaiover.md, docs/users/flai-reference.md, docs/operators/settings.md
