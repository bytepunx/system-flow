---
id: I-0407
title: flai stream sync panics while continuing a rebase whose stop the agent resolved, and leaves the rebase stopped at the next commit
class: defect
status: open
count: 1
cost: 5m
first_reported: 2026-11-18T13:08:06Z
last_reported: 2026-11-18T13:08:06Z
updated: 2026-11-18T13:08:06Z
---

# I-0407 flai stream sync panics while continuing a rebase whose stop the agent resolved, and leaves the rebase stopped at the next commit

## Description
flai stream sync panics while continuing a rebase whose stop the agent resolved, and leaves the rebase stopped at the next commit

## Instances

### 2026-11-18T13:08:06Z
Story: S-0809.
S-0809, sent back by TH-1315: the rebase of story/S-0809 onto main (after S-0807's acceptance) stopped on flai/internal/verify/select.go, story.go, and story_test.go. I took main's side with git checkout --ours, git add-ed the three files, and ran scripts/flai.sh stream sync S-0809 (the tree's flai, built from the branch). It exited with a Go panic: the goroutine trace ran through storygit.Sync (flai/internal/storygit/sync.go:549, its call to Rebase) from newStreamSyncCmd (flai/cmd/stream.go:302); the panic's own message scrolled out of the 15 lines I kept, so its first frame is not recorded. Git was then left mid-rebase, stopped on the next commit's conflicts in design/system/flai-cli.md and design/system/workflow.md, which git rebase had met after the continue. Running flai stream sync again reported those two paths as unresolved, and the sync after they were resolved finished cleanly. So the continue went on, but the panic hid the second stop's message: an agent that reads only the exit sees a crash, not the paths to resolve.

## Remediation
