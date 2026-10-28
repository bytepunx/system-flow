---
id: S-0508
type: story
nature: feature
title: flai serve runs a story's agent in a container through a container harness adapter, with the worktree, the story's cache folder, and the login mounted, the keys and the MCP token by environment, and the image and network the host sets
status: backlog
parent: E-0029
owner: alex
created: 2026-10-26T04:44:13Z
updated: 2026-10-28T14:24:21Z
transitions: []
tags: [cli, docker]
topics: [harness, docker]
touches: [flai/internal/harness/container.go, flai/internal/harness/container_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/internal/config/config.go, flai/internal/config/config_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/internal/protected/protected.go, flai/internal/protected/protected_test.go, flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, flai/internal/serve/start.go, flai/internal/serve/start_test.go, flai/internal/serve/stop.go, flai/internal/serve/stop_test.go, flai/internal/serve/restart.go, docs/operators/settings.md, docs/users/flai-reference.md, docs/users/flai.md, docs/operators/index.md, docs/operators/runbooks/install.md, design/system/flai-cli.md, design/system/agent-adapters.md, design/tech/docker.md, design/issues/summary.md, design/issues]
after: [S-0506, S-0507]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 215
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 13
          output: 110
          cache_read: 4374448
          cache_write: 13240
          cost: 0.9619
cost_of_delay:
  value: 27.78
  by: planner-S-0508
  at: 2026-10-26T12:40:34Z
forecast:
  duration: 2h30m
  delivery: 2026-10-29T02:31:00Z
  basis: "Its own forecast of 2h30m; 23rd in the pull order with an in-progress limit of 5, behind S-0502, S-0415, S-0522, S-0419, S-0423, S-0418, S-0421, S-0416, S-0422, S-0420, S-0417, S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0424, S-0425, S-0506 and S-0507."
  by: flai
  at: 2026-10-28T14:24:21Z
finalized:
  by: alex
  at: 2026-10-28T12:54:28Z
---
# S-0508 flai serve runs a story's agent in a container through a container harness adapter, with the worktree, the story's cache folder, and the login mounted, the keys and the MCP token by environment, and the image and network the host sets

## Goal

The `container` harness adapter ADR-0203 decides, under [ADR-0130](../../../design/adrs/0130-every-harness-meets-flai-through-neutral-contracts-an-adapter-s-capabilities-a.md)'s contracts, so that `flai serve` starts a story's agent as a container rather than a process and nothing above the adapter changes. `harness: container` on a story's agent or the project's default names it, and `config.inner` names the harness pilot runs inside, `claude-code` by default. `Start` builds the `docker run` argument list: the image the host set for the inner harness, `--rm`, `-i` with pilot's start request on standard input, `--user` the operator's uid and gid, `--network` as the host set it, the story's worktree mounted read-write at its host path, `.flai-cache/agents/<S-nnnn>` mounted as the cache folder, the login the host named mounted read-only, and the keys, `FLAI_MCP_URL`, and `FLAI_MCP_TOKEN` passed as `-e NAME` so that the docker client takes each value from `flai serve`'s environment and no value is in argv; never the Docker socket, the home, `~/.flai`, or the main checkout. `Capabilities` are the inner harness's, `Reader` is the inner harness's, since pilot relays its log unchanged, and a stop stops the container by name. The host sets the adapter with `flai serve agent harness container`: the image per inner harness, the network, and the login path. The proof is one story of this repository worked end to end by `flai serve` on `harness: container`.

## Acceptance criteria

- [ ] `harness: container` on a story's agent, the project's default, or a strategic agent selects the adapter through `harness.For`, `config.inner` names the inner harness with `claude-code` the default and an unknown one refused before the start, `Capabilities` are the inner harness's, and `Reader` is the inner harness's reader
- [ ] `flai serve agent harness container` sets on the host, and `flai serve settings` shows, the image per inner harness (`--image claude-code=flai-agent-claude-code:local`), `--network` (the default bridge with Docker's host gateway unless set), and `--login <path>` (the folder with the operator's Claude login, mounted read-only); each key has a row in `docs/operators/settings.md`
- [ ] `Start` gives `docker run --rm -i --name flai-agent-<S-nnnn>-<session> --user <uid>:<gid> --network <n> --add-host host.docker.internal:host-gateway`, the three mounts and no other, `-e` for `FLAI_AGENT`, `FLAI_SESSION`, `FLAI_STORY`, `FLAI_STARTED_BY`, `FLAI_MCP_URL`, `FLAI_MCP_TOKEN`, and each provider key's `key_env`, by name only, the image, and pilot's start request on standard input naming `inner`, the variables, and the cache folder; a test checks the argv holds no key, no token, and no socket mount
- [ ] `flai serve` refuses a start, with the message naming what to do, when `docker` is not on `PATH`, when the image for the inner harness is not present (`make agent-images`), or when the login path is not set for an inner harness that needs one; a stop from the dashboard ([ADR-0058](../../../design/adrs/0058-the-agent-host-action-lets-the-dashboard-stop-a-story-s-agent-and-a-stop.md)) stops the container by its name, and a restart after a thread is answered starts a container on the same cache folder and session so the inner harness resumes
- [ ] The protected list gains the files that change what a containerised agent may do: `flai/images/`, `scripts/agent-image.sh`, and `scripts/agent-container.sh`, so a story changing them is accepted by the operator alone, as [ADR-0106](../../../design/adrs/0106-a-story-whose-branch-changes-a-path-claude-code-protects-is-accepted-by-the.md) has it
- [ ] One story of this repository is worked end to end by `flai serve` with `harness: container` on its agent, from `story_start` to review, its usage measured and its card's status shown as for a host run, and each limitation met is an issue under `design/issues`
- [ ] `docs/users/flai.md` § agents and § harnesses, `docs/operators/index.md`, `docs/operators/runbooks/install.md` (Docker Engine and `make agent-images`), `design/system/flai-cli.md`, `design/system/agent-adapters.md` (the container harness as built), and `design/tech/docker.md` describe the adapter, its settings, its mounts, and what it never mounts

