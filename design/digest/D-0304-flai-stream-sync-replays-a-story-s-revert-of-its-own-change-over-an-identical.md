---
id: D-0304
title: flai stream sync replays a story's revert of its own change over an identical change on main, undoing main's
kind: cause
sources: [I-0124]
topics: [cli, flai, go, flaiover, dashboard, sveltekit, code, all]
paths: [flai/internal/storygit/undone.go, flai/internal/storygit/undone_test.go, flai/internal/storygit/sync.go, flai/internal/storygit/sync_test.go, flai/cmd/stream_restore.go, flai/cmd/stream_restore_test.go, flai/cmd/stream_sync.go, flai/cmd/stream_sync_test.go, flai/cmd/stream.go, flai/internal/preview/accept.go, flai/cmd/branch.go, flai/cmd/accept.go, flai/cmd/accept_undone_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flaiover/src/lib/server/agent.ts, docs/users/flai-reference.md, docs/operators/settings.md, design/adrs/README.md, design/system/flai-cli.md, design/system/workflow.md, design/system/strategic-agents.md, docs/users/flai.md, design/issues/I-0124-flai-stream-sync-replays-a-story-s-revert-of-its-own-change-over-an-identical-change-on-main-undoing-main-s.md, design/issues/summary.md, design/adrs/0155-a-rebased-story-branch-that-undoes-main-s-change-to-a-path-is-recorded-at-sync.md, flai/cmd/move.go, flai/internal/preview/orchestrator.go, design/issues/I-0146-testthegateisaskedoneveryattempt-fails-now-and-then-the-state-reads-refused-before-the-second-refusal-is-logged.md]
date: 2026-10-14
status: absorbed
absorbed_by: "design/adrs/0155-a-rebased-story-branch-that-undoes-main-s-change-to-a-path-is-recorded-at-sync.md#decision"
asked_again: []
cited: []
---

# D-0304 flai stream sync replays a story's revert of its own change over an identical change on main, undoing main's

I-0124 was closed on 2026-10-14.

The cause:

```text
flai stream sync replays a story's revert of its own change over an identical change on main, undoing main's
```

The fix:

```text
Fixed by S-0371 (ADR-0155): after a clean rebase storygit.Rebase compares the lines main added since the branch's old merge base with the lines the rebased branch removes, and records each path that removes any in .flai-cache/undone/S-nnnn.json. flai stream sync and flai task done warn on every recorded path still undone; the acceptance preview lists it as a blocker with its lines as undoes_main, and flai accept refuses it, its own rebase checking too, until the path is restored with flai stream restore (host write stream.restore) or kept with --keep-removal and a reason (accept.run keep_removals), which the orchestrator may not give. Tests reproduce I-0124 in storygit (sync_test.go), flai stream sync, flai stream restore, and flai accept.
```

Source: I-0124.
