---
id: T-3194
type: task
nature: feature
title: The image job runs the operator's build command when the host names one, and a start refuses a missing image when the automatic job is off
status: backlog
parent: S-0759
owner: alex
created: 2026-11-12T16:05:51Z
updated: 2026-11-12T16:05:51Z
transitions: []
stream: S-0759
tags: [cli]
touches: [flai/internal/serve/imagejob.go, flai/internal/serve/imagejob_test.go, flai/internal/serve/container.go, flai/internal/serve/container_test.go]
after: [T-3190, T-3192]
---
# T-3194 The image job runs the operator's build command when the host names one, and a start refuses a missing image when the automatic job is off

## Work

- `flai/internal/serve/imagejob.go`: a job runs `Ensure`, which runs T-3192's command when the host names one. The journal's entries say `command` as the job's way.
- `flai/internal/serve/container.go`: with `auto_image` off for the project, a start that finds its image missing or stale refuses as before S-0751, naming `flai image ensure` and `--auto-image`. No job is requested.

Tests cover:

- a job through the command;
- a start refused with the setting off;
- a start held and gone on with it on.

It waits for T-3190, whose setting it reads, and T-3192, whose command `Ensure` runs.

## Done when

- `go test ./internal/serve/ -run 'ImageJob|Container'` passes, under `-race` too.

## Notes
