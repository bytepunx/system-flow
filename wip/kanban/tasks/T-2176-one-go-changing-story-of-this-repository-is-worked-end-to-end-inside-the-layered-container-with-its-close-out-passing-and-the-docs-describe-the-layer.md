---
id: T-2176
type: task
nature: feature
title: One Go-changing story of this repository is worked end to end inside the layered container with its close-out passing, and the docs describe the layer
status: backlog
parent: S-0510
owner: alex
created: 2026-10-26T04:51:51Z
updated: 2026-10-26T04:51:51Z
transitions: []
stream: S-0510
tags: [cli, docker]
touches: [docs/users/flai.md, docs/operators/index.md, docs/operators/runbooks/install.md, design/system/flai-cli.md, design/tech/docker.md, design/issues/summary.md]
after: [T-2165, T-2171, T-2172]
---
# T-2176 One Go-changing story of this repository is worked end to end inside the layered container with its close-out passing, and the docs describe the layer

## Work

With the layer built and named in `system-flow.yaml`, set `harness: container` on one small story of this repository that changes Go, and let `flai serve` work it to review on the layered image: `flai test` on its changed paths and `scripts/close-out.sh` must run every selected tier inside the container (gofmt, vet, golangci-lint, the Go tests, the markdown lint, and the flaiover tiers when selected). Record each limitation met as an issue with `flai issue new`, and the image's size and build time in the story's notes. Then write the documentation: `docs/users/flai.md` (the layer and `agent.image.layer`), `docs/operators/index.md` and `docs/operators/runbooks/install.md` (building the layer with `make agent-images`), `design/system/flai-cli.md` (the check and the adapter's choice of image), and `design/tech/docker.md` (this repository's layer and what it carries). Waits for T-2165, T-2171, and T-2172, which build, check, and run the layer.

## Done when

The chosen Go story reached review with its close-out passing inside the container, each limitation met has an issue, the five documents describe the layer, and `flai check --strict` and the markdown lint pass.

## Notes

A tier that still cannot run inside the layer is a finding for the layer's Dockerfile, fixed in this story if the tool is small, otherwise an issue.
