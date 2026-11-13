---
id: T-3192
type: task
nature: feature
title: flai image ensure runs the operator's build command in place of its own pull and build, with the image named in variables and the step limit applied
status: backlog
parent: S-0759
owner: alex
created: 2026-11-12T16:05:43Z
updated: 2026-11-12T16:05:43Z
transitions: []
stream: S-0759
tags: [cli]
touches: [flai/internal/agentimage/command.go, flai/internal/agentimage/command_test.go, flai/internal/agentimage/ensure.go, flai/internal/agentimage/ensure_test.go]
after: [T-3190]
---
# T-3192 flai image ensure runs the operator's build command in place of its own pull and build, with the image named in variables and the step limit applied

## Work

`flai/internal/agentimage/command.go` runs the host's `build` for one image:

- from the project's root, with no shell;
- with the image's name, its inner harness, and whether it is the layer in variables named as S-0733's ADR names them (by this plan's assumption `FLAI_IMAGE`, `FLAI_IMAGE_INNER`, and `FLAI_IMAGE_LAYER`);
- through S-0740's runner, so its output is streamed and logged and a silent step is stopped. It is not retried, since flai cannot tell a network failure in another's command.

Make `Ensure` call it in place of the pull and the build when the host names one, then read the image's state again. Report an image the command did not make as a failure naming the command.

It waits for T-3190, whose `build` it reads. It shares no file with T-3191, so the two run together.

## Done when

- `go test ./internal/agentimage/ -run 'Command|Ensure'` passes, with a command that makes the image, one that exits 1, one that exits 0 and makes nothing, and one that falls silent.

## Notes
