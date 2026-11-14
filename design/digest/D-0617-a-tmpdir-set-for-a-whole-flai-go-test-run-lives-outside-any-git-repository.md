---
id: D-0617
title: A TMPDIR set for a whole flai go test run lives outside any git repository, under /tmp from mktemp -d, never under the worktree's bin/, since about 150 tests fail when their temp folders sit inside a git repository and under a system-flow.yaml
kind: decision
sources: [S-0777, I-0371, S-0781]
topics: [flai, go, testing, code, cli]
paths: [flai, flai/cmd/cmd_test.go, flai/internal/mcpserver/server_test.go, flai/cmd/story_start_test.go]
date: 2026-11-13
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0617 A TMPDIR set for a whole flai go test run lives outside any git repository, under /tmp from mktemp -d, never under the worktree's bin/, since about 150 tests fail when their temp folders sit inside a git repository and under a system-flow.yaml

S-0777 ran flai's whole short suite with TMPDIR a symlink, to reproduce macOS's symlinked temp dir on Linux (D-0616). Its start prompt asks that files written outside the project go under the worktree's git-ignored `bin/S-0777/`. A first run with TMPDIR there was thrown away. From S-0777's Decisions (`wip/archive/agents/S-0777.md`):

```text
The symlinked TMPDIR for the whole-suite run lives under /tmp (mktemp -d), not under bin/S-0777/ as the start prompt asks for files written outside the project: a TMPDIR inside the worktree puts every test's temp folder inside a git repository and under a system-flow.yaml, and tests that expect neither (TestLoadAndFind, TestReadWithoutGitLeavesTheCommitsOutAndWarns, and about 150 more) fail for that alone. A first run there was discarded for it. go test writes its temp folders under /tmp in any ordinary run; the scratch folder is removed after the story's runs.
```

That discarded run also leaked branches into the main repository (I-0371, fixed by S-0781). A temp project with no `git init` of its own let git walk up into S-0777's worktree, so `flai story start` in the tests made `story/S-0001`, `S-0002`, `S-0003`, and `S-0005` in the main repository. The main checkout's `flai check` then warned `story.unaccepted` on them. I-0371's fix:

```text
the TestMain of flai/cmd and flai/internal/mcpserver set GIT_CEILING_DIRECTORIES to the temporary folder, with TestStoryStartInATemporaryFolderInsideARepositoryMakesNoBranch reproducing the leak (T-3343); the four branches were deleted and their worktrees pruned on the operator's word on TH-1190 (T-3345)
```

A later story that sets TMPDIR for a go test run makes it with `mktemp -d` under `/tmp`, removes it after, and records the exception to the `bin/<story>/` rule in its Decisions. Failures under a TMPDIR inside the worktree say nothing about the change. A new test package whose tests make temp projects without `git init` sets `GIT_CEILING_DIRECTORIES` to its temp folder in its TestMain, as flai/cmd and flai/internal/mcpserver do.

Sources: S-0777 (wip/archive/agents/S-0777.md, Decisions); I-0371 and S-0781 (wip/archive/agents/S-0781.md).
