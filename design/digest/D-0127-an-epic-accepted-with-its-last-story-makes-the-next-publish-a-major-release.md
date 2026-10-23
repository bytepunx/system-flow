---
id: D-0127
title: "An epic accepted with its last story makes the next publish a major release: keep that: option 2 - we need to update documentation to reflect this change"
kind: answer
sources: [TH-0082]
topics: [orchestration, planning, analysis, flai, cli, go, flaiover, dashboard, sveltekit, template, conventions, code, all]
paths: [flai/internal/workitem, flai/internal/itemedit, flai/internal/release, flai/cmd/release.go, flai/cmd/move.go, flai/cmd/move_cancel.go, flai/cmd/accept.go, flai/cmd/edit.go, flai/cmd/check.go, flai/cmd/accept_epic_test.go, flai/cmd/release_pending_test.go, flai/cmd/accept_research_test.go, flai/cmd/accept_test.go, flai/cmd/stream_diff_test.go, flai/cmd/workitems_test.go, flai/cmd/check_stats_test.go, flai/internal/mcpserver, design/system/workflow.md, flai/cmd/edit_test.go, flai/cmd/hostapi_reads_test.go, flai/cmd/move_cancel_test.go, flai/internal/preview, flai/internal/check, design/system/work-hierarchy.md, design/system/flai-cli.md, design/adrs, docs/users, design/conventions/git.md, template/root/design/conventions/git.md, template/template.yaml, template/CHANGELOG.md, design/tech/ci.md, design/issues/summary.md, design/issues/I-0057-flai-check-strict-stops-a-story-s-close-out-on-wip-findings-outside-the-story.md, design/issues/I-0062-flai-adr-new-numbers-from-the-story-s-worktree-only-so-parallel-story-branches-take-the-same-adr-number.md]
date: 2026-10-03
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0127 An epic accepted with its last story makes the next publish a major release: keep that: option 2 - we need to update documentation to reflect this change

The operator, alex, answered TH-0082 at 2026-10-03T17:51:17Z.

What was asked:

```text
S-0200 makes an epic go to done, and be archived, in the acceptance of its last open story, as its criteria say. An accepted epic contributes a major bump to the next publish (git.md; release.LevelFor returns Major for any epic). No epic has been accepted in this repository yet (wip/archive/epics is empty, flai is at 1.29.1), so once this ships, accepting the last story of E-0016, or of any epic, makes the next `flai release --pending` take flai to 2.0.0, and the dashboard and template too where the epic's stories touched them.
```

The answer:

```text
option 2 - we need to update documentation to reflect this change
```

Source: TH-0082.
