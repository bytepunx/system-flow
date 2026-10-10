---
id: TH-0395
title: "S-0289 plan: MD056 in flai's mdlint, then close I-0077"
anchor:
  path: wip/kanban/stories/S-0289-flai-s-markdown-lint-does-not-split-a-table-row-on-pipes-inside-a-code-span-as-markdownlint-does.md
  item: S-0289
status: resolved
participants: [planner-S-0289, orchestrator]
created: 2026-10-09T18:23:20Z
updated: 2026-10-09T18:23:31Z
---

# TH-0395 S-0289 plan: MD056 in flai's mdlint, then close I-0077

On wip/kanban/stories/S-0289-flai-s-markdown-lint-does-not-split-a-table-row-on-pipes-inside-a-code-span-as-markdownlint-does.md.

## Entries

### 2026-10-09T18:23:20Z planner-S-0289
S-0289 has two tasks in the backlog, plus touches, a 22m forecast, and a cost of delay value of 15.00 USD a week. The story's `### Planning` note says where each figure and touch came from.

## Proposed fix

flai's `mdlint` gains MD056, table column count. I reproduced I-0077's row, from before commit adf27c49, on today's `mdlint` and on markdownlint-cli2 0.20.0:

| Linter | Findings on the row |
|--------|---------------------|
| markdownlint-cli2 0.20.0 | MD038 on the code span the split leaves; MD056, Expected: 2; Actual: 9 |
| flai `mdlint` | MD038 on the same code span |

flai already splits a row on every unescaped pipe, inside a code span too. What it lacks is MD056, so the row passed.

## Tasks and layers

| Layer | Task | After | Touches |
|-------|------|-------|---------|
| 1 | T-1433 MD056 in `mdlint`, a `tables.md` fixture checked against markdownlint-cli2, and the rule named in `flai-cli.md` and `docs/users/flai.md` | none | 6 files under `flai/internal/mdlint`, `design/system/flai-cli.md`, `docs/users/flai.md` |
| 2 | T-1434 `flai issue close I-0077` with a reason naming MD056 | T-1433 | I-0077's file, `design/issues/summary.md` |

T-1434 waits because its reason names the fix T-1433 makes.

## Assumptions

1. MD055 (table pipe style) and MD058 (blanks around tables) stay out. The configuration turns both on, but I-0077 names neither.
2. The docs ship in T-1433, not in a task of their own, because `documentation.md` asks for them in the same commit as the behaviour.
3. ADR-0061 stays as it is, as it did for MD007 and MD038. `flai-cli.md` records the added rule.
4. The wip markdown on main has no MD056 that markdownlint-cli2 does not report too, since CI lints it with markdownlint-cli2. T-1433's Done when checks this.

## Proposals

No split, merge, or drop. The story's title says flai does not split the row, but it does today. You may want the title to say that flai has no MD056. I left the title alone because the story is finalized.

### 2026-10-09T18:23:30Z orchestrator
Approved under `plan_backlog_stories`. The plan fits S-0289:

- Criterion 1: T-1433 adds MD056 to `mdlint`, with a `tables.md` fixture checked against markdownlint-cli2. It reproduces I-0077's row, which today passes flai and fails markdownlint.
- Criterion 2: T-1434 closes I-0077 with `flai issue close` naming MD056. It waits for T-1433 because its reason names that fix.
- The touches are all files, and each one names its source. Leaving out `inline.go` is explained.
- The figures stand as flai gave them: 22m from `flai forecast`, in line with the four earlier mdlint fixes, and 15.00 USD a week from `flai cod` on flai's 6m input from I-0077.

Leaving out MD055 and MD058 keeps to I-0077, which is the story's scope.

The title stays as it is: a finalized story's words are the operator's (`strategic-agents.md`). The story's agent can say in `## Decisions` that the fix is MD056, not the split.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-10-09T18:23:31Z orchestrator
Resolved: Plan approved by the orchestrator under plan_backlog_stories.
