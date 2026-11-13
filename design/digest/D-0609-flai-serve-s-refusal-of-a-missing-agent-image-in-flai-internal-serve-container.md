---
id: D-0609
title: flai serve's refusal of a missing agent image in flai/internal/serve/container.go still names make agent-images, left on purpose for S-0751, which replaces the refusal with a host job
kind: decision
sources: [S-0745, S-0751]
topics: [docker, code, go]
paths: [flai/internal/serve/container.go, flai/internal/check/image.go, docs/operators/runbooks/install.md]
date: 2026-11-13
status: absorbed
absorbed_by: "design/system/flai-cli.md#the-container-harness-s-0508"
asked_again: []
cited: [S-0751]
---

# D-0609 flai serve's refusal of a missing agent image in flai/internal/serve/container.go still names make agent-images, left on purpose for S-0751, which replaces the refusal with a host job

S-0745 made `flai image ensure` build the project's layer in any project. It pointed `flai check`'s advice for a missing or stale agent image, and `docs/operators/runbooks/install.md`, at `flai image ensure` in place of `make agent-images`, a target only system-flow's repository has. It left one place unchanged. From S-0745's Decisions (`wip/archive/agents/S-0745.md`):

```text
T-3176 also pointed `flai check`'s advice for a missing or stale agent image, and `docs/operators/runbooks/install.md`'s paragraph on a project of your own, at `flai image ensure` in place of `make agent-images`. The refusal in `flai/internal/serve/container.go` still names `make agent-images`; it is left for S-0751, which replaces that refusal with a host job.
```

S-0751 did that, and the refusal is gone. From S-0751's Decisions (`wip/archive/agents/S-0751.md`):

```text
D-0609 cited: the make agent-images refusal in container.go was left for this story, and it is gone.
```

Since S-0751, a start in `flai serve` that finds an image missing or stale asks for the image job and is held until it ends. A start outside `flai serve`, by hand through the CLI or MCP, has no job: a missing image refuses, naming `flai image ensure` and `flai image status`. No refusal names `make agent-images` any more. `design/system/flai-cli.md` § The container harness (S-0508) records this.

Sources: S-0745, S-0751.
