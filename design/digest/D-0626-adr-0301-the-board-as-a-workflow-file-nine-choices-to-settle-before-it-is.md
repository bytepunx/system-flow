---
id: D-0626
title: "ADR-0301 (the board as a workflow file): nine choices to settle before it is accepted: Recommended: all as recommended"
kind: answer
sources: [TH-1198, S-0736]
topics: [cli, template, conventions, flai, go, code, all]
paths: [design/adrs/README.md, design/adrs/0003-work-item-hierarchy.md, design/adrs/0004-workflow-states-and-transitions.md, design/system/workflow.md, design/system/work-hierarchy.md, design/workflows/README.md, design/adrs/0301-the-board-is-a-workflow-file-design-workflows-board-md-whose-front-matter-names.md]
date: 2026-11-14
status: absorbed
absorbed_by: "design/adrs/0301-the-board-is-a-workflow-file-design-workflows-board-md-whose-front-matter-names.md#choices-for-the-operator"
asked_again: []
cited: []
---

# D-0626 ADR-0301 (the board as a workflow file): nine choices to settle before it is accepted: Recommended: all as recommended

The operator, alex, answered TH-1198 at 2026-11-14T17:54:49Z.

What was asked:

```text
ADR-0301 is written, proposed, on `story/S-0736`. It fixes the workflow file, `design/workflows/board.md`: states with a kind each (`queue`, `active`, `waiting`, `accepted`, `cancelled`), transitions with `by`, `rules`, and `effects` (`commit`, `pull`, `accept`, `send-back`, `cancel`, `reopen`), a 14-rule vocabulary, timers (between two points, or a sum of states), aliases for renamed or removed states, and a baseline that is today's board. The baseline file is in its §7.
```

The answer:

```text
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
```

Source: TH-1198.
