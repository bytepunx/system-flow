---
id: ADR-0272
title: "A checked creation holds a checked-write lock from its write until it keeps or removes the item, and flai serve's look reads the items it judges holds by under that lock"
status: accepted
date: 2026-11-09
supersedes: []
superseded_by: []
refines: [ADR-0213, ADR-0230]
---

# ADR-0272 A checked creation holds a checked-write lock from its write until it keeps or removes the item, and flai serve's look reads the items it judges holds by under that lock

## Context

A checked creation ([ADR-0213](0213-a-checked-creation-of-a-work-item-refuses-only-a-new-finding-inside-the-new.md)) writes the new item into its folder, runs the check with the item in place, and removes it when the check refuses it. For those few seconds any other reader sees an item that may never exist.

On 2026-11-07, I-0301 showed what that costs. At 08:18:43Z planner-S-0568 called `item_new` for a task of S-0568 with touches under `flai/internal/guard/`, and `itemnew.Create` wrote it as T-2550. flai serve's look at 08:18:46Z read it, so it judged S-0568's claim by that task's touches. S-0568 declared no touches of its own and three stories were in progress. Read without the task, it would have been held `no-touches`; read with it, no hold applied, and the look started agent-S-0568. The check refused T-2550 (no `## Notes`) and removed it by 08:18:48Z. The agent's `story_start` at 08:18:50Z read the items without it and refused the story as held. Judging the hold again just before the start would not have helped: a read a moment later still sees the task.

## Decision

A checked creation holds a checked-write lock from its write until it keeps or removes the item, and flai serve's look reads the items it judges holds by under that lock.

- `(*workitem.Repo).LockCheckedWrites` is an operating system lock on `.flai-cache/items-checked.lock`, waited on for at most `itemsLockWait`, 30 seconds. It is a lock apart from the items lock that `Save` takes, so that its holder may still save.
- `itemnew.Create` takes it just before it writes the new item and holds it through the check and any removal. It lets go once the check has kept the item, before `Then` runs, or once the item is removed. A lock it cannot take refuses the creation, and nothing is written.
- flai serve's look takes it to read the ready stories, the items, and the holds (`readyChecked`). Under it, it marks the item folders to be looked at again, so that a snapshot kept before a removal is not taken as fresh. The look's starts run after it lets go. A lock it cannot take is warned of and ends the look, as a failed read does.

## Consequences

- A look never judges a claim by an item that a checked creation may yet remove, so it starts no agent that `story_start` then refuses for that reason.
- A look waits at most as long as a creation's check runs, a few seconds. Two creations at once take turns.
- Once a minute the look stats every active item file. It parses only the files that changed.
- A checked edit of an existing item still writes, checks, and restores without this lock. No instance shows the window there, and it is recorded as an issue.
- Other readers, such as `inbox`, `wait_for_work`, and the dashboard, still read without the lock. They start no agent, and the next read corrects them.

## Alternatives considered

- **Judge the hold again just before the start.** This was the planner's first plan. It does not close the window: the refused item is still on disk moments later.
- **Hold the items lock across the creation.** `Save`, which the creation calls, takes the items lock. A second flock on the same file from the same process waits on itself until the timeout.
- **Check the new item before it is written, in memory.** The check reads files: the markdown rule and the narrative rules among them. Running them on an overlay is a much larger change.
- **Start a story only once two looks agree on its claim.** That holds every start back a minute, and it hides the window rather than closing it.
