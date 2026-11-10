---
title: "Runbook: run a story's agent in a container"
updated: 2026-11-10
status: active
---

# Run a story's agent in a container

Building the agent images and running one story's Claude Code agent in a container by hand, or a probe that tries an image without working the story. This is the first step of E-0029, flai running agents in Docker containers. Since S-0508 `flai serve` starts a story's agent this way on `harness: container`, as [docs/users/flai.md](../../users/flai.md#in-a-container-the-container-harness) describes; this runbook is the hand-run form. The scripts came from the experiment S-0505, and its results document records what was measured.

## The images

Two images, and a third when the project names a layer, built on this machine from this repository's tree:

| Image | What it holds |
|-------|---------------|
| `flai-agent-bedrock:local` | Alpine, `git`, `ca-certificates`, `tini`, and `flai` built from the tree. The base of every agent image. Dockerfile: `flai/images/bedrock/Dockerfile`. |
| `flai-agent-claude-code:local` | Claude Code from npm, at the version the Dockerfile pins, with Node and ripgrep, on the bedrock image. Auto-update is off: a new Claude Code is a rebuild. Dockerfile: `flai/images/claude-code/Dockerfile`. |
| `flai-agent-<key>-claude-code:local` | The project's layer, on the claude-code image, when `system-flow.yaml` names its Dockerfile under `agent.image.layer`. This repository's is `flai-agent-sf-claude-code:local`, from `flai/images/system-flow/Dockerfile`: Go, golangci-lint, ShellCheck, and pnpm, the tools the close-out's tiers need ([docker.md](../../../design/tech/docker.md)). |

They are local: nothing pushes them, and no registry publishes them. They are not stamped and not signed, unlike flai's releases ([release-key.md](release-key.md)). None holds a login, a key, or a token: the agent's credentials reach the container when it runs. The first two hold no Go, golangci-lint, pnpm, or markdownlint, so an agent in them cannot run this repository's close-out tiers that need them; the layer image holds the tiers' tools.

## Build

From the repository root, with Docker running:

```bash
make agent-images
```

It runs `scripts/agent-image.sh`, which builds `flai-agent-bedrock:local` and then `flai-agent-claude-code:local` on top of it, and prints `built flai-agent-bedrock:local and flai-agent-claude-code:local`. When the manifest names a layer, it then builds the layer image over them and prints `built flai-agent-<key>-claude-code:local`; `scripts/agent-image.sh --layer` builds the layer alone. `scripts/agent-image.sh --tag` prints the image a story's agent runs in and builds nothing. The flai inside is built from the checked-out tree, labelled with the last `flai/v*` tag and the commit. An unchanged tree rebuilds from Docker's cache. Build again after a change to flai that the agent should have, or to move Claude Code to another version.

## What the run needs on the host

- Docker. S-0505 ran on Docker Desktop with WSL2 integration; a native Docker Engine host was not tried. There the MCP server must listen on the bridge's gateway, `flai serve mcp bind <gateway>`, since the engine routes `host.docker.internal` to it; `flai serve` refuses its own container starts until it does ([ADR-0269](../../../design/adrs/0269-flai-serve-refuses-a-container-start-whose-http-mcp-server-listens-on-loopback.md)), but this script makes no such check.
- The images, built as above.
- The story started, so that its worktree exists: `flai story start S-nnnn`.
- The project's HTTP MCP server running: `flai host` runs one for every project `flai serve` serves, or start it with `flai mcp start` ([MCP over HTTP](../index.md#mcp-over-http)). `flai mcp status` says where it listens.
- `flai mcp token` answering the server's bearer token. The script reads it with the installed `flai` on `PATH`, not this tree's `bin/flai`; set `FLAI_BIN` to use another.
- A Claude Code login on the host, `~/.claude/.credentials.json`: sign in with `claude` on the host first.

## Run

```bash
scripts/agent-container.sh S-nnnn                   # the story's agent, working the live story
scripts/agent-container.sh --probe S-nnnn           # a trial of the image, which works no item
scripts/agent-container.sh [--probe] S-nnnn -- cmd  # cmd in place of pilot, same mounts
```

The script takes no prompt. pilot, the image's entrypoint, builds it from the start request the script feeds it, as `flai serve` does ([flai-cli.md](../../../design/system/flai-cli.md#pilot-the-agent-images-entrypoint-s-0507)).

| Form | Agent name | Session | `FLAI_STORY` | Prompt |
|------|------------|---------|--------------|--------|
| `S-nnnn` | `agent-S-nnnn` | `FLAI_AGENT_SESSION`, or a new one | `S-nnnn` | the story's agent's: it works the live story |
| `--probe S-nnnn` | `probe-S-nnnn` | always a new one | none | read and report: it calls flai's read tools, runs `flai version`, `git status`, and `git log -1`, reports one line per check, and ends |

A run without `--probe` is the story's agent. It moves the story's tasks, commits through `task_done`, and writes the narrative, as `flai serve` would. Use `--probe` to try an image: its agent works no item and never calls a writing tool, so the live story is left as it was ([I-0224](../../../design/issues/I-0224-scripts-agent-container-sh-can-run-only-a-story-s-own-agent-under-the-story-s-own-name-so-a-trial-run-of-the-image-acts-on-the-live-story.md), S-0682). The script says on standard error that the run is a probe. The mounts are the story's in both forms, so a probe sees what the story's agent would.

A capped trial, on the smallest model and a budget of two cents:

```bash
FLAI_AGENT_MODEL=haiku FLAI_AGENT_MAX_BUDGET_USD=0.02 scripts/agent-container.sh --probe S-nnnn
```

With `--` the script runs the command in place of pilot, with `FLAI_AGENT` and `FLAI_SESSION` set as in the table, and `FLAI_STORY` only without `--probe`.

Without `--`, pilot runs `claude -p` headless in the agent image ([Settings](#settings)) with the arguments `flai serve` gives a story's agent, as the operator's uid, and the container is removed when it ends (`--rm`). The output is Claude Code's stream JSON on standard output.

The script and pilot write these to `.flai-cache/agents/S-nnnn/`, which is git-ignored:

- `mcp.json`, mode 600, which pilot writes. It names the HTTP MCP server, with its URL, token, and Host as variables that Claude Code expands as it reads the file, so the token is not written. Its `timeout` of 1800000 milliseconds keeps a slow `task_done`, `story_start`, or `test` from being dropped at 60 seconds ([ADR-0258](../../../design/adrs/0258-the-mcp-file-flai-writes-for-its-http-server-names-a-per-server-timeout-of.md)). An image built before that change writes it without the `timeout` and keeps the 60-second limit until `make agent-images` builds it again.
- `claude/.credentials.json`, a symbolic link to `/run/claude-login/.credentials.json`, the login inside the container. It does not resolve on the host, and nothing from the login is copied.
- `flai/`, an empty folder for the container's flai. On its first run that flai writes its config there, `flai/config.json`, and keeps its cache in `flai/cache`.

### What it mounts

| Host path | In the container | Mode | For |
|-----------|------------------|------|-----|
| `.flai-cache/worktrees/S-nnnn` | the same path, the working directory | read-write | the story's worktree |
| `.git` of the main checkout | the same path | read-only | git in the worktree reads the objects, refs, and `config` from it |
| `.git/worktrees/S-nnnn` of the main checkout, the story's registration | the same path, over the read-only `.git` | read-write | the worktree's own `index`, `HEAD`, logs, and locks, which its `.git` file names |
| `wip/` of the main checkout, as the manifest's `layout.wip` names it | the same path | read-only | the narrative, the board, the items, the threads, and the messages that `flai verify` checks ([ADR-0253](../../../design/adrs/0253-a-container-that-runs-a-story-s-agent-mounts-the-main-checkout-s-wip-and-system.md)) |
| `system-flow.yaml` of the main checkout | the same path | read-only | by it flai knows the main checkout is there, so `flai verify` runs rather than refusing |
| `.flai-cache/verify` of the main checkout, made first so it is yours | the same path | read-write | the report `flai verify` writes, which `flai verify --last` reads on the host |
| `.flai-cache/agents/S-nnnn` | the same path | read-write | in `claude/`, the `CLAUDE_CONFIG_DIR`, Claude Code's state: settings, transcript, caches; in `flai/`, the container's flai config and cache |
| `~/.claude` (`FLAI_CLAUDE_LOGIN`) | `/run/claude-login` | read-only | the login; only `.credentials.json` is read |

Nothing else is mounted. The Docker socket is never mounted, so the agent cannot start containers or reach the host's Docker. The two git mounts are [ADR-0246](../../../design/adrs/0246-a-container-that-runs-a-story-s-agent-mounts-the-main-checkout-s-git-read-only.md)'s: git in the worktree reads its history and refreshes its index, so `git status`, `diff`, `log`, and `show` work, and a `git commit`, `stash`, or new branch is refused by git on the read-only objects and refs. Commits are the host's, through `task_done` and the close-out. The script refuses to start when the story's registration, `.git/worktrees/S-nnnn`, is missing, and names its path.

### The close-out's and the checks' containers

Under `flai serve`, a story whose agent runs on `harness: container` is closed out, and checked, in containers of their own from the same image, so nothing from its worktree runs on the host ([ADR-0274](../../../design/adrs/0274-the-close-out-of-a-story-whose-agent-runs-in-a-container-runs-its-verify-in.md)). This script does not start them; `flai serve` does. Their mounts beside the agent's:

| Mount | The agent's container | The close-out's and the checks' |
|-------|-----------------------|---------------------------------|
| the story's worktree, the working directory | read-write | read-write |
| `.git` of the main checkout | read-only | read-only |
| the story's registration, `.git/worktrees/S-nnnn` | read-write | read-write |
| `wip/` and `system-flow.yaml` of the main checkout | read-only | read-only |
| `.flai-cache/verify` of the main checkout | read-write | read-write |
| `.flai-cache/agents/S-nnnn` | read-write | not mounted |
| the login | read-only | not mounted |

The close-out's container is `flai-closeout-S-nnnn-<stamp>`, with `-r2` and `-r3` when main moved and it ran again, and runs `flai verify S-nnnn --record-issues`. Each check's is `flai-check-S-nnnn-<stamp>-<n>`. Neither gets the MCP server's URL or token, a provider's key, or `host.docker.internal`. The close-out's commit, its sync check, and its sync run on the host through git and `flai serve`'s own flai, with the main checkout's hooks only. `docker ps --filter name=flai-closeout-` and `--filter name=flai-check-` list them while they run.

The script sets two variables in the container so that its flai writes nothing under the container's home, `/home/agent`:

| Variable | Points at | For |
|----------|-----------|-----|
| `FLAI_CONFIG` | `.flai-cache/agents/S-nnnn/flai/config.json` | the config file the container's flai reads, and writes with defaults on its first run |
| `FLAI_CACHE_DIR` | `.flai-cache/agents/S-nnnn/flai/cache` | the `cache_dir` that first run writes into that config |

Both are host paths, the same inside the container, in the mounted agent folder. They do not touch the host's flai config: the script's own `flai mcp status` and `flai mcp token` run on the host, with the `FLAI_CONFIG` that `scripts/env.sh` sets there.

### Settings

Environment variables the script reads:

| Variable | Default | Effect |
|----------|---------|--------|
| `FLAI_BIN` | the `flai` on `PATH` | the flai that answers `flai mcp status` and `flai mcp token` |
| `FLAI_AGENT_IMAGE` | what `scripts/agent-image.sh --tag` prints: the layer image when the manifest names `agent.image.layer`, else `flai-agent-claude-code:local` | the image the run uses. The default is the one `flai serve`'s container harness runs (S-0670). The script refuses an image not built, naming `scripts/agent-image.sh`; `scripts/agent-container.sh --image` prints its choice and runs nothing |
| `FLAI_CLAUDE_LOGIN` | `~/.claude` | the folder holding the login, mounted read-only |
| `FLAI_AGENT_MODEL` | Claude Code's own | `claude --model`, such as `haiku` |
| `FLAI_AGENT_NETWORK` | `bridge` | `bridge` reaches the host as `host.docker.internal`; `host` uses Docker's host network |
| `FLAI_AGENT_MAX_BUDGET_USD` | none | the agent's `max_budget_usd`, such as `0.02`, for a trial run |
| `FLAI_AGENT_SESSION` | a new UUID | the session the agent logs under; a second run with the same one resumes it. Ignored with `--probe`, which always starts a new one |

On the bridge network the agent sends the server's own `Host` header, `127.0.0.1:<port>`, whatever address it dials, because the server refuses a request on its loopback listener that names another host. On Docker Desktop the `host` network is the VM's, not the machine's, and the server is not reachable from it.

## Stop

The agent stops when `claude -p` finishes. To stop it sooner, interrupt the script, or stop the container:

```bash
docker ps --filter ancestor=flai-agent-claude-code:local
docker stop <container id>
```

The container is removed when it stops. Whatever the agent left in the worktree stays there: read `git status` in the worktree before you go on with the story.

## Clean up

After the story, remove its cache folder. It holds the agent's transcripts, a probe's among them, and `mcp.json`:

```bash
rm -rf .flai-cache/agents/S-nnnn
```

Remove the images when you no longer run agents in containers:

```bash
docker image rm flai-agent-sf-claude-code:local flai-agent-claude-code:local flai-agent-bedrock:local
```

If the MCP token may have been copied elsewhere, rotate it with `flai mcp token --rotate`; every agent connected over HTTP then needs the new one.

## Limitations

S-0505's results document records what was tried, what worked, and what is still open, among them the network attempts, the login's OAuth refresh, which was not exercised, and git inside the container, which ADR-0246's two git mounts have since settled. The close-out runs inside the container past `flai verify`'s own steps, since the main checkout's `wip/`, `system-flow.yaml`, and `.flai-cache/verify` are mounted (ADR-0253, I-0218). Its tiers need the project's layer image, `flai-agent-sf-claude-code:local` here (S-0510), which this script runs since S-0670: in `flai-agent-claude-code:local` the close-out stops at the first tier that needs `go`. Its commit step stops inside the container whenever anything is uncommitted, such as an issue `--record-issues` recorded, since commits are the host's: commit through `task_done` and ask for the close-out from the host, whose job verifies in a container and commits on the host ([above](#the-close-outs-and-the-checks-containers)). Read S-0505's results before relying on this runbook for long runs: [S-0505's results](../../../design/experiments/S-0505-claude-code-works-one-story-of-this-repository-inside-a-hand-built-container-from-a-bedrock-image-and-a-claude-code-layer-and-the-results-settle-what-adr-0203-leaves-open.md).
