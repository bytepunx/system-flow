---
id: ADR-0275
title: "A checked edit of a work item holds the checked-write lock from its first save until it keeps or restores the edit, as a checked creation does"
status: accepted
date: 2026-11-10
supersedes: []
superseded_by: []
refines: [ADR-0272]
---

# ADR-0275 A checked edit of a work item holds the checked-write lock from its first save until it keeps or restores the edit, as a checked creation does

## Context

[ADR-0272](0272-a-checked-creation-holds-a-checked-write-lock-from-its-write-until-it-keeps-or.md) gave a checked creation the checked-write lock, `.flai-cache/items-checked.lock`, from its write until it keeps or removes the item. flai serve's look reads the items it judges holds by under that lock. Its Consequences left the checked edit out: no instance showed the window there, and it was recorded as I-0328.

A checked edit has the same window. `itemedit.Apply`, behind `flai edit`, `flai criteria`, and the MCP tools `item_edit` and `criteria_tick`, saves the item, its parents, a retitled story's narrative, and `wip/agents/index.md`. It then runs the check and, on a refusal, puts every file back. `itemedit.ApplyForecasts`, behind flai serve's replan, saves many stories with `SaveAll`, checks them, and puts each refused one back with `Restore`. Between the save and the restore, a look could read touches the check then took out, and start or hold a story by them.

## Decision

A checked edit holds the checked-write lock from its first save until it keeps or restores the edit, as a checked creation does.

- `itemedit.Apply` takes `(*workitem.Repo).LockCheckedWrites` after the check that runs before the write and just before its first save. It holds it through the parents' saves, the narrative and relink writes, the index write, the check after the write, and any restore. It lets go once the check has kept the edit and the kept file is read back, before the notice and the commit, or once the restore has put the files back. A lock it cannot take refuses the edit, and nothing is written.
- `itemedit.ApplyForecasts` takes it just before `SaveAll` and holds it through the index write, the check, and each `Restore` of a refused item. A lock it cannot take refuses the batch, and nothing is written. The replan warns and tries again on its next round.
- The lock order is the checked-write lock first. `Save`, `SaveAll`, and `Restore` take the items lock inside it. No caller of either function holds the checked-write lock already, and nothing inside the window takes it.

## Consequences

- A look never judges a claim by touches a refused edit wrote.
- A look waits at most as long as an edit's check runs, a few seconds. Two checked writes, creations or edits, take turns. The replan's batch waits for them too, and they for it.
- A lock not taken within `itemsLockWait`, 30 seconds, refuses the edit with nothing changed. The caller tries again.
- `flai touches` saves with no check and no restore, so it has no window. It takes no lock.
- Other readers, such as `inbox`, `wait_for_work`, and the dashboard, still read without the lock. They start no agent, and the next read corrects them.
- Every checked write of a work item now follows one rule.

## Alternatives considered

- **Hold the lock only around the check.** The window opens at the first save, not at the check. A look between the save and the check would still see the edit.
- **Check the edit in memory before it is written.** The check reads files, the markdown rule and the narrative rules among them. ADR-0272 rejected this for a creation as a much larger change, and the same holds for an edit.
- **Leave `ApplyForecasts` out.** A forecast changes no claim, so the look would judge no hold wrongly by it. But it is the same window, and leaving it out keeps two rules for checked writes where one will do.
