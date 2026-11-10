---
id: S-0713
type: story
nature: improvement
title: make agent-images fails at the first base image when Docker Desktop's credential helper desktop.exe cannot start from WSL
status: in-progress
owner: alex
created: 2026-11-10T10:11:03Z
updated: 2026-11-10T10:52:33Z
transitions:
  - to: ready
    at: 2026-11-10T10:25:05Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-10T10:52:33Z
    by: agent-S-0713
tags: [scripts, docker]
touches: [scripts/agent-image.sh, scripts/agent-image-test.sh, scripts/smoke.sh, scripts/README.md, docs/operators/runbooks/agent-container.md, design/issues/I-0339-make-agent-images-fails-at-the-first-base-image-when-docker-desktop-s.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 104
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-10
      ceremony: 2
      work: 12
  models:
    - model: claude-opus-5-5
      input: 34
      output: 147
      cache_read: 1760128
      cache_write: 170152
      cost: 0.7106
  strategic:
    - kind: orchestrator
      seconds: 143
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 32
          output: 277
          cache_read: 1129071
          cache_write: 131899
          cost: 0.3446
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-11-10T10:11:03Z
  value: 12.5
  by: planner-S-0713
  at: 2026-11-10T10:23:28Z
forecast:
  duration: 30m
  delivery: 2026-11-10T15:18:00Z
  basis: flai forecast gave 17m (107 s per unit over 30 done medium improvement stories, size 9), raised to 30m, the median of the three agent-image and agent-container remediations before it (S-0670 9m, S-0672 29m, S-0682 66m), with the delivery moved 13m later to match.
  by: planner-S-0713
  at: 2026-11-10T10:23:28Z
finalized:
  by: orchestrator
  at: 2026-11-10T10:24:22Z
---
# S-0713 make agent-images fails at the first base image when Docker Desktop's credential helper desktop.exe cannot start from WSL

## Goal

This story remediates [I-0339](../../../design/issues/I-0339-make-agent-images-fails-at-the-first-base-image-when-docker-desktop-s.md), "make agent-images fails at the first base image when Docker Desktop's credential helper desktop.exe cannot start from WSL". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0339 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0339 is closed with `flai issue close I-0339 --reason` saying what fixed it

## Tasks
- T-2955 scripts/agent-image.sh builds with an anonymous Docker configuration when the configured credential helper fails
- T-2956 scripts/agent-image-test.sh reproduces I-0339 without Docker, and the smoke tier runs it
- T-2957 The agent-container runbook and scripts/README.md say what make agent-images does when the credential helper fails
- T-2958 I-0339 is closed with the reason that scripts/agent-image.sh builds with no login when the credential helper fails

## Notes

Cost of delay inputs set by flai from I-0339. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-10T09:13:13Z, 0 days before this story; under one cycle counts as one).

### Planning

Proposed solution, from I-0339's one instance: before its first `docker build`, `scripts/agent-image.sh` reads the Docker configuration the build would use (`$DOCKER_CONFIG/config.json`, else `~/.docker/config.json`). When it names a `credsStore`, the script asks that helper (`docker-credential-<store> list`) under a short timeout. When the helper is missing, fails, or does not answer, the script builds with `DOCKER_CONFIG` pointed at a scratch folder whose `config.json` is `{}`, so the public base images are pulled anonymously, and prints one line naming the helper and that it did so. The operator's configuration is not changed. A working helper, or none named, leaves the build as it is today. The CLI plugins folder of the original configuration is linked into the scratch one, so `docker build` still finds buildx when it lives there.

Touches, each with where it came from:

| Touch | From |
|-------|------|
| `scripts/agent-image.sh` | design: I-0339's instance names it; co-change lists nothing else that the fix needs |
| `scripts/agent-image-test.sh` | layout: a new test beside `scripts/agent-container-test.sh`, which reproduces I-0279 the same way, with no Docker |
| `scripts/smoke.sh` | layout: it runs `agent-container-test.sh`, so it runs the new test too |
| `scripts/README.md` | co-change (2 of 5 commits of `agent-image.sh`): its rows for `agent-image.sh`, `smoke.sh`, and the new test |
| `docs/operators/runbooks/agent-container.md` | design: its build section says what `make agent-images` does |
| `design/issues/I-0339-make-agent-images-fails-at-the-first-base-image-when-docker-desktop-s.md` | criterion 2 |
| `design/issues/summary.md` | criterion 2: closing the issue rewrites its row |

Co-change also listed `flai/images/system-flow/Dockerfile` (3 of 5) and `Makefile` (2 of 5). Neither is a touch: the fix changes no image and no target.

No folder touch and no `after`: the story reads no other story's output. S-0672's `--check` mode in `scripts/agent-image.sh` builds nothing, so the fallback does not concern it; a trial merge settles the shared file.

Forecast: `flai forecast` gave 17m, at 107 s per unit of size over 30 done medium improvement stories, size 9. Raised to 30m, the median of the three agent-image and agent-container remediations before it: S-0670 took 9m, S-0672 29m, S-0682 66m. This one's test needs a fake credential helper and a fake `docker` that records the configuration, which is more than S-0670's. The delivery is flai's, 2026-11-10T15:05Z, moved 13m later to match.

Cost of delay: 12.50 USD a week, as `flai cod` works it out from the operator's input, 5m lost per 168h cycle at 150 USD an hour. It stands: I-0339 has one instance, and the workaround it records took 5m.
