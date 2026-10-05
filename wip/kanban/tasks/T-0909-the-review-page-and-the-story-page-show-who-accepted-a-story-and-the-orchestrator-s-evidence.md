---
id: T-0909
type: task
nature: feature
title: The review page and the story page show who accepted a story, and the orchestrator's evidence
status: backlog
parent: S-0221
owner: alex
created: 2026-10-05T04:48:07Z
updated: 2026-10-05T04:48:07Z
transitions: []
stream: S-0221
tags: [dashboard]
touches: [flaiover/src/lib/components/Review.svelte, flaiover/src/lib/components/Review.svelte.test.ts, flaiover/src/routes/items]
after: [T-0907]
---
# T-0909 The review page and the story page show who accepted a story, and the orchestrator's evidence

## Work

In `flaiover/src/lib/components/Review.svelte`, the message after an acceptance names who accepted: `S-nnnn is accepted by <by>`, read from the acceptance's `--json`.

On the story page, `flaiover/src/routes/items/[id]/+page.svelte`, a done story says who accepted it, read from its done transition's `by`. When that is the orchestrator, the page links to the evidence: the `### Accepted by the orchestrator` section T-0907 writes in its Notes.

The history list already shows each transition's `by`, so it stays as it is.

Add tests:

- In `Review.svelte.test.ts`, the accepted message names the acceptor.
- In `flaiover/src/routes/items/item.svelte.test.ts`, a story accepted by the orchestrator shows it.

This task waits for T-0907, which sets the transition's `by` and the `--json` this reads.

## Done when

- The review page and the story page name who accepted a story, and an orchestrator's acceptance shows its evidence.
- The new tests and `scripts/flaiover-test.sh` pass.

## Notes
