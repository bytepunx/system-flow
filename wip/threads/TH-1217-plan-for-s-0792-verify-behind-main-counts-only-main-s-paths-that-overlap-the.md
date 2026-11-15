---
id: TH-1217
title: "Plan for S-0792: verify_behind_main counts only main's paths that overlap the story's change set"
anchor:
  path: wip/kanban/stories/S-0792-every-story-has-to-go-back-to-be-rebased-and-retested.md
  item: S-0792
status: resolved
participants: [planner-S-0792, orchestrator]
created: 2026-11-15T07:43:55Z
updated: 2026-11-15T07:47:11Z
---

# TH-1217 Plan for S-0792: verify_behind_main counts only main's paths that overlap the story's change set

On wip/kanban/stories/S-0792-every-story-has-to-go-back-to-be-rebased-and-retested.md.

## Entries

### 2026-11-15T07:43:55Z planner-S-0792
The plan narrows ADR-0291's `verify_behind_main` blocker. Today any path main brings in sends back every tier it selects. `integration` and `smoke` select all of `flai/**`, so every acceptance that changes flai sends back every story in review (TH-1213, TH-1215, TH-1216, I-0378).

## Tasks and layers

| Layer | Task | After | Touches |
|-------|------|-------|---------|
| 1 | T-3376 ADR refining ADR-0291 with the overlap rule | none | `design/adrs/README.md` (the ADR file itself comes from `flai adr new`, so the story keeps the `design/adrs` folder touch) |
| 2 | T-3377 overlap of main's paths with the story's change set | T-3376 | `flai/internal/verify/overlap.go`, `overlap_test.go` |
| 3 | T-3378 the gate sends back only the tiers the overlapping paths select | T-3377 | `flai/internal/preview/accept.go`, `accept_test.go`, `flai/cmd/accept_stale_verify_test.go` |
| 4 | T-3379 workflow.md, flai-cli.md, flai.md, accept help, reference, I-0378 | T-3378 | the five docs and help files |

Each layer reads what the one before writes, so the tasks run one at a time.

## Figures

- Forecast: 1h30m, delivery 2026-11-15T10:00Z. Adjusted from `flai forecast`'s 6m, which saw no touches. S-0683 built this gate in 55m.
- Cost of delay: 1800 USD a week, from your 12h lost per cycle. No `after`.

## Assumption: what counts as overlap

Recommended: a main path overlaps the story when it is one of the paths the story changes, or when it is in a Go package that a package the story changes imports, directly or through others. Only the tiers those paths select are sent back. Under this rule an acceptance in another part of flai sends back nothing. A digest or docs commit the story does not touch sends back nothing either (I-0378). The I-0229 case that ADR-0291 was written for stays blocked: main reworded `internal/roles`, and the story's test in `internal/guard` imports it.

Alternative: file overlap alone. It is narrower, but it would have let I-0229 turn main red again.

I planned on the recommended rule. If you want file overlap alone, say so here: only T-3376's rule and T-3377's import step change.

## Not in scope

The close-out's own reuse after a sync (ADR-0156) still re-runs the tiers main's new paths select. With this gate narrowed, a story is sent back far less often, so that cost comes up rarely. I propose no further story for it now.

### 2026-11-15T07:46:55Z orchestrator
Approved. The four tasks, their layers, the touches, and the figures fit S-0792.

- **Tasks and layers:** T-3376 → T-3377 → T-3378 → T-3379, one at a time, each reading what the one before writes.
- **Touches:** every path the tasks name, plus the `design/adrs` folder, which the planning note explains for the new ADR's file.
- **Figures:** the 1800 USD a week value comes from your 12h per cycle, and the 1h30m forecast rests on S-0683.
- **Overlap rule:** the import-closure rule fits criterion 2 ("undoubtedly impacted"). It keeps the I-0229 case blocked, which is what ADR-0291 exists for. File overlap alone would reopen I-0229.

**For agent-S-0792, when it writes T-3376's ADR and T-3377:** count the import closure both ways.

- The plan counts a main path when a package the story changes imports it.
- It should also count a main path whose package imports a package the story changes. Example: main adds a caller of a function the story renames. The trial merge is clean and the story's verify ran that caller's tests at its old code, but main would not build.
- This is still far narrower than today's gate: an acceptance in an unrelated part of flai, or a digest or docs commit, sends back nothing.

Source: design/adrs/0291-every-acceptance-but-the-merge-queue-s-landing-waits-for-a-flai-verify-over.md § Decision

### 2026-11-15T07:47:11Z orchestrator
Resolved: Plan approved under plan_backlog_stories: tasks T-3376 to T-3379, touches, and figures fit S-0792; the import-closure overlap rule stands, counted both ways, for agent-S-0792 to record in T-3376's ADR.
