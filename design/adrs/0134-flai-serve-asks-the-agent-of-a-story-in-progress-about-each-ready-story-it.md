---
id: ADR-0134
title: "flai serve asks the agent of a story in progress about each ready story it holds on overlap alone, and that agent may share the overlapping paths with a split of the work, so the overlap no longer holds"
status: accepted
date: 2026-10-08
supersedes: []
superseded_by: []
refines: [ADR-0046, ADR-0096, ADR-0120, ADR-0121]
---

# ADR-0134 flai serve asks the agent of a story in progress about each ready story it holds on overlap alone, and that agent may share the overlapping paths with a split of the work, so the overlap no longer holds

## Context

A ready story whose claim overlaps the claim of a story in progress is held until that story moves to review, is cancelled, or is sent back ([ADR-0046](0046-a-ready-story-whose-claim-overlaps-an-open-story-s-is-held-yellow-and-with-its.md), [ADR-0096](0096-a-story-in-review-holds-nothing-an-overlap-inside-the-manifest-s-shared-paths.md)). Often the holding agent will not change the paths at all, or will change them in a section the held story does not touch. Only the holding agent knows. Today it is never asked: the held story waits out the holding story's whole time in progress. The held story has no agent until it starts, so it cannot ask for itself. Since [ADR-0120](0120-agents-of-two-open-stories-message-each-other-in-conversations-kept-under-wip.md) and [ADR-0121](0121-flai-tells-two-stories-of-a-trial-merge-conflict-or-a-grown-overlap-with-a.md) agents of two open stories talk in a conversation per pair, but only between stories in progress or in review.

The operator decided on TH-0318 and TH-0319, by finalizing S-0334 as written, that a share lifts the hold on its own, without their confirmation.

## Decision

**flai serve asks the agent of each story in progress that holds a ready story on overlap alone, in the pair's conversation, and that agent answers by narrowing its touches, by sharing the overlapping paths with a stated split of who changes what, or by saying why the hold stands. An overlap on shared paths does not hold.** This refines ADR-0046, ADR-0096, ADR-0120, and ADR-0121.

- **The ask.** At each look, for each ready story held on overlap alone (code `overlap`, not `after` or `no-touches`), flai serve messages each story in progress that holds it, by `flai`, from the held story to the holding one, `about` the overlapping paths, with the held story's goal and the three answers. It asks once per pair while the held story stays in ready: an ask already in the pair's open conversation since the held story last entered ready is not repeated. A message that cannot be written is logged and the look goes on.
- **A ready story in a conversation.** A conversation may now have a ready story on one side, for the ask and the share alone. Nobody answers for the held story until it starts: the conversation awaits the holding story, and a reply it writes awaits no one while the other side is ready, so it does not count as awaiting another story's reply. Sending, replying, and escalating are otherwise between stories in progress or in review, as before.
- **Narrowing.** The holding agent narrows its touches with `flai touches`; the hold clears as it does today.
- **Sharing.** `flai message share <conversation> --paths <path>… "<split>"`, and the MCP tool `message_share`, record in the conversation's front matter, under `shares`, the holding story, the held story, the paths, the split, who shared, and when, and add an entry with the split. Each path must lie inside the two stories' overlap. Only the holding story's agent, writing for the holding story, or the operator, the holding story's owner, may share. `flai guard` refuses a sub-agent the share, as it refuses every other message write.
- **What a share does.** An overlap between the held story and the holding one whose narrower path lies wholly inside a shared path does not hold, as an overlap inside the manifest's shared paths does not (ADR-0096). `flai board`, `inbox`, `wait_for_work`, `flai story start`, and the launcher read the shares through the same holds, so they agree. A hold that a share clears in part names only the paths, and the stories, that still hold. A share with one story clears nothing held by a third.
- **When a share ends.** A share is in force while its conversation is stored open, both stories are in ready, in progress, or in review, and neither has moved to backlog, done, or cancelled since it was made. So it ends when either story leaves the open columns, or when the held story goes back to backlog before it starts.
- **Starting on a share.** When flai serve starts a story whose hold a share cleared, the prompt it gives the agent names the conversation, the shared paths, and the split. The agent's first `inbox` lists the conversation under `messages`.
- **Sync and acceptance.** A share changes no claim. The trial merge at `flai stream sync` still compares the two branches, and a conflict between them adds to their conversation as ADR-0121 has it; the notice at acceptance still reports the overlap.

## Consequences

- A hold the holding agent judges needless ends while that agent is still working, not when it reaches review, so the held story starts sooner.
- A wrong split is caught where it would be without the share: at sync, by the trial merge, and at acceptance, by the notice.
- The holding agent is woken by the ask, as by any message to its story, and spends a turn answering it.
- Conversations under `wip/messages` carry a new front matter key, `shares`; an older flai keeps it as unknown front matter and ignores it, so it holds as before.

## Alternatives considered

- The share waits for the operator's confirmation: the operator decided against it on TH-0318 and TH-0319, since the trial merge and the acceptance notice still catch a split that does not hold.
- The held story's agent asks once it starts: the held story has no agent while it is held.
- Record the share on the stories' front matter: the conversation is where the two agents agree, and a share there ends with it.
- Add the shared paths to the manifest's `claims.shared`: that list is the operator's and holds for every story, not one pair.
