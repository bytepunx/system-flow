---
id: T-3041
type: task
nature: remediation
title: The runbook and docker.md say how a story's container stamps go builds, and I-0255 is closed
status: done
parent: S-0686
owner: alex
created: 2026-11-10T10:34:15Z
updated: 2026-11-11T07:50:28Z
transitions:
  - to: ready
    at: 2026-11-11T07:50:18Z
    by: system-flow
  - to: in-progress
    at: 2026-11-11T07:50:18Z
    by: system-flow
  - to: done
    at: 2026-11-11T07:50:28Z
    by: agent-S-0686
stream: S-0686
tags: [docker, agent-container, docs]
touches: [docs/operators/runbooks/agent-container.md, design/tech/docker.md, design/issues/I-0255-go-build-s-vcs-stamping-fails-with-git-exit-status-128-inside-the-agent.md, design/issues/summary.md]
after: [T-3036, T-3039]
usage:
  source: log
  seconds: 10
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 2
      output: 303
      cache_read: 86674
      cache_write: 3553
      cost: 0.0518
---
# T-3041 The runbook and docker.md say how a story's container stamps go builds, and I-0255 is closed

## Work

Write down what T-3039 changed, then close the issue. It waits for T-3039, since it describes the fix that task built. It also changes I-0255, which T-3032 wrote to.

- In `design/tech/docker.md`, update the "Agent images" row with the layer's Go environment as it now stands: `GOFLAGS` gone, or kept and why.
- In `docs/operators/runbooks/agent-container.md`, under `## Limitations` or `### What it mounts`, write down the cause and what an operator sees when it comes back: `go build: error obtaining VCS status: exit status 128`. Say how to read git's message with `GIT_TRACE`.
- Close the issue with `flai issue close I-0255 --reason "<what fixed it>"`. That regenerates `design/issues/summary.md`.

## Done when

- Both documents describe the layer's Go environment as built, and the markdown lint passes on them.
- I-0255 is closed with a reason that names the fix and T-3036's smoke case, and `design/issues/summary.md` lists it as closed.

## Notes
