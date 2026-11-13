---
id: T-3192
type: task
nature: feature
title: flai image ensure runs the operator's build command in place of its own pull and build, with the image named in variables and the step limit applied
status: backlog
parent: S-0759
owner: alex
created: 2026-11-12T16:05:43Z
updated: 2026-11-13T16:15:42Z
transitions: []
stream: S-0759
tags: [cli]
touches: [flai/internal/agentimage/command.go, flai/internal/agentimage/command_test.go, flai/internal/agentimage/ensure.go, flai/internal/agentimage/ensure_test.go, flai/cmd/image.go, flai/cmd/image_test.go]
after: [T-3190]
---
# T-3192 flai image ensure runs the operator's build command in place of its own pull and build, with the image named in variables and the step limit applied

## Work

`flai/internal/agentimage/command.go` runs the operator's build command (ADR-0295 point 7):

- in the project's main checkout, with no shell;
- with the variables ADR-0295 names: `FLAI_IMAGE_VERSION`, `FLAI_IMAGE_INNER`, `FLAI_IMAGE_BEDROCK`, `FLAI_IMAGE_HARNESS`, `FLAI_IMAGE_LAYER` (empty when the manifest names no layer), `FLAI_IMAGE_RECIPE_BEDROCK`, `FLAI_IMAGE_RECIPE_HARNESS`, `FLAI_IMAGE_RECIPE_LAYER`, and `FLAI_IMAGE_STALE` (the images found missing or stale, of `bedrock`, `harness`, `layer`, space-separated);
- through S-0740's `Runner`, so its output is streamed and logged as flai's own build's is, and a step silent for the stall limit is stopped. It is not retried: flai cannot tell a network failure inside another's command.

`Ensurer.Ensure` gains the command (a field, nil for none). When set, it runs it in place of the pull, the build, and the layer's build, when an image is missing or stale (in `build` mode always), then judges the images again with `Check`. A non-zero exit, or an image still missing or stale, fails, naming the command and the image. The operator's `--image` stays untouched. A development flai's build from the tree is left as it is: S-0760 makes this repository's script its command.

`flai image ensure` and `flai image build` (`flai/cmd/image.go`) read the project's command from the host's configuration, `AgentStart.ImageCommand(root)` from T-3190, and pass it on. Test that in `image_test.go`.

It waits for T-3190, whose field it reads. It shares no file with T-3191, so the two run together.

## Done when

- `flai test` on the changed files passes, with a command that makes the image, one that exits 1, one that exits 0 and makes nothing, and one that falls silent, and the variables checked.

## Notes
