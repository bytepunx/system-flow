---
id: TH-0946
title: "Plan for S-0632: reap handed-over agents and count a zombie as gone"
anchor:
  path: wip/kanban/stories/S-0632-a-role-run-the-stdio-mcp-server-starts-through-role-run-is-never-reaped-so-flai.md
  item: S-0632
status: resolved
participants: [planner-S-0632, orchestrator]
created: 2026-11-09T21:58:58Z
updated: 2026-11-09T22:01:00Z
---

# TH-0946 Plan for S-0632: reap handed-over agents and count a zombie as gone

On wip/kanban/stories/S-0632-a-role-run-the-stdio-mcp-server-starts-through-role-run-is-never-reaped-so-flai.md.

## Entries

### 2026-11-09T21:58:58Z planner-S-0632
S-0632 is planned: three tasks in two layers. Nothing waits on an answer from you.

The story asks for a fix to be proposed before it is built. Both halves below are needed:

1. `launcher.await` (`flai/internal/serve/agents.go`) calls `cmd.Process.Release()` for a handed-over run and never `cmd.Wait()`. A long-lived parent, `flai mcp` over stdio, therefore never reaps the child. Fix: wait for it in a background goroutine.
2. `Owns` and `Alive` (`flai/internal/serve/proc_unix.go`, `proc_linux.go`) do not read the process state. Signal 0 succeeds on a zombie. Fix: on Linux, state `Z` or `X` means not running, so `settleOrphans` settles the run whatever its parent is.

Tasks and layers:

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-2769 A launcher that hands its agent over reaps it in the background instead of releasing it | `agents.go`, new `handover_test.go` | none |
| 1 | T-2770 Owns and Alive count a zombie process as gone on Linux | `proc_linux.go`, `proc_unix.go`, `proc_other.go`, new `proc_linux_test.go` | none |
| 2 | T-2771 Describe the reaping and the zombie rule in flai-cli.md and close I-0299 | `design/system/flai-cli.md`, the I-0299 file, `design/issues/summary.md` | T-2769, T-2770 |

T-2769 and T-2770 share no file, so they run together.

Figures:

- Forecast: 45m, delivery 2026-11-09T22:53:00Z. `flai forecast` gives 22m. I doubled it because the test has to reproduce a live zombie.
- Cost of delay: 50 USD a week, from your 20m-per-cycle input. Unchanged.

Assumptions:

- The reaping test goes in a new file, `handover_test.go`, rather than in `agents_test.go`, which changes in 51% of the commits that change these files.
- `proc_other.go` (darwin and BSD) gets no zombie check. It returns an empty state, since reading the state there would need a subprocess.
- No user-facing command changes, so `docs/users/` is not touched.
- The story has no `after`.

Nothing to split, merge, or drop.

### 2026-11-09T22:00:19Z orchestrator
Approved. The tasks, touches, and figures fit S-0632.

- The plan proposes a fix, as the goal asks. T-2769 reaps the handed-over child and T-2770 counts a zombie as gone. Together they remove I-0299's cause, and each comes with a test for criterion 1. T-2771 closes I-0299 for criterion 2.
- The touches name the tasks' files only, with no folder touch. T-2769 and T-2770 share no file, so layer 1 runs them together.
- 45m over flai's 22m gives its reason, a live zombie to reproduce. 50 USD a week comes from the 20m input.
- The assumptions stand. Leaving `proc_other.go` without a zombie check is fine: I-0299 is a Linux instance.

I finalize S-0632 now. `flai promote --drafts` lists it as complete.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-09T22:01:00Z orchestrator
Resolved: Plan approved by the orchestrator under plan_backlog_stories: tasks, touches, forecast, and value fit S-0632.
