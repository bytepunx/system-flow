---
id: S-0667
type: story
nature: improvement
title: A docker run from a story's agent's shell mounts the main checkout as an empty folder, so a hand call in the agent image cannot read the project
status: ready
owner: alex
created: 2026-11-10T08:34:40Z
updated: 2026-11-10T09:53:50Z
transitions:
  - to: ready
    at: 2026-11-10T08:58:35Z
    by: orchestrator
tags: [docker, agent-container]
touches: [scripts/agent-container.sh, docs/operators/runbooks/agent-container.md, design/tech/docker.md, design/issues/I-0276-a-docker-run-from-a-story-s-agent-s-shell-mounts-the-main-checkout-as-an-empty.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 21
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 4
          output: 21
          cache_read: 676495
          cache_write: 8724
          cost: 0.1857
cost_of_delay:
  inputs:
    time_lost_per_cycle: 20m
    by: flai
    at: 2026-11-10T08:34:40Z
  value: 50
  by: planner-S-0667
  at: 2026-11-10T08:49:29Z
forecast:
  duration: 1h30m
  delivery: 2026-11-10T12:18:00Z
  basis: "Its own forecast of 1h30m; 2nd in the pull order with an in-progress limit of 5, behind S-0668, S-0669, S-0672, S-0673, S-0684 and S-0665."
  by: flai
  at: 2026-11-10T09:53:50Z
finalized:
  by: alex
  at: 2026-11-10T08:51:21Z
---
# S-0667 A docker run from a story's agent's shell mounts the main checkout as an empty folder, so a hand call in the agent image cannot read the project

## Goal

This story remediates [I-0276](../../../design/issues/I-0276-a-docker-run-from-a-story-s-agent-s-shell-mounts-the-main-checkout-as-an-empty.md), "A docker run from a story's agent's shell mounts the main checkout as an empty folder, so a hand call in the agent image cannot read the project". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0276 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0276 is closed with `flai issue close I-0276 --reason` saying what fixed it

## Tasks
- T-2926 The cause of I-0276's empty mount is found on the host and a fix is proposed
- T-2927 A hand-run container sees the main checkout's files, and refuses with the cause when a mount comes up empty
- T-2929 The runbook and docker.md say why a hand docker run can mount an empty folder, and I-0276 is closed

## Notes

Cost of delay inputs set by flai from I-0276. time_lost_per_cycle 20m: 20m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-07T05:43:57Z, 3.1 days before this story; under one cycle counts as one).

### Planning

Planned by planner-S-0667 on 2026-11-10. The story declared no touches. `flai touches suggest` needs a start path, so it was given `scripts/agent-container.sh` and the runbook. It found 10 commits that changed them.

| Touch | Source | Why |
|-------|--------|-----|
| `scripts/agent-container.sh` | layout | The repository's hand-run form of a container start (S-0505), and the place a check for an empty mount goes (T-2927) |
| `docs/operators/runbooks/agent-container.md` | design, co-change | The runbook for hand-run containers, linked from `design/system/strategic-agents.md` § The live run (S-0567) (T-2929) |
| `design/tech/docker.md` | co-change (30%) | Where the fact about the host's Docker that the cause rests on belongs (T-2929) |
| `design/issues/I-0276-…empty.md` | design | The issue the goal names, bumped by T-2926 and closed by T-2929 |
| `design/issues/summary.md` | layout | `flai issue close` regenerates it |

Co-change also listed `design/system/flai-cli.md`, `docs/operators/index.md`, `docs/users/flai.md`, and `scripts/README.md` at 40%. They are left out: no task changes a command or adds a script. A story's agent widens the touches if T-2926's fix needs one of them. There are no folder touches.

The cause is not known yet. S-0567's agent ran on the host on `harness: claude-code`, and its `docker run -v <checkout>:<checkout>:ro` came up empty inside and outside the sandbox. A container `flai serve` started for the planner, with the same path, had the files. T-2926 compares the two, by docker context, `DOCKER_HOST`, and daemon, before anything is built.

Figures:

- Forecast: `flai forecast` gave 13m: 107 s a unit over 30 done improvement stories, times size 7 (2 criteria, 5 touches). Raised to 1h30m, since T-2926 runs docker live on the host from several shells before T-2927 can start. S-0567's live run was forecast at 1h. Delivery is moved from flai's 11:09Z by the extra 77 minutes.
- Cost of delay: `flai cod` gave 50.00 USD a week from flai's input: 20m lost per 168h cycle at 150 USD an hour. Kept as given: I-0276 has one occurrence, and the cost of a hand container that cannot read the project is about that.
- No `after`. No other story writes output this one reads. S-0670 and S-0682 change `scripts/agent-container.sh` too, which a trial merge judges.
