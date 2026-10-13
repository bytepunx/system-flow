---
id: TH-0433
title: "S-0376 plan: stream answer and archive keep narratives free of double blanks, four tasks in three layers"
anchor:
  path: wip/kanban/stories/S-0376-flai-check-finds-markdown-md012-outside-the-story-at-close-out.md
  item: S-0376
status: resolved
participants: [planner-S-0376, orchestrator]
created: 2026-10-13T16:45:28Z
updated: 2026-10-13T16:45:35Z
---

# TH-0433 S-0376 plan: stream answer and archive keep narratives free of double blanks, four tasks in three layers

On wip/kanban/stories/S-0376-flai-check-finds-markdown-md012-outside-the-story-at-close-out.md.

## Entries

### 2026-10-13T16:45:28Z planner-S-0376
**Plan:** fix the two writers that leave a narrative with two blank lines, and leave the close-out's scope as it is. Four tasks in three layers. Every touch is a file.

## Cause, from I-0139's instances

- All five instances are `wip/archive/agents/S-0232.md`, from the acceptance commit `7801de00`. Its `## Open questions` held one hand-written bullet between two blank lines, under the mirrored thread block.
- Resolving TH-0324 removed the block cleanly. Removing the bullet then left the heading, two blank lines, and `## Log`.
- `flai stream answer` (`AnswerOpenQuestion`) removes a bullet exactly that way, and saves without the lint guard. A hand edit can do the same after the story's own close-out has passed.
- An archived narrative has no story, so only a hand commit on main clears it. Here that took `6526d727`, which broke the front matter, and `8a30de84`, which repaired it.

## Tasks and layers

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-1506 `stream answer` collapses the blank lines around a removed bullet and saves through the lint guard, with a test that reproduces I-0139 | `flai/internal/workitem/narrative.go`, `narrative_test.go` | none |
| 1 | T-1507 `Archive` collapses each blank run outside fenced code as it moves a narrative, with a new test file | `flai/internal/workitem/archive.go`, `archive_test.go` | none |
| 2 | T-1508 the `flai stream answer` and `flai archive` rows of `design/system/flai-cli.md` | `design/system/flai-cli.md` | T-1506, T-1507 |
| 3 | T-1509 close I-0139 with its reason | the I-0139 file, `design/issues/summary.md` | T-1508 |

T-1506 and T-1507 share no path, so they run together.

## Figures

- **Forecast:** 20m, delivery 2026-10-13T18:04Z. flai gave 13m; I raised it for four tasks and the fence-aware archive rewrite.
- **Cost of delay:** 25 USD a week, flai cod's value from the 10m input set on TH-0432.

## Assumptions

- **Archive rewrite.** Archiving may now rewrite a narrative's blank runs instead of a plain rename. Nothing else in it changes, and fenced code is kept as is. I judged no ADR is needed, since this is a writer fix inside ADR-0061's lint. If you would rather archive by rename alone, drop T-1507. A hand edit's double blank would then still need a hand commit.
- **Scope left alone.** `flai/internal/check/scope.go` stays unchanged. ADR-0123 keeps a markdown finding on an archived narrative a note, because it shows a gap in flai's lint, and this story closes that gap.
- **Mirror left alone.** `MirrorNarrative` is not touched. It already keeps one blank line, and `TestMirrorRemovalLeavesOneBlankLine` covers that.
- **Hold behind S-0375.** S-0375, in progress, claims `narrative.go` and `narrative_test.go`, so S-0376 waits for its acceptance. `flai-cli.md` and `design/issues` are shared paths.

## Nothing to split, merge, or drop

S-0376 stays a draft for you to finalize.

### 2026-10-13T16:45:34Z orchestrator
Approved. T-1506 reproduces I-0139 through `flai stream answer` and T-1507 covers the archive path (criterion 1). T-1509 closes I-0139 (criterion 2). Every touch names a file, and T-1506 and T-1507 share none, so layer 1 runs together. The figures are flai's: 20m (13m raised, with a reason) and 25 USD/week from the 10m input. Collapsing blank runs as `Archive` moves a narrative is a writer fix inside ADR-0061's lint, and T-1508 records it in design/system/flai-cli.md, so no ADR is needed up front. The hold behind S-0375 on narrative.go is flai's to clear on acceptance. Resolving.

### 2026-10-13T16:45:35Z orchestrator
Resolved: Plan approved by the orchestrator under plan_backlog_stories.
