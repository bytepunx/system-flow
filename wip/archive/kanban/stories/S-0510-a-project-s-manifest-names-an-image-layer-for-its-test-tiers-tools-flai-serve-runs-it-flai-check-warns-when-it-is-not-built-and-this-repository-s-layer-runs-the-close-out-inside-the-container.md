---
id: S-0510
type: story
nature: feature
title: A project's manifest names an image layer for its test tiers' tools, flai serve runs it, flai check warns when it is not built, and this repository's layer runs the close-out inside the container
status: done
parent: E-0029
owner: alex
created: 2026-10-26T04:46:53Z
updated: 2026-11-03T02:39:15Z
transitions:
  - to: ready
    at: 2026-11-02T18:12:39Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-02T18:16:07Z
    by: agent-S-0510
  - to: review
    at: 2026-11-03T01:26:52Z
    by: agent-S-0510
  - to: done
    at: 2026-11-03T02:39:15Z
    by: alex
blocked:
  - from: 2026-11-02T19:32:01Z
    until: 2026-11-02T19:34:36Z
    reason: awaiting the operator on TH-0814
  - from: 2026-11-02T19:35:01Z
    until: 2026-11-02T19:42:28Z
    reason: awaiting the operator on TH-0814
  - from: 2026-11-02T19:42:56Z
    until: 2026-11-02T19:52:45Z
    reason: awaiting the operator on TH-0814
  - from: 2026-11-02T19:53:10Z
    until: 2026-11-02T20:04:55Z
    reason: awaiting the operator on TH-0814
  - from: 2026-11-02T20:05:11Z
    until: 2026-11-02T20:31:49Z
    reason: awaiting the operator on TH-0814
  - from: 2026-11-02T20:32:25Z
    until: 2026-11-02T20:35:57Z
    reason: awaiting the operator on TH-0814
  - from: 2026-11-02T20:36:11Z
    until: 2026-11-02T21:37:57Z
    reason: awaiting the operator on TH-0814
  - from: 2026-11-02T21:38:08Z
    until: 2026-11-02T23:25:31Z
    reason: awaiting the operator on TH-0814
tags: [cli, docker]
topics: [harness, docker, manifest]
touches: [flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/check/check.go, flai/internal/check/image.go, flai/internal/check/image_test.go, flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, scripts/agent-image.sh, Makefile, scripts/README.md, flai/images/system-flow/Dockerfile, system-flow.yaml, design/system/project-manifest.md, docs/operators/settings.md, template/root/system-flow.yaml.tmpl, template/CHANGELOG.md, docs/users/flai.md, docs/operators/index.md, docs/operators/runbooks/install.md, design/system/flai-cli.md, design/tech/docker.md, design/issues/summary.md, design/issues, flai/internal/manifest/agent.go, flai/internal/manifest/agent_test.go, flai/internal/workitem/front-matter-fields.txt, flai/internal/workitem/store_test.go, flai/internal/serve/container.go, flai/internal/serve/container_test.go, flai/cmd/check.go, flai/cmd/check_stats_test.go, flai/internal/check/story.go, flai/internal/verify/story.go, design/issues/I-0255-go-build-s-vcs-stamping-fails-with-git-exit-status-128-inside-the-agent.md, scripts/env.sh, scripts/flai-build.sh, design/issues/I-0259-a-story-worktree-s-flaiover-node-modules-installed-inside-the-alpine-agent.md]
after: [S-0508]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 5201
  priced_by: harness
  turns:
    - day: 2026-11-02
      ceremony: 5
      test_runs: 5
      hand_edits: 3
      work: 148
    - day: 2026-11-03
      ceremony: 4
      work: 43
  models:
    - model: claude-opus-5-5
      input: 764
      output: 190568
      cache_read: 84941918
      cache_write: 1563893
      cost: 32.1654
  strategic:
    - kind: orchestrator
      seconds: 341
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 73
          output: 575
          cache_read: 3505965
          cache_write: 297280
          cost: 0.8695
        - model: claude-sonnet-5-5
          input: 10
          output: 66
          cache_read: 144168
          cache_write: 60576
          cost: 0.1441
cost_of_delay:
  value: 27.78
  by: planner-S-0510
  at: 2026-10-28T08:06:14Z
forecast:
  duration: 2h30m
  delivery: 2026-11-02T21:18:00Z
  basis: "Its own forecast of 2h30m; 5th in the pull order with an in-progress limit of 5, behind S-0424, S-0518, S-0517, S-0554, S-0425 and S-0509."
  by: flai
  at: 2026-11-02T18:14:02Z
finalized:
  by: orchestrator
  at: 2026-10-28T08:07:13Z
---
# S-0510 A project's manifest names an image layer for its test tiers' tools, flai serve runs it, flai check warns when it is not built, and this repository's layer runs the close-out inside the container

## Goal

