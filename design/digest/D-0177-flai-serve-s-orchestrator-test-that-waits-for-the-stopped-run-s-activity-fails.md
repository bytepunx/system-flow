---
id: D-0177
title: flai serve's orchestrator test that waits for the stopped run's activity fails under the close-out's full parallel test run
kind: cause
sources: [I-0090]
topics: [flai, cli, go, code, all]
paths: [design/issues/I-0074-two-stories-that-each-record-or-close-an-issue-always-conflict-in-design-issues-summary-md-whose-updated-line-both-rewrite.md, design/issues/I-0078-flai-check-finds-item-archive-outside-the-story-at-close-out.md, design/issues/I-0090-flai-serve-s-orchestrator-test-that-waits-for-the-stopped-run-s-activity-fails-under-the-close-out-s-full-parallel-test-run.md, design/issues/summary.md, flai/internal/serve/orchestrate_test.go]
date: 2026-10-06
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0177 flai serve's orchestrator test that waits for the stopped run's activity fails under the close-out's full parallel test run

I-0090 was closed on 2026-10-06.

The cause:

```text
flai serve's orchestrator test that waits for the stopped run's activity fails under the close-out's full parallel test run
```

The fix:

```text
S-0292: the orchestrator stub in flai/internal/serve/orchestrate_test.go wrote what it was given before it printed its stream, so under load the stop's TERM could land before its usage reached the run's log, and no activity was logged. The stub now prints its stream first, so given returning means the usage is in the log. With a 0.5 s delay before the stream the test failed in the old order with the same message and passes in the new one.
```

Source: I-0090.
