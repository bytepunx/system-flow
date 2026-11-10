---
id: S-0682
type: story
nature: improvement
title: scripts/agent-container.sh can run only a story's own agent, under the story's own name, so a trial run of the image acts on the live story
status: done
owner: alex
created: 2026-11-10T08:37:06Z
updated: 2026-11-10T10:10:48Z
transitions:
  - to: ready
    at: 2026-11-10T08:49:30Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-10T08:49:35Z
    by: agent-S-0682
  - to: review
    at: 2026-11-10T09:41:35Z
    by: agent-S-0682
  - to: in-progress
    at: 2026-11-10T09:54:27Z
    by: flai
  - to: review
    at: 2026-11-10T10:08:14Z
    by: agent-S-0682
  - to: done
    at: 2026-11-10T10:10:48Z
    by: alex
tags: [containers, pilot, scripts]
touches: [flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/internal/pilot/request.go, flai/internal/pilot/request_test.go, flai/internal/pilot/run.go, flai/internal/pilot/run_test.go, scripts/agent-container.sh, docs/operators/runbooks/agent-container.md, design/system/flai-cli.md, design/tech/docker.md, design/issues/I-0224-scripts-agent-container-sh-can-run-only-a-story-s-own-agent-under-the-story-s-own-name-so-a-trial-run-of-the-image-acts-on-the-live-story.md, design/issues/summary.md, design/issues/I-0338-flai-task-done-s-sync-reports-undoes-main-for-paths-the-story-never-changed.md, design/issues/I-0339-make-agent-images-fails-at-the-first-base-image-when-docker-desktop-s.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 3958
  priced_by: harness
  turns:
    - day: 2026-11-10
      ceremony: 2
      hand_edits: 5
      work: 127
  models:
    - model: claude-opus-5-5
      input: 394
      output: 72877
      cache_read: 49243052
      cache_write: 735669
      cost: 16.2937
  strategic:
    - kind: orchestrator
      seconds: 103
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 19
          output: 139
          cache_read: 690120
          cache_write: 74409
          cost: 0.2082
        - model: claude-sonnet-5-5
          input: 4
          output: 53
          cache_read: 12769
          cache_write: 30314
          cost: 0.0323
cost_of_delay:
  inputs:
    time_lost_per_cycle: 3m
    by: flai
    at: 2026-11-10T08:37:06Z
  value: 7.5
  by: planner-S-0682
  at: 2026-11-10T08:45:44Z
forecast:
  duration: 45m
  delivery: 2026-11-10T15:47:00Z
  basis: "Its own forecast of 45m; 24th in the pull order with an in-progress limit of 5, behind S-0685, S-0241, S-0659, S-0660, S-0661, S-0662, S-0663, S-0665, S-0666, S-0667, S-0668, S-0669, S-0670, S-0671, S-0672, S-0673, S-0674, S-0675, S-0676, S-0677, S-0678, S-0679, S-0680 and S-0681."
  by: flai
  at: 2026-11-10T08:49:28Z
finalized:
  by: orchestrator
  at: 2026-11-10T08:48:59Z
---
# S-0682 scripts/agent-container.sh can run only a story's own agent, under the story's own name, so a trial run of the image acts on the live story

## Goal

This story remediates [I-0224](../../../design/issues/I-0224-scripts-agent-container-sh-can-run-only-a-story-s-own-agent-under-the-story-s-own-name-so-a-trial-run-of-the-image-acts-on-the-live-story.md), "scripts/agent-container.sh can run only a story's own agent, under the story's own name, so a trial run of the image acts on the live story". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0224 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0224 is closed with `flai issue close I-0224 --reason` saying what fixed it

## Tasks
- T-2896 A start request marked probe gets a prompt that reads and reports and works no story
- T-2897 pilot starts a probe without the story's FLAI_STORY and refuses a probe that names a role
- T-2898 scripts/agent-container.sh --probe runs the image as probe-S-nnnn with a probe start request
- T-2899 The runbook, flai-cli.md, and docker.md say how to probe an agent image without acting on its story
- T-2900 Probe the rebuilt image on S-0682, check its narrative is untouched, and close I-0224

## Notes

Cost of delay inputs set by flai from I-0224. time_lost_per_cycle 3m: 5m per occurrence × 1 occurrence ÷ 1.8 cycles of 168h (first reported 2026-10-28T23:10:55Z, 12.4 days before this story).

### Planning

The proposed remedy, from I-0224's one instance (S-0507's T-2174, whose capped trial called `story_start`, `prime`, `inbox`, and `stream_state` on the live story): a probe. `scripts/agent-container.sh --probe S-nnnn` writes `"probe":true` in pilot's start request and names the agent `probe-S-nnnn` in a new session. `harness.Prompt` answers a probe with a prompt that only reads and reports, and pilot starts it without `FLAI_STORY`. The mounts stay the story's, so the trial sees what the story's agent would. A throwaway story of its own was the issue's other option; it leaves an item behind per trial and was not taken.

Touches, all files, no folder:

| Touch | Source |
|-------|--------|
| `flai/internal/harness/harness.go`, `harness_test.go` | layout: `harness.Request` and `Prompt` live there, and pilot decodes the request into it with unknown fields refused |
| `flai/internal/pilot/request.go`, `request_test.go`, `run.go`, `run_test.go` | layout: `Check` and `childEnv`, which sets `FLAI_STORY` |
| `scripts/agent-container.sh` | design: named by the issue and the goal |
| `docs/operators/runbooks/agent-container.md` | design: the runbook of the script; its Run section still shows a `prompt.md` form the script lost at S-0507 |
| `design/system/flai-cli.md`, `design/tech/docker.md` | design: both describe pilot and the script |
| `design/issues/I-0224-….md`, `design/issues/summary.md` | criterion 2 |

`flai touches suggest` listed nothing, since the story declared no touches.

No `after`. S-0670 also changes `scripts/agent-container.sh` (the image it runs), but this story reads nothing S-0670 writes, so the trial merge judges the two.

Forecast 45m, raised from `flai forecast`'s 24m (101 s a unit, size 14). S-0517, S-0520, and S-0521, container remediations of this script and pilot, took 39m, 49m, and 43m of agent time. This story also rebuilds the images for its trial and runs an integration close-out, which the size band does not count.

Cost of delay 7.50 USD a week, as `flai cod` gives it from flai's input of 3m lost per cycle at 150 USD an hour. The figure stands: the input is the issue's own measure.
- 2026-11-10T09:54:27Z: moved to in-progress: the rebase of story/S-0682 onto main stopped on conflicts in design/tech/docker.md, docs/operators/runbooks/agent-container.md, which the acceptance of S-0682 by alex cannot settle, so it aborted the rebase; its agent settles it in the worktree and syncs it with flai stream sync S-0682
