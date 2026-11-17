---
id: T-3499
type: task
nature: improvement
title: Close I-0387, saying the section merge now settles a line two stories change in different words
status: backlog
parent: S-0817
owner: alex
created: 2026-11-17T02:32:00Z
updated: 2026-11-17T02:32:00Z
transitions: []
stream: S-0817
tags: [issues]
touches: [design/issues/I-0387-stories-of-one-epic-rewrite-the-same-section-of-one-design-document-so.md, design/issues/summary.md]
after: [T-3498, T-3497]
---
# T-3499 Close I-0387, saying the section merge now settles a line two stories change in different words

## Work

Run `flai issue close I-0387 --reason "<reason>"`. The reason names the ADR T-3495 wrote, the merge T-3496 built, and the reproductions T-3498 added. It says the first instance, E-0033's stories on `design/system/server-performance.md`, has not stopped an acceptance since ADR-0287, and that the word merge settles the second, S-0797's table row and sentence.

Then tick both criteria of S-0817 with `flai criteria tick`, once each is verified.

Waits for T-3498 and T-3497: the reason cites the tests and the documents they write.

## Done when

- I-0387 is closed with that reason, and `design/issues/summary.md` no longer lists it.
- `flai check --strict` is clean.

## Notes
