---
id: T-1586
type: task
nature: improvement
title: Close I-0099 as fixed by the rebase's row merge of the ADR index
status: backlog
parent: S-0306
owner: alex
created: 2026-10-14T23:04:17Z
updated: 2026-10-14T23:04:17Z
transitions: []
stream: S-0306
tags: []
touches: [design/issues/I-0099-two-stories-that-each-add-an-adr-conflict-in-design-adrs-readme-md-whose-index-both-append-a-row-to-and-flai-opens-a-thread-for-each-sync-that-finds-it.md, design/issues/summary.md]
after: [T-1584, T-1585]
---
# T-1586 Close I-0099 as fixed by the rebase's row merge of the ADR index

## Work

Close I-0099 with `flai issue close I-0099 --reason --commit` from the worktree, naming the ADR and the merge that fixed it.

Waits for T-1584 and T-1585: the reason names the ADR and the behaviour they land.

## Done when

- I-0099 is closed and `summary.md` no longer lists it.

## Notes
