---
title: Workflow and board policies
updated: 2026-10-03
status: active
topics: [all]
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

Limits apply to stories. Tasks inherit their story's column budget. Epics are not limited; there should rarely be more than a handful open. `flai check` warns when a limit is exceeded. It does not block, because an agent finishing a story is more valuable than a hard stop, but the warning is recorded in the narrative. `flai check --strict` fails on ready or in-progress over its limit, but passes over review over its limit: only the operator's acceptance clears it, so it must not stop a story's close-out. It counts in the result's `advisory` and the summary line says so. A full review holds the pull instead (S-0243, [ADR-0073](../adrs/0073-a-full-review-holds-the-pull-and-flai-check-strict-passes-over-review-over-its.md)).

## Pull policy

Work is pulled, not pushed. An agent starting a session:

1. Reads `wip/agents/index.md` to find streams with open narratives it was working on.
2. If none, reads `wip/kanban/board.md` and pulls the highest ordered `ready` story into `in-progress`, respecting the in-progress limit. While review is at or over its limit, it pulls nothing and waits for acceptance (S-0243, [ADR-0073](../adrs/0073-a-full-review-holds-the-pull-and-flai-check-strict-passes-over-review-over-its.md)).
3. Opens or resumes the narrative for that story in `wip/agents/<story-id>.md`.
4. If the story has no tasks, reads its goal, acceptance criteria, and notes and writes them, each with `## Work` and `## Done when`. If the story does not say enough to do that, blocks it with the reason, opens a thread on it saying what is missing, and pulls the next story.
5. Works tasks in order, transitioning each task as it goes.

## Definition of ready (story)

- `## Goal` and `## Acceptance criteria` are filled.
- Parent epic is not `cancelled`.
- The story is not a `draft` (S-0199).

A draft is a story an agent wrote, by the planner or from an issue, that the operator has not yet read ([ADR-0074](../adrs/0074-work-items-carry-planning-data-a-story-s-draft-flag-an-epic-s-or-story-s-cost.md)). `flai move` refuses a draft to ready with "finalize it first". The operator finalizes it: `flai move <story> ready --yes` finalizes it as it moves, `flai edit <story> --no-draft` finalizes it in place, and the dashboard's `item.move` with `finalize` runs the move with `--yes`. No agent finalizes until the orchestrator's permission to finalize exists (S-0218): MCP's `item_move` refuses a draft to ready, and `item_edit` refuses `draft: false`; the agent says in a thread or its narrative that the story is ready to be finalized. `flai check` warns (`story.draft`) on a draft in ready or later.

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
| backlog to ready | Human or agent during refinement; a draft story, the operator only | Definition of ready met; a draft is finalized first or by the move (`--yes`) |
| ready to in-progress | Agent pulling work | Narrative opened; tasks written next if the story has none |
| in-progress to review | Agent | At least one task exists, no open question in the narrative, nothing uncommitted in the story's worktree (S-0140: `flai move` and `item_move` refuse it, naming the paths), acceptance criteria self-checked, narrative summary current |
| review to done | Human, or agent if the story is tagged `auto-accept` | Definition of done met. For a story this transition is acceptance, however it is made: `flai accept`, `flai move <story> done`, a card dropped on done, or the item page button all run the same flow (S-0046). `flai accept` does the acceptance: rebase and merge the story branch, move to done, archive, commit. Its preview lists anything uncommitted in the story's worktree (`worktree_uncommitted`) and blocks on it before anything is merged; `flai serve agent commit <story>`, or **Have an agent commit them** in the dashboard's confirmation, starts the story's agent to commit it (S-0140). Nothing is tagged or pushed at acceptance (S-0087, ADR-0032): publishing what has accumulated is `flai release --pending`, a deliberate step of its own, run by hand or from the board's Publish action, see `design/conventions/git.md`. Publishing is the one way accepted work reaches the remote: `git fetch`, then `flai release --pending` or Publish, which tags and pushes the branch and the tags, never forcing, and refuses a clone behind its remote's tags or branch. The operator publishes, or an agent the operator asks; no agent pushes accepted work on its own ([ADR-0067](../adrs/0067-accepted-work-reaches-the-remote-only-when-it-is-published-and-agents-publish.md)) |
| review to in-progress | Human | Reason appended to story notes |
| ready to backlog, in-progress to ready | Human | A step back one column, from the board's lane menu or `flai move` ([ADR-0055](../adrs/0055-a-story-moves-back-one-column-from-ready-in-progress-review-or-cancelled-and.md)). A reason is optional and appended to the notes when given. A story back in ready goes last in the ready order; one back in backlog leaves the order |
| cancelled to backlog | Human | Reopens a cancelled item; refused while its parent is cancelled, which is moved back first. What its cancellation cancelled under it stays cancelled. The item is not completed again until it closes ([ADR-0055](../adrs/0055-a-story-moves-back-one-column-from-ready-in-progress-review-or-cancelled-and.md)) |
| done to anything | Nobody | Done is final: acceptance merged and archived the work |
| backlog, ready, or in-progress to cancelled | Human | Reason appended to the item's notes. Everything open under the item is cancelled with it ([ADR-0028](../adrs/0028-cancelling-an-item-cancels-everything-open-under-it.md)) |
| review to cancelled | Only a parent's cancellation | An item in review is accepted or sent back; it is cancelled only with its epic or story |

