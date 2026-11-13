---
id: S-0751
type: story
nature: feature
title: flai serve brings a missing or stale agent image back as a host job that holds the start, in place of refusing it
status: done
parent: E-0036
owner: alex
created: 2026-11-12T16:01:16Z
updated: 2026-11-13T16:13:05Z
transitions:
  - to: ready
    at: 2026-11-13T07:30:13Z
    by: alex
  - to: in-progress
    at: 2026-11-13T14:38:04Z
    by: agent-S-0751
  - to: review
    at: 2026-11-13T16:11:23Z
    by: agent-S-0751
  - to: done
    at: 2026-11-13T16:13:05Z
    by: orchestrator
tags: [cli]
topics: [docker]
touches: [flai/internal/serve/imagejob.go, flai/internal/serve/imagejob_test.go, flai/internal/serve/container.go, flai/internal/serve/container_test.go, flai/internal/serve/agents_start.go, flai/internal/serve/agents.go, flai/internal/serve/verifyjob_container.go, flai/internal/serve/verifyjob_test.go, flai/internal/serve/checks.go, flai/internal/check/image.go, flai/internal/check/image_test.go, design/system/flai-cli.md, docs/users/flai.md, docs/operators/runbooks/agent-container.md, flai/internal/serve/serve.go, flai/internal/serve/agents_launch.go, flai/cmd/serve.go, flai/cmd/image.go, flai/internal/serve/plan.go, flai/internal/serve/orchestrate.go, flai/internal/serve/analyze.go, flai/internal/serve/curate.go, flai/internal/serve/role.go, flai/internal/serve/agents_activity.go, flai/internal/serve/stop.go, flai/internal/serve/restart.go, flai/internal/serve/plan_test.go, flai/internal/serve/checks_test.go, flai/internal/serve/kluster.go, design/issues/I-0363-flai-s-markdown-lint-reports-md038-on-a-code-span-holding-escaped-backticks.md, design/issues/summary.md]
after: [S-0740, S-0745]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 5617
  priced_by: harness
  turns:
    - day: 2026-11-13
      ceremony: 7
      work: 64
  models:
    - model: claude-opus-5-5
      input: 542
      output: 278630
      cache_read: 48398723
      cache_write: 1186625
      cost: 22.2069
  strategic:
    - kind: orchestrator
      seconds: 21
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 5
          output: 25
          cache_read: 187590
          cache_write: 26002
          cost: 0.0612
forecast:
  duration: 50m
  delivery: 2026-11-13T16:09:00Z
  basis: "Its own forecast of 50m; 3rd in the pull order with an in-progress limit of 3, behind S-0662, S-0663, S-0771, S-0760 and S-0759."
  by: flai
  at: 2026-11-13T14:36:34Z
finalized:
  by: alex
  at: 2026-11-13T07:27:49Z
---
# S-0751 flai serve brings a missing or stale agent image back as a host job that holds the start, in place of refusing it

## Goal

When a start finds its image missing today, `containerImage` (`flai/internal/serve/container.go`) refuses it. The run fails with "is not on this host ... build it with make agent-images". On 2026-11-12 that refused the planners of S-0723, S-0725, and S-0730, and `flai serve` waited for the operator to rebuild by hand. A flai upgrade makes the images stale with no word at all.

After this story, a start that finds its image missing or stale starts the ensure of S-0740 and S-0745 as a host job, one at a time on the host. The start is held until the job ends. The container never gets the Docker socket (ADR-0203), so this is the host's work, not an agent's. When the job passes, the held start goes on. When it fails, the run fails with the cause and the build's log. The run records what the job did, and the journal says it.

## Acceptance criteria

- [x] A story's start, a strategic agent's start, the close-out container, and the checks container that find their image missing or stale ask for an image job in place of refusing. The image is not the operator's own (`--image`), which still refuses as today
- [x] One image job runs on the host at a time. A second start that needs the same images waits for the job in flight rather than starting another, and every held start goes on, in its order, when the job passes
- [x] While held, the run is recorded as waiting for its image in `agents.json`, with the job's id. The serve journal logs the job's start, its end, its way (pull or build), its duration, and its log file
- [x] When the job fails, each held run fails with the cause, the failing step, and the log file's path, in place of naming `make agent-images`. The next start that needs the image starts a new job, and a failure is not retried in a loop
- [x] A strategic agent waits, or starts on the harness image alone, as S-0733's ADR decides for the layer
- [x] After a flai upgrade, the first start that finds its images stale renews them through the same job, shown by a test whose label names another flai version
- [x] `flai check`'s `agent.image` and the start's remaining refusals name `flai image ensure`. `design/system/flai-cli.md` § The container harness (S-0508), `docs/users/flai.md` § In a container, and `docs/operators/runbooks/agent-container.md` describe the job

