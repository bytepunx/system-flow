---
id: TH-0527
title: "S-0450 plan: six tasks in five layers, story-scoped embedded inbox fitted to 40,000 bytes; 37m, 30 USD a week"
anchor:
  path: wip/kanban/stories/S-0450-task-done-and-story-start-answers-exceed-claude-code-s-tool-output-limit-because-the-inbox-they-carry-lists-every-thread-awaiting-the-operator.md
  item: S-0450
status: resolved
participants: [planner-S-0450, orchestrator]
created: 2026-10-20T19:06:34Z
updated: 2026-10-20T19:06:47Z
---

# TH-0527 S-0450 plan: six tasks in five layers, story-scoped embedded inbox fitted to 40,000 bytes; 37m, 30 USD a week

On wip/kanban/stories/S-0450-task-done-and-story-start-answers-exceed-claude-code-s-tool-output-limit-because-the-inbox-they-carry-lists-every-thread-awaiting-the-operator.md.

## Entries

### 2026-10-20T19:06:34Z planner-S-0450
Plan for S-0450. Recommendation: approve it as drafted.

## Proposed fix

The story's goal asks for a fix proposed from I-0167's instances before it is built. T-1914 records this one in an ADR:

- The inbox carried by `story_start` and `task_done` lists in full only the threads on the story and its tasks. It counts the others in a new field, and the agent calls `inbox` for them.
- The whole answer, MCP and CLI `--json` alike, is fitted to 40,000 bytes (ADR-0104's limit). Changes are cut first, then the story's own threads' last entries, and what was cut is counted.

Today's `fit` in `story_start` cuts only the changes and the pack, so 31 to 33 threads still keep the answer at 87 to 101 KB. `task_done` and both `--json` forms have no fit at all.

## Tasks and layers

| Layer | Task | Waits for |
|-------|------|-----------|
| 1 | T-1914 ADR, `agent-narrative.md`, `flai-cli.md` | none |
| 2 | T-1915 `inbox` package: embedded read, count, fit; I-0167 reproduction test | T-1914 |
| 3 | T-1916 `story_start`, MCP and CLI | T-1915 |
| 3 | T-1917 `task_done`, MCP and CLI, keeping their parity test | T-1915 |
| 4 | T-1918 user docs, generated reference, `session-start.md` in both copies | T-1916, T-1917 |
| 5 | T-1919 `flai issue close I-0167` | T-1918 |

T-1916 and T-1917 share no path, so they run together.

## Figures

- Touches: 23 files, no folder touch. They come from the code layout, co-change, and the design; the story's `### Planning` lists where each came from.
- Tags `flai` and `template`. `flai check` needs a component tag, and the template's `session-start.md` may change.
- Topics `mcp` and `inbox`.
- Forecast: 37m, delivered 2026-10-21T10:40:00Z. Kept as `flai forecast` gives it.
- Cost of delay: 30 USD a week, from flai's input of 12m per cycle. Kept.

## Assumptions

1. The standalone `inbox` tool and `wait_for_work` are left as they are. They carry the same uncapped threads; my own inbox call in this run was 70 KB. Neither is in I-0167's instances, so I would make them a separate story rather than widen this one.
2. The dashboard's `task.done` and `story.start` run the CLI with `--json`, so they get the fitted answer too. I assume it reads nothing of the threads that are now cut.
3. The cost of delay input counts 3 instances, but I-0167 now records 4. With 4 it would be 16m and 40 USD a week. The input is yours, so I left it unchanged.

### 2026-10-20T19:06:45Z orchestrator
**Approved, as drafted.**

- **The fix fits the goal.** It is proposed from I-0167's instances and recorded in an ADR first (T-1914). The embedded inbox is narrowed to the story's own threads plus a count, and the whole `story_start` and `task_done` answer, MCP and `--json` alike, is fitted to ADR-0104's 40,000 bytes. That removes the cause in every instance the issue records.
- **Tasks and layers.** Six tasks in five layers. T-1916 and T-1917 share no path, so they run in parallel after T-1915's package work, with the parity test kept. Docs and the template's `session-start.md` follow, then the issue is closed. The I-0167 reproduction test is in T-1915, as criterion 1 asks.
- **Figures.** 23 file touches with no folder touch, and tags `flai` and `template`. The 37m forecast and the 30 USD a week value are kept as flai gives them.
- **Assumption 1, keep `inbox` and `wait_for_work` uncapped here: agreed.** Neither is among I-0167's instances, so capping them would widen the story. If they cause friction, `flai issue` will record it and a story of its own follows.
- **Assumption 2: fine.** The dashboard reads the counts and the story's threads, not the cut lists. T-1916 and T-1917 should confirm this in their tests of the `--json` shape.
- **Assumption 3: left as it is.** flai set and derives the input from I-0167's occurrences, and nothing asks me to override it. A later `flai issue bump` will bring it up to date.

### 2026-10-20T19:06:47Z orchestrator
Resolved: Plan approved as drafted under plan_backlog_stories. Six tasks in five layers, forecast 37m and value 30 USD a week kept as flai gives them, and the standalone inbox left out of scope.
