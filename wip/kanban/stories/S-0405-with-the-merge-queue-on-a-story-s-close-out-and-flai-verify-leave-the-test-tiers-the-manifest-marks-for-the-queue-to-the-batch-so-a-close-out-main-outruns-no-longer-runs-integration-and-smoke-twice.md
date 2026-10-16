---
id: S-0405
type: story
nature: improvement
title: With the merge queue on, a story's close-out and flai verify leave the test tiers the manifest marks for the queue to the batch, so a close-out main outruns no longer runs integration and smoke twice
status: backlog
parent: E-0022
owner: alex
created: 2026-10-16T04:59:12Z
updated: 2026-10-16T20:56:07Z
transitions: []
tags: [cli]
topics: [conventions, template]
touches: [design/adrs, flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/verify/manifest.go, flai/internal/verify/select.go, flai/internal/verify/story.go, flai/internal/verify/story_test.go, flai/cmd/verify.go, flai/cmd/verify_test.go, flai/internal/mergequeue/batch.go, flai/internal/mergequeue/batch_test.go, scripts/close-out.sh, template/root/scripts/close-out.sh, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/template.yaml, template/CHANGELOG.md, system-flow.yaml, design/system/project-manifest.md, design/system/workflow.md, design/system/flai-cli.md, docs/operators/settings.md, docs/users/flai.md]
after: [S-0399]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
forecast:
  duration: 41m
  delivery: 2026-10-17T03:15:00Z
  basis: "Its own forecast of 41m; 15th in the pull order with an in-progress limit of 5, behind S-0241, S-0392, S-0393, S-0394, S-0395, S-0396, S-0397, S-0398, S-0399, S-0400, S-0401, S-0402, S-0403 and S-0404."
  by: flai
  at: 2026-10-16T20:56:07Z
---
# S-0405 With the merge queue on, a story's close-out and flai verify leave the test tiers the manifest marks for the queue to the batch, so a close-out main outruns no longer runs integration and smoke twice

## Goal

A story's close-out stops paying for the most expensive tiers when the queue will run them anyway on the batch. A test tier marked for the queue is left out of the close-out and of `flai verify` at review while the queue is on, and the batch runs it once for every story it holds, so that a close-out that main outruns, as in I-0148, no longer runs integration and smoke twice past the agent's tool limit.

## Acceptance criteria

- [ ] A test tier in `system-flow.yaml` may carry `queue: true`, validated by the manifest's checks. This project's `integration` and `smoke` tiers carry it.
- [ ] While the `merge-queue` host action is on, `flai verify` leaves the tiers marked `queue` out, names each in its step as left to the queue, and passes without them; with it off, it runs them as today. `--all` and `flai test --all` still run every tier.
- [ ] The batch's verification always runs the tiers marked `queue` that the batch's paths select, so no story lands without them.
- [ ] An ADR, refining ADR-0110 and ADR-0156, records it. The close-out rule in `work-management.md`, in the repository and the template, and both copies of `scripts/close-out.sh` say what the close-out leaves to the queue; `template/CHANGELOG.md` names the version of flai it needs.
- [ ] `design/system/project-manifest.md`, `docs/operators/settings.md`, and `docs/users/flai.md` describe the tier key.

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
