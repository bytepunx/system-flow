---
id: ADR-0121
title: "flai tells two stories of a trial-merge conflict or a grown overlap with a message in their conversation, and either escalates to the operator on a thread only when they do not agree"
status: accepted
date: 2026-10-07
supersedes: []
superseded_by: []
refines: [ADR-0046, ADR-0120]
---

# ADR-0121 flai tells two stories of a trial-merge conflict or a grown overlap with a message in their conversation, and either escalates to the operator on a thread only when they do not agree

## Context

flai finds two kinds of collision between open stories on its own ([ADR-0046](0046-a-ready-story-whose-claim-overlaps-an-open-story-s-is-held-yellow-and-with-its.md)). `flai stream sync` trial-merges the story's branch with every other open story's and, on a conflict, opens a thread per pair on the syncing story. A write that grows a story's claim until it overlaps another in-progress story's (S-0244) records an `overlapped` change for both, which tells each agent to coordinate on a thread. Both reach the operator: a thread awaits the designer in `inbox` and on the dashboard, and is mirrored into the narrative's `## Open questions`. The operator has nothing to decide until the two agents disagree about who changes what. Since [ADR-0120](0120-agents-of-two-open-stories-message-each-other-in-conversations-kept-under-wip.md) the two agents have a channel of their own, and since S-0333 flai keeps one conversation per pair for the notices it sends when a task closes (`messages.Notify`).

## Decision

**flai tells the two stories of a trial-merge conflict or a grown overlap with a message in the pair's conversation, not a thread, and either story's agent escalates to the operator on a thread only when the two do not agree.** This refines ADR-0046 and ADR-0120.

- **Conflict at sync.** A trial merge that conflicts with another open story's branch adds a message from the syncing story to the other, by `flai`, `about` the conflicting paths, naming both branches and the paths, to the pair's open conversation, or starts one when none is open (`messages.Notify`). The conversation then awaits the other story. A later sync that finds the same paths adds nothing; one that finds other paths adds a message.
- **Clean merge.** A later sync that finds the two merging cleanly closes the pair's open conversation when the newest conflict message flai wrote in it is not yet answered by a clean merge, with the reason that they merge cleanly. A story sent back out of in progress or review closes it the same way at the next sync of the other; an accepted, cancelled, or archived story's conversations close as ADR-0120 says.
- **Grown overlap.** A write that grows a story's claim into another in-progress story's adds a message from the story whose claim grew to the other, by the writer, `about` the paths gained, to the pair's conversation, saying the claims now overlap and asking which of them changes those paths first. The `overlapped` change both stories get stays, so the dashboard and an older flai still see it. An overlap inside the shared paths tells nothing (S-0295). A message that cannot be written is logged, and the write that grew the claim stands.
- **Escalation.** `flai message escalate <conversation> "<reason>"`, and the MCP tool `message_escalate`, run by the agent of one of the conversation's two stories, open a thread on that story that names both stories, links the conversation, and quotes the reason: what they could not agree. The conversation gets a message from that story naming the thread, and stays open. A third story, a closed conversation, and an empty reason are refused. `flai guard` refuses a sub-agent the escalation, as it refuses the other message writes.
- **Old threads.** A conflict thread already open is still resolved by the sync that finds the pair merging cleanly or the other story gone, as before. No sync opens a conflict thread or adds to one.
- **Output.** `flai stream sync` names each conflicting pair's conversation, and `--json` gives it as `conversation` on the branch; `thread` names only an old conflict thread that sync resolved.

## Consequences

- The agents settle a conflict or an overlap between themselves while both are working, and the operator hears of it only through an escalation, which says what the two could not agree.
- An agent learns of a conflict another story's sync found from `inbox` `messages` and from `wait_for_events`, which wakes on the message, rather than from a thread on the other story.
- A conversation can now hold the agents' messages, task-close notices, conflict messages, and overlap messages for one pair, and a clean merge may close it while the agents still discuss something else; a new message then starts another.
- The conflict thread's text moves into the message; a dashboard that lists conflict threads lists only the old ones.

## Alternatives considered

- Keep the conflict thread and add the message beside it: the operator would still be asked about every conflict.
- One conversation per kind of notice: the agents would read the same pair's collision in several places.
- Escalate automatically when a conflict survives a number of syncs: a count cannot tell disagreement from work in progress.
