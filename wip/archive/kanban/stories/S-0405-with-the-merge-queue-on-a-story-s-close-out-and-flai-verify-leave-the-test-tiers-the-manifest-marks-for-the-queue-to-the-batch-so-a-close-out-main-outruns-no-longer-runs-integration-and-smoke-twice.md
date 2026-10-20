---
id: S-0405
type: story
nature: improvement
title: With the merge queue on, a story's close-out and flai verify leave the test tiers the manifest marks for the queue to the batch, so a close-out main outruns no longer runs integration and smoke twice
status: done
parent: E-0022
owner: alex
created: 2026-10-16T04:59:12Z
updated: 2026-10-19T21:52:13Z
transitions:
  - to: ready
    at: 2026-10-19T04:54:02Z
    by: alex
  - to: in-progress
    at: 2026-10-19T20:34:01Z
    by: agent-S-0405
  - to: review
    at: 2026-10-19T21:50:42Z
    by: agent-S-0405
  - to: done
    at: 2026-10-19T21:52:13Z
    by: orchestrator
tags: [cli]
topics: [conventions, template]
touches: [design/adrs, flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/verify/manifest.go, flai/internal/verify/select.go, flai/internal/verify/story.go, flai/internal/verify/story_test.go, flai/cmd/verify.go, flai/cmd/verify_test.go, flai/internal/mergequeue/batch.go, flai/internal/mergequeue/batch_test.go, scripts/close-out.sh, template/root/scripts/close-out.sh, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/template.yaml, template/CHANGELOG.md, system-flow.yaml, design/system/project-manifest.md, design/system/workflow.md, design/system/flai-cli.md, docs/operators/settings.md, docs/users/flai.md, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, flai/internal/verify/manifest_test.go, flai/internal/verify/verify.go, flai/cmd/mcp.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/verify.go, flai/internal/mcpserver/verify_test.go, flai/internal/mergequeue/eligible.go, flai/internal/mergequeue/eligible_test.go, flai/internal/verify/select_test.go, flaiover/src/lib/review.test.ts, flaiover/src/lib/review.ts, docs/operators/index.md, docs/users/flai-reference.md, flai/cmd/mcp_http.go, design/issues/I-0157-the-host-s-installed-flai-is-again-behind-the-project-s-releases-after-i-0149-closed-so-the-mcp-task-done-refuses-every-close.md, design/issues/summary.md]
after: [S-0399]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 4630
  priced_by: harness
  turns:
    - day: 2026-10-19
      ceremony: 1
      hand_edits: 2
      work: 75
  models:
    - model: claude-haiku-5-5
      input: 64
      output: 39414
      cache_read: 3512975
      cache_write: 195265
      cost: 0.3131
    - model: claude-opus-5-5
      input: 358
      output: 117485
      cache_read: 22453728
      cache_write: 929128
      cost: 12.1927
  strategic:
    - kind: orchestrator
      seconds: 2569
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 61
          output: 1142
          cache_read: 9654900
          cache_write: 23404
          cost: 2.4265
        - model: claude-sonnet-5-5
          input: 8
          output: 40
          cache_read: 151933
          cache_write: 44753
          cost: 0.161
forecast:
  duration: 41m
  delivery: 2026-10-19T21:19:00Z
  basis: "Its own forecast of 41m; 2nd in the pull order with an in-progress limit of 5, behind S-0404."
  by: flai
  at: 2026-10-19T20:33:45Z
finalized:
  by: alex
  at: 2026-10-19T04:53:47Z
---
# S-0405 With the merge queue on, a story's close-out and flai verify leave the test tiers the manifest marks for the queue to the batch, so a close-out main outruns no longer runs integration and smoke twice

## Goal

A story's close-out stops paying for the most expensive tiers when the queue will run them anyway on the batch. A test tier marked for the queue is left out of the close-out and of `flai verify` at review while the queue is on, and the batch runs it once for every story it holds, so that a close-out that main outruns, as in I-0148, no longer runs integration and smoke twice past the agent's tool limit.

## Acceptance criteria

