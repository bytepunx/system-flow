---
id: T-3194
type: task
nature: feature
title: The image job runs the operator's build command when the host names one, and a start refuses a missing image when the automatic job is off
status: done
parent: S-0759
owner: alex
created: 2026-11-12T16:05:51Z
updated: 2026-11-13T17:19:28Z
transitions:
  - to: ready
    at: 2026-11-13T16:55:38Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T16:55:38Z
    by: system-flow
  - to: done
    at: 2026-11-13T17:19:28Z
    by: agent-S-0759
stream: S-0759
tags: [cli]
touches: [flai/internal/serve/imagejob.go, flai/internal/serve/imagejob_test.go, flai/internal/serve/container.go, flai/internal/serve/container_test.go, flai/cmd/image_test.go, flai/cmd/serve_agent.go, flai/cmd/serve_agent_image_test.go, flai/internal/serve/agents.go, flai/internal/serve/checks.go]
after: [T-3190, T-3192]
usage:
  source: log
  seconds: 1428
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 110
      output: 625
      cache_read: 7795038
      cache_write: 157732
      cost: 2.98
---
# T-3194 The image job runs the operator's build command when the host names one, and a start refuses a missing image when the automatic job is off

## Work

- `flai/internal/serve/imagejob.go`: a job runs `Ensure` with the project's command from the host's configuration (`AgentStart.ImageCommand(root)`, T-3190), so T-3192's runner runs it in place of the pull and build. The job's record and journal say the command ran.
- `flai/internal/serve/container.go`: with `AutoImage(root)` false, a start that finds its image missing or stale requests no job and refuses as before S-0751, naming the image, `flai image ensure`, and `flai serve agent image auto on`. This holds for a story's agent, a strategic agent, the close-out's container, and the checks' containers, wherever S-0751 holds them for a job.

Tests cover:

- a job through the command;
- a start refused with the setting off;
- a start held and gone on with it on.

It waits for T-3190, whose settings it reads, and T-3192, whose command `Ensure` runs.

## Done when

- `flai test` on the changed files passes; run `flai test race` too, since the job is concurrent.

## Notes
