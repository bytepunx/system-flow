---
id: T-3209
type: task
nature: improvement
title: scripts/agent-image.sh serves as this project's build command, building from the tree the one image flai asks for, and make agent-images still builds all three
status: backlog
parent: S-0760
owner: alex
created: 2026-11-12T16:06:44Z
updated: 2026-11-12T16:06:55Z
transitions: []
stream: S-0760
tags: [cli]
touches: [scripts/agent-image.sh, scripts/agent-image-test.sh, Makefile]
---
# T-3209 scripts/agent-image.sh serves as this project's build command, building from the tree the one image flai asks for, and make agent-images still builds all three

## Work

Teach `scripts/agent-image.sh` the variables S-0759 passes to a build command (by that story's assumption `FLAI_IMAGE`, `FLAI_IMAGE_INNER`, and `FLAI_IMAGE_LAYER`). Given them, it builds from the tree only the one image named:

- bedrock and claude-code from the tree's `flai` and `pilot`, as its default mode does (lines 202 to 206);
- the layer as `--layer` does (lines 189 to 194).

It exits non-zero when it cannot. Keep `--tag`, `--layer`, and `--check` as they are, and its build with no variables, which `make agent-images` (`Makefile`, lines 69 to 70) runs. Cover the new mode in `scripts/agent-image-test.sh` with a fake `docker`.

Record in the task's Notes the command the operator runs once on the host to name it, `flai serve agent harness container --build scripts/agent-image.sh` for this project. Setting it is the operator's: `flai guard` refuses it to agents.

It shares no file with T-3207 or T-3212, so the three run together.

## Done when

- `scripts/agent-image-test.sh` passes, and ShellCheck passes on both scripts.
- `make agent-images` still builds the three images in its order.

## Notes
