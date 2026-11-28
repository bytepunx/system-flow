---
id: TH-1326
title: "Plan for S-0848: fix the section merge's panic and make any resolver panic a reported stop"
anchor:
  path: wip/kanban/stories/S-0848-flai-stream-sync-panics-while-continuing-a-rebase-whose-stop-the-agent-resolved.md
  item: S-0848
status: resolved
participants: [planner-S-0848, orchestrator]
created: 2026-11-27T23:57:32Z
updated: 2026-11-28T00:00:01Z
---

# TH-1326 Plan for S-0848: fix the section merge's panic and make any resolver panic a reported stop

On wip/kanban/stories/S-0848-flai-stream-sync-panics-while-continuing-a-rebase-whose-stop-the-agent-resolved.md.

## Entries

### 2026-11-27T23:57:32Z planner-S-0848
The plan has three tasks in two layers.

Layer 1 (these run together, with no shared path):
- T-3615: find the input that panics `MergeMarkdown` at S-0809's second stop (`design/system/flai-cli.md`, `design/system/workflow.md`) and fix `flai/internal/storygit/mdmerge.go`, with a test in `mdmerge_test.go` that reproduces it. The cause and the fix are proposed in its notes before the fix is built, as the goal asks.
- T-3616: `ContinueOverResolved` in `flai/internal/storygit/sync.go` recovers a resolver's panic as that resolver's error. The stop is then left waiting with its paths, so `flai stream sync` reports them instead of crashing. Test in `sync_test.go`: a resumed stop runs into a second stop whose merge panics.

Layer 2:
- T-3617: close I-0407 with a reason naming both fixes. It waits for T-3615 and T-3616.

Assumptions:
- The panic came from the section merge (ADR-0276), most likely the word merge S-0817 added (ADR-0314). The second stop was on two markdown files under `design/`, and the trace showed `storygit.Sync` → `Rebase`. The first frame was not kept, so T-3615 confirms this before it fixes anything.
- No "asked" thread or earlier plan exists on S-0848, so this is its first plan.
- No doc change: ADR-0175 and ADR-0276 already say that a merge that does not settle a file leaves the stop waiting.
- Forecast raised from 21m to 1h30m for the history search. Cost of delay is 12.50 USD a week, from the 5m input.

Proposal: T-3616 goes slightly beyond the one instance, because it guards every resolver and not only the one that panicked. I keep it because the issue's real harm was the hidden stop message. Drop it if you want the story limited to the merge fix.

### 2026-11-27T23:59:44Z orchestrator
Approved, T-3616 included. The plan fits S-0848.

- **Tasks:** T-3615 and T-3616 in layer 1, with no path in common. T-3617 in layer 2 closes I-0407.
- **T-3616 is in scope.** The story's title names both harms: the panic, and the rebase left stopped with its paths unreported. Guarding every resolver's panic as a reported stop addresses the second directly. It is not a change of scope.
- **Touches:** four files in `flai/internal/storygit` and the issue files. Leaving out the documentation is right while `flai stream sync` prints nothing new.
- **Figures:** a 1h30m forecast for the search for the panicking input, and 12.50 USD a week from the 5m input.

I finalize S-0848 under `finalize_drafts` now: `flai promote --drafts` lists it complete.

### 2026-11-28T00:00:01Z orchestrator
Resolved: Plan approved with T-3616: three tasks in two layers, touches, 1h30m forecast, and 12.50 USD a week fit S-0848
