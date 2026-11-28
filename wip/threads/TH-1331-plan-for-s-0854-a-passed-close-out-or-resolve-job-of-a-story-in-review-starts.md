---
id: TH-1331
title: "Plan for S-0854: a passed close-out or resolve job of a story in review starts the orchestrator"
anchor:
  path: wip/kanban/stories/S-0854-a-story-in-review-whose-close-out-or-resolve-job-passed-is-not-a-reason-that.md
  item: S-0854
status: resolved
participants: [planner-S-0854, claude-concierge]
created: 2026-11-28T00:00:08Z
updated: 2026-11-28T00:01:51Z
---

# TH-1331 Plan for S-0854: a passed close-out or resolve job of a story in review starts the orchestrator

On wip/kanban/stories/S-0854-a-story-in-review-whose-close-out-or-resolve-job-passed-is-not-a-reason-that.md.

## Entries

### 2026-11-28T00:00:08Z planner-S-0854
S-0854 is planned: 5 tasks in 4 layers. Touches, a forecast of 2h, and a cost of delay of 7,200 USD a week are written. The reasons are under `### Planning` in the story's Notes.

## Proposed fix

The close-out job and the rebase-and-resolve job already append a `verified` event with verdict `passed`. The orchestrator's triggers do not name that event, so nothing starts it for the story. The fix:

- Add `verified` to the trigger events.
- Add a condition, `passed-in-review`. It holds when the verdict is `passed`, the job ran something (not answered from the record), and the story is in review.
- Add `{event: verified, when: passed-in-review}` to the orchestrator's role file and to flai serve's built-in triggers.
- A matched event becomes the reason `review:S-nnnn`.

## Tasks and layers

| Layer | Task | Waits for | Touches |
| --- | --- | --- | --- |
| 1 | T-3641 Write the ADR, then update `design/system/strategic-agents.md`, the roles README, and the convention, each with its template copy | none | `design/adrs/*`, design and convention files |
| 2 | T-3642 Add the `verified` event and the `passed-in-review` condition to the roles vocabulary | T-3641 | `flai/internal/roles/{roles,conditions}{,_test}.go` |
| 3 | T-3645 Match the event in flai serve and prime it as `review:S-nnnn`, with a test that reproduces I-0413 | T-3642 | `flai/internal/serve/orchestrate_reasons{,_test}.go`, `orchestrate_test.go` |
| 3 | T-3646 Add the trigger to the three copies of the orchestrator's role file | T-3642 | role files, `flai/internal/harness/harness_test.go` |
| 4 | T-3648 Update the user docs and the template changelog, and close I-0413 | T-3645, T-3646 | `docs/users/flai.md`, `template/CHANGELOG.md`, I-0413, `summary.md` |

T-3645 and T-3646 run together: they share no file.

## Assumptions

- A verdict answered from the record (`recorded: true`) starts nothing. Otherwise the orchestrator's own `verify` call on a story it cannot accept would start it again, in a loop.
- `failed`, `stopped`, and `interrupted` verdicts start nothing. The story's agent or the operator acts on those.
- Coalescing and the dedupe key stay as they are. A story's move to review and its verified event in one window give one reason.
- The ADR's number is a guess (`0322-…`). `flai task done` corrects the touch, so no `design/adrs` folder touch is kept.
- The story names no `after`. S-0845 and S-0824 also touch the resolve job and the orchestrator's role file, and a trial merge judges that.
- The cost of delay value is worked out from your input, `time_lost_per_cycle` 48h, carried over from I-0413. No input of yours is changed.

Nothing here needs your answer unless you want a different fix.

### 2026-11-28T00:01:51Z claude-concierge
For the operator, who delegated completing drafts and queueing ready work in this session: the plan stands as drafted, assumptions included.

Resolved: plan reviewed and the story finalized on the operator's delegation
