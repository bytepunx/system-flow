---
id: TH-0764
title: "S-0523 is in progress: which permission lets the orchestrator set cost of delay inputs on an epic or a started story (TH-0694 question 2)?"
anchor:
  path: wip/kanban/stories/S-0523-the-orchestrator-cannot-set-a-cost-of-delay-input-the-operator-decided-on-an-epic-or-a-story-past-the-backlog.md
  item: S-0523
status: open
participants: [agent-S-0523, orchestrator]
created: 2026-10-30T13:06:33Z
updated: 2026-10-30T13:07:33Z
---

# TH-0764 S-0523 is in progress: which permission lets the orchestrator set cost of delay inputs on an epic or a started story (TH-0694 question 2)?

On wip/kanban/stories/S-0523-the-orchestrator-cannot-set-a-cost-of-delay-input-the-operator-decided-on-an-epic-or-a-story-past-the-backlog.md.

## Entries

### 2026-10-30T13:06:33Z agent-S-0523
S-0523 has started. Most of its work does not depend on your answer, and I am doing that now. One choice is yours. TH-0694 asked it two days ago, and it is still open there.

On TH-0579 you wrote "you should be allowed to update the cost of delay". Two ways do that. Both stay inputs only: no value, no removal, nothing else changed in the same call, and never on a done, cancelled, or archived item.

1. **Recommended: a new permission, `orchestration.permissions.set_cost_of_delay`, off by default.** With it on, the orchestrator adds or changes the inputs of any open epic or story, including inputs you already set, but only to a figure you decided on a thread. `plan_backlog_stories` keeps what it does today: first inputs on a backlog story that has none. It adds a manifest key, a settings toggle, and rows in the docs. You turn it on here once S-0523 is accepted.
   - Why: changing inputs you set is a new power over money. The convention says today "never overwrite the operator's inputs". A switch of its own keeps that consent explicit for every project made from the template, and you can turn it off without losing story planning.
2. **Alternative: widen `plan_backlog_stories`**, the planner's option on TH-0694. No new key. The power arrives for every project that already has that permission on.

Until you answer, I build the parts both options share:
- tests that reproduce I-0187's refusals;
- flai's rule admitting an epic and a started story;
- the guard admitting an epic ID.

The permission gate is one place in each, so either answer costs little rework. With option 1, I add the key and write the ADR for it. With option 2, I write the ADR refining ADR-0119.

### 2026-10-30T13:07:33Z orchestrator (recommendation)
**Recommended for the operator:** option 1. Add a new permission, `orchestration.permissions.set_cost_of_delay`, off by default. Leave `plan_backlog_stories` as ADR-0119 § Decision item 4 has it: inputs on a backlog story only.

## Why this is yours to decide

- It is a decision not yet recorded. ADR-0119 limits my cost of delay writes to a backlog story. Either option widens that.
- It widens my own powers over money figures. I should not settle that myself.

## Why option 1

- **You turn it on yourself.** The new key is off by default. Confirming it grants me nothing until you set it in `system-flow.yaml`. Option 2 grants the power at once to every project that already has `plan_backlog_stories` on.
- **It keeps consent separate.** The convention says "never overwrite the operator's inputs". Option 1 lets me change inputs you set. A separate switch keeps that a choice you make on its own. You can also turn it off without losing story planning.
- **The cost is small.** It adds one manifest key, one settings toggle, and docs rows. The story's agent already builds the gate as one place.

## The same question on TH-0694

TH-0694 question 2 asks the same thing. My recommendation there, from 2026-10-28T14:34:37Z, is the same answer and still awaits you. One answer on either thread settles both.

The thread stays open for your answer.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md § Decision
