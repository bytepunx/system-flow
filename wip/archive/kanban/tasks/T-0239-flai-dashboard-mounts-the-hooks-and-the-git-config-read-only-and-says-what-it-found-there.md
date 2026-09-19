---
id: T-0239
type: task
nature: feature
title: flai dashboard mounts the hooks and the git config read-only, and says what it found there
status: done
parent: S-0064
owner: alex
created: 2026-09-19T10:18:10Z
updated: 2026-09-19T10:24:51Z
transitions:
  - to: ready
    at: 2026-09-19T10:22:07Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T10:22:07Z
    by: system-flow
  - to: done
    at: 2026-09-19T10:24:51Z
    by: system-flow
stream: S-0064
tags: []
---

# T-0239 flai dashboard mounts the hooks and the git config read-only, and says what it found there

## Work
`flai dashboard` adds read-only mounts over the read-write clone for what T-0238 found can be closed that way, creating an empty hooks directory on the host when there is none so there is something to mount. It says in its start message that the hooks and the git config are read-only in the container. Where the clone is a linked worktree or `.git` is not a directory it says so and mounts nothing extra. Tests with the fake runner for the mounts and the fallbacks.

## Done when
- The mounts are in the docker arguments and tested
- `make flai-test` passes

## Notes
