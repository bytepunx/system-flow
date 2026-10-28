---
id: D-0376
title: Folder-wide touches and after hubs on planned stories held up to ten ready stories a day while in-progress stayed at two of five
kind: cause
sources: [I-0152]
topics: [planning, flai, cli, go, template, conventions, code, all]
paths: [flai/internal/workitem/create.go, flai/internal/workitem/create_test.go, flai/internal/workitem/hold.go, flai/internal/workitem/hold_test.go, flai/internal/itemedit/itemedit.go, flai/internal/itemedit/planning_test.go, flai/internal/itemedit/folder_test.go, flai/internal/itemnew/itemnew.go, flai/internal/itemnew/itemnew_test.go, flai/cmd/touches.go, flai/cmd/touches_test.go, flai/internal/mcpserver/items_write.go, flai/internal/mcpserver/items_write_test.go, flai/internal/check/check.go, flai/internal/check/check_test.go, flai/internal/check/scope.go, flai/internal/check/scope_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, design/adrs, design/adrs/README.md, design/system/strategic-agents.md, design/system/workflow.md, design/system/flai-cli.md, design/conventions/strategic-agents.md, design/conventions/work-management.md, template/root/design/conventions/strategic-agents.md, template/root/design/conventions/work-management.md, ".claude/agents/planner.md", template/root/.claude/agents/planner.md, template/template.yaml, template/CHANGELOG.md, docs/users/flai.md, docs/users/flai-reference.md, design/issues/I-0152-folder-wide-touches-and-after-hubs-on-planned-stories-held-up-to-ten-ready-stories-a-day-while-in-progress-stayed-at-two-of-five.md, design/issues/summary.md, flai/cmd/edit.go, flai/cmd/check.go, flai/cmd/check_stats_test.go, design/issues/I-0197-testservesettingsareshowntothedashboard-in-flai-cmd-fails-now-and-then-in-the-close-out-s-go-test-tier-on-a-branch-that-does-not-touch-serve.md]
date: 2026-10-23
status: absorbed
absorbed_by: "design/adrs/0188-a-task-takes-no-folder-touch-the-planner-s-folder-touch-on-a-story-needs-the.md#decision"
asked_again: []
cited: []
---

# D-0376 Folder-wide touches and after hubs on planned stories held up to ten ready stories a day while in-progress stayed at two of five

I-0152 was closed on 2026-10-23.

The cause:

```text
Folder-wide touches and after hubs on planned stories held up to ten ready stories a day while in-progress stayed at two of five
```

The fix:

```text
Fixed by S-0468 (ADR-0188). flai refuses a task a folder touch through flai task new, flai edit --touches, flai touches, item_new, and item_edit, so ADR-0096's narrowing always has files to narrow to. It refuses the planner's folder touch on a story unless the story's ### Planning note names the folder. flai check gives the advisory warnings touches.folder, on a story in progress whose claim keeps a folder, and wip.after-hub, on a story three or more open stories name in after. The planner's convention, planPrompt, and planner.md say that after names only a story whose accepted change is needed on main, never an order of preference.
```

Source: I-0152.
