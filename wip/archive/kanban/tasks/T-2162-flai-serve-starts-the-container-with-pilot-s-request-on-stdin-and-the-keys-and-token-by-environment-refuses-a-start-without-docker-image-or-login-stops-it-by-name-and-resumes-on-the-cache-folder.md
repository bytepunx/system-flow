---
id: T-2162
type: task
nature: feature
title: flai serve starts the container with pilot's request on stdin and the keys and token by environment, refuses a start without docker, image, or login, stops it by name, and resumes on the cache folder
status: done
parent: S-0508
owner: alex
created: 2026-10-26T04:49:23Z
updated: 2026-11-02T07:25:41Z
transitions:
  - to: ready
    at: 2026-11-02T07:10:20Z
    by: system-flow
  - to: in-progress
    at: 2026-11-02T07:10:21Z
    by: system-flow
  - to: done
    at: 2026-11-02T07:25:41Z
    by: agent-S-0508
stream: S-0508
tags: [cli, docker]
touches: [flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, flai/internal/serve/start.go, flai/internal/serve/start_test.go, flai/internal/serve/stop.go, flai/internal/serve/stop_test.go, flai/internal/serve/restart.go, flai/internal/serve/claudecheck.go, flai/internal/serve/claudecheck_test.go, flai/internal/serve/container.go, flai/internal/serve/container_test.go]
after: [T-2155, T-2156]
usage:
  source: log
  seconds: 920
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 106
      output: 38510
      cache_read: 8793387
      cache_write: 292237
      cost: 4.496
---
# T-2162 flai serve starts the container with pilot's request on stdin and the keys and token by environment, refuses a start without docker, image, or login, stops it by name, and resumes on the cache folder

## Work

In `agents.go`, make `startOf` and `spawn` feed `Start.Stdin` to the process when the adapter set it, pass `FLAI_MCP_URL` and `FLAI_MCP_TOKEN` (from S-0506) and the provider keys in the docker client's environment so each `-e NAME` carries its value, and create `.flai-cache/agents/<S-nnnn>` before the start so the mount exists. Add to `startAllowed`, for the container adapter, the refusals with a message naming what to do: `docker` not on `PATH`, the image for the inner harness absent from `docker image inspect` (naming `make agent-images`), and no login path set for an inner harness that needs one. Before the start, read the inner harness's version from its image, for claude-code `docker run --rm --entrypoint claude <image> --version`, as ADR-0106's check reads the host's `claude` (TH-0666): record it on the run, and when `claude-checks.json` holds no passed record for that version, log a warning naming the version and the image rather than start the headless check; the check inside the image is left to an issue T-2169 files. In `stop.go`, when the run's harness is `container`, stop the container by its name (`docker stop flai-agent-<S-nnnn>-<session>`) after signalling the `docker run` process, so a stop from the dashboard (ADR-0058) ends the agent inside. In `restart.go`, a restart after a thread is answered starts a container on the same cache folder and session, so the inner harness resumes as `resumes()` has it. Test each with a fake `docker` on `PATH` in `agents_test.go`, `start_test.go`, and `stop_test.go`. Waits for T-2155, whose `Start` gives the argv and `Stdin`, and T-2156, whose settings name the image and the login.

## Done when

A story on `harness: container` is started by `flai serve` with the request on the container's standard input and no key or token in argv, each refusal is tested, a stop ends the container, a restart resumes, the image's Claude Code version is read and recorded, and `flai test flai/internal/serve` passes.

## Notes

`FLAI_STORY` is set by `spawn` today and reaches the container through `-e FLAI_STORY`; the strategic roles' variables are S-0509's.

The image's version step was added by S-0508's agent on pull, as TH-0666 settled it (the orchestrator's recommendation, confirmed by flai serve on 2026-10-28).
