---
id: T-3039
type: task
nature: remediation
title: "A story's container stamps go builds again, and the layer drops GOFLAGS=-buildvcs=false"
status: backlog
parent: S-0686
owner: alex
created: 2026-11-10T10:34:01Z
updated: 2026-11-10T10:34:01Z
transitions: []
stream: S-0686
tags: [docker, agent-container, go]
touches: [flai/images/system-flow/Dockerfile, scripts/agent-container.sh, flai/internal/harness/container.go, flai/internal/harness/container_test.go]
after: [T-3036]
---
# T-3039 A story's container stamps go builds again, and the layer drops GOFLAGS=-buildvcs=false

## Work

Build the fix T-3032 proposed. It waits for T-3036, so that the smoke case shows the fix works and does not just pass by chance.

- Fix the cause where T-3032 found it. For a mount or a `safe.directory` the container lacks, fix both starts: `scripts/agent-container.sh`'s `run` and `flai/internal/harness/container.go` (`safeMain` and the mounts), with a case in `container_test.go`. If T-3032 found the cause is gone, change no start.
- Remove `GOFLAGS=-buildvcs=false` and its I-0255 comment from `flai/images/system-flow/Dockerfile`. Keep the other `ENV` lines. If the cause lies in Go or git and cannot be fixed here, keep the flag, rewrite the comment with the cause, and say why in the narrative's `## Decisions`.
- Rebuild the layer with `scripts/agent-image.sh --layer`, then check it with `scripts/agent-image.sh --check S-0686`.

## Done when

- T-3036's case passes in the rebuilt layer's image, with no `GOFLAGS=-buildvcs=false` set there, unless the narrative records why the flag stays.
- `TestSubprojectEditAddsSetsAndRemovesCommittingNothing` and `TestSubprojectEditWorksInTheStorysWorktree` pass inside the container.
- `flai test` passes on the changed paths.

## Notes
