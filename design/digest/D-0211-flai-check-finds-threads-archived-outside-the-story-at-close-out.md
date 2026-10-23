---
id: D-0211
title: "flai check finds `threads.archived` outside the story at close-out"
kind: cause
sources: [I-0073]
topics: [cli, threads, flai, go, code, all]
paths: [flai/internal/threads/archived.go, flai/internal/threads/archived_test.go, flai/internal/check/check.go, flai/internal/check/check_test.go, flai/cmd/accept.go, flai/internal/preview/accept.go, flai/cmd/accept_threads_test.go, flai/cmd/archive.go, flai/cmd/archive_test.go, design/system/flai-cli.md, design/system/workflow.md, docs/users/flai.md, docs/users/flai-reference.md, design/adrs, design/issues/I-0073-flai-check-finds-threads-archived-outside-the-story-at-close-out.md, design/issues/summary.md, design/issues/I-0076-flai-check-finds-wip-overlap-outside-the-story-at-close-out.md, design/issues/I-0078-flai-check-finds-item-archive-outside-the-story-at-close-out.md]
date: 2026-10-07
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0211 flai check finds `threads.archived` outside the story at close-out

I-0073 was closed on 2026-10-07.

The cause:

```text
flai check finds `threads.archived` outside the story at close-out
```

The fix:

```text
S-0277: flai accept and flai archive resolve the threads still open or answered on what they archive (ADR-0109), tested in TestAcceptResolvesTheThreadsOnWhatItArchives, TestFinishingAnAcceptanceResolvesTheThreadsLeftOpen, and TestArchiveResolvesTheThreadsOnWhatItArchives; the threads.archived finding names flai thread resolve, and the two left open (TH-0203, TH-0240) are resolved
```

Source: I-0073.
