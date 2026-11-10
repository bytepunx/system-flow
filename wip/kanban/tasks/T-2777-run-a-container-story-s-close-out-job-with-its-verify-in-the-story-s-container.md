---
id: T-2777
type: task
nature: improvement
title: Run a container story's close-out job with its verify in the story's container and its commit and sync on the host through flai alone
status: done
parent: S-0628
owner: alex
created: 2026-11-09T21:58:30Z
updated: 2026-11-10T06:41:42Z
transitions:
  - to: ready
    at: 2026-11-10T06:15:10Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T06:15:11Z
    by: system-flow
  - to: done
    at: 2026-11-10T06:41:42Z
    by: agent-S-0628
stream: S-0628
tags: [flai, serve, security]
touches: [flai/internal/serve/verifyjob.go, flai/internal/serve/verifyjob_test.go, flai/internal/serve/container.go, flai/internal/serve/verifyjob_container.go]
after: [T-2775]
usage:
  source: log
  seconds: 1591
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 74
      output: 16255
      cache_read: 12911339
      cache_write: 130936
      cost: 3.7074
---
# T-2777 Run a container story's close-out job with its verify in the story's container and its commit and sync on the host through flai alone

## Work

In `flai/internal/serve/verifyjob.go`, `command` and `runJob` decide the job's run from the story's agent, the story's own over the project's default (`Agent.With`, ADR-0251). When the harness is `container`, the job never execs anything from the worktree on the host. Its steps:

1. It runs the image's flai, `flai verify <story> --record-issues`, through T-2775's argv. The image is the one `containerImageHost` in `flai/internal/serve/container.go` decides for the story's agent: the project's layer, else the host's image. Expose that choice for the job rather than deciding it twice.
2. It commits what is outstanding on the host with git and the request's message, as `scripts/close-out.sh`'s commit step does.
3. It runs the host's own `flai verify <story> --sync-only`, with this flai's executable (`p.exe()`).
4. It syncs and verifies again in the container, at most three rounds in all, as the script does (ADR-0156).

The container can write the worktree's `.git` file and its registration under `.git/worktrees/<story>`, including `config.worktree`. So the host's git commands do not take the repository from them. They name the git dir and work tree from the main checkout's own record, and set `core.hooksPath` to the main checkout's `.git/hooks`. A `core.hooksPath` or `core.fsmonitor` the container wrote then runs nothing on the host.

The same refusals as `containerAllowed` apply when docker or the image is missing. Each names its fix, and the verdict is `failed` with that reason, not a host run. A story on any other harness keeps today's `sh scripts/close-out.sh` run. Waits for T-2775, whose argv it runs.

## Done when

- A test in `flai/internal/serve/verifyjob_test.go` reproduces I-0295 with a fake `closeOutRunner`. For a story whose agent's harness is `container`, the job runs `docker` with the story's mounts and `flai verify --record-issues`. It never runs `sh scripts/close-out.sh`, or any program from the worktree, on the host. The test fails on today's code.
- A test writes a `core.hooksPath` hook into the registration's `config.worktree`, and the host's commit step does not run it.
- Tests cover the commit and `--sync-only` steps on the host, a missing image refused with its fix, and a host-harness story keeping `sh scripts/close-out.sh`.
- `flai test` on the three paths passes.

## Notes
