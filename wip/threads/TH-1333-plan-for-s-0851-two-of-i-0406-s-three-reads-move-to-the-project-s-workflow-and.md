---
id: TH-1333
title: "Plan for S-0851: two of I-0406's three reads move to the project's workflow, and kluster.go's legacy check keeps the baseline"
anchor:
  path: wip/kanban/stories/S-0851-flai-serve-still-reads-three-state-questions-by-the-baseline-workflow-directly.md
  item: S-0851
status: open
participants: [planner-S-0851, orchestrator]
created: 2026-11-28T00:04:46Z
updated: 2026-11-28T00:07:08Z
---

# TH-1333 Plan for S-0851: two of I-0406's three reads move to the project's workflow, and kluster.go's legacy check keeps the baseline

On wip/kanban/stories/S-0851-flai-serve-still-reads-three-state-questions-by-the-baseline-workflow-directly.md.

## Entries

### 2026-11-28T00:04:46Z planner-S-0851
S-0851 is planned: three tasks in two layers, with touches, a forecast, and a cost of delay value written.

## Proposed solution

| Read | Proposal |
|------|----------|
| `acceptance()`, `flai/internal/serve/curate_schedule.go` | read `repo.Flow()`, which the scheduler's look already holds; S-0748, which its comment waits for, is done |
| `StoryLeft`, `flai/internal/serve/components.go` | read the story's project's workflow, which `readStory` already opens; this is also defect I-0402, a dev server stopped in a state the workflow adds |
| `klusterKept`, `flai/internal/serve/kluster.go` | no change: the baseline judges only a legacy send-back record without `To`, which a flai that ran only the baseline wrote |

## Tasks and layers

| Layer | Task | Touches |
|-------|------|---------|
| 1 | T-3652 The curate scheduler matches an acceptance by the project's workflow | `curate_schedule.go`, `curate_schedule_test.go`, `orchestrate_reasons_test.go` |
| 1 | T-3653 StoryLeft reads the kinds of the story's project's workflow | `components.go`, `components_test.go`, `design/system/component-library.md` |
| 2 | T-3654 Close I-0406 and I-0402 (after T-3652, T-3653) | the two issue files, `design/issues/summary.md` |

T-3652 and T-3653 share no path, so they run together.

## Figures

- Forecast: 1h, delivery 2026-11-28T08:17:00Z. `flai forecast` gave 23m. I raised it because of the two Go changes with renamed-workflow tests, and because the close-out runs the integration and smoke tiers.
- Cost of delay value: 17.50 USD a week, from `flai cod` on the 7m input. It stands as given.

## Assumptions

- `klusterKept`'s baseline fallback is deliberate. T-3654 checks the git history before closing. If a flai that wrote records without `To` could already run a project's workflow file, T-3654 says so here instead of closing.
- `builtinTriggersIn` in `orchestrate_reasons.go` reads `baselineFlow` only to translate baseline state names to the project's. It is outside I-0406, as D-0676 records.
- `flai/cmd/serve_actions.go` stays unchanged: T-3653 picks the form of `ComponentsOptions` that leaves it alone.

## Proposal

T-3654 also closes I-0402, since T-3653 fixes its cause. The criteria name only I-0406. If you would rather keep I-0402 open, I will drop that step from T-3654.

### 2026-11-28T00:07:08Z orchestrator (recommendation)
**Recommended:** approve the plan as it stands, and let T-3654 close I-0402 too.

**Why close I-0402 here:**

- `design/system/workflow.md` § A state added in the workflow file alone (S-0756) records I-0402 as one of the parts of E-0034's outcome that do not follow yet: flai serve stops the dev server of a story in the added state.
- I-0402's cause is the `StoryLeft` read that T-3653 fixes, and T-3653's new `TestStoryLeft` case is that instance: a story in `testing` has not left.
- Leaving I-0402 open after its cause is gone would give the analyzer a stale issue. A later story would then only close it.

**The plan fits the story:**

- T-3652 and T-3653 each carry a test on a renamed workflow (criterion 1). T-3654 closes I-0406 (criterion 2).
- Keeping `klusterKept` on the baseline fits: it judges only a legacy record with no `To`. T-3654 checks that in git history before it closes.
- T-3652 and T-3653 share no path, so layer 1 runs together.
- The 1h forecast states its reason. The 17.50 USD a week value comes from your 7m input.

**Why I recommend rather than answer:** the criteria name only I-0406, so closing I-0402 adds to the story's scope. That call is yours.

**Alternative:** keep I-0402 open. The planner then drops that step from T-3654.

**What happens next:** once this is settled, I finalize S-0851 under `finalize_drafts`. `flai promote --drafts` already lists it complete.

Source: design/system/workflow.md § A state added in the workflow file alone (S-0756)
