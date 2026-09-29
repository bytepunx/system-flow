---
id: ADR-0055
title: "A story moves back one column from ready, in progress, review, or cancelled, and done is final; an item reopened from cancelled is not completed"
status: accepted
date: 2026-09-29
supersedes: []
superseded_by: []
refines: [ADR-0004]
topics: [cli, dashboard]
---

# ADR-0055 A story moves back one column from ready, in progress, review, or cancelled, and done is final; an item reopened from cancelled is not completed

## Context

The board's lanes get a right-click menu (S-0167) that moves stories forward and back. flai allowed one move back, review to in-progress, which sends work back for a reason. A story pulled too early could not go back to ready, a ready story that was not ready after all could not go back to backlog, and a cancelled story could not be brought back: each needed a hand edit of front matter that a command owns. Creation is backlog and is not recorded as a transition, so no transition to backlog was valid at all. `completed` was the first transition to done, else the first to cancelled, so an item once cancelled counted as completed for ever.

## Decision

Any item moves back one column from ready, in-progress, review, or cancelled, and done is final.

| From | Back to | Reason |
|------|---------|--------|
| ready | backlog | optional |
| in-progress | ready | optional |
| review | in-progress | required, as before |
| cancelled | backlog | optional; refused while the item's parent is cancelled, which is moved back first |

Done is final: acceptance merged and archived the work. A transition to backlog is valid anywhere but first, where creation already implies it. A story moved to ready goes last in the ready order, as any story entering ready does; a story moved to backlog leaves the order.

`completed` is the time of an item's last transition while it is done or cancelled, and nothing while it is open. An item reopened from cancelled has no lead time or cycle time until it closes again. `committed` and `started` stay the first transitions to ready and in-progress, so a story sent back keeps the time it first started.

## Consequences

- The metrics contract in `design/system/metrics.md` changes, which is why this is an ADR.
- A story moved back to ready is a story entering ready: `flai serve` may start an agent for it (ADR-0043). An agent still working the story finds its move to review refused, since ready does not go to review.
- Cancelling a story cancelled its open tasks (ADR-0028); moving it back does not reopen them. The agent that pulls it again writes or reopens the tasks it needs.
- An older flai rejects an item with a later transition to backlog as invalid, so `flai check` from an older installed flai, such as the one the MCP server runs, reports it until that flai is upgraded.

## Alternatives considered

- A single "reopen" that returns a cancelled item to the column it was cancelled from: it would skip the refinement that cancelled work needs before it is ready again.
- Moving back any number of columns in one step: the menu moves one column at a time, and so does flai, so each step is a transition with its own rules.
- Leaving `completed` as the first closing transition: a story reopened and delivered would report the cancellation as its completion.
