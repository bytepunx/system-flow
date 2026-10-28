---
id: D-0337
title: The merge queue's send-back thread quotes only the first line of a failing plain tier's finding and keeps no log of the batch, so the step that failed cannot be told
kind: cause
sources: [I-0172]
topics: [cli, flai, go, code, all]
paths: [flai/internal/verify/text.go, flai/internal/verify/text_test.go, flai/cmd/verify.go, flai/cmd/verify_test.go, flai/internal/verify/run.go, flai/internal/verify/run_test.go, flai/internal/verify/story.go, flai/internal/mergequeue/batch.go, flai/internal/mergequeue/batch_test.go, flai/internal/mergequeue/state.go, flai/internal/mergequeue/bisect.go, flai/internal/mergequeue/bisect_test.go, design/system/flai-cli.md, docs/users/flai.md, docs/operators/index.md, docs/operators/runbooks/backup.md, design/issues/I-0172-the-merge-queue-s-send-back-thread-quotes-only-the-first-line-of-a-failing-plain-tier-s-finding-and-keeps-no-log-of-the-batch-so-the-step-that-failed-cannot-be-told.md, design/issues/summary.md, flai/internal/verify/verify.go, flai/internal/mergequeue/land_test.go, flai/internal/serve/mergequeue_test.go]
date: 2026-10-21
status: absorbed
absorbed_by: "design/system/flai-cli.md#the-merge-queue-pass"
asked_again: []
cited: []
---

# D-0337 The merge queue's send-back thread quotes only the first line of a failing plain tier's finding and keeps no log of the batch, so the step that failed cannot be told

I-0172 was closed on 2026-10-21.

The cause:

```text
The merge queue's send-back thread quotes only the first line of a failing plain tier's finding and keeps no log of the batch, so the step that failed cannot be told
```

The fix:

```text
S-0457: the merge queue's send-back thread now quotes the failing tier's findings in full, as flai verify prints them (Result.findings, cut at 8,000 characters), and names the folder that keeps each batch's and bisection step's tier output under .flai-cache/merge-queue/batches/<batch id>/, the newest 20 kept; TestASentBackThreadQuotesAPlainTiersFailureInFull reproduces I-0172 with a plain smoke tier whose first line is its heading.
```

Source: I-0172.
