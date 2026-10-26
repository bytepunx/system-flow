---
id: T-2169
type: task
nature: feature
title: One story of this repository is worked end to end by flai serve on harness container, and the docs describe the adapter, its settings, its mounts, and what it never mounts
status: backlog
parent: S-0508
owner: alex
created: 2026-10-26T04:50:37Z
updated: 2026-10-26T04:50:37Z
transitions: []
stream: S-0508
tags: [cli, docker]
touches: [docs/users/flai.md, docs/operators/index.md, docs/operators/runbooks/install.md, design/system/flai-cli.md, design/system/agent-adapters.md, design/tech/docker.md, design/issues/summary.md]
after: [T-2162, T-2157]
---
# T-2169 One story of this repository is worked end to end by flai serve on harness container, and the docs describe the adapter, its settings, its mounts, and what it never mounts

## Work

Set `harness: container` on one small story of this repository (`flai edit --harness container`, or ask the operator to), with the host settings from T-2156 in place and the images built, and let `flai serve` start it: from `story_start` to review, with its usage measured, its card's status shown, and a stop and restart tried once. Record each limitation met as an issue with `flai issue new`, naming the story of E-0029 that carries it on. Then write the documentation: `docs/users/flai.md` § agents and § harnesses (the `container` harness, `config.inner`, what it mounts and never mounts, that Claude Code inside runs on the operator's login mounted read-only), `docs/operators/index.md` and `docs/operators/runbooks/install.md` (Docker Engine as a requirement for the harness and `make agent-images`), `design/system/flai-cli.md` (the adapter in the package map and the serve start path), `design/system/agent-adapters.md` (the container harness as built, in the table that compares it), and `design/tech/docker.md`. Waits for T-2162 and T-2157, which finish the behaviour it documents and proves.

## Done when

The chosen story reached review on `harness: container` with usage on its item, each limitation met has an issue, the six documents describe the adapter as the tests check it, and `flai check --strict` and the markdown lint pass.

## Notes

The proof story should change documentation or shell only: the Go tiers need the project layer of S-0510.
