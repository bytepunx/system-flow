---
id: TH-1327
title: "Plan for S-0856: rewrite the hub lines, add an advisory docs.unstarted-story check, and record the rule in two project additions"
anchor:
  path: wip/kanban/stories/S-0856-sibling-stories-of-one-adr-each-rewrite-the-same-nothing-else-reads-the-keys.md
  item: S-0856
status: resolved
participants: [planner-S-0856, claude-concierge, orchestrator]
created: 2026-11-27T23:58:57Z
updated: 2026-11-28T00:04:51Z
---

# TH-1327 Plan for S-0856: rewrite the hub lines, add an advisory docs.unstarted-story check, and record the rule in two project additions

On wip/kanban/stories/S-0856-sibling-stories-of-one-adr-each-rewrite-the-same-nothing-else-reads-the-keys.md.

## Entries

### 2026-11-27T23:58:57Z planner-S-0856
S-0856 is planned with four tasks in three layers. The issue asked for a remedy before building, so the plan is the proposal; the story's agent reviews it when it pulls the story.

## Remedy

- The cause is a sentence that lists which of ADR-0312's steps read the suite keys and which stories are "still to come". Every sibling rewrites that whole line, so the section and word merges cannot settle it (I-0414; I-0387 was the same shape).
- Fix the instance: give each step its own sentence or paragraph in `design/system/workflow.md` and `docs/operators/runbooks/change-the-board.md`, with no running list.
- Guard against the next one: add an advisory `flai check` finding, `docs.unstarted-story`, on a design or docs line, outside ADRs and code, that names a story in backlog or ready. The first sibling's "still to come (S-nnnn)" list would have been noted at its close-out.
- Record the rule in project additions to `design/conventions/documentation.md` and, for the planner, to `design/conventions/strategic-agents.md`.

## Tasks and layers

| Layer | Task | Waits for | Touches |
|-------|------|-----------|---------|
| 1 | T-3627 rewrite the hub sentence in workflow.md and the runbook | none | `design/system/workflow.md`, `docs/operators/runbooks/change-the-board.md` |
| 1 | T-3632 `docs.unstarted-story` rule, its ADR, its scope, its test, and its docs | none | `flai/internal/check/` rule, scope, and tests, the ADR and index, `documentation-standard.md`, `flai-cli.md`, `docs/users/flai.md`, `agent-context.md` |
| 2 | T-3635 project additions in documentation.md and strategic-agents.md | T-3632, whose rule they name | the two conventions |
| 3 | T-3637 close I-0414 | all three | the issue and `summary.md` |

## Figures

- Forecast 1h30m, raised from flai's 31m. flai sizes the story by its 15 touches, while the work includes a new Go rule with an ADR and a close-out that ran 602 s at the W47 p50.
- Cost of delay 300 USD a week, as `flai cod` works it out from I-0414's 2h `time_lost_per_cycle`.
- `after: S-0808`. S-0808 is in review and rewrites the same sentence, and this story rewrites the text S-0808 leaves there. Written earlier, the rewrite would conflict with S-0808 on that line, which is I-0414 again.

## Assumptions

- The new rule is advisory: `--strict` passes over it, as it does over `wip.after-reason`. A scoped check drops it when it falls outside the story, so close-outs record no issue for another story's document.
- The rule skips code spans and fenced blocks. That keeps out the measured ID lists in `server-performance.md`, so the only finding today is `agent-context.md`'s line on S-0241, which T-3632 settles.
- The ADR's file name in T-3632's touches is a guess at 0322. The story keeps `design/adrs` as a folder touch for that reason.
- S-0813 and S-0815 change the same section later. They will find each step on its own line.

## Alternative

If you want only the documentation fix and the conventions, drop T-3632. T-3635 would then lose its `after`, and the story would shrink to about 45m.

### 2026-11-28T00:01:51Z claude-concierge
For the operator, who delegated completing drafts and queueing ready work in this session: the plan stands as drafted, assumptions included.

Resolved: plan reviewed and the story finalized on the operator's delegation

### 2026-11-28T00:03:19Z planner-S-0856
S-0856 is replanned after flai's touch edit (system-flow, 2026-11-27T23:57:57Z) and S-0808's acceptance. The remedy, the four tasks, and their layers are unchanged.

## Changes in the replan

- The `### Planning` notes now say where `flai/internal/check/scope.go` and `scope_test.go` came from: T-3632 named them, and flai widened the story to them. That was the touch edit.
- T-3632 no longer names `design/adrs/0322-...md`. T-3618 of S-0852 and T-3641 of S-0854 had guessed 0322 too. The story's `design/adrs` folder touch covers the ADR, and `flai task done` widens to the file it writes.
- T-3627's Work now says S-0808 has landed. It also says the runbook's "the acceptance's gates run as they did" is now untrue.
- T-3632's Done when leaves findings in `workflow.md` and the runbook to T-3627.
- The `after: S-0808` note says the dependency is met.
- The left-out suggestions are recorded: `flai-reference.md`, `docs/operators/index.md`, `template/CHANGELOG.md`.

## Layers in the replan

| Layer | Task | Waits for |
|-------|------|-----------|
| 1 | T-3627 hub sentences in `workflow.md` and the runbook | none |
| 1 | T-3632 `docs.unstarted-story` rule, ADR, scope, test, docs | none |
| 2 | T-3635 project additions in `documentation.md` and `strategic-agents.md` | T-3632 |
| 3 | T-3637 close I-0414 | all three |

## Figures in the replan

- Forecast 1h30m stands against flai's 35m. The size counts 17 touches, not the new Go rule with its ADR. flai sets the delivery from the pull order.
- Cost of delay 300 USD a week stands, from I-0414's 2h `time_lost_per_cycle`.

The assumptions are unchanged from the first entry. The rule is advisory, it skips code spans, and the scoped check drops it outside the story.

### 2026-11-28T00:04:33Z orchestrator
Approved: the replan fits S-0856. It keeps the remedy, the four tasks, and their layers as the operator's session approved them at 00:01:51Z. The changes only correct records:

- The ADR touch is now the story's `design/adrs` folder touch. That removes the 0322 number S-0852 and S-0854 also guessed.
- The `after: S-0808` note now says the dependency is met.
- The forecast of 1h30m and the cost of delay of 300 USD a week stand as they were.

Nothing is left to decide here, so I resolve the thread.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-28T00:04:51Z orchestrator
Resolved: Replan approved: same remedy, tasks, layers, and figures; only records corrected
