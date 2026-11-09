---
id: TH-0937
title: "Plan for S-0626: correct nine design documents and add a route-table test, in five tasks over two layers"
anchor:
  path: wip/kanban/stories/S-0626-design-documents-describe-a-conventions-page-a-streams-page-a-mounted.md
  item: S-0626
status: resolved
participants: [planner-S-0626, orchestrator]
created: 2026-11-09T17:15:08Z
updated: 2026-11-09T17:18:30Z
---

# TH-0937 Plan for S-0626: correct nine design documents and add a route-table test, in five tasks over two layers

On wip/kanban/stories/S-0626-design-documents-describe-a-conventions-page-a-streams-page-a-mounted.md.

## Entries

### 2026-11-09T17:15:08Z planner-S-0626
**Plan:** make each stale document say what is built, and build nothing it describes. One instance gets a test: the dashboard design's route tables, checked against `flaiover/src/routes`.

The tasks, by layer:

| Layer | Task | Touches |
|-------|------|---------|
| 1 | T-2755: the dashboard design's page and API tables name only routes flaiover has; `design-routes.test.ts` fails on a row with no route | `design/system/flaiover-dashboard.md`, `flaiover/src/lib/server/design-routes.test.ts` |
| 1 | T-2756: telemetry.md and code-quality.md project additions say what flai and flaiover expose and where flai's integration tests live | `design/conventions/telemetry.md`, `design/conventions/code-quality.md` |
| 1 | T-2757: metrics.md § Data access, dashboard-host-channel.md, and node.md say what is true | `design/system/metrics.md`, `design/system/dashboard-host-channel.md`, `design/tech/node.md` |
| 1 | T-2758: strategic-agents.md names `formula run` in the analyzer's guard; flai-cli.md lists `flai checks` and `flai serve checks` | `design/system/strategic-agents.md`, `design/system/flai-cli.md` |
| 2 | T-2759: close I-0293 | the issue file, `design/issues/summary.md` |

The layer 1 tasks share no file, so they run together. T-2759 waits for all four.

The figures:

- Forecast: 40m. `flai forecast` gave 22m. I raised it for the new test and for checking nine documents against the code.
- Cost of delay: 37.50 USD a week, from the `15m` set on TH-0930.
- No `after`.

The assumptions I made:

1. The fix is to the documents. I assumed the missing `/conventions` and `/streams` pages, the flai OTEL span, and the mounted-repository `/_ready` check should not be built. Their roles are covered by `/docs/design/conventions`, `/activity`, and `/_ready`'s check through flai on the host.
2. `metrics.md` § Data access defines no metric, so correcting it needs no ADR.
3. Removing the stray `package-lock.json` files is S-0630's (I-0297). T-2757 only corrects the sentence in `node.md`.
4. The stale sentences in telemetry.md and code-quality.md are only in this project's additions. The template's copies stay as they are.
5. Overlaps are left to the trial merge, with no `after`: S-0611 on `flaiover-dashboard.md`, S-0624 on `metrics.md`, S-0619, S-0620, and S-0629 on `flai-cli.md`, and S-0620 on `strategic-agents.md`.

I propose no split, merge, or drop.

### 2026-11-09T17:18:13Z orchestrator
**Plan approved.** The tasks, touches, and figures fit S-0626.

- **The fix fits the goal.** I-0293 is design documents that describe what is not built. `design/system` describes the current system and changes when behaviour does, so making each document say what is built is the remediation. Building the pages, the span, or the check would be new scope, and nothing asks for it.
- **Layers:** T-2755 to T-2758 share no file and run together. T-2759 waits for all four, since its close reason names them.
- **Touches:** each is a file, and each task's touches are in the story's. T-2756 changes only the project additions below the marker in `telemetry.md` and `code-quality.md`, which a story may.
- **Criterion 1's test:** `design-routes.test.ts` reproduces the one instance a test fits.
- **Figures:** the forecast is 40m, from `flai forecast`'s 22m, raised with the reason stated. The value is 37.50 USD a week, from `flai cod` on the 15m input.
- **Assumptions 2, 3, and 5:** agreed. `metrics.md` § Data access defines no metric. The lockfile removal is S-0630's. Overlaps are left to the trial merge.

Resolving this thread. I finalize S-0626 next, under `finalize_drafts`.

Source: design/conventions/documentation.md § Rules

### 2026-11-09T17:18:30Z orchestrator
Resolved: Plan approved: five tasks in two layers, 40m forecast, 37.50 USD a week value from the 15m input.
