---
id: T-2159
type: task
nature: feature
title: scripts/agent-container.sh starts a container by hand with the worktree, the cache folder, and the login mounted, and Claude Code inside it signs in and reaches the project's HTTP MCP server
status: done
parent: S-0505
owner: alex
created: 2026-10-26T04:48:54Z
updated: 2026-10-28T12:21:45Z
transitions:
  - to: ready
    at: 2026-10-28T12:15:25Z
    by: agent-S-0505
  - to: in-progress
    at: 2026-10-28T12:15:26Z
    by: agent-S-0505
  - to: done
    at: 2026-10-28T12:21:45Z
    by: agent-S-0505
stream: S-0505
tags: [cli, docker]
touches: [scripts/agent-container.sh, scripts/README.md, design/experiments/S-0505-claude-code-works-one-story-of-this-repository-inside-a-hand-built-container-from-a-bedrock-image-and-a-claude-code-layer-and-the-results-settle-what-adr-0203-leaves-open.md]
after: [T-2149]
usage:
  source: log
  seconds: 379
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 59
      output: 19190
      cache_read: 3626285
      cache_write: 81566
      cost: 1.6759
---
# T-2159 scripts/agent-container.sh starts a container by hand with the worktree, the cache folder, and the login mounted, and Claude Code inside it signs in and reaches the project's HTTP MCP server

## Work

Write `scripts/agent-container.sh S-nnnn [prompt-file]`, POSIX `sh` under `set -eu`, that runs `docker run --rm -i --user "$(id -u):$(id -g)"` on `flai-agent-claude-code:local` with the story's worktree (`.flai-cache/worktrees/S-nnnn`) mounted read-write at its host path and as the working directory, `.flai-cache/agents/S-nnnn` mounted read-write and set as `CLAUDE_CONFIG_DIR`, the operator's Claude login folder mounted read-only at a path the script takes from `FLAI_CLAUDE_LOGIN` (default `~/.claude`), and no other mount, never the Docker socket. Before the run it writes `.flai-cache/agents/S-nnnn/mcp.json` naming the project's HTTP MCP server (the address `flai host status` shows, the token from `flai mcp token`, and `X-Flai-Agent: agent-S-nnnn`), and starts `claude -p` headless with `--mcp-config` and `--strict-mcp-config` naming that file, `--allowedTools` as `flai serve` passes them, and `FLAI_AGENT`, `FLAI_SESSION`, and `FLAI_STORY` set. Try the network modes in turn, the default bridge with `--add-host host.docker.internal:host-gateway` first, then `--network host`, and keep the first that reaches the listener; record each attempt. Record in the results document's "What was done" what `claude` read from the login mount, what it wrote under `CLAUDE_CONFIG_DIR`, and whether it signed in without a browser. Waits for T-2149, which builds the image.

## Done when

`scripts/agent-container.sh S-nnnn` with a one-line prompt file makes Claude Code inside the container call `mcp__flai__inbox` and answer, the results document records the login's behaviour and the network mode that worked, and `scripts/README.md` names the script.

## Notes

The worktree's `.mcp.json` names the stdio server and is a protected path: do not write it. `--mcp-config` with `--strict-mcp-config` makes Claude Code read only the file in the cache folder.