## Tasks
- T-2155 The container adapter: harness.For selects it, config.inner names the inner harness, Capabilities and Reader are the inner harness's, and Start gives the docker run argv with three mounts and -e names only
- T-2156 flai serve agent harness container sets the image per inner harness, the network, and the login path on the host, shown by flai serve settings and indexed in settings.md
- T-2157 The protected list gains flai/images, scripts/agent-image.sh, and scripts/agent-container.sh
- T-2162 flai serve starts the container with pilot's request on stdin and the keys and token by environment, refuses a start without docker, image, or login, stops it by name, and resumes on the cache folder
- T-2169 One story of this repository is worked end to end by flai serve on harness container, and the docs describe the adapter, its settings, its mounts, and what it never mounts

## Notes

Waits on S-0506 for the HTTP server the container reaches and the token by environment, and on S-0507 for pilot and the images. What this story leaves to others: the strategic agents (S-0509, whose mounts differ since they run in the main checkout), a project's own image layer for its test tiers (S-0510), and the signed, published image, which is outside E-0029.

### Planning

Planned by planner-E-0029 on 2026-10-26 with the epic; revisited by planner-S-0508 on 2026-10-26 (TH-0666), which kept the five tasks, their touches, and their `after`, and wrote the cost of delay value.

Touches. The first run's `flai touches suggest S-0508`, seeded with `flai/internal/harness/adapters.go`, `flai/internal/serve/agents.go`, `flai/cmd/serve_actions.go`, and `flai/internal/config/config.go` (108 commits), named `serve_actions_test.go`, `serve/agents_test.go`, `harness_test.go`, `harness.go`, `flai-cli.md`, `flai-reference.md`, `settings.md`, `docs/users/flai.md`, `config_test.go`, `serve/restart.go`, `docs/operators/index.md`, `serve/start.go`, and `start_test.go`, each kept (co-change). Run again on the 27 touches (1284 of 2427 commits), it names only paths outside the story: `flaiover-dashboard.md`, `docs/users/flaiover.md`, `design/adrs/README.md`, `workflow.md`, `template/CHANGELOG.md`, `hostapi/writes.go`, `project-manifest.md`, `cmd/root.go`, `mcpserver/*`, `check/check.go`, and `manifest/manifest.go`; none added, since the story writes no ADR, no manifest key, no new command, no dashboard or hostapi change, and the strategic paths are S-0509's. The rest are the tasks' (layout):

- `flai/internal/harness/container.go`, `container_test.go`: the new adapter, beside `adapters.go`, registered in the `Adapters` map of `harness.go`, which gains `Stdin` on `Start`. T-2155. `adapters.go` itself is S-0507's and is not touched here.
- `flai/internal/protected/protected.go`, `protected_test.go` (design, criterion 5). T-2157.
- `flai/internal/serve/stop.go`, `stop_test.go`: the stop by container name. T-2162. `serve/usage.go` needs no change: `logReader` resolves the reader through `harness.For`, so the container adapter's `Reader()` gives the inner harness's.
- `docs/operators/runbooks/install.md`, `design/system/agent-adapters.md`, `design/tech/docker.md` (design, criterion 7). T-2169.
- `design/issues/summary.md` and the folder `design/issues` (layout, folder kept): `flai issue new` in T-2169 writes issue files whose numbers and slugs no task can name; the folder is kept for them alone, and `flai task done` widens the touches with the committed paths (ADR-0144).

Overlaps on 2026-10-26: no story in progress touches these paths; S-0498, in review, touches `config.go`, `config_test.go`, and `settings.md`, and a story in review holds nothing (ADR-0096); its change lands before this story starts, since this one waits on S-0506 and S-0507.

Forecast. `flai forecast S-0508` gives 1h1m (median 107 s per unit over 93 done feature stories on claude-opus-5-5 in the large band; size 34: 7 criteria, 27 touches) and delivery 2026-10-28T02:56Z, 26th in the pull order. Kept at 2h30m, the first run's figure: the adapter and the serve start path change flai's largest file (`serve/agents.go`) with the stop, the restart, and three refusals, and criterion 6 runs a whole story of this repository under `flai serve` inside a container, a session of its own of 30 to 60 minutes with a stop and a restart tried. The delivery written is flai's play-out start, 2026-10-28T02:56Z less its 1h1m, plus the 2h30m: 2026-10-28T04:25Z.

Cost of delay. `flai cod S-0508` gives 27.78 USD a week: this story's share of E-0029's 150 USD a week (`penalty_per_week`, set by the operator on 2026-10-26T12:24Z after TH-0657), 2h30m of the 13h30m forecast over the epic's six open stories without inputs. Written as flai gives it. It moves with a sibling's forecast; the next planner run writes it again.
