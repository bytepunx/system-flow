---
id: S-0686
type: story
nature: remediation
title: go build's VCS stamping fails with git exit status 128 inside the agent container, though git answers each command there by hand
status: done
owner: alex
created: 2026-11-10T08:37:12Z
updated: 2026-11-11T08:46:40Z
transitions:
  - to: ready
    at: 2026-11-11T01:44:14Z
    by: alex
  - to: in-progress
    at: 2026-11-11T05:19:26Z
    by: agent-S-0686
  - to: review
    at: 2026-11-11T08:00:16Z
    by: agent-S-0686
  - to: done
    at: 2026-11-11T08:46:40Z
    by: alex
blocked:
  - from: 2026-11-11T05:26:33Z
    until: 2026-11-11T07:42:21Z
    reason: pause for urgent fix
tags: [docker, agent-container, go]
touches: [flai/images/system-flow/Dockerfile, scripts/agent-container.sh, scripts/agent-vcs-test.sh, scripts/smoke.sh, scripts/README.md, flai/internal/harness/container.go, docs/operators/runbooks/agent-container.md, design/tech/docker.md, design/issues/I-0255-go-build-s-vcs-stamping-fails-with-git-exit-status-128-inside-the-agent.md, design/issues/summary.md, design/issues/I-0348-a-go-build-in-a-story-s-worktree-stamps-vcs-revision-of-the-main-checkout-s.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1543
  priced_by: harness
  turns:
    - day: 2026-11-11
      ceremony: 4
      work: 58
  models:
    - model: claude-opus-5-5
      input: 138
      output: 33162
      cache_read: 10355922
      cache_write: 389650
      cost: 5.8522
  strategic:
    - kind: orchestrator
      seconds: 48
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 20
          output: 126
          cache_read: 753078
          cache_write: 113890
          cost: 0.2409
cost_of_delay:
  inputs:
    time_lost_per_cycle: 14m
    by: flai
    at: 2026-11-10T08:37:12Z
  value: 35
  by: planner-S-0686
  at: 2026-11-10T10:34:53Z
forecast:
  duration: 1h30m
  delivery: 2026-11-11T06:53:00Z
  basis: "Its own forecast of 1h30m; 2nd in the pull order with an in-progress limit of 3, behind S-0690, S-0691, S-0693 and S-0678."
  by: flai
  at: 2026-11-11T05:07:32Z
finalized:
  by: orchestrator
  at: 2026-11-10T11:10:11Z
---
# S-0686 go build's VCS stamping fails with git exit status 128 inside the agent container, though git answers each command there by hand

## Goal

This story remediates [I-0255](../../../design/issues/I-0255-go-build-s-vcs-stamping-fails-with-git-exit-status-128-inside-the-agent.md), "go build's VCS stamping fails with git exit status 128 inside the agent container, though git answers each command there by hand". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0255 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0255 is closed with `flai issue close I-0255 --reason` saying what fixed it

## Tasks
- T-3032 The cause of I-0255's exit status 128 is found in the layer's image and a fix is proposed
- T-3036 The smoke tier's agent-container test runs a VCS-stamped go build in the layer's image
- T-3039 A story's container stamps go builds again, and the layer drops GOFLAGS=-buildvcs=false
- T-3041 The runbook and docker.md say how a story's container stamps go builds, and I-0255 is closed

## Notes

Cost of delay inputs set by flai from I-0255. time_lost_per_cycle 14m: 15m per occurrence × 1 occurrence ÷ 1.1 cycles of 168h (first reported 2026-11-02T19:02:40Z, 7.6 days before this story).

### Planning

Planned by planner-S-0686 on 2026-11-10. The orchestrator asked for a plan again, naming a thread `asked`. No thread of that name exists, and none is on S-0686 or I-0255. So this is the story's first plan, made from I-0255. The story declared no touches. `flai touches suggest` needs a start path, so it was given the layer's Dockerfile, `scripts/agent-container-test.sh`, and `flai/internal/harness/container.go`. It found 14 commits that changed them.

| Touch | Source | Why |
|-------|--------|-----|
| `flai/images/system-flow/Dockerfile` | layout | Holds the workaround `GOFLAGS=-buildvcs=false` with its I-0255 comment, which T-3039 removes |
| `scripts/agent-container.sh` | layout | The hand-run container start: its `run` sets the mounts and `safe.directory` (T-3039) |
| `scripts/agent-container-test.sh` | layout | The smoke tier's container test, which already reproduces I-0279 (T-3036) |
| `flai/internal/harness/container.go` | layout | The container harness: `safeMain` and the mounts a served story's container gets (T-3039) |
| `flai/internal/harness/container_test.go` | co-change (50%) | Its tests, changed with `container.go` in half the commits (T-3039) |
| `docs/operators/runbooks/agent-container.md` | design | The runbook for the agent images and their mounts (T-3041) |
| `design/tech/docker.md` | design | Its "Agent images" row describes the layer's contents (T-3041) |
| `design/issues/I-0255-…agent.md` | design | The issue the goal names: T-3032 writes to it and T-3041 closes it |
| `design/issues/summary.md` | layout | `flai issue close` regenerates it |

Co-change also listed `scripts/agent-image.sh` (21%), `flai/internal/harness/harness.go`, `harness_test.go`, and `scripts/README.md` (14%). They are left out. `agent-image.sh --check` sets the same `safe.directory` as `agent-container.sh`, so a story's agent widens the touches if T-3032's fix reaches it. No script is added. There are no folder touches.

The cause is not known yet. I-0255's one instance, S-0510's on 2026-11-02, came before the main checkout's mounts and their `safe.directory` (ADR-0253 of 2026-11-06, S-0628). It may no longer reproduce. Go's error leaves out git's own message, and exit 128 fits more than one cause, such as dubious ownership or a registration that is not mounted (I-0215). So T-3032 reproduces with `GIT_TRACE` written to a file before anything is built.

Figures:

- Forecast: `flai forecast` gave 25m: 134 s a unit over 54 done remediation stories, times size 11 (2 criteria, 9 touches). Raised to 1h30m, since T-3032 and T-3039 run the layer's image live and T-3039 rebuilds it. S-0667's plan for a like container cause was forecast at 1h30m. Delivery is moved from flai's 17:06Z by the extra 65 minutes.
- Cost of delay: `flai cod` gave 35.00 USD a week from flai's input: 14m lost per 168h cycle at 150 USD an hour. Kept as given. I-0255 has one occurrence, and the layer's workaround holds meanwhile.
- No `after`. No other story writes output this one reads. S-0667 and S-0713 change `scripts/agent-container.sh` and the runbook too, and a trial merge judges that.
