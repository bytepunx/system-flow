---
id: TH-1198
title: "ADR-0301 (the board as a workflow file): nine choices to settle before it is accepted"
anchor:
  path: wip/kanban/stories/S-0736-an-adr-fixes-the-workflow-file-s-format-and-place-under-design-workflows-its.md
  item: S-0736
status: resolved
participants: [agent-S-0736, orchestrator, alex]
created: 2026-11-14T17:47:40Z
updated: 2026-11-14T17:55:34Z
---

# TH-1198 ADR-0301 (the board as a workflow file): nine choices to settle before it is accepted

On wip/kanban/stories/S-0736-an-adr-fixes-the-workflow-file-s-format-and-place-under-design-workflows-its.md.

## Entries

### 2026-11-14T17:47:40Z agent-S-0736
ADR-0301 is written, proposed, on `story/S-0736`. It fixes the workflow file, `design/workflows/board.md`: states with a kind each (`queue`, `active`, `waiting`, `accepted`, `cancelled`), transitions with `by`, `rules`, and `effects` (`commit`, `pull`, `accept`, `send-back`, `cancel`, `reopen`), a 14-rule vocabulary, timers (between two points, or a sum of states), aliases for renamed or removed states, and a baseline that is today's board. The baseline file is in its §7.

Read it at `design/adrs/0301-the-board-is-a-workflow-file-design-workflows-board-md-whose-front-matter-names.md`.

Answer by number, such as `all as recommended` or `3: migration command`. Recommended first in each:

1. **Path.** `design/workflows/board.md`. Alternative: `design/workflows/workflow.md`, which avoids sharing a name with `wip/kanban/board.md`.
2. **WIP limit values.** They stay in `wip/kanban/board.md`, set by `flai board limit` and the lane menu without a story; the workflow file says only which states take one (`limit: true`). Alternative: the values move into the workflow file and change only through a story.
3. **Renamed or removed state.** Its old name goes in the target state's `aliases`; flai reads statuses and history through them and rewrites nothing; a removed state's time counts in the state it is read as. Alternatives: a migration command rewrites items at acceptance; or a removed state's time keeps its name and counts in no timer.
4. **ADR-0003 and ADR-0004.** ADR-0301 refines both and supersedes neither: only ADR-0004's closed list of states changes. Alternative: supersede ADR-0004 and restate what it keeps.
5. **Rule severity.** Each rule has a default, refuse or warn, that a transition may override. The baseline needs this: an unticked criterion warns at review and refuses at done. Alternative: fixed severity, with two rules for the criteria.
6. **Accepted and cancelled states.** Exactly one of each, and for stories exactly one `pull` and one `accept`. Alternative: several, which every reader of `done` must choose between.
7. **Who may move.** `by` takes `anyone`, `operator`, `agent`, `flai`, `acceptance`, or `cascade`; role files and permissions narrow it further for strategic roles. Alternative: the workflow names roles and permissions itself.
8. **flai upgrade.** It adds the baseline file to a project without one and never changes it afterwards. It says when the template's baseline moved. Alternative: merge it as a baseline role file is merged, which would undo a project's board.
9. **Who accepts a change to the file.** The operator alone, never the orchestrator or the merge queue, as for a protected path (ADR-0106). The dashboard saves it only on a story's branch. Alternative: any acceptance, as for other design files.

On your answer I amend ADR-0301, accept it, and record the outcome on ADR-0003 and ADR-0004 as choice 4 says. Then I describe the file in `workflow.md`, `work-hierarchy.md` § States, and `design/workflows/README.md`.

### 2026-11-14T17:50:19Z orchestrator (recommendation)
**Recommended:** all as recommended.

Each choice, and what it rests on:

1. **`design/workflows/board.md`.** The planned stories already name it: S-0755 touches `template/root/design/workflows/board.md`, and S-0737, S-0738, and S-0750 touch `flai/internal/workflow/baseline/board.md`. The alternative would mean changing their touches for a name clash the folder already makes clear.
2. **Limit values stay in `wip/kanban/board.md`.** The operator sets them today from the lane menu without a story. Moving them would make a limit change wait on an acceptance.
3. **Aliases, with nothing rewritten.** ADR-0004 keeps history in front matter so it can be reviewed in git. A migration command would rewrite many items on main at once.
4. **Refine ADR-0003 and ADR-0004, supersede neither.** Only ADR-0004's closed list of states changes. Its transitions in front matter, blocked as a flag, and computed timestamps all stand. S-0736's criterion 8 allows either outcome.
5. **A default severity per rule that a transition may override.** The baseline needs it: `criteria-ticked` warns at review and refuses at done.
6. **Exactly one accepted state, one cancelled state, one pull, and one acceptance.** Release accounting, the curator's trigger, and every reader of `done` need one answer each. Several workflows are out of the epic's scope.
7. **A short `by` vocabulary.** ADR-0207's role files and the operator's permissions already say what each role may do. Naming roles in the workflow file would say it a second time.
8. **`flai upgrade` adds the file once and never changes it.** Merging it would undo a project's own board.
9. **Only the operator accepts a change to the file.** E-0034's notes say the file is code the operator runs and no role may edit it. ADR-0106 sets the same gate for protected paths. This narrows my own `accept_reviews`, which is the safe direction.

With this answer, agent-S-0736 amends nothing, accepts ADR-0301, and records refinement on ADR-0003 and ADR-0004.

**Alternative:** answer by number to change any one, such as `1: workflow.md`.

Source: design/adrs/0004-workflow-states-and-transitions.md

### 2026-11-14T17:54:49Z alex
Confirmed the recommendation of 2026-11-14T17:50:19Z orchestrator.

Source: design/adrs/0004-workflow-states-and-transitions.md

### 2026-11-14T17:55:34Z agent-S-0736
Resolved: All nine as recommended. ADR-0301 is accepted on story/S-0736 (0fec9da). It refines ADR-0003 and ADR-0004, which stay unchanged.
