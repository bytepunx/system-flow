---
id: S-0510
type: story
nature: feature
title: A project's manifest names an image layer for its test tiers' tools, flai serve runs it, flai check warns when it is not built, and this repository's layer runs the close-out inside the container
status: backlog
parent: E-0029
owner: alex
created: 2026-10-26T04:46:53Z
updated: 2026-10-26T04:57:31Z
transitions: []
tags: [cli, docker]
topics: [harness, docker, manifest]
touches: [flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/check/check.go, flai/internal/check/image.go, flai/internal/check/image_test.go, flai/internal/harness/container.go, flai/internal/harness/container_test.go, flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, scripts/agent-image.sh, Makefile, scripts/README.md, flai/images/system-flow/Dockerfile, system-flow.yaml, design/system/project-manifest.md, docs/operators/settings.md, template/root/system-flow.yaml.tmpl, template/CHANGELOG.md, docs/users/flai.md, docs/operators/index.md, docs/operators/runbooks/install.md, design/system/flai-cli.md, design/tech/docker.md, design/issues/summary.md, design/issues]
after: [S-0508]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
forecast:
  duration: 2h30m
  delivery: 2026-10-27T09:00:00Z
  basis: "flai forecast gave 56m (median 107 s per unit over 93 done feature stories on claude-opus-5-5 in the large band, size 31: 6 criteria, 25 touches) and delivery 2026-10-27T07:08Z, 29th in the pull order by a cycle factor of 1.17; raised to 2h30m because a layer with Go, golangci-lint, Node, and pnpm takes minutes per build and is built several times, and criterion 5 runs a Go-changing story of this repository inside the container with its close-out's Go tiers, a session of its own of 30 to 60 minutes; delivery is flai's start of 2026-10-27T06:02Z plus 2h30m times the same cycle factor."
  by: planner-E-0029
  at: 2026-10-26T04:57:31Z
---
# S-0510 A project's manifest names an image layer for its test tiers' tools, flai serve runs it, flai check warns when it is not built, and this repository's layer runs the close-out inside the container

## Goal

ADR-0203's limitation: only the worktree's files are the container's to write, and a test that needs a tool outside the image fails inside it until the project's layer carries the tool. The bedrock and claude-code images carry no Go toolchain, no golangci-lint, and no pnpm, so a story of this repository that changes Go or Svelte cannot run `flai test` or `scripts/close-out.sh` inside the container; S-0505's results name each tier that failed for that reason. This story gives a project one image layer of its own: the manifest names a Dockerfile under `agent.image.layer`, `scripts/agent-image.sh` builds it over the harness image as `flai-agent-<project key>-<inner>:local`, `flai serve` runs that image for the project's agents when it is named, and `flai check` warns, as `agent.image`, when the manifest names a layer that is not built, since `docker image inspect` finds no such image. This repository's layer, `flai/images/system-flow/Dockerfile`, carries Go at the version `flai/go.mod` names, golangci-lint v2, Node, and pnpm, so that `flai test` and the close-out run inside the container, which one Go-changing story worked end to end on `harness: container` proves.

## Acceptance criteria

- [ ] `agent.image.layer` in `system-flow.yaml` names a Dockerfile path relative to the repository root, read by `flai/internal/manifest`, documented in `design/system/project-manifest.md` and `docs/operators/settings.md`, and left as a comment in `template/root/system-flow.yaml.tmpl` with `template/CHANGELOG.md` saying so
- [ ] `scripts/agent-image.sh --layer` builds the named Dockerfile with the harness image as its base argument and tags it `flai-agent-<project key>-<inner>:local`, `make agent-images` builds it after the two base images when the manifest names one, and the container adapter runs that image in place of the harness image for the project's agents when it is named
- [ ] `flai check` reports `agent.image` when the manifest names a layer whose image is not built or whose Dockerfile does not exist, as a warning that names `make agent-images`, and stays silent when the manifest names none or `docker` is not on `PATH`, saying so as a note
- [ ] `flai/images/system-flow/Dockerfile` carries Go at `flai/go.mod`'s version, golangci-lint v2 at the version `scripts/install-tools.sh` installs, Node 24, and pnpm, and `system-flow.yaml` names it, so that `scripts/flai.sh test` and `scripts/close-out.sh` run every tier the manifest's `tests` select inside the container
- [ ] One story of this repository that changes Go is worked end to end by `flai serve` on `harness: container` with this layer, its close-out passing inside the container, and each limitation met is an issue under `design/issues`
- [ ] `docs/users/flai.md`, `docs/operators/index.md`, `docs/operators/runbooks/install.md`, `design/system/flai-cli.md`, and `design/tech/docker.md` describe the layer, how it is built, and what this repository's carries

## Tasks
- T-2164 agent.image.layer in the manifest names a project's Dockerfile, read by flai and documented in project-manifest.md, settings.md, and the template
- T-2165 scripts/agent-image.sh --layer builds a project's layer over the harness image, make agent-images builds it when named, and this repository's layer carries Go, golangci-lint, Node, and pnpm
- T-2171 flai check reports agent.image when the manifest names a layer that is not built or whose Dockerfile is missing, and notes when docker is not on PATH
- T-2172 The container adapter runs the project's layer image in place of the harness image when the manifest names one
- T-2176 One Go-changing story of this repository is worked end to end inside the layered container with its close-out passing, and the docs describe the layer

## Notes

Waits on S-0508 for the adapter that runs the image. The layer is the project's risk, as ADR-0203 says; the warning on an unsigned layer it foresees waits for the signed image, which is outside E-0029. A layer's tools are a project's: the template ships none.

### Planning

Planned by planner-E-0029 on 2026-10-26 with the epic.

Touches. `flai touches suggest`, seeded with `flai/internal/manifest/manifest.go` and `flai/internal/check/check.go` (89 commits), names `manifest_test.go`, `settings.md`, `docs/users/flai.md`, `flai-cli.md`, `project-manifest.md`, and `docs/operators/index.md`, each kept (co-change), and `check_test.go`, `manifest/settings.go`, `cmd/root.go`, `flaiover-dashboard.md`, `workitem/item.go`, `flai-reference.md`, and `mcpserver/server.go`, left out because the key is not settable from the dashboard, no command changes, and the check has its own file and test. The rest are the tasks' (layout):

- `flai/internal/check/image.go`, `image_test.go`: the new check, one file per check as the folder has them. T-2171.
- `flai/internal/harness/container.go`, `container_test.go`, `flai/internal/serve/agents.go`, `agents_test.go`: S-0508's files, which gain the layer's image. T-2172.
- `scripts/agent-image.sh`, `Makefile`, `scripts/README.md`, `flai/images/system-flow/Dockerfile`, `system-flow.yaml`: the build and this repository's layer. T-2165.
- `template/root/system-flow.yaml.tmpl`, `template/CHANGELOG.md` (design, criterion 1): the template's manifest comment and its changelog. T-2164.
- `docs/operators/runbooks/install.md`, `design/tech/docker.md` (design, criterion 6). T-2176.
- `design/issues/summary.md` and the folder `design/issues` (layout, folder kept): `flai issue new` in T-2176 writes issue files no task can name; the folder is kept for them alone, and `flai task done` widens the touches with the committed paths (ADR-0144).

Forecast and cost of delay: see the second Planning paragraph, written once the figures are set.
