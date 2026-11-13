---
id: T-3194
type: task
nature: feature
title: The image job runs the operator's build command when the host names one, and a start refuses a missing image when the automatic job is off
status: backlog
parent: S-0759
owner: alex
created: 2026-11-12T16:05:51Z
updated: 2026-11-13T16:15:46Z
transitions: []
stream: S-0759
tags: [cli]
touches: [flai/internal/serve/imagejob.go, flai/internal/serve/imagejob_test.go, flai/internal/serve/container.go, flai/internal/serve/container_test.go]
after: [T-3190, T-3192]
---
# T-3194 The image job runs the operator's build command when the host names one, and a start refuses a missing image when the automatic job is off

## Work

- `flai/internal/serve/imagejob.go`: a job runs `Ensure` with the project's command from the host's configuration (`AgentStart.ImageCommand(root)`, T-3190), so T-3192's runner runs it in place of the pull and build. The job's record and journal say the command ran.
- `flai/internal/serve/container.go`: with `AutoImage(root)` false, a start that finds its image missing or stale requests no job and refuses as before S-0751, naming the image, `flai image ensure`, and `flai serve agent image auto on`. This holds for a story's agent, a strategic agent, the close-out's container, and the checks' containers, wherever S-0751 holds them for a job.

Tests cover:

- a job through the command;
- a start refused with the setting off;
- a start held and gone on with it on.

It waits for T-3190, whose settings it reads, and T-3192, whose command `Ensure` runs.

## Done when

- `flai test` on the changed files passes; run `flai test race` too, since the job is concurrent.

## Notes
