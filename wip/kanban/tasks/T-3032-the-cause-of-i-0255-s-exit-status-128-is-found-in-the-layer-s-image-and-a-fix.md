---
id: T-3032
type: task
nature: research
title: The cause of I-0255's exit status 128 is found in the layer's image and a fix is proposed
status: done
parent: S-0686
owner: alex
created: 2026-11-10T10:33:36Z
updated: 2026-11-11T05:25:13Z
transitions:
  - to: ready
    at: 2026-11-11T05:24:53Z
    by: agent-S-0686
  - to: in-progress
    at: 2026-11-11T05:24:53Z
    by: agent-S-0686
  - to: done
    at: 2026-11-11T05:25:13Z
    by: agent-S-0686
stream: S-0686
tags: [docker, agent-container, go]
touches: [design/issues/I-0255-go-build-s-vcs-stamping-fails-with-git-exit-status-128-inside-the-agent.md]
usage:
  source: log
  seconds: 20
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 5
      output: 1692
      cache_read: 542297
      cache_write: 19899
      cost: 0.3015
---
# T-3032 The cause of I-0255's exit status 128 is found in the layer's image and a fix is proposed

## Work

Reproduce I-0255 before anything is built. It waits for no task.

- Build or reuse the layer's image (`flai/images/system-flow/Dockerfile`) and start it on a story worktree with `scripts/agent-container.sh --probe`, in the mounts a story's container has today.
- In it, run the build `builtFlai` runs (`flai/internal/mcpserver/subproject_test.go`), `go build -o <tmp>/flai github.com/bytepunx/system-flow/flai` from `flai/`, with `GOFLAGS=` and `-buildvcs=true`, `-x`, and `GIT_TRACE` and `GIT_TRACE_SETUP` written to a file. The file holds git's own message, which Go's error text leaves out.
- Run it again as the container harness starts a story (`flai/internal/harness/container.go`), with its `GIT_CONFIG_*` and mounts, and once with the image's default user, 65532.
- Check whether the cause is gone. I-0255's instance (S-0510, 2026-11-02) came before the main checkout's mounts and their `safe.directory` (ADR-0253, S-0628). If the build stamps cleanly now, say which change fixed it.
- Write the cause, git's message, and the fix you propose as a new instance of I-0255. Write it in the story's narrative `## Decisions` too. If the fix needs a path outside the story's touches, widen them. If the fix changes the container's mounts or what a story's agent may write, ask on a thread on S-0686 first.

## Done when

- git's own message behind `exit status 128` is recorded in I-0255, with the command and the mounts that reproduce it, or the record says it no longer reproduces and names the change that fixed it.
- The proposed fix is named in the narrative's `## Decisions`, with the paths it changes.

## Notes
