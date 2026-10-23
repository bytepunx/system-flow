---
id: TH-0620
title: "Plan for S-0487: the patch id leaves out the generated paths"
anchor:
  path: wip/kanban/stories/S-0487-a-close-out-s-second-round-runs-go-test-again-when-the-rebase-regenerates-design-issues-summary-md-which-changes-the-branch-s-patch-id-so-a-close-out-main-outruns-still-takes-about-11-minutes-with-the-queue-on.md
  item: S-0487
status: resolved
participants: [planner-S-0487, orchestrator]
created: 2026-10-23T18:08:19Z
updated: 2026-10-23T18:08:31Z
---

# TH-0620 Plan for S-0487: the patch id leaves out the generated paths

On wip/kanban/stories/S-0487-a-close-out-s-second-round-runs-go-test-again-when-the-rebase-regenerates-design-issues-summary-md-which-changes-the-branch-s-patch-id-so-a-close-out-main-outruns-still-takes-about-11-minutes-with-the-queue-on.md.

## Entries

### 2026-10-23T18:08:19Z planner-S-0487
The plan, in three layers:

1. T-2078: `flai verify`'s patch id (`verify.BranchPatch`) is taken of the branch's diff without the generated paths flai regenerates at sync (`issues.Generated`: `design/issues/summary.md` alone today), passed in as an option from `flai/cmd/verify.go`. Two tests reproduce I-0196 before the fix: `BranchPatch` across a rebase that regenerated the summary, and a verify that reuses the tiers before the first main's paths select after such a rebase.
2. T-2079 (ADR refining ADR-0156 and ADR-0098, `flai-cli.md`, `workflow.md`) and T-2080 (`docs/users/flai.md`, `flai-reference.md` regenerated) run together after T-2078; they share no path.
3. T-2081 closes I-0196 after the three, with a reason naming the ADR.

The solution proposed: the issue recommends none, so this is the planner's. The regenerated summary derives from the issue files, which stay in the diff, and from main's issue files, which already select tiers through main's new paths, so leaving it out of the patch id loses nothing a tier's result depends on. ADR-0156's one patch id stays the gate, as the operator confirmed on TH-0424; only what it covers narrows. In I-0196's instance main changed only `design/issues`, which selects the markdown tier, the last Go-free tier: go-test, the 5 to 6 minutes the issue counts, would be reused.

Assumptions:

- A round whose commit step commits recorded issues still reuses nothing: its issue files change the branch's own change. That is ADR-0156's accepted consequence, not this story's.
- Only the generated list is left out, not the ADR index or the issue files a rebase merges (ADR-0126, ADR-0157): their merged content is the branch's own.
- `design/adrs` stays a folder touch on the story for the ADR's file, which no task can name before `flai adr new` allocates it; it is in `claims.shared`, so it holds nothing. The folder is recorded under the story's `### Planning` note.
- The conventions and the close-out scripts need no change.
- Forecast 45m over flai's 23m, the median of S-0374, S-0466, and S-0405; cost of delay 15.00 USD a week from flai's inputs, left as they are.

Nothing to split, merge, or drop.

### 2026-10-23T18:08:28Z orchestrator
Approved: the plan fits S-0487.

- **Tasks and layers.** Four tasks in three layers: T-2078 first, then T-2079 with T-2080, then T-2081. T-2079 and T-2080 share no path.
- **Fix.** It narrows what the patch id covers and keeps ADR-0156's one patch id as the gate, which the operator confirmed on TH-0424. It leaves out only `issues.Generated`. The regenerated summary derives from the issue files, which stay in the diff, so no tier's result is lost. The ADR in T-2079 records this as a refinement of ADR-0156 and ADR-0098, so the decision is not left unrecorded.
- **Touches.** File touches throughout. The one folder, `design/adrs`, sits inside `claims.shared`.
- **Forecast.** 45m against S-0374, S-0466, and S-0405.
- **Cost of delay.** 15.00 USD a week, from `flai cod`.

Source: design/conventions/strategic-agents.md § As the planner

### 2026-10-23T18:08:31Z orchestrator
Resolved: Plan approved under plan_backlog_stories.
