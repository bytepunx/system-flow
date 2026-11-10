---
id: D-0547
title: "TestAPassTellsReleaseOfAMovedTagOrBranch in flai/internal/serve/mergequeue_test.go failed once under load while six task splits ran their tests at once, and passed when run again; no issue records it"
kind: cause
sources: [S-0631, T-2797]
topics: [testing, go]
paths: [flai/internal/serve/mergequeue_test.go]
date: 2026-11-10
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0547 TestAPassTellsReleaseOfAMovedTagOrBranch in flai/internal/serve/mergequeue_test.go failed once under load while six task splits ran their tests at once, and passed when run again; no issue records it

From S-0631's Decisions (`wip/archive/agents/S-0631.md`):

```text
T-2797's first `flai test` failed once on `TestAPassTellsReleaseOfAMovedTagOrBranch` in `serve/mergequeue_test.go`, a package T-2797 did not touch, while the six splits ran their tests at once. It passed when run again, and so did every later run with the splits together. Load made it flaky; the split did not cause it.
```

No issue was opened for it. So a later story whose `flai test` or close-out fails on this test, in a package it does not change, on a loaded host, runs it again before it investigates, and records an issue if it fails a second time, with this entry as the first instance.

Sources: S-0631 (T-2797).
