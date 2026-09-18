---
id: T-0164
type: task
nature: improvement
title: flai dashboard gives the container the host's global git excludes
status: done
parent: S-0051
owner: alex
created: 2026-09-18T21:00:23Z
updated: 2026-09-18T21:05:14Z
transitions:
  - to: ready
    at: 2026-09-18T21:04:13Z
    by: alex
  - to: in-progress
    at: 2026-09-18T21:04:14Z
    by: alex
  - to: done
    at: 2026-09-18T21:05:14Z
    by: alex
stream: S-0051
tags: []
touches: [flai/cmd]
---

# T-0164 flai dashboard gives the container the host's global git excludes

## Work
In `flai/cmd/dashboard.go`, find the host's global excludes file: `git config --global --get core.excludesFile`, else `$XDG_CONFIG_HOME/git/ignore`, else `~/.config/git/ignore`, expanding `~`. When it exists, mount it read-only in the container and point git at it with `GIT_CONFIG_COUNT`, `GIT_CONFIG_KEY_0=core.excludesFile`, `GIT_CONFIG_VALUE_0=<target>`, beside the identity arguments. When there is none, pass nothing. Tests with the fake runner for: configured file, default file, none.

## Done when
- The three cases are pinned by tests on the `docker run` arguments
- `go test -race ./cmd/...` passes

## Notes
