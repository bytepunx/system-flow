---
id: T-3073
type: task
nature: improvement
title: Record in an ADR that a verified batch lands over main's wip-only commits and that every drop is logged and counted
status: backlog
parent: S-0675
owner: alex
created: 2026-11-10T10:36:24Z
updated: 2026-11-10T10:36:24Z
transitions: []
stream: S-0675
tags: [adr, merge-queue]
touches: [design/adrs/README.md]
---
# T-3073 Record in an ADR that a verified batch lands over main's wip-only commits and that every drop is logged and counted

## Work

Write one ADR with `flai adr new` and accept it. It records two decisions, from the evidence under `### Planning` in S-0675's Notes:

1. A verified merge queue batch is not dropped when the main branch moved from the commit the batch expects (`expectedMain`: the head its landing left, or the commit it was built on) only by commits that descend from that commit and that change only paths under the manifest's `wip` folder (`layout.wip`) that no story of the batch changes. The batch lands on the new main without being verified again, as `flai verify`'s sync step passes over such commits (ADR-0135). The same rule holds between two landings of one batch. Any other move drops the batch, as today: a commit outside `wip/`, a `wip` path a story of the batch changes, or a main that does not descend from the expected commit.
2. Every drop appends a `dropped` line to `.flai-cache/merge-queue/events.jsonl`, with the batch, its stories, and the drop's reason. `flai stats` counts them per week as `queue.weeks[].dropped`. `metrics.md` changes only with an ADR, so this one carries it.

Give the evidence in Context: 37 of the 94 batches that passed between 2026-10-21 and 2026-11-02 never landed, and the drops sampled were flai's `chore: replan forecasts after reordered` commits, which change only `wip/`. Say in Consequences that one-story batches stay, since stories reach review hours apart and `merge_queue.coalesce` already joins those enqueued close together. Set `refines` to ADR-0173, ADR-0177, and ADR-0227, and give the ADR topics with `flai adr topics`.

This task waits for nothing. The other tasks build what it decides.

## Done when

- The ADR is accepted under `design/adrs/`, and `design/adrs/README.md` lists it.
- It states both decisions, the cases that still drop a batch, and the `dropped` event's fields.
- `flai check --strict` reports nothing on it.

## Notes
