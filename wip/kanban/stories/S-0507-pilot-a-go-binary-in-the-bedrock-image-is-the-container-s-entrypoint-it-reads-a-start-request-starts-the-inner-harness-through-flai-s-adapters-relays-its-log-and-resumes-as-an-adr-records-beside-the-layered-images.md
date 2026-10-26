---
id: S-0507
type: story
nature: feature
title: "pilot, a Go binary in the bedrock image, is the container's entrypoint: it reads a start request, starts the inner harness through flai's adapters, relays its log, and resumes, as an ADR records beside the layered images"
status: backlog
parent: E-0029
owner: alex
created: 2026-10-26T04:43:39Z
updated: 2026-10-26T04:57:43Z
transitions: []
tags: [cli, docker]
topics: [harness, docker]
touches: [flai/internal/harness/harness.go, flai/internal/harness/adapters.go, flai/internal/harness/harness_test.go, flai/internal/pilot/request.go, flai/internal/pilot/request_test.go, flai/internal/pilot/run.go, flai/internal/pilot/run_test.go, flai/pilot/main.go, flai/images/bedrock/Dockerfile, flai/images/claude-code/Dockerfile, scripts/agent-image.sh, scripts/agent-container.sh, design/adrs/README.md, design/adrs, design/system/flai-cli.md, design/tech/docker.md, design/tech/go.md, docs/operators/index.md]
after: [S-0505]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
forecast:
  duration: 2h
  delivery: 2026-10-27T18:10:00Z
  basis: "Its own forecast of 2h; 26th in the pull order with an in-progress limit of 5, behind S-0498, S-0504, S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0414, S-0415, S-0416, S-0417, S-0418, S-0419, S-0420, S-0421, S-0422, S-0423, S-0424, S-0425, S-0499, S-0502, S-0503, S-0505 and S-0506."
  by: flai
  at: 2026-10-26T04:57:43Z
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

Planned by planner-E-0029 on 2026-10-26 with the epic.

Touches. `flai touches suggest S-0507`, seeded with `flai/internal/harness/adapters.go` and `harness.go` (72 commits), names `harness_test.go` (kept, co-change) and then the template's conventions, `template/CHANGELOG.md`, `serve/agents.go`, `docs/users/flai.md`, and the guard, which this story does not change: pilot adds no command and changes no convention, and `flai serve`'s start path is S-0508's. The rest are the tasks' (layout):

- `flai/internal/pilot/request.go`, `request_test.go`, `run.go`, `run_test.go`, `flai/pilot/main.go`: the new package and the second `main` of the flai module, placed as `flai/main.go` and `flai/internal/*` are. T-2161, T-2168.
- `flai/images/bedrock/Dockerfile`, `flai/images/claude-code/Dockerfile`, `scripts/agent-image.sh`, `scripts/agent-container.sh`: S-0505's files, which gain pilot. T-2174.
- `design/system/flai-cli.md`, `design/tech/docker.md`, `design/tech/go.md`, `docs/operators/index.md` (design, criterion 6). T-2153, T-2178.
- `design/adrs/README.md` and the folder `design/adrs` (layout, folder kept): `flai adr new` in T-2153 names the ADR's file by its number and title, which no task can name; the folder is kept for that file alone, and `flai task done` widens the touches with the committed path (ADR-0144).

Forecast and cost of delay: see the second Planning paragraph, written once the figures are set.