`flai move <id> <state>` performs a transition, appends it to `transitions`, updates `status` and `updated`, and validates the rule for that transition. A `--reason` is recorded under the item's `## Notes`. Tasks may skip `review` and go from `in-progress` to `done` directly. A move to `cancelled` takes everything open under the item with it, in one operation: an epic's stories that are not done or cancelled and their open tasks, a story's open tasks. Each gets its own `cancelled` transition with the same actor and time and a note naming the item whose cancellation caused it; done, cancelled, and archived items are left alone; every move is validated before any file is written. `flai move` lists the items first, asks on a terminal unless `--yes`, and changes nothing under `--dry-run`; the board shows the same list in its confirmation. Branches, worktrees, and narratives are never removed by a cancellation; the command names what it left. `flai check` reports an open item under a cancelled parent as `item.parent-cancelled`. A story moved from `review` to `done` is accepted rather than merely transitioned; a story found `done` but unarchived was never accepted, `flai check` flags it as `story.unaccepted`, and `flai accept` completes it. The board's `order` list is the pull order: ready stories first, then backlog stories in the order they should be refined.

### Natures that do not release ([ADR-0025](../adrs/0025-research-is-accepted-without-a-release.md), [ADR-0066](../adrs/0066-an-experiment-story-is-accepted-like-any-other-and-records-its-results-in-a.md))

Acceptance of a `research` or `experiment` story is the same act as any other (merge, done, archive, commit) since S-0087: no acceptance tags or pushes any more. What is still true: publishing (`flai release --pending`, [ADR-0032](../adrs/0032-accepting-a-story-merges-it-publishing-is-a-deliberate-batched-step-over.md)) never gives a component a bump on a research story's or an experiment's account, whatever it touched; the plan for one (`flai release <id>`) says it releases nothing and names each component whose files it touched as landing on main unreleased, and the next story that delivers to the component releases them.

An `experiment` story records its results in a document of its own, `design/experiments/<S-nnnn>-<slug>.md`, named as the story's file is, with the hypothesis, the success measure, what was done, the results, and a recommendation to adopt, adapt, or drop ([repository-layout.md](repository-layout.md#design)). Acceptance, from `flai accept`, `flai move <story> done`, or the board, refuses an experiment whose branch has no such document before anything is merged, and names the document it expects; without a branch it looks in the main checkout. The document is found by the story's ID, so any `<S-nnnn>-*.md` in the folder counts. `flai check` validates it (`experiments.*`).

### The pull order (S-0057)

Only the `ready` and `backlog` columns have an order, and only stories are in it; the other columns are read by what happened to the items in them, and tasks follow their story. A column reads as the stories `order` names, in its order, then the stories it does not name, by ID. So a backlog nobody has prioritised still has one reading, and placing a story says something only about the stories placed above it.

`flai order <story> --before <other> | --after <other> | --top | --bottom` places a story within its column. It refuses an epic or a task, a story in any other state, a reference in another column (changing the column is `flai move`), and a story relative to itself, each with the reason. On every write the list is normalised: ready stories, then backlog stories, and names that are no longer ready or backlog stories are dropped. Backlog stories that were never placed and still come last by ID stay unnamed, so one placement does not make `board.md` name the whole backlog. Ready stories are always named.

