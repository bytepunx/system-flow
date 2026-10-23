---
id: D-0366
title: TestWaitForEventsAcrossAFolder can wake on the new story's creation before its move to ready, and finds no event
kind: cause
sources: [I-0186]
topics: [flai, cli, go, code, all]
paths: [flai/internal/workitem/itemlock.go, flai/internal/workitem/itemlock_test.go, flai/internal/workitem/repo.go, flai/internal/itemedit/itemedit.go, flai/internal/itemedit/lock_test.go, design/system/work-hierarchy.md, design/adrs, design/issues/I-0175-a-story-s-in-progress-transition-from-flai-story-start-was-lost-when-a-concurrent-write-rewrote-the-item-file-leaving-it-ready-while-its-agent-worked-it.md, design/issues/summary.md, flai/cmd/accept_resume_test.go, flai/cmd/check_stats_test.go, flai/cmd/issue_test.go, flai/internal/workitem/item.go, design/issues/I-0176-a-usage-write-to-a-story-overwrites-its-move-to-in-progress-made-meanwhile-so-a-started-story-reads-as-ready-again.md, flai/internal/itemnew/itemnew_test.go, flai/internal/statsread/statsread_test.go, flai/internal/mcpserver/folder_test.go, design/issues/I-0186-testwaitforeventsacrossafolder-can-wake-on-the-new-story-s-creation-before-its-move-to-ready-and-finds-no-event.md]
date: 2026-10-21
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0366 TestWaitForEventsAcrossAFolder can wake on the new story's creation before its move to ready, and finds no event

I-0186 was closed on 2026-10-21.

The cause:

```text
TestWaitForEventsAcrossAFolder can wake on the new story's creation before its move to ready, and finds no event
```

The fix:

```text
S-0479: TestWaitForEventsAcrossAFolder writes the story before the wait (backlogStoryIn) and moves it to ready only during the wait, so the first change the wait sees is the move, which is an event.
```

Source: I-0186.
