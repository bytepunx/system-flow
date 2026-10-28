---
id: S-0507
type: story
nature: feature
title: "pilot, a Go binary in the bedrock image, is the container's entrypoint: it reads a start request, starts the inner harness through flai's adapters, relays its log, and resumes, as an ADR records beside the layered images"
status: backlog
parent: E-0029
owner: alex
created: 2026-10-26T04:43:39Z
updated: 2026-10-28T21:35:17Z
transitions: []
tags: [cli, docker]
topics: [harness, docker]
touches: [flai/internal/harness/harness.go, flai/internal/harness/adapters.go, flai/internal/harness/harness_test.go, flai/internal/pilot/request.go, flai/internal/pilot/request_test.go, flai/internal/pilot/run.go, flai/internal/pilot/run_test.go, flai/pilot/main.go, flai/images/bedrock/Dockerfile, flai/images/claude-code/Dockerfile, scripts/agent-image.sh, scripts/agent-container.sh, scripts/README.md, design/adrs/README.md, design/adrs, design/system/flai-cli.md, design/tech/docker.md, design/tech/go.md, docs/operators/index.md]
after: [S-0505]
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
      seconds: 228
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 21
          output: 284
          cache_read: 6625262
          cache_write: 19897
          cost: 1.4569
cost_of_delay:
  value: 22.22
  by: planner-S-0507
  at: 2026-10-26T12:35:04Z
forecast:
  duration: 2h
  delivery: 2026-10-29T05:54:00Z
  basis: "Its own forecast of 2h; 19th in the pull order with an in-progress limit of 5, behind S-0416, S-0419, S-0423, S-0418, S-0421, S-0422, S-0420, S-0417, S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0424, S-0425 and S-0506."
  by: flai
  at: 2026-10-28T21:35:17Z
finalized:
  by: orchestrator
  at: 2026-10-26T12:35:35Z
---
# S-0507 pilot, a Go binary in the bedrock image, is the container's entrypoint: it reads a start request, starts the inner harness through flai's adapters, relays its log, and resumes, as an ADR records beside the layered images

## Goal

E-0029 names the common tooling inside the container `pilot`: a second Go binary built from the flai module, shipped in the bedrock image as its entrypoint, so that `flai serve` talks to a container the way it talks to a child process today and the image, not flai, decides which harness runs inside. pilot reads one JSON start request on standard input, the `Request` the adapters take (`flai/internal/harness/harness.go`) with the fields a container adds: the inner harness, the names of the variables that hold the MCP URL and token, and the mounted cache folder. It writes an MCP client configuration naming the host's HTTP server, starts the inner harness with the argv and environment that harness's adapter gives, Claude Code first, with `CLAUDE_CONFIG_DIR` in the cache folder and `--mcp-config` naming the file it wrote, relays the inner harness's standard output line by line unchanged to its own, so that `flai serve`'s readers for that harness, usage, stream, and activity, read a container run as a host run, and exits with the inner harness's code. A second start with the same session resumes. An ADR refining [ADR-0203](../../../design/adrs/0203-flai-may-run-a-story-s-agent-in-a-container-from-a-baseline-image-it-builds.md) records the two refinements E-0029 makes: layered images, a bedrock layer and one layer per harness, and pilot as the entrypoint in place of `flai agent run`, built locally within this epic.

## Acceptance criteria

