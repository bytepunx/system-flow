---
id: D-0645
title: "S-0739 cannot be enqueued: its branch is 37 commits behind main, and the close-out stops at sync: done"
kind: answer
sources: [TH-1226, TH-1227, TH-1228]
topics: [orchestration, acceptance, verify, git, flai, all]
paths: [flai/internal/preview/accept.go, flai/cmd/accept.go]
date: 2026-11-16
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0645 S-0739 cannot be enqueued: its branch is 37 commits behind main, and the close-out stops at sync: done

When a story in review cannot be enqueued because its branch is behind main and its close-out stops at `sync`, the operator brings the branch up to main, and the orchestrator then asks for the close-out again and enqueues the story once it passes. The orchestrator does not rebase a story branch itself.

On 2026-11-15 three stories in review were in this state: S-0739 (37 commits behind main, TH-1226), S-0742 (51 behind, TH-1227), and S-0747 (51 behind, TH-1228). Their agents had ended. Each accept dry run had one blocker, `verify_behind_main` ([ADR-0291](../adrs/0291-every-acceptance-but-the-merge-queue-s-landing-waits-for-a-flai-verify-over.md)). The close-out the orchestrator asked for stopped at `sync`: "the branch does not contain main, ... which change paths outside wip/". `queue_add` then refused the story as `verify_failed` at step sync. Review was full (5 of 5), which starved the board (TH-1224). The orchestrator asked on each thread:

```text
**Recommended:** run `flai stream sync S-0739` in its worktree, settle what it reports, and ask for the close-out again (the MCP tool `verify`). Then I enqueue it.
...
Rebasing the story branch is the story's work, not mine.
```

The operator, alex, answered each with "done" between 2026-11-16T00:27:35Z and 00:28:30Z, having synced the three branches. The orchestrator found each branch 0 commits behind main and asked for its close-out again, for example: "story/S-0739 now contains main (0 commits behind, head 8d1e1ba). Close-out job S-0739-20261116T002837Z is running at that head; the queue takes the story once it passes."

Sources: TH-1226, TH-1227, TH-1228. D-0646 and D-0647 are merged here. D-0633 covers a different send-back: a story whose verify tier main's new paths sent back, which its own agent syncs and closes out again.
