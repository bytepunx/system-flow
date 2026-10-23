---
id: D-0324
title: A story asked to share a path agrees in a settled reply but records no share, so the asking story's task stays refused on the overlap
kind: cause
sources: [I-0168]
topics: [orchestration, cli, analysis, flai, go, flaiover, dashboard, sveltekit, code, all]
paths: [design/adrs, flai/internal/mergequeue/land.go, flai/internal/mergequeue/land_test.go, flai/internal/storygit/storygit.go, flai/internal/storygit/storygit_test.go, flai/internal/serve/mergequeue.go, flai/internal/serve/mergequeue_test.go, flai/cmd/queue.go, flai/cmd/queue_test.go, flaiover/src/lib/components/MergeQueue.svelte, flaiover/src/lib/components/MergeQueue.svelte.test.ts, design/system/workflow.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, docs/operators/index.md, docs/users/flai.md, design/issues/I-0157-the-host-s-installed-flai-is-again-behind-the-project-s-releases-after-i-0149-closed-so-the-mcp-task-done-refuses-every-close.md, design/issues/summary.md, design/issues/I-0168-a-story-asked-to-share-a-path-agrees-in-a-settled-reply-but-records-no-share-so-the-asking-story-s-task-stays-refused-on-the-overlap.md, flai/cmd/accept.go, flai/cmd/accept_batch_dirty.go, flai/cmd/accept_batch_dirty_test.go, flai/cmd/accept_markdown.go, flai/internal/mergequeue/eligible.go, flai/internal/mergequeue/progress.go, design/components/flaiover/src/lib/components/MergeQueue.examples.json, flai/cmd/board.go, flai/cmd/queue_held_test.go]
date: 2026-10-20
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0324 A story asked to share a path agrees in a settled reply but records no share, so the asking story's task stays refused on the overlap

I-0168 was closed on 2026-10-20.

The cause:

```text
A story asked to share a path agrees in a settled reply but records no share, so the asking story's task stays refused on the overlap
```

The fix:

```text
recorded in error by S-0402: S-0403's agent did try to share, and flai refused because a share lifts the hold on a ready story only, so no share can be recorded between two stories in progress; that gap is I-0169, recorded by S-0403
```

Source: I-0168.