- [ ] An ADR refining ADR-0203 records the layered images (bedrock: a small Linux base, `git`, `flai`, and `pilot`; one layer per inner harness, `claude-code` first), pilot as the container's entrypoint in place of `flai agent run`, that the images are built locally by `scripts/agent-image.sh` within E-0029 with the signed, published image left to a later story, and that pilot relays the inner harness's log unchanged so flai's readers stay per harness
- [ ] `go build ./pilot` in `flai/` builds the `pilot` binary from `flai/pilot/main.go` and `flai/internal/pilot`, `flai/images/bedrock/Dockerfile` builds it beside `flai` and sets it as the image's `ENTRYPOINT`, the claude-code layer keeps that entrypoint, and `scripts/agent-image.sh` still builds both images
- [ ] pilot reads one JSON start request on standard input, `harness.Request` plus `inner` (the harness to run), `mcp_url_env` and `mcp_token_env` (the variables holding the host's MCP URL and token), and `cache` (the mounted folder), refuses a malformed request or an unknown inner harness with one line on standard error and exit code 2, and never writes the environment or a token to any log
- [ ] For `inner: claude-code`, pilot writes an MCP client configuration under the cache folder naming the HTTP server with the bearer token and `X-Flai-Agent`, sets `CLAUDE_CONFIG_DIR` to a folder under the cache, and starts `claude` with the argv the claude-code adapter's `Start` gives for the request, `--mcp-config` and `--strict-mcp-config` naming the file, `--session-id` on a first start and `--resume` on a later one with the same session, and `--agent` for a strategic role, so that a run of this repository's planner or a story's agent inside the image is driven as `flai serve` drives it on the host
- [ ] pilot copies the inner harness's standard output to its own line by line unchanged, writes its own diagnostics to standard error prefixed `pilot:`, forwards `SIGTERM` to the inner harness, and exits with its exit code; a test with a fake inner harness checks the relay, the exit code, the resume, and the refusals
- [ ] `design/system/flai-cli.md` names the `pilot` binary and the `flai/internal/pilot` package, `design/tech/docker.md` lists the layered images and the entrypoint, `design/tech/go.md` names the second binary, and `docs/operators/index.md` says what pilot is and that it reads no file of the operator's

## Tasks
- T-2153 An ADR refining ADR-0203 records the layered images, pilot as the container's entrypoint, and that E-0029 builds its images locally
- T-2154 harness.Request gains the container fields and the claude-code adapter writes an HTTP MCP client configuration beside the stdio one
- T-2161 flai/internal/pilot reads and checks the start request, refusing a malformed one or an unknown inner harness with one line on stderr and exit code 2
- T-2168 pilot starts the inner harness through its adapter with the MCP file, CLAUDE_CONFIG_DIR, session and resume, relays its stdout unchanged, forwards SIGTERM, and exits with its code
- T-2174 The bedrock image builds pilot beside flai and sets it as the entrypoint, the claude-code layer keeps it, and a container started with a request on stdin runs Claude Code through pilot
- T-2178 flai-cli.md, docker.md, go.md, and the operators' index describe pilot, the flai/internal/pilot package, and the layered images

## Notes

The claude-code adapter's `Start` writes a stdio MCP configuration today (`flai/internal/harness/adapters.go`, `mcpConfig`); pilot needs an HTTP one, so `Request` gains the MCP fields and the adapter writes either shape, tested. The relay keeps the readers per harness, as [ADR-0130](../../../design/adrs/0130-every-harness-meets-flai-through-neutral-contracts-an-adapter-s-capabilities-a.md) has them; neutral usage events of pilot's own are not built here, since no inner harness yet needs them. The guard runs inside the container as the hook `.claude/settings.json` in the worktree names, with `flai` from the image; pilot sets `FLAI_AGENT` and `FLAI_SESSION` from the request so the guard and the MCP server see one name.

### Planning

Planned by planner-E-0029 on 2026-10-26 with the epic. Revisited by planner-S-0507 on 2026-10-26 against the story's criteria and the code: the six tasks, their `after`, and their layers stand as drafted; one touch is added to T-2174 and the story; the forecast is kept at 2h; and the cost of delay value is written now that E-0029 has its input. The plan and what the code shows are on the story's plan thread.

Touches. `flai touches suggest S-0507` on the declared touches (590 of 2427 commits changed them) lists `docs/users/flai.md` (50%), `design/system/flaiover-dashboard.md` (18%), `docs/users/flai-reference.md` (17%), `docs/users/flaiover.md` (16%), `docs/operators/settings.md` (12%), `template/CHANGELOG.md` (12%), and the conventions; none is added. pilot is a second binary, not a flai command, so the generated reference and the user guide do not change (the container section of `docs/users/flai.md` is S-0508's, where the user-facing path exists); it adds no setting and changes no convention or dashboard behaviour. Every other touch is declared and a task's, with these sources:

- `flai/internal/harness/harness.go`, `adapters.go`, `harness_test.go` (declared, co-change: the first two changed together in 72 commits, the test with them): `Request` gains the container fields and `mcpConfig` the HTTP shape. T-2154.
- `flai/internal/pilot/request.go`, `request_test.go`, `run.go`, `run_test.go`, `flai/pilot/main.go` (layout): the new package and the second `main` of the flai module, placed as `flai/main.go` and `flai/internal/*` are; the manifest's Go tiers select `flai/**/*.go`, so they are tested and linted without a manifest change, and `flai/.goreleaser.yaml` builds `flai` alone, so no release change. T-2161, T-2168.
- `flai/images/bedrock/Dockerfile`, `flai/images/claude-code/Dockerfile`, `scripts/agent-image.sh`, `scripts/agent-container.sh` (layout): S-0505's files, which gain pilot; the story's `after` names S-0505 for them. The root `.dockerignore` already admits `flai/**`, so `flai/pilot` is in the build context and it is not touched. T-2174.
- `scripts/README.md` (layout, added): T-2174 changes what `scripts/agent-container.sh` does, from starting `claude` to feeding pilot a start request, and the index of the scripts describes each one (S-0505, criterion 6), so it changes in the same task. T-2174.
- `design/system/flai-cli.md`, `design/tech/docker.md`, `design/tech/go.md`, `docs/operators/index.md` (design, criterion 6). T-2153 (`docker.md`), T-2178.
- `design/adrs/README.md` and the folder `design/adrs` (layout, folder kept): `flai adr new` in T-2153 names the ADR's file by its number and title, which no task can name; the folder is kept for that file alone, and `flai task done` widens the touches with the committed path (ADR-0144).

Forecast. `flai forecast S-0507` gives 43m (median 107 s per unit over 93 done feature stories on claude-opus-5-5 in the large band; size 24: 6 criteria, 18 touches) and delivery 2026-10-28T00:12Z, 25th in the pull order. Kept at 2h, the epic planner's figure, for its reason: the story adds a package and a second binary with process tests (relay, exit code, signals, resume), an ADR, and image rebuilds that each run a Go build and a Claude Code install, and no story in the band added a binary. Delivery 2026-10-28T13:43Z is flai's delivery for 43m played out for 2h as flai's replan of 2026-10-26T04:57Z did (43m to 2h moved it 13h31m, from 04:39Z to 18:10Z); flai's next replan plays it out again from the pull order.

Cost of delay. `flai cod S-0507` gives 22.22 USD a week: the story's share of E-0029's 150 USD a week of penalty, set by the operator on 2026-10-26 after TH-0657, 2h of the 13h30m forecast over the epic's six open stories without inputs. Written as it stands: the share by forecast duration is flai's rule, and nothing argues for inputs of the story's own, since pilot brings no revenue and loses no work per cycle apart from the epic's.
