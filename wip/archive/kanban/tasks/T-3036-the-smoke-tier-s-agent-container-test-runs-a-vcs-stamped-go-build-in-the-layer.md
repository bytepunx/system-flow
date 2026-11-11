---
id: T-3036
type: task
nature: remediation
title: The smoke tier's agent-container test runs a VCS-stamped go build in the layer's image
status: done
parent: S-0686
owner: alex
created: 2026-11-10T10:33:48Z
updated: 2026-11-11T05:26:30Z
transitions:
  - to: ready
    at: 2026-11-11T05:25:30Z
    by: agent-S-0686
  - to: in-progress
    at: 2026-11-11T05:25:30Z
    by: agent-S-0686
  - to: done
    at: 2026-11-11T05:26:30Z
    by: agent-S-0686
stream: S-0686
tags: [docker, agent-container, go, test]
touches: [scripts/agent-vcs-test.sh, scripts/smoke.sh, scripts/README.md]
after: [T-3032]
usage:
  source: log
  seconds: 60
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 6
      output: 1800
      cache_read: 576950
      cache_write: 21170
      cost: 0.3208
---
# T-3036 The smoke tier's agent-container test runs a VCS-stamped go build in the layer's image

## Work

Add a case to `scripts/agent-container-test.sh`, which the smoke tier runs and which already reproduces I-0279. The case reproduces I-0255. It waits for T-3032, since the case runs the command and mounts that T-3032 found.

- In a story's mounts, it runs the build `builtFlai` runs, `go build` of `github.com/bytepunx/system-flow/flai` from `flai/`, with `GOFLAGS=` and `-buildvcs=true`. Use the layer's image.
- It fails when the build exits non-zero, and it prints the build's error and git's own message from a `GIT_TRACE` file.
- It skips with the reason, as the script's other cases do, when Docker or the layer's image is missing.
- If T-3032 found that the cause no longer occurs, the case passes now and guards against a return. Otherwise it fails until T-3039's fix. Say which in the task's `## Notes`.

## Done when

- `scripts/agent-container-test.sh` holds the case, and ShellCheck passes on it.
- Run in the layer's image, the case fails as I-0255 describes before the fix, or passes and guards against a return when the cause is gone. The run's result is in `## Notes`.

## Notes
