---
id: D-0307
title: Two stories that each add an ADR conflict in design/adrs/README.md, whose index both append a row to, and flai opens a thread for each sync that finds it
kind: cause
sources: [I-0099]
topics: [flai, cli, go, code, all]
paths: [design/adrs/0157-flai-stream-sync-and-flai-accept-merge-design-adrs-readme-md-s-index-row-by-row.md, flai/internal/adrindex/adrindex.go, flai/internal/adrindex/adrindex_test.go, design/adrs/README.md, design/system/flai-cli.md, docs/users/flai.md, flai/cmd/stream_sync.go, flai/internal/storygit/sync.go, flai/internal/storygit/sync_test.go, design/issues/I-0099-two-stories-that-each-add-an-adr-conflict-in-design-adrs-readme-md-whose-index-both-append-a-row-to-and-flai-opens-a-thread-for-each-sync-that-finds-it.md, design/issues/summary.md, docs/users/flai-reference.md, flai/cmd/stream.go, flai/cmd/stream_sync_test.go, flai/internal/issues/generated.go, design/issues/I-0148-a-close-out-whose-run-main-outruns-takes-two-rounds-of-the-integration-and-smoke-tiers-about-24-minutes-past-the-bash-tool-s-10-minute-limit-so-the-agent-s-call-is-killed-mid-round.md]
date: 2026-10-14
status: absorbed
absorbed_by: "design/adrs/0157-flai-stream-sync-and-flai-accept-merge-design-adrs-readme-md-s-index-row-by-row.md#decision"
asked_again: []
cited: []
---

# D-0307 Two stories that each add an ADR conflict in design/adrs/README.md, whose index both append a row to, and flai opens a thread for each sync that finds it

I-0099 was closed on 2026-10-14.

The cause:

```text
`flai adr new` adds a row for the new ADR to `design/adrs/README.md`. Two stories in progress that each add an ADR therefore change the same lines of that file, and git cannot merge them. flai's trial merge at sync reports the pair as conflicting and opens a thread to the two agents, which says to narrow or to wait; neither applies, since every row belongs. Whichever story is accepted second resolves the index by hand at its sync or its rebase, keeping both rows.

Since ADR-0096 lets stories in review stop holding others and the shared paths include `design/adrs`, more stories run at once, and most add an ADR. The thread churn is the same as the issue summary's before S-0278.
```

The fix:

```text
Fixed by S-0306 (ADR-0157): flai stream sync, flai task done, and flai accept merge design/adrs/README.md row by row when a rebase stops on it, keeping both sides' rows in ADR number order, and the trial merge leaves the index out of a pair's conflicts, so two stories that each add an ADR neither stop each other's rebase nor get a message about it.
```

Source: I-0099.
