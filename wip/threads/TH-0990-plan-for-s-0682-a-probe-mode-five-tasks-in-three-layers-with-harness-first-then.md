---
id: TH-0990
title: "Plan for S-0682: a --probe mode, five tasks in three layers, with harness first, then pilot and the script, then docs and the trial"
anchor:
  path: wip/kanban/stories/S-0682-scripts-agent-container-sh-can-run-only-a-story-s-own-agent-under-the-story-s.md
  item: S-0682
status: resolved
participants: [planner-S-0682, orchestrator]
created: 2026-11-10T08:47:06Z
updated: 2026-11-10T08:48:55Z
---

# TH-0990 Plan for S-0682: a --probe mode, five tasks in three layers, with harness first, then pilot and the script, then docs and the trial

On wip/kanban/stories/S-0682-scripts-agent-container-sh-can-run-only-a-story-s-own-agent-under-the-story-s.md.

## Entries

### 2026-11-10T08:47:06Z planner-S-0682
S-0682 is planned: a proposed remedy, five draft tasks in three layers, 12 file touches, forecast 45m, cost of delay 7.50 USD a week.

## Remedy proposed

`scripts/agent-container.sh --probe S-nnnn` writes `"probe":true` in pilot's start request and names the agent `probe-S-nnnn`, in a new session. `harness.Prompt` gives a probe a prompt that only reads and reports. pilot starts it without `FLAI_STORY`. The mounts stay the story's. I did not take the issue's other option, a throwaway story per trial, because it leaves an item behind each time.

## Tasks and layers

| Layer | Task | Touches |
|-------|------|---------|
| 1 | T-2896 `Request.Probe` and the probe prompt | `flai/internal/harness/harness.go`, `harness_test.go` |
| 2 | T-2897 pilot: no `FLAI_STORY` for a probe, a probe with a role or an item refused, and the test that reproduces I-0224 | `flai/internal/pilot/request.go`, `request_test.go`, `run.go`, `run_test.go` |
| 2 | T-2898 the script's `--probe` | `scripts/agent-container.sh` |
| 3 | T-2899 the runbook, `flai-cli.md` § pilot, `docker.md` | the three documents |
| 3 | T-2900 rebuild, probe S-0682, close I-0224 | `design/issues/I-0224-….md`, `summary.md` |

T-2897 and T-2898 wait for T-2896. pilot refuses unknown fields, so `probe` must exist in `harness.Request` first. T-2899 and T-2900 wait for both layer-2 tasks. Tasks in the same layer share no file.

## Assumptions

- A guard rule that holds a probe to read-only calls is out of scope. The prompt, the name `probe-S-nnnn`, and the missing `FLAI_STORY` are the fix. If you want the guard to enforce it, say so and I will add a task.
- T-2900's capped trial uses haiku with a 0.02 USD limit, a paid call. Its agent asks you on a thread before running it. The free `-- env` check runs without asking.
- `scripts/agent-container.sh` is on the container harness's protected list (S-0508). T-2898's write asks the owner unless auto-approve is on, and the story is yours to accept.
- No `after`. S-0670 also changes the script, but this story reads nothing from it.
- The runbook's Run section still shows a `prompt.md` form that the script lost at S-0507. T-2899 fixes it, since the section changes anyway.
- This run was started for a thread named "asked", but no such thread exists on S-0682 and the story had no plan before. I planned it from scratch.

### 2026-11-10T08:48:42Z orchestrator
Approved as written: T-2896, then T-2897 and T-2898, then T-2899 and T-2900, 45m, 7.50 USD a week.

- The probe remedy fits I-0224's one instance. A trial run gets a read-and-report prompt, the name `probe-S-nnnn`, and no `FLAI_STORY`, so it no longer acts on the live story. That is criterion 1, and T-2897's test reproduces it.
- The layers follow what each task reads. pilot refuses unknown fields, so `harness.Request` must change first. The 12 touches are files, and no two tasks in a layer share one.
- A guard rule that holds a probe to read-only calls stays out of scope, as you proposed. Nobody asked for it.
- T-2900's paid trial is asked of the operator on its own thread before it runs. That is the right place for that spend decision.

### 2026-11-10T08:48:55Z orchestrator
Resolved: Plan approved as written by the orchestrator under plan_backlog_stories.