- [x] A test tier in `system-flow.yaml` may carry `queue: true`, validated by the manifest's checks. This project's `integration` and `smoke` tiers carry it.
- [x] While the `merge-queue` host action is on, `flai verify` leaves the tiers marked `queue` out, names each in its step as left to the queue, and passes without them; with it off, it runs them as today. `--all` and `flai test --all` still run every tier.
- [x] The batch's verification always runs the tiers marked `queue` that the batch's paths select, so no story lands without them.
- [x] An ADR, refining ADR-0110 and ADR-0156, records it. The close-out rule in `work-management.md`, in the repository and the template, and both copies of `scripts/close-out.sh` say what the close-out leaves to the queue; `template/CHANGELOG.md` names the version of flai it needs.
- [x] `design/system/project-manifest.md`, `docs/operators/settings.md`, and `docs/users/flai.md` describe the tier key.

## Tasks
- T-1653 An ADR records that while the queue runs, the close-out and flai verify leave the tiers marked for the queue to the batch
- T-1654 A test tier in the manifest may carry queue: true, validated by the manifest's checks
- T-1665 While the merge-queue host action is on, flai verify leaves the tiers marked queue out and names each in its step
- T-1666 The batch's verification always runs the tiers marked queue that the batch's paths select
- T-1674 The close-out, the close-out rule in work-management.md, and this project's manifest leave the integration and smoke tiers to the queue
- T-1681 The manifest design, the workflow, the settings index, the CLI docs, and the template's changelog describe the queue tier key

## Notes

- Drafted by the planner for E-0022 as an addition, asked on TH-0465 question 5. The epic's outcome names I-0148; the queue alone does not remove it, since each close-out still runs the tiers.
- After S-0399, whose batch verification runs the tiers this story leaves out of the close-out.

### Planning

- **Touches.** None were declared; every touch is the planner's.
  - Layout, from the explorer's map:
    - the tiers: `TestTier` in `flai/internal/manifest/manifest.go` and `verify.Tier` in `flai/internal/verify/manifest.go`
    - selection: `flai/internal/verify/select.go`; the run and its reuse: `story.go`
    - the close-out script and its template copy at `template/root/scripts/close-out.sh`.
  - Design: ADR-0110, ADR-0156, and I-0148; the close-out rule in `work-management.md` and its template copy.
  - Co-change: `flai touches suggest` added `design/system/flai-cli.md`, `design/system/workflow.md` (the definition of done), and `template/template.yaml`.
  - `system-flow.yaml` is touched to mark the two tiers. CLAUDE.md gives flai the keys outside `name`, `description`, and `dashboard`, so T-1674 sets them through flai where a command can, and records a hand edit otherwise.
  - Folder touch kept: `design/adrs`, for the new ADR.
- **Forecast.** flai's 41m is kept: a tier key, a selection rule, and both copies of a script and a convention.
- **Cost of delay.** Waits on TH-0465 question 1.

### Accepted by the orchestrator

- Verified: 0d592a0e062b89f55abcace63e7be1834a3547b4
- At: 2026-10-19T21:52:13Z

Verdict: every criterion 1 to 5 is met by changed files and ticked. Queue tiers are left out only while the merge queue can batch; with the merge-queue host action off every tier runs as today. The batch always runs them. The baseline close-out rule is changed in both copies of work-management.md, with template 1.0.102 in its CHANGELOG. flai verify passed every tier at head 0d592a0, integration and smoke included.
- 1: flai/internal/manifest/manifest.go, flai/internal/manifest/settings.go, flai/internal/verify/manifest.go, flai/internal/verify/verify.go, system-flow.yaml, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings_test.go, flai/internal/verify/manifest_test.go
- 2: flai/internal/verify/select.go, flai/internal/verify/story.go, flai/internal/verify/select_test.go, flai/internal/verify/story_test.go, flai/cmd/verify.go, flai/cmd/verify_test.go, flai/cmd/mcp.go, flai/cmd/mcp_http.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/verify.go, flai/internal/mcpserver/verify_test.go, flai/internal/mergequeue/eligible.go, flai/internal/mergequeue/eligible_test.go, flaiover/src/lib/review.ts, flaiover/src/lib/review.test.ts
- 3: flai/internal/mergequeue/batch.go, flai/internal/mergequeue/batch_test.go
- 4: design/adrs/0174-while-the-merge-queue-can-batch-flai-verify-and-the-close-out-leave-the-test.md, design/adrs/README.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, scripts/close-out.sh, template/root/scripts/close-out.sh, template/CHANGELOG.md, template/template.yaml
- 5: design/system/project-manifest.md, docs/operators/settings.md, docs/users/flai.md, design/system/workflow.md, design/system/flai-cli.md, docs/users/flai-reference.md, docs/operators/index.md
