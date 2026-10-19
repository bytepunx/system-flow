---
id: ADR-0176
title: "The merge queue lands each story of a passing batch through flai accept --batch, which takes the batch's verification of the story's rebased head in place of the story head's own, and flai serve sends back to in-progress the story a red batch is bisected to"
status: accepted
date: 2026-10-19
supersedes: []
superseded_by: []
refines: [ADR-0032, ADR-0093, ADR-0173]
---

# ADR-0176 The merge queue lands each story of a passing batch through flai accept --batch, which takes the batch's verification of the story's rebased head in place of the story head's own, and flai serve sends back to in-progress the story a red batch is bisected to

## Context

[ADR-0173](0173-a-merge-queue-in-flai-serve-with-no-model-runs-behind-the-merge-queue-host.md) put a merge queue in `flai serve`. It builds a batch of the enqueued stories on the main branch and verifies the batch once. S-0399 built and verified the batch. S-0400 ends a verified batch, with no person in between, in acceptances or in a send-back.

Two rules stood in the way:

- [ADR-0093](0093-with-accept-reviews-on-the-orchestrator-accepts-a-story-in-review-through-flai.md) accepts a story as the orchestrator only when the verifier passed at the story branch's head (`--verified` equal to that head). A batch verifies a copy of each story rebased onto the batch, which is a different commit.
- [ADR-0032](0032-accepting-a-story-merges-it-publishing-is-a-deliberate-batched-step-over.md) and the workflow let only a person move a story from review back to in-progress.

The acceptance steps are unexported in `package cmd`, and `flai serve` cannot import that package.

## Decision

The queue lands each story of a passing batch with `flai accept --batch`, which takes the batch's verification of the story's rebased head in place of the story head's own, and `flai serve` sends back to in-progress the story a red batch is bisected to.

- **Landing.** For a batch that passed, `flai serve` accepts each of its stories in queue order. Each runs as `flai accept <story> --by <whoever enqueued it> --batch <batch ID>`, a process of this flai's own binary in the main checkout. That is the acceptance of today: merge, done, archive, commit. So each story's transitions, acceptance commit, archive, and `flai stats` figures are those of an acceptance by hand.
- **The verification it takes.** The verified commit is the story's rebased head inside the batch that passed. The preview blocks the acceptance unless all of these hold:
  - the queue's batch is the one named, is not dropped, and passed, or its bisection lets the story land;
  - the story is the next of the batch to land;
  - the story's branch is still at the head the batch rebased a copy of;
  - the main branch is where the landing left it: the batch's base, then the head the last landing left.
  Each acceptance rebases the story's branch onto main as the batch did. The acceptance commits between landings change only `wip/` paths, so what lands is what the batch verified.
- **The orchestrator.** For a story the orchestrator enqueued, the landing is the orchestrator's acceptance. It takes `--batch` in place of `--verified` and needs no `--evidence`: the queue's eligibility has checked the story's verify, its criteria, and its threads. Its notes get `### Accepted by the orchestrator`, naming the batch's verified commit, with the evidence when one is given. S-0403 makes the orchestrator enqueue with its evidence.
- **The queue's record.** `flai accept --batch` records the story as landed in the queue's state and leaves the batch in place. `flai accept` without it is an acceptance by hand: it takes the story out of the queue and drops the batch that holds it.
- **Bisection.** For a batch that failed, the queue verifies main with the first half of its stories, each prefix checked out at its last story's rebased head in the batch's worktree, and halves again. It stops at the first story with which the batch fails. Main itself counts as passing. Each step is kept in the queue's state and shown. The stories before the failing one land as above. The stories after it stay enqueued for the next batch.
- **Send-back.** `flai serve` moves the failing story from review to in-progress, by `flai`. The reason names the failing tier and its first findings, and the move appends it to the story's notes, as every send-back does. It opens a thread on the story to its agent, saying what failed and at which commit, and takes the story out of the queue. The story's agent is started again as for any story back in progress. This makes `flai serve` a second maker of that move, for the queue's send-back only, beside the forge's send-back ([ADR-0167](0167-flai-serve-reads-each-pull-request-s-required-checks-and-reviews-into-kept.md)).
- **Main moving.** When main moves other than by the batch's own landing, as by an acceptance by hand, the queue drops the batch and builds it again from the new main.

## Consequences

- An acceptance by the queue reads as one by hand, so the review page, the board, and `flai stats` need nothing new.
- Each landing starts a process, so a batch of four lands in four acceptances, one after another.
- A landing that fails stops the batch. The queue drops the batch, and the stories not yet landed are batched again from the main branch as it is.
- `flai serve` moves a story back to in-progress without a person. The send-back's reason, its thread, and the queue's outcome say why.

## Alternatives considered

- **Move the acceptance into an internal package that `flai serve` calls.** Rejected for now: it moves about a thousand lines. A subprocess keeps one acceptance flow and needs no second.
- **Fast-forward main to the batch's head and accept each story afterwards.** Rejected: each story's own acceptance commit would no longer follow its merge, so the history would not read as an acceptance by hand.
- **Rebuild each bisection prefix from main.** Rejected: the batch already holds each prefix's rebased head, and a rebuild could land a different tree from the one verified.
- **Send back every story of a red batch.** Rejected: one failure would cost every other story a cycle.
