---
id: D-0359
title: A story's in-progress transition from flai story start was lost when a concurrent write rewrote the item file, leaving it ready while its agent worked it
kind: cause
sources: [I-0175]
topics: [flai, cli, go, code, all]
paths: [flai/internal/workitem/itemlock.go, flai/internal/workitem/itemlock_test.go, flai/internal/workitem/repo.go, flai/internal/itemedit/itemedit.go, flai/internal/itemedit/lock_test.go, design/system/work-hierarchy.md, design/adrs, design/issues/I-0175-a-story-s-in-progress-transition-from-flai-story-start-was-lost-when-a-concurrent-write-rewrote-the-item-file-leaving-it-ready-while-its-agent-worked-it.md, design/issues/summary.md, flai/cmd/accept_resume_test.go, flai/cmd/check_stats_test.go, flai/cmd/issue_test.go, flai/internal/workitem/item.go, design/issues/I-0176-a-usage-write-to-a-story-overwrites-its-move-to-in-progress-made-meanwhile-so-a-started-story-reads-as-ready-again.md, flai/internal/itemnew/itemnew_test.go, flai/internal/statsread/statsread_test.go, flai/internal/mcpserver/folder_test.go, design/issues/I-0186-testwaitforeventsacrossafolder-can-wake-on-the-new-story-s-creation-before-its-move-to-ready-and-finds-no-event.md]
date: 2026-10-21
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0359 A story's in-progress transition from flai story start was lost when a concurrent write rewrote the item file, leaving it ready while its agent worked it

I-0175 was closed on 2026-10-21.

The cause:

```text
A story's in-progress transition from flai story start was lost when a concurrent write rewrote the item file, leaving it ready while its agent worked it
```

The fix:

```text
S-0479 (ADR-0184): Repo.Save takes one items lock, reads the file again, and merges a copy read before another writer's save by front-matter key and body section, refusing a part both changed, so the replanner's forecast write through itemedit.Apply keeps a move made meanwhile; a refused edit takes out only its own change. TestSaveKeepsAMoveSavedSinceTheCopyWasRead and TestRestoreKeepsAMoveSavedAfterTheEdit reproduce it.
```

Source: I-0175.
