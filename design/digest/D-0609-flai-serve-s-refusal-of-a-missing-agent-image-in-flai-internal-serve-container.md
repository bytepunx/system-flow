---
id: D-0609
title: flai serve's refusal of a missing agent image in flai/internal/serve/container.go still names make agent-images, left on purpose for S-0751, which replaces the refusal with a host job
kind: decision
sources: [S-0745]
topics: [docker, code, go]
paths: [flai/internal/serve/container.go, flai/internal/check/image.go, docs/operators/runbooks/install.md]
date: 2026-11-13
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0609 flai serve's refusal of a missing agent image in flai/internal/serve/container.go still names make agent-images, left on purpose for S-0751, which replaces the refusal with a host job

S-0745 made `flai image ensure` build the project's layer in any project. It pointed `flai check`'s advice for a missing or stale agent image, and `docs/operators/runbooks/install.md`, at `flai image ensure` in place of `make agent-images`, a target only system-flow's repository has. It left one place unchanged. From S-0745's Decisions (`wip/archive/agents/S-0745.md`):

```text
T-3176 also pointed `flai check`'s advice for a missing or stale agent image, and `docs/operators/runbooks/install.md`'s paragraph on a project of your own, at `flai image ensure` in place of `make agent-images`. The refusal in `flai/internal/serve/container.go` still names `make agent-images`; it is left for S-0751, which replaces that refusal with a host job.
```

So the two errors in `flai/internal/serve/container.go` that refuse a start on a missing harness or layer image still tell the operator to run `make agent-images` or `scripts/agent-image.sh --layer`. A story other than S-0751 that touches that file should leave the wording to S-0751, rather than reword a refusal that S-0751 removes. S-0751 should drop the `make agent-images` advice with it.

Source: S-0745.
