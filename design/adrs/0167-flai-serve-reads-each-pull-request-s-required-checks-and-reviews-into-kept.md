---
id: ADR-0167
title: "flai serve reads each pull request's required checks and reviews into kept state under .flai-cache that flai show, item_get, and the dashboard read without the forge, and on a failing required check or a requested change opens one thread and sends the story back to in-progress, starting its agent"
status: accepted
date: 2026-10-17
supersedes: []
superseded_by: []
refines: [ADR-0055, ADR-0108, ADR-0166]
topics: [release, dashboard]
---

# ADR-0167 flai serve reads each pull request's required checks and reviews into kept state under .flai-cache that flai show, item_get, and the dashboard read without the forge, and on a failing required check or a requested change opens one thread and sends the story back to in-progress, starting its agent

## Context

[ADR-0166](0166-under-pull-request-acceptance-flai-serve-pushes-a-story-s-branch-and-opens-or.md) has flai serve push each story in review and open its one pull request. Nothing reads the pull request back. CI runs its checks on it, and reviewers review it, on the forge. The story's agent hears of a failed check or a requested change only when the operator relays it.

The forge token stays with flai serve on the host. An agent does not hold it. So an agent cannot ask the forge how its pull request stands, and neither can the dashboard without the host.

[ADR-0055](0055-a-story-moves-back-one-column-from-ready-in-progress-review-or-cancelled-and.md) lets a story move back from review to in-progress, with a reason. Until now only the operator made that move. [ADR-0108](0108-flai-serve-restarts-a-story-s-agent-that-ended-with-its-story-in-progress-up-to.md) has flai serve start a story's agent again after a failure, up to `agent.auto_restarts` times. flai serve also starts an agent again in its last session when a thread it waits on is answered.

## Decision

Under pull-request acceptance, flai serve reads each pull request's required checks and reviews into kept state under `.flai-cache`. Every reader reads that state, never the forge. A failing required check or a requested change opens one thread, sends the story back to in-progress, and starts its agent.

- **Reading.** Under `acceptance.mode: pull-request`, with the host action `forge` on, flai serve reads each story in review's pull request through the forge contract. It does this in the same once-a-minute forge pass that pushes stories in review. It reads the checks on the head commit, keeps the required ones, and reads the reviews.
- **Kept state.** flai keeps the latest read per story in `.flai-cache/forge/status/S-nnnn.json`, one file per story, as ADR-0166 keeps its push state. The file holds:
  - the pull request's number, URL, and head commit;
  - each required check's name, state, log link, and the time it entered that state;
  - each review's author, state, comment, and link;
  - an event for each change, with its time: a new head commit, a check that appeared or changed state, a review that is new or changed.

  The events are appended, so a reader can tell when the story waited on checks and when on reviewers. S-0396 reads them. The file never holds the token.
- **Readers read only the kept state.** `flai show S-nnnn` and its `--json`, the MCP `item_get`, and `flai forge status` read that file. So do the dashboard's story and review pages, through the host read method `forge.status` and the route `/api/items/[id]/forge`. None of them calls the forge. An agent without the token sees the same state, and no reader spends the forge's rate limit.
- **Send back.** A required check that failed at the head commit, or a review that requests changes, sends the story back:
  - flai serve opens one thread on the story, by `flai serve`. For a check, it names the check and links its log. For a review, it quotes the review's comment and links it.
  - flai then moves the story from review to in-progress, by `flai`, with a reason that names the same.

  This makes flai a new maker of the review to in-progress move. ADR-0055 allows the move; until now the operator alone made it. A pending check does nothing, and the story waits in review.
- **One send back per failure.** A check is keyed by its name and the head commit. The same check failing again on a new commit sends the story back again. A review is keyed by its URL. A review that requests changes stands across the fix's new commit until the reviewer reviews again, so it sends the story back once.
- **The agent is started.** flai serve then starts the story's agent again, in its last session, as it does when a thread the agent waits on is answered. This refines ADR-0108:
  - the prompt names the thread, and the agent's first `inbox` lists it;
  - the start does not count toward `agent.auto_restarts`, because it answers a thread and does not retry a failure;
  - the story keeps the place in in-progress it had.

  The agent fixes the branch, closes out, and moves the story to review again. flai serve pushes the fix under ADR-0166. flai serve resolves the send-back thread once the story is back in review at a different head commit.

## Consequences

- The story's agent hears from CI and from reviewers without the operator relaying.
- Every fix goes through the close-out and the move to review, as every other change does.
- The review column holds only stories that wait on someone else, and the review limit counts only those.
- A story can go back and forth between review and in-progress several times. Each round adds two transitions and one thread.
- The kept state is as fresh as the last forge pass. With `forge` off, or in `local` mode, flai serve does not update it.
- A reviewer who requests changes and does not review again keeps the thread's cause standing. flai sends the story back once for that review, not on every push.

## Alternatives considered

- **Leave the story in review while its agent changes the branch.** Rejected. A fix would skip the close-out and the move to review that every other change goes through. The review column would hold a story that is being worked, and the review limit would count it.
- **Have readers call the forge.** Rejected. Every agent and the dashboard would need the token, which ADR-0166 keeps out of every agent's environment. Each reader would also spend the forge's rate limit.