ADR-0203's limitation: only the worktree's files are the container's to write, and a test that needs a tool outside the image fails inside it until the project's layer carries the tool. The bedrock and claude-code images carry no Go toolchain, no golangci-lint, and no pnpm, so a story of this repository that changes Go or Svelte cannot run `flai test` or `scripts/close-out.sh` inside the container; S-0505's results name each tier that failed for that reason. This story gives a project one image layer of its own: the manifest names a Dockerfile under `agent.image.layer`, `scripts/agent-image.sh` builds it over the harness image as `flai-agent-<project key>-<inner>:local`, `flai serve` runs that image for the project's agents when it is named, and `flai check` warns, as `agent.image`, when the manifest names a layer that is not built, since `docker image inspect` finds no such image. This repository's layer, `flai/images/system-flow/Dockerfile`, carries Go at the version `flai/go.mod` names, golangci-lint v2, Node, and pnpm, so that `flai test` runs every per-path tier inside the container. The close-out inside the container, and a Go-changing story worked end to end on `harness: container` with the layer, moved to S-0565 on TH-0814.

## Acceptance criteria

- [x] `agent.image.layer` in `system-flow.yaml` names a Dockerfile path relative to the repository root, read by `flai/internal/manifest`, documented in `design/system/project-manifest.md` and `docs/operators/settings.md`, and left as a comment in `template/root/system-flow.yaml.tmpl` with `template/CHANGELOG.md` saying so
- [x] `scripts/agent-image.sh --layer` builds the named Dockerfile with the harness image as its base argument and tags it `flai-agent-<project key>-<inner>:local`, `make agent-images` builds it after the two base images when the manifest names one, and the container adapter runs that image in place of the harness image for the project's agents when it is named
- [x] `flai check` reports `agent.image` when the manifest names a layer whose image is not built or whose Dockerfile does not exist, as a warning that names `make agent-images`, and stays silent when the manifest names none or `docker` is not on `PATH`, saying so as a note
- [x] `flai/images/system-flow/Dockerfile` carries Go at `flai/go.mod`'s version, golangci-lint v2 at the version `scripts/install-tools.sh` installs, Node 24, and pnpm, and `system-flow.yaml` names it, so that `scripts/flai.sh test`, run by hand inside the layered image on this story's worktree, runs every per-path tier the manifest's `tests` select, each tier's outcome recorded in this story's notes
- [x] `docs/users/flai.md`, `docs/operators/index.md`, `docs/operators/runbooks/install.md`, `design/system/flai-cli.md`, and `design/tech/docker.md` describe the layer, how it is built, and what this repository's carries

## Tasks
- T-2164 agent.image.layer in the manifest names a project's Dockerfile, read by flai and documented in project-manifest.md, settings.md, and the template
- T-2165 scripts/agent-image.sh --layer builds a project's layer over the harness image, make agent-images builds it when named, and this repository's layer carries Go, golangci-lint, Node, and pnpm
- T-2171 flai check reports agent.image when the manifest names a layer that is not built or whose Dockerfile is missing, and notes when docker is not on PATH
- T-2172 The container adapter runs the project's layer image in place of the harness image when the manifest names one
- T-2176 flai test runs every per-path tier inside the layered container on this story's worktree, each limitation met is an issue, and the docs describe the layer

## Notes

Waits on S-0508 for the adapter that runs the image. The layer is the project's risk, as ADR-0203 says; the warning on an unsigned layer it foresees waits for the signed image, which is outside E-0029. A layer's tools are a project's: the template ships none.

### Moved to S-0565 (TH-0814)

On 2026-11-02 the operator confirmed, on TH-0814, the orchestrator's recommendation, as D-0466 records for S-0508. The former criterion 5 (one Go-changing story of this repository worked end to end by `flai serve` on `harness: container` with this layer, its close-out passing inside the container, each limitation an issue) and criterion 4's close-out half (`scripts/close-out.sh` running inside the container) moved to S-0565 as its criterion 3, and S-0565's `after` gained S-0510 and S-0517. The reasons: this host's `flai serve` is the installed release, which gains the layer only once S-0510 is accepted and published, and the close-out runs git, which does not work inside the container until S-0517 (I-0215). Criterion 4 here is the hand run of `flai test` inside the layer, below. The title still says the close-out; S-0565 carries that part.

### In the container (T-2176)

On 2026-11-02, `make agent-images` built `flai-agent-bedrock:local`, `flai-agent-claude-code:local`, and the layer `flai-agent-sf-claude-code:local`, 1.56 GB. The layer holds Go 1.26.0, golangci-lint 2.5.0 built with go1.26.0, ShellCheck 0.10.0, Node v24.22.1, and pnpm 10.34.5.

`scripts/flai.sh test` then ran inside the layer by hand, as the operator's uid. The mounts were this story's worktree read-write, the main checkout's `.git` read-only, and the story's registration read-write. The paths given selected every per-path tier. All passed: gofmt, vet, golangci-lint, go-test with `-race` (4m10s), flaiover-lint, vitest, svelte-check, ShellCheck, and the markdown lint. Three things stopped it at first, and this story fixed each:

