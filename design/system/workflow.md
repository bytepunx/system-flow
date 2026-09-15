---
title: Workflow and board policies
updated: 2026-09-15
status: active
---

# Workflow and board policies

This is the Lean part. The board is `wip/kanban`, the columns are the states in [work-hierarchy.md](work-hierarchy.md), and the policies below are what keeps flow visible and measurable. `wip/kanban/board.md` in each project records that project's WIP limits and any local policy; the defaults are here.

## Columns and WIP limits

| Column | Default WIP limit | Applies to |
|--------|-------------------|------------|
| backlog | none | all |
| ready | 5 | stories |
| in-progress | 2 | stories |
| review | 3 | stories |
| done | none, archived by `flai archive` | all |

Limits apply to stories. Tasks inherit their story's column budget. Epics are not limited; there should rarely be more than a handful open. `flai check` warns when a limit is exceeded. It does not block, because an agent finishing a story is more valuable than a hard stop, but the warning is recorded in the narrative.

## Pull policy

Work is pulled, not pushed. An agent starting a session:

1. Reads `wip/agents/index.md` to find streams with open narratives it was working on.
2. If none, reads `wip/kanban/board.md` and pulls the highest ordered `ready` story into `in-progress`, respecting the WIP limit.
3. Opens or resumes the narrative for that story in `wip/agents/<story-id>.md`.
4. Works tasks in order, transitioning each task as it goes.

## Definition of ready (story)

- `## Goal` and `## Acceptance criteria` are filled.
- Tasks exist, each with `## Work` and `## Done when`.
- Parent epic is not `cancelled`.

## Definition of done (story)

- All acceptance criteria checked.
- All tasks `done` or `cancelled`.
- Any decision made during the story is captured: an ADR if it is architectural, an edit to `design/system` or `design/tech` if it changes the living design, a note in the story otherwise.
- User-facing behaviour changes have a matching edit under `docs/`.
- Narrative in `wip/agents/<story-id>.md` has a closing entry.

## Transitions and who makes them

| Transition | Made by | Required action |
|------------|---------|-----------------|
| backlog to ready | Human or agent during refinement | Definition of ready met |
| ready to in-progress | Agent pulling work | Narrative opened |
| in-progress to review | Agent | Acceptance criteria self-checked, narrative summary current |
| review to done | Human, or agent if the story is tagged `auto-accept` | Definition of done met |
| review to in-progress | Human | Reason appended to story notes |
| any to cancelled | Human | Reason appended to story notes |

`flai move <id> <state>` performs a transition, appends it to `transitions`, updates `status` and `updated`, and validates the rule for that transition. A `--reason` is recorded under the item's `## Notes`. Tasks may skip `review` and go from `in-progress` to `done` directly. The board's `order` list is the pull order: ready stories first, then backlog stories in the order they should be refined. Moving a story to `ready` appends it if absent; starting or cancelling it removes it.

## Blocking

When an agent cannot proceed it appends an open `blocked` interval with a reason, writes the question in the narrative under `## Open questions`, and moves on to another task or story if one is available. When the block clears, the interval is closed. The item never leaves its column because of a block.

## Cadence

There are no sprints. Flow is continuous. The dashboard's charts replace the status meeting: cycle time scatter shows whether stories are getting slower, burn-up shows scope versus completion per epic, and the state-time breakdown shows where time is spent. Reviewing those charts weekly and adjusting WIP limits or definitions is the process improvement loop.
