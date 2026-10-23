---
id: D-0261
title: TestRoundTripRepositoryItems reads the live main checkout and fails a close-out when another agent edits a story mid-run
kind: cause
sources: [I-0079]
topics: [flai, cli, go, code, all]
paths: [flai/internal/workitem/workitem_test.go, design/issues/I-0079-testroundtriprepositoryitems-reads-the-live-main-checkout-and-fails-a-close-out-when-another-agent-edits-a-story-mid-run.md, design/issues/summary.md]
date: 2026-10-08
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0261 TestRoundTripRepositoryItems reads the live main checkout and fails a close-out when another agent edits a story mid-run

I-0079 was closed on 2026-10-08.

The cause:

```text
TestRoundTripRepositoryItems reads the live main checkout and fails a close-out when another agent edits a story mid-run
```

The fix:

```text
Fixed by S-0290 (T-1367): TestRoundTripRepositoryItems now parses, validates, and compares one read of each listed item file, so an item another agent rewrites in the main checkout during a close-out is compared with itself; Save writes atomically, so one read is one whole version. A file gone since the list is skipped. TestRoundTripItemsRewrittenSinceTheList reproduces the race, and fails against the old two-read check.
```

Source: I-0079.
