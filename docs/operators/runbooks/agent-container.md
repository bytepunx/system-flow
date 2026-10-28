---
title: "Runbook: run a story's agent in a container"
updated: 2026-10-28
status: active
---

# Run a story's agent in a container

Building the agent images and running one story's Claude Code agent in a container by hand. This is the first step of E-0029, flai running agents in Docker containers; `flai serve` does not start agents this way yet. The scripts came from the experiment S-0505, and its results document records what was measured.

## The images

Two images, built on this machine from this repository's tree:

| Image | What it holds |
|-------|---------------|
| `flai-agent-bedrock:local` | Alpine, `git`, `ca-certificates`, `tini`, and `flai` built from the tree. The base of every agent image. Dockerfile: `flai/images/bedrock/Dockerfile`. |
| `flai-agent-claude-code:local` | Claude Code from npm, at the version the Dockerfile pins, with Node and ripgrep, on the bedrock image. Auto-update is off: a new Claude Code is a rebuild. Dockerfile: `flai/images/claude-code/Dockerfile`. |

They are local: nothing pushes them, and no registry publishes them. They are not stamped and not signed, unlike flai's releases ([release-key.md](release-key.md)). Neither holds a login, a key, or a token: the agent's credentials reach the container when it runs. Neither holds Go, golangci-lint, pnpm, or markdownlint, so an agent in the container cannot run this repository's close-out tiers that need them.

## Build

From the repository root, with Docker running:

```bash
make agent-images
```

It runs `scripts/agent-image.sh`, which builds `flai-agent-bedrock:local` and then `flai-agent-claude-code:local` on top of it, and prints `built flai-agent-bedrock:local and flai-agent-claude-code:local`. The flai inside is built from the checked-out tree, labelled with the last `flai/v*` tag and the commit. An unchanged tree rebuilds from Docker's cache. Build again after a change to flai that the agent should have, or to move Claude Code to another version.

## What the run needs on the host

- Docker. S-0505 ran on Docker Desktop with WSL2 integration; a Docker Engine host should work the same, but was not tried.
- The images, built as above.
- The story started, so that its worktree exists: `flai story start S-nnnn`.
- The project's HTTP MCP server running: `flai host` runs one for every project `flai serve` serves, or start it with `flai mcp start` ([MCP over HTTP](../index.md#mcp-over-http)). `flai mcp status` says where it listens.
- `flai mcp token` answering the server's bearer token. The script reads it with the installed `flai` on `PATH`, not this tree's `bin/flai`; set `FLAI_BIN` to use another.
- A Claude Code login on the host, `~/.claude/.credentials.json`: sign in with `claude` on the host first.

## Run

```bash
scripts/agent-container.sh S-nnnn prompt.md          # the prompt from a file
scripts/agent-container.sh S-nnnn < prompt.md        # or on standard input
scripts/agent-container.sh S-nnnn -- git status      # a command in place of claude, same mounts
```

It runs `claude -p` headless in `flai-agent-claude-code:local` with the arguments `flai serve` gives a story's agent, as the operator's uid, and removes the container when it ends (`--rm`). The output is Claude Code's stream JSON on standard output. The agent's name on the MCP server is `agent-S-nnnn`.

Before the run it writes two things to `.flai-cache/agents/S-nnnn/`, which is git-ignored:

- `mcp.json`, mode 600, naming the HTTP MCP server with the token. It holds a secret: do not copy it or commit it.
- `.credentials.json`, a symbolic link to `/run/claude-login/.credentials.json`, the login inside the container. It does not resolve on the host, and nothing from the login is copied.

### What it mounts

| Host path | In the container | Mode | For |
|-----------|------------------|------|-----|
| `.flai-cache/worktrees/S-nnnn` | the same path, the working directory | read-write | the story's worktree |
| `.flai-cache/agents/S-nnnn` | the same path, as `CLAUDE_CONFIG_DIR` | read-write | Claude Code's state: settings, transcript, caches |
| `~/.claude` (`FLAI_CLAUDE_LOGIN`) | `/run/claude-login` | read-only | the login; only `.credentials.json` is read |
| `.git` of the main checkout, only with `FLAI_AGENT_GIT_MOUNT=1` | the same path | read-write | git in the worktree, whose `.git` file points there |

Nothing else is mounted. The Docker socket is never mounted, so the agent cannot start containers or reach the host's Docker. Without `FLAI_AGENT_GIT_MOUNT=1`, git in the worktree fails with `not a git repository`, and the agent cannot commit or close a task.

### Settings

Environment variables the script reads:

| Variable | Default | Effect |
|----------|---------|--------|
| `FLAI_BIN` | the `flai` on `PATH` | the flai that answers `flai mcp status` and `flai mcp token` |
| `FLAI_CLAUDE_LOGIN` | `~/.claude` | the folder holding the login, mounted read-only |
| `FLAI_AGENT_MODEL` | Claude Code's own | `claude --model`, such as `haiku` |
| `FLAI_AGENT_NETWORK` | `bridge` | `bridge` reaches the host as `host.docker.internal`; `host` uses Docker's host network |
| `FLAI_AGENT_GIT_MOUNT` | unset | `1` also mounts the main checkout's `.git` |
| `FLAI_SESSION` | a new UUID | the session ID the agent logs under |

On the bridge network the agent sends the server's own `Host` header, `127.0.0.1:<port>`, whatever address it dials, because the server refuses a request on its loopback listener that names another host. On Docker Desktop the `host` network is the VM's, not the machine's, and the server is not reachable from it.

## Stop

The agent stops when `claude -p` finishes. To stop it sooner, interrupt the script, or stop the container:

```bash
docker ps --filter ancestor=flai-agent-claude-code:local
docker stop <container id>
```

The container is removed when it stops. Whatever the agent left in the worktree stays there: read `git status` in the worktree before you go on with the story.

## Clean up

After the story, remove its cache folder. It holds the agent's transcript and `mcp.json` with the MCP token:

```bash
rm -rf .flai-cache/agents/S-nnnn
```

Remove the images when you no longer run agents in containers:

```bash
docker image rm flai-agent-claude-code:local flai-agent-bedrock:local
```

If a token in an `mcp.json` may have been copied elsewhere, rotate it with `flai mcp token --rotate`; every agent connected over HTTP then needs the new one.

## Limitations

S-0505's results document records what was tried, what worked, and what is still open, among them the network attempts, the login's OAuth refresh, which was not exercised, and git inside the container. Read it before relying on this runbook for long runs: [S-0505's results](../../../design/experiments/S-0505-claude-code-works-one-story-of-this-repository-inside-a-hand-built-container-from-a-bedrock-image-and-a-claude-code-layer-and-the-results-settle-what-adr-0203-leaves-open.md).