- The worktree's `bin/flai` was built with cgo on the host, so it did not run on the image's musl. `scripts/flai-build.sh` now builds it with `CGO_ENABLED=0`.
- `scripts/env.sh` kept the npm and pnpm caches in the main checkout's `.flai-cache`, which the container does not mount, and npm failed with EACCES. They now go to the home when that folder is not writable.
- `go build`'s VCS stamping failed with git's exit status 128, though git answers each of its commands by hand there. The layer sets `GOFLAGS=-buildvcs=false`, and I-0255 records it.

Not run inside: `scripts/close-out.sh`, which commits, and only the host commits (S-0517, S-0520); and the `--all` tiers `integration` and `smoke`. Both are S-0565's now, with the live run by `flai serve`.

### Planning

Planned by planner-E-0029 on 2026-10-26 with the epic; revisited by planner-S-0510 on 2026-10-28, which kept the five tasks, their touches, and their `after`, and wrote the cost of delay value from the operator's input on E-0029.

Touches. The first run's `flai touches suggest`, seeded with `flai/internal/manifest/manifest.go` and `flai/internal/check/check.go` (89 commits), named `manifest_test.go`, `settings.md`, `docs/users/flai.md`, `flai-cli.md`, `project-manifest.md`, and `docs/operators/index.md`, each kept (co-change). Run again on the 25 touches (1255 of 2440 commits), it names only paths outside the story: `docs/users/flai-reference.md` (generated from the command help, and no command changes), `flaiover-dashboard.md` and `docs/users/flaiover.md` (the key is not settable from the dashboard), `design/adrs/README.md` (no ADR: the layer is ADR-0203's), `workflow.md`, the conventions, `template/template.yaml` (its version is bumped at release, not by the story), `cmd/serve_actions.go` and `harness.go` (S-0508's: the host's image setting stays as it is, the layer only wins over it), and `hostapi`, `mcpserver`, `cmd/root.go`; none added. `settings.md` is checked by `TestSettingsIndexIsComplete` in `flai/cmd/settings_doc_test.go`, which walks the `Manifest` struct's yaml tags, so T-2164's row is enough and `manifest/settings.go`, the `flai manifest set` catalog, needs no entry. `scope.go` needs no change: the `agent.image` findings fall on `system-flow.yaml`, which the generic scope rule marks inside or outside by the branch's diff. The rest are the tasks' (layout):

- `flai/internal/check/image.go`, `image_test.go`: the new check beside `design.go` and `gitcompat.go`, registered in `check.go`. T-2171.
- `flai/internal/harness/container.go`, `container_test.go`, `flai/internal/serve/agents.go`, `agents_test.go`: S-0508's files, which gain the layer's image. T-2172.
- `scripts/agent-image.sh`, `Makefile`, `scripts/README.md` (S-0505's and S-0507's, which gain `--layer`), `flai/images/system-flow/Dockerfile`, `system-flow.yaml`: the build and this repository's layer. T-2165. The manifest reader ignores keys it does not know, so naming `agent.image.layer` in `system-flow.yaml` before T-2164 lands breaks nothing.
- `template/root/system-flow.yaml.tmpl`, `template/CHANGELOG.md` (design, criterion 1): the template's manifest comment and its changelog. T-2164.
- `docs/operators/runbooks/install.md`, `design/tech/docker.md` (design, criterion 6). T-2176.
- `design/issues/summary.md` and the folder `design/issues` (layout, folder kept): `flai issue new` in T-2176 writes issue files whose numbers and slugs no task can name; the folder is kept for them alone, and `flai task done` widens the touches with the committed paths (ADR-0144).

Overlaps on 2026-10-28: no story is in progress or in review. S-0508 and S-0509 share `container.go`, `container_test.go`, `agents.go`, and `agents_test.go` with this story, and S-0505, S-0507, and S-0508 share `agent-image.sh`, `Makefile`, `scripts/README.md`, `docker.md`, and the docs; `after` names S-0508, which waits on S-0507 and so on S-0505, so each lands before this story starts, and the pull order holds the rest.

Forecast. `flai forecast S-0510` gives 56m (median 107 s per unit over 93 done feature stories on claude-opus-5-5 in the large band; size 31: 6 criteria, 25 touches) and delivery 2026-10-30T01:13Z, 28th in the pull order. Kept at 2h30m, the first run's figure: a layer with Go, golangci-lint, Node, and pnpm takes minutes per build and is built several times while the Dockerfile is worked out, and criterion 5 works a Go-changing story of this repository inside the container with its close-out's Go tiers, a session of its own of 30 to 60 minutes. The delivery written is flai's play-out start, 2026-10-30T01:13Z less its 56m, plus the 2h30m: 2026-10-30T02:47Z.

Cost of delay. `flai cod S-0510` gives 27.78 USD a week: this story's share of E-0029's 150 USD a week (`penalty_per_week`, set by the operator on 2026-10-26T12:24Z after TH-0657), 2h30m of the 13h30m forecast over the epic's six open stories without inputs. Written as flai gives it. It moves with a sibling's forecast; the next planner run writes it again.
