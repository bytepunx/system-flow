---
id: T-2176
type: task
nature: feature
title: flai test runs every per-path tier inside the layered container on this story's worktree, each limitation met is an issue, and the docs describe the layer
status: done
parent: S-0510
owner: alex
created: 2026-10-26T04:51:51Z
updated: 2026-11-02T19:11:49Z
transitions:
  - to: ready
    at: 2026-11-02T18:54:46Z
    by: agent-S-0510
  - to: in-progress
    at: 2026-11-02T18:54:47Z
    by: agent-S-0510
  - to: done
    at: 2026-11-02T19:11:49Z
    by: agent-S-0510
stream: S-0510
tags: [cli, docker]
touches: [docs/users/flai.md, docs/operators/index.md, docs/operators/runbooks/install.md, design/system/flai-cli.md, design/tech/docker.md, design/issues/summary.md, scripts/flai-build.sh, docs/operators/runbooks/agent-container.md, scripts/env.sh]
after: [T-2165, T-2171, T-2172]
usage:
  source: log
  seconds: 1022
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 210
      output: 57030
      cache_read: 24771170
      cache_write: 304162
      cost: 8.148
---
# T-2176 flai test runs every per-path tier inside the layered container on this story's worktree, each limitation met is an issue, and the docs describe the layer

## Work

With the layer built and named in `system-flow.yaml`, run `scripts/flai.sh test` inside `flai-agent-sf-claude-code:local` by hand, on this story's worktree, with paths that select every per-path tier: gofmt, vet, golangci-lint, the Go tests, flaiover-lint, vitest, svelte-check, ShellCheck, and the markdown lint. Mount the worktree read-write, the main checkout's `.git` read-only, and the story's registration read-write, as S-0517 does. Fix in this story what a small change fixes. Record each limitation met as an issue with `flai issue new`, and the image's size in the story's notes. Then write the documentation: `docs/users/flai.md` (the layer and `agent.image.layer`), `docs/operators/index.md` and `docs/operators/runbooks/install.md` (building the layer with `make agent-images`), `design/system/flai-cli.md` (the check and the adapter's choice of image), and `design/tech/docker.md` (this repository's layer and what it carries). Waits for T-2165, T-2171, and T-2172, which build, check, and run the layer.

The live run by `flai serve` on `harness: container` (criterion 5) is asked about on TH-0814. It needs the installed release to carry the layer, so it is not this task's: it becomes a task of its own here, or moves to S-0565, as the answer says.

## Done when

Every per-path tier passes inside the container, each limitation met has an issue, the five documents describe the layer, and `flai check --strict` and the markdown lint pass.

## Notes

A tier that still cannot run inside the layer is a finding for the layer's Dockerfile, fixed in this story if the tool is small, otherwise an issue. Retitled by agent-S-0510 on 2026-11-02, after TH-0814 asked whether the live run moves to S-0565.