`flai move` keeps the list true as stories change column: a story moved to `ready` is named after the ready stories already named and before any backlog story, wherever the list had it; a story moved anywhere else is removed, so one sent back to `backlog` returns to the unplaced ones. `flai board`, the MCP `board` and `inbox` tools, and the dashboard's board all lay out `backlog` and `ready` in this sequence, and the dashboard's drag within a column calls `flai order` (ADR-0016). An agent's MCP `inbox` reports a reordering of the stories it had already seen as a change.

Ready means work starts (S-0079). An agent that is running and idle holds `wait_for_work` (S-0097), which names the story to pull as soon as one is ready, not held (below), the in-progress limit leaves room, and review is under its limit; `inbox` lists ready stories to one that ends its turns. When none is, and the operator has enabled it on the host, `flai serve` starts each ready story's agent in this order: once per story entering ready, and again when its agent is changed while it is in ready (S-0116), only while the in-progress limit leaves room and review is under its limit, and not while the story is held, whoever is attending ([ADR-0043](../adrs/0043-flai-serve-starts-a-ready-story-s-agent-whenever-the-in-progress-limit-has-room.md), [ADR-0046](../adrs/0046-a-ready-story-whose-claim-overlaps-an-open-story-s-is-held-yellow-and-with-its.md)). A full review holds the pull (S-0243, [ADR-0073](../adrs/0073-a-full-review-holds-the-pull-and-flai-check-strict-passes-over-review-over-its.md)): flai serve starts no agent for a ready story and says why (`review is full (5 of 5): accept or send back a story`), `inbox` and `wait_for_work` report `can_pull` false with the reason in `pull_hold`, and `wait_for_work` waits with `waiting_for: review`, where a full in-progress limit gives `room`. Stories in progress, and agents already started, go on. A story in ready or in progress whose agent dropped or failed gets a new one on the operator's word: `flai serve agent restart`, or Retry on its page. For a story in ready while the in-progress limit is full, review is full, or the story is held, that queues the agent, and flai serve starts it once there is room and it is clear (S-0118, [ADR-0044](../adrs/0044-a-retry-of-a-failed-agent-on-a-ready-story-with-the-in-progress-limit-full-is.md)). A story in ready gets its agent at once on the operator's word, whatever the launcher's own rules say, past the limit, a full review, and a hold with a warning: `flai serve agent start`, or Start agent on its page (S-0115). See the operators' guide for what enabling it means.

## Blocking

When an agent cannot proceed it appends an open `blocked` interval with a reason, writes the question in the narrative under `## Open questions`, and moves on to another task or story if one is available. When the block clears, the interval is closed. The item never leaves its column because of a block.

## Cadence

There are no sprints. Flow is continuous. The dashboard's charts replace the status meeting: cycle time scatter shows whether stories are getting slower, burn-up shows scope versus completion per epic, and the state-time breakdown shows where time is spent. Reviewing those charts weekly and adjusting WIP limits or definitions is the process improvement loop.

## Branches and collisions (ADR-0019)

Each story is worked on `story/S-nnnn` in a worktree under `.flai-cache/worktrees/`; `wip/` is written in the main checkout so the board is live. Agents work a story task by task ([ADR-0069](../adrs/0069-agents-commit-each-task-on-the-story-branch-after-flai-stream-sync-and-never.md)): when a task is done they commit its changes, with its docs and work item updates, on `story/S-nnnn`; run `flai stream sync`, resolving each conflict it lists, `git add` and `git rebase --continue`; then run the tests for what the task changed and commit any fix they need. Before moving the story to review they commit what is outstanding, sync again, and close out with `scripts/close-out.sh`, which refuses a branch that does not contain the main branch. `flai stream sync` is how agents rebase a story branch; they never start a `git rebase` or `git merge` by hand. It never stashes: it refuses, touching nothing, a worktree with uncommitted changes, naming each path, or with a rebase already in progress; when the rebase stops on conflicts it lists each path with how to continue and how to abort (`git rebase --abort` puts the branch back as it was before the sync), and exits non-zero. `flai accept` rebases and merges. Stories and tasks may declare `touches`; `flai check` warns when two in-progress items overlap and the dashboard badges the documents. The designer edits on `main` through the dashboard, so their intent reaches the next sync.

