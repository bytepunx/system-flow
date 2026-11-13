---
id: T-3181
type: task
nature: feature
title: The host API reads the project's agent images and image jobs, reads a job's log from an offset, and starts a new job behind the agent host action
status: done
parent: S-0757
owner: alex
created: 2026-11-12T16:03:36Z
updated: 2026-11-13T16:35:56Z
transitions:
  - to: ready
    at: 2026-11-13T16:14:37Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T16:14:38Z
    by: system-flow
  - to: done
    at: 2026-11-13T16:35:56Z
    by: agent-S-0757
stream: S-0757
tags: [cli]
touches: [flai/internal/hostapi/images.go, flai/internal/hostapi/images_test.go, flai/internal/serve/imagejob.go, flai/internal/serve/imagejob_test.go, flai/internal/hostapi/guardreads_test.go, flai/internal/hostapi/writes_actions.go, flai/internal/hostapi/writes_methods.go, flai/internal/serve/serve.go, flaiover/src/lib/server/agent.ts]
usage:
  source: log
  seconds: 1278
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 142
      output: 54472
      cache_read: 8161210
      cache_write: 306851
      cost: 4.5182
---
# T-3181 The host API reads the project's agent images and image jobs, reads a job's log from an offset, and starts a new job behind the agent host action

## Work

Add to `flai/internal/hostapi/images.go`, registered as `kluster.go` registers its methods (line 1497):

- `images.get`: the project's images, with their state from `agentimage` status, and the image jobs of S-0751's runner, running and the last ended per image set. Each job has its id, images, way, current step, start, end, outcome, cause, and log size.
- `images.log`: a job's log from an offset, as `klusterLogTail` (line 70) does.
- `images.ensure`: a write behind `ActionAgent` that starts a new job for the project. It is refused, naming the running job, while one for the same images runs.

Give the runner in `flai/internal/serve/imagejob.go` what these read: the jobs it keeps, each job's current step, and a start by request. Test the three methods and the refusal in `images_test.go`, and what the runner keeps in `imagejob_test.go`.

It is the first layer.

## Done when

- `go test ./internal/hostapi/ ./internal/serve/ -run 'Image'` passes, under `-race` too.
- `contract_test.go`'s list of methods names the three.

## Notes
