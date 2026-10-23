---
id: D-0288
title: A host crash during flai story start leaves the story's branch ref, worktree, and narrative as empty files, and no agent can repair the ref
kind: cause
sources: [I-0137]
topics: [git, narrative, flai, cli, go, code, all]
paths: [flai/internal/atomicfile/atomicfile.go, flai/internal/atomicfile/atomicfile_test.go, flai/internal/workitem/narrative.go, flai/internal/workitem/narrative_test.go, flai/internal/storygit/storygit.go, flai/internal/storygit/open.go, flai/internal/storygit/open_test.go, flai/cmd/stream.go, flai/internal/storystart/start_test.go, design/system/flai-cli.md, design/system/agent-narrative.md, docs/users/flai.md, docs/operators/runbooks/restore.md, design/issues/I-0137-a-host-crash-during-flai-story-start-leaves-the-story-s-branch-ref-worktree-and-narrative-as-empty-files-and-no-agent-can-repair-the-ref.md, design/issues/summary.md, flai/cmd/stream_reopen_test.go]
date: 2026-10-13
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0288 A host crash during flai story start leaves the story's branch ref, worktree, and narrative as empty files, and no agent can repair the ref

I-0137 was closed on 2026-10-13.

The cause:

```text
A host crash during flai story start leaves the story's branch ref, worktree, and narrative as empty files, and no agent can repair the ref
```

The fix:

```text
S-0375 fixed it. atomicfile.WriteFile syncs the file before the rename and the folder after it, so a host crash leaves the old content or the new, never an empty narrative (T-1500). flai stream open and flai story start write an empty or blank narrative afresh from the template (T-1501). storygit.Open removes the story's own branch ref when it resolves to no commit, which git refuses to, and makes a worktree folder git does not know again when its branch is missing or broken or its files are all empty, refusing any other; the branch is made with core.fsync=reference (T-1502). TestStartRecoversWhatACrashedStartLeft reproduces S-0337's crash state and passes (T-1503).
```

Source: I-0137.
