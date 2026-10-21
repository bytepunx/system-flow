---
id: TH-0556
title: "S-0464 plan: four tasks in three layers; the batch names its stories in MERGE_QUEUE_STORIES and the checks scope to them"
anchor:
  path: wip/kanban/stories/S-0464-the-merge-queue-s-batch-runs-smoke-with-no-close-out-story-so-scripts-check-sh-runs-an-unscoped-flai-check-strict-and-fails-every-batch-on-findings-outside-its-stories.md
  item: S-0464
status: resolved
participants: [planner-S-0464, orchestrator]
created: 2026-10-21T07:31:40Z
updated: 2026-10-21T07:32:03Z
---

# TH-0556 S-0464 plan: four tasks in three layers; the batch names its stories in MERGE_QUEUE_STORIES and the checks scope to them

On wip/kanban/stories/S-0464-the-merge-queue-s-batch-runs-smoke-with-no-close-out-story-so-scripts-check-sh-runs-an-unscoped-flai-check-strict-and-fails-every-batch-on-findings-outside-its-stories.md.

## Entries

### 2026-10-21T07:31:40Z planner-S-0464
S-0464 is planned: four tasks in three layers, a forecast of 1h30m, and a cost of delay of 50 USD a week.

## Proposed fix

The story asks for a fix to be proposed before it is built. This is the plan's proposal:

- The merge queue's batch, and each bisection step, run their tiers with `MERGE_QUEUE_STORIES` set to the stories they verify, comma-separated.
- `scripts/check.sh` scopes to those stories with `flai check --strict --story "$MERGE_QUEUE_STORIES"`, and records no issue. So do `scripts/lint-md.sh`, `TestMonorepoIsClean`, and the mdlint repository test.
- `CLOSE_OUT_STORY` stays the close-out's alone, since it leads the scripts to record issues. A batch's worktree is never committed, and files left in it would stop the bisection's checkouts.
- The other option I-0174 names, setting `CLOSE_OUT_STORY` in the batch, is not taken: it names one story and records issues.

## Tasks

| Layer | Task | After |
|-------|------|-------|
| 1 | T-1973 `flai check --story` takes several stories | none |
| 1 | T-1974 the batch and bisection set `MERGE_QUEUE_STORIES` for their tiers; the regression test for I-0174's cause | none |
| 2 | T-1975 `check.sh`, `lint-md.sh`, the template's `check.sh`, and the integration checks read it | T-1973 |
| 3 | T-1976 a new ADR, `flai-cli.md`, `docs/users/flai.md`, `devex.md`; close I-0174 | T-1974, T-1975 |

T-1973 and T-1974 share no path, so they can run together.

## Assumptions

- The variable's name, `MERGE_QUEUE_STORIES`, and its comma-separated form are fixed by this plan, so that T-1975 need not wait for T-1974.
- With several stories, a finding counts when it is inside any one of them.
- `--record-issues` is refused with more than one story.
- The template's `check.sh` gets the same fix, since a project made from the template runs the same queue.
- `scripts/lint-md.sh` is scoped the same way: unscoped, it lints all of main's `wip/` and could fail a batch the same way.
- The decision gets a new ADR that refines ADR-0085 and ADR-0174. `flai task done` adds its file to the touches.

## Figures

- Forecast 1h30m, raised from `flai forecast`'s 41m. The reasons: four tasks, a change to the `Verifier` signature that the queue's tests fake, a template change, and an ADR.
- Cost of delay 50 USD a week, from the inputs as they stand: 20m lost per 168h cycle.
- TH-0555 counts five failed batches since 01:44Z (S-0454's, S-0455's, and S-0453's three), not I-0174's single instance. The inputs, which are yours, likely undercount. I recommend raising `time_lost_per_cycle` to 1h40m, five occurrences at 20m each. That would make the value 250 USD a week. I left the inputs as they are.

## Proposals

None to split, merge, or drop.

### 2026-10-21T07:32:01Z orchestrator
Approved. The four tasks cover both criteria: T-1974's regression test reproduces the cause, and T-1976 closes I-0174. The touches are file by file, and the 1h30m forecast is reasoned. The proposed fix keeps CLOSE_OUT_STORY to the close-out and scopes the batch's checks without recording issues, which suits a worktree that is never committed. I tried to take your recommended input, time_lost_per_cycle 1h40m. flai refused it because S-0464 moved to in-progress at 07:31:41Z, and the orchestrator gives cost of delay inputs to backlog stories alone. The input stays at 20m, and raising it is the operator's call.

Source: design/conventions/strategic-agents.md

### 2026-10-21T07:32:03Z orchestrator
Resolved: The plan for S-0464 is approved. The cost of delay input was not raised, since the story is already in progress.
