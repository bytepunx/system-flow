---
title: Workflow and board policies
updated: 2026-09-24
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
| done | none; a story reaches done only by acceptance, which archives it | all |

Limits apply to stories. Tasks inherit their story's column budget. Epics are not limited; there should rarely be more than a handful open. `flai check` warns when a limit is exceeded. It does not block, because an agent finishing a story is more valuable than a hard stop, but the warning is recorded in the narrative.

## Pull policy

Work is pulled, not pushed. An agent starting a session:

1. Reads `wip/agents/index.md` to find streams with open narratives it was working on.
2. If none, reads `wip/kanban/board.md` and pulls the highest ordered `ready` story into `in-progress`, respecting the WIP limit.
3. Opens or resumes the narrative for that story in `wip/agents/<story-id>.md`.
4. If the story has no tasks, reads its goal, acceptance criteria, and notes and writes them, each with `## Work` and `## Done when`. If the story does not say enough to do that, blocks it with the reason, opens a thread on it saying what is missing, and pulls the next story.
5. Works tasks in order, transitioning each task as it goes.

## Definition of ready (story)

- `## Goal` and `## Acceptance criteria` are filled.
- Parent epic is not `cancelled`.

Tasks are not part of ready. The agent that pulls the story writes them once it is `in-progress`, and a story cannot enter `review` without at least one ([ADR-0021](../adrs/0021-story-ready-without-tasks.md)). `flai move` and `flai check` enforce both halves: nothing about tasks is checked on the way into `ready` or `in-progress`, and `story.tasks` is raised for `review` and `done` stories.

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
| ready to in-progress | Agent pulling work | Narrative opened; tasks written next if the story has none |
| in-progress to review | Agent | At least one task exists, acceptance criteria self-checked, narrative summary current |
| review to done | Human, or agent if the story is tagged `auto-accept` | Definition of done met. For a story this transition is acceptance, however it is made: `flai accept`, `flai move <story> done`, a card dropped on done, or the item page button all run the same flow (S-0046). `flai accept` does the acceptance: rebase and merge the story branch, move to done, archive, commit. Nothing is tagged or pushed at acceptance (S-0087, ADR-0032): publishing what has accumulated is `flai release --pending`, a deliberate step of its own, run by hand or from the board's Publish action, see `design/conventions/git.md` |
| review to in-progress | Human | Reason appended to story notes |
| backlog, ready, or in-progress to cancelled | Human | Reason appended to the item's notes. Everything open under the item is cancelled with it ([ADR-0028](../adrs/0028-cancelling-an-item-cancels-everything-open-under-it.md)) |
| review to cancelled | Only a parent's cancellation | An item in review is accepted or sent back; it is cancelled only with its epic or story |

`flai move <id> <state>` performs a transition, appends it to `transitions`, updates `status` and `updated`, and validates the rule for that transition. A `--reason` is recorded under the item's `## Notes`. Tasks may skip `review` and go from `in-progress` to `done` directly. A move to `cancelled` takes everything open under the item with it, in one operation: an epic's stories that are not done or cancelled and their open tasks, a story's open tasks. Each gets its own `cancelled` transition with the same actor and time and a note naming the item whose cancellation caused it; done, cancelled, and archived items are left alone; every move is validated before any file is written. `flai move` lists the items first, asks on a terminal unless `--yes`, and changes nothing under `--dry-run`; the board shows the same list in its confirmation. Branches, worktrees, and narratives are never removed by a cancellation; the command names what it left. `flai check` reports an open item under a cancelled parent as `item.parent-cancelled`. A story moved from `review` to `done` is accepted rather than merely transitioned; a story found `done` but unarchived was never accepted, `flai check` flags it as `story.unaccepted`, and `flai accept` completes it. The board's `order` list is the pull order: ready stories first, then backlog stories in the order they should be refined.

### Natures that do not release ([ADR-0025](../adrs/0025-research-is-accepted-without-a-release.md))

Acceptance of a `research` story is the same act as any other (merge, done, archive, commit) since S-0087: no acceptance tags or pushes any more, so there is nothing left to say a research story does differently at that point. What is still true: publishing (`flai release --pending`, [ADR-0032](../adrs/0032-accepting-a-story-merges-it-publishing-is-a-deliberate-batched-step-over.md)) never gives a component a bump on a research story's account, whatever it touched. An `experiment` story is refused by acceptance, before anything is merged, and stays on its branch.

### The pull order (S-0057)

Only the `ready` and `backlog` columns have an order, and only stories are in it; the other columns are read by what happened to the items in them, and tasks follow their story. A column reads as the stories `order` names, in its order, then the stories it does not name, by ID. So a backlog nobody has prioritised still has one reading, and placing a story says something only about the stories placed above it.

`flai order <story> --before <other> | --after <other> | --top | --bottom` places a story within its column. It refuses an epic or a task, a story in any other state, a reference in another column (changing the column is `flai move`), and a story relative to itself, each with the reason. On every write the list is normalised: ready stories, then backlog stories, and names that are no longer ready or backlog stories are dropped. Backlog stories that were never placed and still come last by ID stay unnamed, so one placement does not make `board.md` name the whole backlog. Ready stories are always named.

`flai move` keeps the list true as stories change column: a story moved to `ready` is named after the ready stories already named and before any backlog story, wherever the list had it; a story moved anywhere else is removed, so one sent back to `backlog` returns to the unplaced ones. `flai board`, the MCP `board` and `inbox` tools, and the dashboard's board all lay out `backlog` and `ready` in this sequence, and the dashboard's drag within a column calls `flai order` (ADR-0016). An agent's MCP `inbox` reports a reordering of the stories it had already seen as a change.

Ready means work starts (S-0079). An agent that is running and idle holds `wait_for_work` (S-0097), which names the story to pull as soon as one is ready and the in-progress limit leaves room; `inbox` lists ready stories to one that ends its turns. When none is, and the operator has enabled it on the host, `flai serve` starts each ready story's agent in this order: once per story entering ready, and again when its agent is changed while it is in ready (S-0116), only while the in-progress limit leaves room, whoever is attending ([ADR-0043](../adrs/0043-flai-serve-starts-a-ready-story-s-agent-whenever-the-in-progress-limit-has-room.md)). A story in ready or in progress whose agent dropped or failed gets a new one on the operator's word: `flai serve agent restart`, or Restart agent on its page. See the operators' guide for what enabling it means.

## Blocking

When an agent cannot proceed it appends an open `blocked` interval with a reason, writes the question in the narrative under `## Open questions`, and moves on to another task or story if one is available. When the block clears, the interval is closed. The item never leaves its column because of a block.

## Cadence

There are no sprints. Flow is continuous. The dashboard's charts replace the status meeting: cycle time scatter shows whether stories are getting slower, burn-up shows scope versus completion per epic, and the state-time breakdown shows where time is spent. Reviewing those charts weekly and adjusting WIP limits or definitions is the process improvement loop.

## Branches and collisions (ADR-0019)

Each story is worked on `story/S-nnnn` in a worktree under `.flai-cache/worktrees/`; `wip/` is written in the main checkout so the board is live. Agents run `flai stream sync` at every task transition and `flai accept` rebases and merges. Stories and tasks may declare `touches`; `flai check` warns when two in-progress items overlap and the dashboard badges the documents. The designer edits on `main` through the dashboard, so their intent reaches the next sync.