[ADR-0046](../adrs/0046-a-ready-story-whose-claim-overlaps-an-open-story-s-is-held-yellow-and-with-its.md) (S-0124) turns `touches` into a claim that scheduling acts on. The survey behind it is [agent-coordination.md](agent-coordination.md). S-0128 built the hold:

- **Claim.** A story's claim is its `touches` plus those of its tasks that are not done or cancelled. An entry that is a sub-project's name or tag in `system-flow.yaml` is read as that sub-project's path, so `cli` and `flai` are both `flai`.
- **Open.** A story is open while it is in progress or in review. To flai serve, a story in ready whose agent it has started is open too, so one look never starts two stories that overlap.
- **Held.** A ready story is held when its claim overlaps an open story's by the `/`-prefix rule `flai check` uses (`overlap`), when its own claim or an open story's is empty (`no-touches`), or when a story it names in `after:` is not done (`after`, S-0130).
- **After.** `after: [S-nnnn, …]` names stories that must be done first, whatever they touch ([work-hierarchy.md](work-hierarchy.md#identifiers)). The reason names each with its state and says what clears it: `held (after): waits for S-0129 (in progress); starts when S-0129 is done`. A cancelled story named keeps the hold, and the reason says to drop it from `after:` if the story no longer needs it; a story that does not exist keeps it too, and `flai check` reports it. A story held both ways is held `after`, and its reason goes on `; also held (overlap): …`. Named stories are looked up in the archive too, so an accepted story clears the hold. A story can wait only for a story named in its own `after:`, and `flai check` refuses a cycle, so the pull order does not deadlock.
- **Skip-ahead.** flai serve does not start a held story's agent, and `wait_for_work` does not offer it. Both take the next ready story in pull order that is not held, within the in-progress limit and while review is under its limit. The held story keeps its place and is taken first once it is clear.
- **The reason.** It names the path, the open story and its state, and what clears it: `held (overlap): touches flai/cmd/serve, inside flai/cmd which S-0128 (in progress) touches; starts when S-0128 is accepted, cancelled, or sent back`. Every open story that holds it is named. `flai board` prints it under the card (`HELD`), and `flai board --json`, `board.get`, and the MCP `board`, `inbox`, and `wait_for_work` answers carry it as `held: {code, reason}` on the ready card. `wait_for_work` times out with `waiting_for: held` when every ready story with room is held. `agent.status` reports a held story as `waiting` with the reason and `hold`, whether or not it has had an agent.
- **The operator's word.** `flai move <story> in-progress` and `item_move` warn and move it. `flai serve agent start` warns and starts it, and tells the agent in its prompt, and the journal, that the operator started it and the hold it went past, so the agent does not take it for flai serve's own start or narrow its touches only to clear the hold (S-0182, I-0050). `flai serve agent restart` of a held story in ready queues it.
- **An answered agent.** flai serve starts an agent that ended asking again once its question is answered, in its session, at once while its story is in progress or in review: an open story holds others and is not held. A story sent back to ready meanwhile waits for its hold, the limit, and a full review like any other, and one in backlog is not started (S-0182).

S-0129 built the dashboard's yellow card: a held story's card and page say why it waits and link the story it waits for ([flaiover-dashboard.md](flaiover-dashboard.md)). Still to come under E-0009, in this order: `after:`, and the trial merge and drift check at `flai stream sync`.

S-0132 built the notice at acceptance. When a story is accepted, flai lists the paths its merge brought into the main branch. Every other open story whose claim covers one of them is told which paths those are. A story with an empty claim is told of every path, because it may change anything. The notice is appended to `.flai-cache/overlaps.jsonl`. The MCP `inbox` and `wait_for_events` report it once, as a change of kind `overlapped` on the open story, with `cause` the accepted story and `to` the paths. `flai accept` prints whom it told and returns the list as `overlaps` in `--json`. On such a notice the agent runs `flai stream sync` and its tests again before it goes on ([work-management.md](../conventions/work-management.md)). A story with no branch merged tells nobody. The notices are best effort: a notice that cannot be worked out is logged, and the acceptance stands.
