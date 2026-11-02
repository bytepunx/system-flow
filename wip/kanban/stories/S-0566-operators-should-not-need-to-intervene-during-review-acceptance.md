---
id: S-0566
type: story
nature: improvement
title: Operators should not need to intervene during review acceptance
status: in-progress
owner: alex
created: 2026-11-02T18:21:16Z
updated: 2026-11-02T20:26:25Z
transitions:
  - to: ready
    at: 2026-11-02T18:21:38Z
    by: alex
  - to: in-progress
    at: 2026-11-02T18:26:18Z
    by: agent-S-0566
tags: [cli]
touches: [flai/cmd, design/adrs/0248-every-acceptance-commits-the-main-checkout-s-uncommitted-files-under-wip-design.md, design/adrs/README.md, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/internal/preview/accept.go, flai/internal/preview/accept_test.go, flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, flai/internal/storygit/sync.go, design/system/flai-cli.md, design/system/strategic-agents.md, design/system/workflow.md, docs/users/flai.md, docs/users/flaiover.md, flaiover/src/lib/components/AcceptConfirm.svelte, flaiover/src/lib/components/AcceptConfirm.svelte.test.ts, flaiover/src/lib/components/Review.svelte, flaiover/src/lib/components/Review.svelte.test.ts, "flaiover/src/routes/api/items/[id]/accept/+server.ts", "flaiover/src/routes/api/items/[id]/accept/accept.test.ts", "flaiover/src/routes/api/items/[id]/move/+server.ts", flaiover/src/routes/board/+page.svelte, "flaiover/src/routes/items/[id]/+page.svelte", docs/users/flai-reference.md, flai/internal/hostapi/writes.go, flai/internal/mergequeue/eligible.go, flai/internal/mergequeue/land.go, flai/internal/mergequeue/land_test.go, flai/internal/mergequeue/preflight.go, flai/internal/mergequeue/preflight_test.go, flai/internal/storygit/storygit.go, flai/internal/storygit/storygit_test.go, design/system/flaiover-dashboard.md]
agent:
  harness: claude-code
  model: claude-fable-5-1
  config:
    effort: high
usage:
  source: log
  seconds: 5149
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-02
      hand_edits: 1
      work: 34
  models:
    - model: claude-fable-5-1
      input: 6194
      output: 2407
      cache_read: 41603546
      cache_write: 1363673
      cost: 65.1532
cost_of_delay:
  inputs:
    penalty_per_week: 600
    by: alex
    at: 2026-11-02T18:21:16Z
---
# S-0566 Operators should not need to intervene during review acceptance

## Goal

Right now, nearly every story that lands in review fails to merge because it is no longer in sync with main. Rebase fails due to conflicts that flai is not resolving, forcing the operator to try and resolve merge conflicts in files they did not author.

flai exists to automate the development process. It should automatically resolve the rebase on main including any merge conflicts that arise from this operation.

## Acceptance criteria
- [x] An operator is never instructed to rebase the story’s worktree onto main since flai should do this
- [x] An operator is never told to resolve merge conflicts during rebase, the story agent should do this
- [x] files in shared folders like wip and design should always get checked in on review acceptance - there doesn’t need to be a checkbox/permission request since flai can do this directly

## Tasks
- T-2436 flai accept sends a story whose rebase it cannot settle back to its agent, whoever accepts
- T-2437 Every acceptance commits the main checkout's uncommitted files under wip, design, and the shared paths, and leaves the rest as they are
- T-2438 The dashboard's acceptance drops the checkbox for uncommitted files, shows the shared set it commits, and shows a send back as such
- T-2439 The design and the user guides say every acceptance commits the shared set, leaves the rest, and sends an unsettled rebase back to the agent

## Notes
