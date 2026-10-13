---
id: S-0384
type: story
nature: feature
title: An agent's proposal for the design system is a story whose branch changes it, and the Design page previews that branch beside the main branch's before the operator accepts it
status: backlog
parent: E-0020
owner: alex
created: 2026-10-13T17:45:03Z
updated: 2026-10-13T17:45:03Z
transitions: []
tags: [dashboard, cli]
topics: [design]
touches: [flai/internal/designsystem/proposals.go, flai/internal/designsystem/proposals_test.go, flai/cmd/design.go, flai/cmd/design_test.go, flai/internal/hostapi/design.go, flai/internal/hostapi/design_test.go, flaiover/src/lib/server/repo.ts, flaiover/src/lib/server/agent.ts, flaiover/src/routes/api/design/+server.ts, flaiover/src/routes/api/design/design.test.ts, flaiover/src/routes/api/design/proposals/+server.ts, flaiover/src/routes/api/design/proposals/proposals.test.ts, flaiover/src/lib/design/diff.ts, flaiover/src/lib/design/diff.test.ts, flaiover/src/lib/components/design/ProposalCompare.svelte, flaiover/src/lib/components/design/ProposalCompare.svelte.test.ts, flaiover/src/lib/components/design/ProposeChange.svelte, flaiover/src/lib/components/design/ProposeChange.svelte.test.ts, flaiover/src/routes/design/+page.svelte, flaiover/src/routes/design/design.svelte.test.ts, design/system/flai-cli.md, design/system/flaiover-dashboard.md, docs/users/flai.md, docs/users/flaiover.md, docs/users/flai-reference.md]
after: [S-0383]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0384 An agent's proposal for the design system is a story whose branch changes it, and the Design page previews that branch beside the main branch's before the operator accepts it

## Goal

The operator works with agents on the design system through the board's own workflow: they ask for a change from the Design page, which becomes a story; the story's agent changes the design system on its branch; and the Design page previews that branch's tokens on the sample components beside the main branch's, so the operator sees the proposal respond in the browser before accepting it, and accepting it adopts it.

## Acceptance criteria

- [ ] `flai design show --story S-nnnn` and the host method `design.get` with a story read the design system from that story's worktree, and refuse a story that is not open or has no worktree.
- [ ] `flai design proposals --json` and the host method `design.proposals` list the open stories whose branches change the design system document, each with its ID, title, and state.
- [ ] The Design page lists those stories, and choosing one previews its tokens on the gallery beside the main branch's, with each token that differs marked and its two values shown.
- [ ] From the Design page the operator writes a request for a change to the design system, which becomes a story under an epic they choose, with topic `design`, its goal from their words, and the design system document in its touches, for a story's agent to work like any other.
- [ ] Once such a story is accepted, the Design page shows its tokens as the project's, with no step of its own.
- [ ] `design/system/flai-cli.md`, `design/system/flaiover-dashboard.md`, `docs/users/flai.md`, `docs/users/flaiover.md`, and the generated reference say how.

## Tasks

## Notes

- Drafted by the planner for E-0020. It assumes S-0380 recommends proposals on story branches over a second store of drafts, so that review, acceptance, and history stay the board's. S-0378, in progress, gives `flai doc show` a `--story`; this story reads the design system the same way and should reuse what S-0378 builds.
