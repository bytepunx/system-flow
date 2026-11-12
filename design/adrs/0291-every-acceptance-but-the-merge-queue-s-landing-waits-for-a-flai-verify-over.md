---
id: ADR-0291
title: "Every acceptance but the merge queue's landing waits for a flai verify over main's new paths: it is refused while the paths main brought in since the story's last passing verify send back a tier that verify ran"
status: accepted
date: 2026-11-12
supersedes: []
superseded_by: []
refines: [ADR-0156, ADR-0221]
---

# ADR-0291 Every acceptance but the merge queue's landing waits for a flai verify over main's new paths: it is refused while the paths main brought in since the story's last passing verify send back a tier that verify ran

## Context

An acceptance runs no tier. `flai accept` rebases the story's branch onto main and merges it, trusting the story's last `flai verify`. That verify ran over main as it was then. When main has moved since, nothing checks the rebased head.

[I-0229](../issues/I-0229-two-stories-accepted-apart-leave-main-red-s-0535-reworded-to-backlog-s-condition-and-s-0416-s-testadefinedroleisheldtoitsfile-still-expects-the-old-words.md) records it:

- S-0416's close-out passed at 6c25aad7, over base 1dfa57c8.
- S-0535 then landed. It reworded the `to-backlog` condition's `Allows` in `flai/internal/roles`.
- S-0416 was accepted by hand. The merge queue had refused it for protected paths (TH-0722). `flai accept` rebased it onto S-0535's main and fast-forwarded.
- S-0416's `TestADefinedRoleIsHeldToItsFile`, in `flai/internal/guard`, still held the old words. Main was red. The next batches sent back S-0538 and S-0519, though neither touched the guard ([D-0413](../digest/D-0413-a-merge-queue-batch-that-fails-a-story-in-a-package-the-story-does-not-change.md)).

Tier selection was not the gap. `internal/guard` imports `internal/roles`, so a verify at the rebased head would have run the guard's tests again: [ADR-0156](0156-a-close-out-that-main-outran-syncs-and-verifies-again-at-most-three-rounds-and.md)'s sent-back tiers. The gap is the acceptance. Its only tier gate, `queue_tiers` ([ADR-0221](0221-while-the-merge-queue-can-batch-a-hand-acceptance-is-refused-until-the-queue.md)), applies only while the queue batches, and compares the report only to the branch's own head, never to how far main moved.

## Decision

Every acceptance but the merge queue's landing waits for a `flai verify` over main's new paths.

1. **Which acceptances.** The operator's, from `flai accept`, `flai move <S-nnnn> done`, the board, and the dashboard; and the orchestrator's. The queue on or off. The queue's landing (`--batch`) is not gated: its batch ran every tier over the main it lands on. While the queue batches, the operator's acceptance enqueues the story ([ADR-0225](0225-while-the-merge-queue-batches-and-is-not-held-the-operator-s-acceptance-of-a.md)), so the gate meets only what still lands by hand.
2. **Main's new paths.** The acceptance reads the story's last `flai verify` report. When it passed, main's new paths are those `git diff --name-only` lists from the report's merge base, or its base commit when it recorded none, to the main branch's head now. Paths under the `wip` folder that the branch does not change are left out, as ADR-0135 and ADR-0156 leave them out.
3. **The tiers they send back.** The paths select tiers as a story's own paths do (`SelectStory`, `{importers}` included). A tier the report ran, passed or reused, that they select is sent back. So is a tier the report ran in place of one they select, as the tier that covered it.
4. **The blocker.** One or more tiers sent back is a blocker, `verify_stale`. It names the tiers, how many paths main brought in and the first of them, and what to do: ask for the story's close-out again, with the MCP tool `verify` or `flai verify <S-nnnn> --job`. The close-out syncs the branch onto main and, by ADR-0156's reuse, runs again only the tiers main's new paths send back. `flai accept --dry-run`, the dashboard's preview, and every acceptance list it, and the acceptance refuses it before anything is merged. A report that cannot be read, or a main whose paths cannot be listed, is the blocker too, naming the error, so the gate never opens unread.
5. **Where it lives.** The blocker is among the acceptance's blockers, not among the orchestrator's own (ADR-0093). The orchestrator's enqueue checks only its own blockers, so enqueueing a story whose verify main has outrun is not refused: the batch runs every tier.
6. **No override.** No flag lands a story past it. Asking for the close-out again is cheap, and an override would bring back the hand acceptance that left main red.
7. **What passes.** No report, or a report that did not pass, adds no `verify_stale`: the other gates and the move to review judge those. A main that moved only `wip` paths, or only paths that select no tier the verify ran, stops nothing.

This refines ADR-0156, whose reuse now also bounds what an acceptance waits for, and ADR-0221, whose `queue_tiers` checks the tiers at the story's head but not main's moves.

## Consequences

- A story accepted after main moved under it lands only once its tiers ran over that main, whoever accepts it. A test that pins words another story changed, as in I-0229, is found before it reaches main.
- In this repository the integration and smoke tiers select `flai/**`. A story whose verify ran them is held after any acceptance that changes flai, until its close-out runs them again, about twelve minutes. An acceptance that changes only documents, or only `wip/`, holds nothing.
- The story's agent may have ended its session in review. The operator or the orchestrator asks for the close-out with the MCP tool `verify`, or the operator runs `flai verify <S-nnnn> --job` on the host; the story stays in review.
- The test that went red no longer copies the condition's words: `TestADefinedRoleIsHeldToItsFile` reads `Allows` from `internal/roles`.

## Alternatives considered

- **`flai accept` runs the sent-back tiers itself before it merges.** Rejected as ADR-0221 rejected it: twelve minutes is past the dashboard's ten-minute timeout on `accept.run` and `item.move`.
- **Gate only the operator's acceptance, as ADR-0221 does.** Rejected: with the queue off, the orchestrator's acceptance lands by hand exactly as S-0416 did.
- **An override flag, as `--without-queue-tiers`.** Rejected: confirmed with no override on TH-1044.
- **Block on any commit main gained outside `wip/`.** Rejected: an acceptance that changes only documents would hold every story in review, though no tier the verify ran reads them.
