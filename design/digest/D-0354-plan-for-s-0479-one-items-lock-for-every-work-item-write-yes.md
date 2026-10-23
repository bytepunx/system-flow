---
id: D-0354
title: "Plan for S-0479: one items lock for every work item write: yes"
kind: recommendation
sources: [TH-0577]
topics: [flai, cli, go, code, all]
paths: [flai/internal/workitem/itemlock.go, flai/internal/workitem/itemlock_test.go, flai/internal/workitem/repo.go, flai/internal/itemedit/itemedit.go, flai/internal/itemedit/lock_test.go, design/system/work-hierarchy.md, design/adrs, design/issues/I-0175-a-story-s-in-progress-transition-from-flai-story-start-was-lost-when-a-concurrent-write-rewrote-the-item-file-leaving-it-ready-while-its-agent-worked-it.md, design/issues/summary.md, flai/cmd/accept_resume_test.go, flai/cmd/check_stats_test.go, flai/cmd/issue_test.go, flai/internal/workitem/item.go, design/issues/I-0176-a-usage-write-to-a-story-overwrites-its-move-to-in-progress-made-meanwhile-so-a-started-story-reads-as-ready-again.md, flai/internal/itemnew/itemnew_test.go, flai/internal/statsread/statsread_test.go, flai/internal/mcpserver/folder_test.go, design/issues/I-0186-testwaitforeventsacrossafolder-can-wake-on-the-new-story-s-creation-before-its-move-to-ready-and-finds-no-event.md]
date: 2026-10-21
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0354 Plan for S-0479: one items lock for every work item write: yes

The operator, alex, answered TH-0577 at 2026-10-21T12:42:12Z.

What was asked:

```text
The plan has four tasks in three layers. The story's agent started S-0479 at 12:09:40Z, while it was being planned.
```

The recommendation it overruled:

```text
Recommended for the operator: yes to the merge.

- Add to S-0479 the criterion "I-0176 is closed with `flai issue close I-0176 --reason` saying what fixed it", with its close step in T-2005.
- Cancel S-0480 as a duplicate.

Both issues are the same lost update: a stale whole-file write from another process puts an older `status` back. I-0175 hit S-0479's own start, and I-0176 hit S-0456's move. T-2001 and T-2002 put every writer under the lock, which fixes both. Two stories would build the lock twice, and they would hold each other on the same files.

The planner says the input would rise to 15m, or 37.50 USD a week. That too is yours, since S-0479 is in progress and my permission covers backlog stories only. If you say no, both stories stay as they are.

Source: design/conventions/strategic-agents.md § As the planner
```

The operator's answer:

```text
yes
```

Source: TH-0577.
