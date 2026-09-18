---
id: T-0152
type: task
nature: remediation
title: flai dashboard mounts the repository at its host path and sets PROJECT_DIR to it
status: done
parent: S-0050
owner: alex
created: 2026-09-18T19:47:23Z
updated: 2026-09-18T19:50:38Z
transitions:
  - to: ready
    at: 2026-09-18T19:49:24Z
    by: alex
  - to: in-progress
    at: 2026-09-18T19:49:24Z
    by: alex
  - to: done
    at: 2026-09-18T19:50:38Z
    by: alex
stream: S-0050
tags: []
touches: [flai/cmd]
---

# T-0152 flai dashboard mounts the repository at its host path and sets PROJECT_DIR to it

## Work
In `flai/cmd/dashboard.go`, mount `--volume <root>:<root>` and set `PROJECT_DIR=<root>`. Decide whether `/project` stays as a second mount and record it in the narrative. When the host path cannot be a container path (not absolute in the container's terms, a Windows drive path being the case), keep today's `/project` mount, warn at start that acceptance of stories with a branch will not work from the board, and name the setting and `flai accept`. Update the printed and JSON output, which name the mount. Extend the dashboard command tests to pin the volume and `PROJECT_DIR` arguments for both cases.

## Done when
- The new argument tests fail without the change and pass with it
- `go test -race ./cmd/...` passes

## Notes
