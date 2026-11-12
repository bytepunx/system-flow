---
id: D-0288
title: A host crash during flai story start leaves the story's branch ref, worktree, and narrative as empty files, and no agent can repair the ref
kind: cause
sources: [I-0137, S-0375]
topics: [git, narrative, flai, cli, go, code, all]
paths: [flai/internal/atomicfile/atomicfile.go, flai/internal/atomicfile/atomicfile_test.go, flai/internal/workitem/narrative.go, flai/internal/workitem/narrative_test.go, flai/internal/storygit/storygit.go, flai/internal/storygit/open.go, flai/internal/storygit/open_test.go, flai/cmd/stream.go, flai/internal/storystart/start_test.go, design/system/flai-cli.md, design/system/agent-narrative.md, docs/users/flai.md, docs/operators/runbooks/restore.md, design/issues/I-0137-a-host-crash-during-flai-story-start-leaves-the-story-s-branch-ref-worktree-and-narrative-as-empty-files-and-no-agent-can-repair-the-ref.md, design/issues/summary.md, flai/cmd/stream_reopen_test.go]
date: 2026-10-13
status: compacted
absorbed_by: ""
asked_again: []
cited: []
---

# D-0288 A host crash during flai story start leaves the story's branch ref, worktree, and narrative as empty files, and no agent can repair the ref