## Tasks
- T-3177 flai serve runs the image ensure as a host job, one at a time, that held starts wait on and that the journal records
- T-3178 A story's or a strategic agent's start that finds its image missing or stale waits on an image job, recorded on its run, in place of refusing
- T-3179 The close-out container and the checks container wait on an image job when their image is missing or stale
- T-3180 flai check's agent.image names flai image ensure, and flai-cli.md, the user guide, and the runbook describe the image job

## Notes

Proposed slice 4a of E-0036. The dashboard's view of the job, and its retry, are S-0757. The switch that turns the automatic job off, and the operator's own build command, are S-0759.

### Planning

- Touches:
  - `flai/internal/serve/imagejob.go`, `imagejob_test.go`: new files of the job, modelled on the close-out job (`verifyjob.go`).
  - `flai/internal/serve/container.go`, `container_test.go`, `agents_start.go`, `agents.go`: layout. `containerAllowed` and `containerImage` (lines 108 to 143) refuse today; `startOf` (line 398) and the failed run (around line 179) are in `agents_start.go`; `AgentRun` (lines 206 to 216) is in `agents.go`.
  - `flai/internal/serve/verifyjob_container.go`, `verifyjob_test.go`, `checks.go`: layout. The close-out container (line 115) and the checks container (line 250) are the first criterion's other two paths.
  - `flai/internal/check/image.go`, `image_test.go`: the "not built" warning (lines 86 to 88), named in the last criterion.
  - `design/system/flai-cli.md`, `docs/users/flai.md`, `docs/operators/runbooks/agent-container.md`: named in the last criterion.
  - Left out from `flai touches suggest`: `docs/operators/index.md` (25%), `docs/users/flai-reference.md` (22%, no command changes), and the dashboard's documents, which are S-0757's.
- `after`:
  - S-0733 is not named, though this story reads the ADR's rule for strategic agents and the job's limits: it waits for S-0733 through S-0740 and S-0735, which names it.
  - S-0740: the job runs that story's `agentimage.Ensure`.
  - S-0745: `Ensure` builds the layer, which the first and fifth criteria need.
- Tasks, in four layers: T-3177 writes the job; T-3178 changes `containerImage` to request it; T-3179 makes the close-out and the checks wait, building on T-3178's new `containerImage`; T-3180 renames the finding and describes the job.
- Forecast 50m: `flai forecast` gives 37m (109 s per unit over 118 done large feature stories, size 20). Raised: a concurrent runner with joined and queued requests, tested under `-race`, and four start paths changed. Delivery is as flai gives it.
- Cost of delay: not set yet. E-0036 has no inputs; TH-1094 asks the operator for them.

### Accepted by the orchestrator

- Verified: a11bf9744d43e017df31b42f30c596e43fe282f8
- At: 2026-11-13T16:13:05Z

Verdict: pass at a11bf9744d43e017df31b42f30c596e43fe282f8
- 1: flai/internal/serve/container.go, flai/internal/serve/agents_start.go, flai/internal/serve/checks.go, flai/internal/serve/verifyjob_container.go
- 2: flai/internal/serve/imagejob.go, flai/internal/serve/imagejob_test.go
- 3: flai/internal/serve/agents.go, flai/internal/serve/agents_activity.go, flai/internal/serve/imagejob.go
- 4: flai/internal/serve/imagejob.go, flai/internal/serve/container.go, flai/internal/serve/container_test.go
- 5: flai/internal/serve/plan.go, flai/internal/serve/role.go, flai/internal/serve/container.go, flai/internal/serve/container_test.go
- 6: flai/internal/serve/container_test.go
- 7: flai/internal/check/image.go, design/system/flai-cli.md, docs/users/flai.md, docs/operators/runbooks/agent-container.md
