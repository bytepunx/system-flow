---
title: Claude Code works a story inside a hand-built container
updated: 2026-10-28
status: active
story: S-0505
---

# Claude Code works a story inside a hand-built container

The results of S-0505, an `experiment` story run from 2026-10-28. It starts Claude Code by hand in `flai-agent-claude-code:local`, built on `flai-agent-bedrock:local`, against one story's worktree, to settle what ADR-0203 leaves open about a container harness of flai's own.

## Hypothesis

T-2177 fills this section.

## Success measure

T-2177 fills this section.

## What was done

### The images

T-2149 built `flai-agent-bedrock:local` and `flai-agent-claude-code:local` with [`scripts/agent-image.sh`](../../scripts/agent-image.sh): Alpine 3.24 with `git`, `tini`, and `flai` from the tree, then Node 24 and Claude Code 2.1.312. The image's user is 65532, `HOME` is `/home/agent` (mode 1777), and `DISABLE_AUTOUPDATER=1`.

### The container, started by hand

T-2159 wrote [`scripts/agent-container.sh`](../../scripts/agent-container.sh) `S-nnnn [prompt-file]`. On this host, Docker Desktop 29.7.2 with WSL2 integration and WSL2 mirrored networking, it builds this command for S-0505:

```sh
docker run --rm -i --user 1000:1000 --add-host host.docker.internal:host-gateway \
  -v <root>/.flai-cache/worktrees/S-0505:<root>/.flai-cache/worktrees/S-0505 \
  -w <root>/.flai-cache/worktrees/S-0505 \
  -v <root>/.flai-cache/agents/S-0505:<root>/.flai-cache/agents/S-0505 \
  -v /home/alex/.claude:/run/claude-login:ro \
  -e CLAUDE_CONFIG_DIR=<root>/.flai-cache/agents/S-0505 -e FLAI_AGENT=agent-S-0505 \
  -e FLAI_SESSION=<uuid> -e FLAI_STORY=S-0505 \
  flai-agent-claude-code:local \
  claude -p --output-format stream-json --verbose \
    --mcp-config <root>/.flai-cache/agents/S-0505/mcp.json --strict-mcp-config \
    --permission-prompt-tool mcp__flai__permission_prompt \
    --permission-mode acceptEdits --allowedTools Bash,mcp__flai --model haiku < <prompt-file>
```

`<root>` is the main checkout. The prompt reaches `claude -p` on standard input. The arguments are the ones `flai serve` passes, from `flai/internal/harness/adapters.go`, except that the MCP server is the project's HTTP server rather than `flai mcp` on stdio. There are three mounts and no others, and the Docker socket is not one of them. `FLAI_AGENT_MODEL` sets `--model`, `FLAI_AGENT_NETWORK` sets the network mode, and `agent-container.sh S-nnnn -- cmd` runs `cmd` in place of `claude` with the same mounts, for timing. Before the run, the script writes `.flai-cache/agents/S-nnnn/mcp.json` with mode 600:

```json
{"mcpServers":{"flai":{"type":"http","url":"http://host.docker.internal:4243/mcp","headers":{"Authorization":"Bearer <flai mcp token>","Host":"127.0.0.1:4243","X-Flai-Agent":"agent-S-0505"}}}}
```

The script takes the URL from `flai --json mcp status`, run with the installed `flai` (the one on `PATH` before `env.sh` puts the tree's `bin/` first, or `FLAI_BIN`) in the main checkout. That command reads `.flai-cache/mcp-http.json` from the project's own cache folder (`Repo.CacheDir`), so the `FLAI_CONFIG` that `scripts/env.sh` sets does not change its answer: it names the server the operator's host started, `http://127.0.0.1:4243/mcp`, flai 1.71.1.

### The network

The server listens on `127.0.0.1:4243` in WSL. Each attempt:

1. **Default bridge with `--add-host host.docker.internal:host-gateway`.** `host.docker.internal` resolves to `192.168.65.254`, Docker Desktop's host gateway, and the TCP connection reaches the WSL listener: without a token the server answers 401. With the token, the MCP SDK answers 403 `Forbidden: invalid Host header "host.docker.internal:4243"`. The go-sdk (v1.8.0) `StreamableHTTPHandler` refuses as DNS rebinding any request that arrives on a loopback listener with a Host header that is not loopback, unless the server sets `DisableLocalhostProtection`, and `flai/internal/mcphttp` does not. Claude Code's first run reported the `flai` server `failed`, with no MCP tools, and stopped because `--permission-prompt-tool mcp__flai__permission_prompt` named a tool it did not have.
2. **Default bridge without `--add-host`.** Docker Desktop resolves `host.docker.internal` to the same `192.168.65.254`, with the same results as attempt 1. The `--add-host` changes nothing here, but it is kept for a Docker Engine host.
3. **`--network host`.** `127.0.0.1:4243` refused the connection. Docker Desktop's host network is its VM's, not WSL's, and nothing listens on the VM's port 4243.
4. **Default bridge, `host.docker.internal`, with `"Host":"127.0.0.1:4243"` in `mcp.json`'s headers.** This works. Claude Code 2.1.312 sends the Host header that the configuration names, the SDK sees a loopback Host, and the server connects (`"status":"connected"`). The script makes this its default (`FLAI_AGENT_NETWORK=bridge`). Node's own `fetch` ignores a Host header, so a client other than Claude Code may not get through this way. Another way in, not tried, is a relay inside the container that listens on `127.0.0.1:4243` and forwards to `host.docker.internal`.

Neither attempt changed the host's flai configuration or Docker Desktop. For an agent in a container to reach the server without overriding the Host header, flai's HTTP MCP server needs a way to accept a Host other than loopback, such as an allowed-hosts list, or to listen on an address the container can reach.

### The login

The script mounts the operator's `~/.claude` read-only at `/run/claude-login` and sets `CLAUDE_CONFIG_DIR` to the story's cache folder. Before the run, the script puts a symbolic link in the cache folder, `.credentials.json -> /run/claude-login/.credentials.json`. The link does not resolve on the host, and nothing secret is copied into the cache folder.

- **What `claude` read from the login mount:** `.credentials.json` alone, through the link. `CLAUDE_CONFIG_DIR` points elsewhere, so nothing else on the mount is on a path Claude Code reads. The operator's `settings.json`, hooks, plugins, skills, history, and `~/.claude.json` were not used. The run's init event lists only the built-in plugins and this worktree's own agents. A write to the mount fails with `Read-only file system`.
- **What it wrote under `CLAUDE_CONFIG_DIR`:** `.claude.json` (mode 600), `backups/.claude.json.backup.<ms>`, `.last-cleanup`, `cache/model-catalog/<id>.json`, `plugins/store/cc-plugin-agents-md_builtin-<hash>.json`, `plugins/synced/<id>` and `.bucket-<id>`, `skills/synced/<id>` and `.bucket-<id>`, `projects/<worktree path, dashed>/<session>.jsonl` (the transcript), `projects/<…>/memory/`, and `sessions/`. All of these are owned by the operator's uid. The cache folder is git-ignored.
- **Sign-in:** it signed in without a browser. The init event reports `apiKeySource: "none"`, so the login is the OAuth subscription login and not a key. The model answered on the first run.
- **OAuth refresh:** not exercised. The access token was valid for another seven hours. The read-only `.credentials.json` kept its size and modification time across the runs, so nothing wrote to it. What a refresh does in the container is still open. Claude Code may fail to write through the link to the read-only mount, or it may replace the link with a file of its own in the cache folder. If the refresh token rotates, the operator's copy on the host could stop working. Testing this means forcing an expiry on the operator's own login, which this task did not do.

### Git inside the container

With the three mounts alone, `git status` in the worktree fails with `fatal: not a git repository`. The worktree's `.git` is a file that names `<root>/.git/worktrees/S-0505`, which is not mounted. With `FLAI_AGENT_GIT_MOUNT=1`, the script also mounts the main checkout's `.git` read-write at its host path, and then `git status` and `git log` work (`On branch story/S-0505`). Git reports no ownership complaint, because the container runs as the operator's uid. T-2166 and T-2173 compare the two.

### The test

The prompt was "Call the flai MCP tool inbox and reply with only the number awaiting_you and the ids of the threads listed.", run with `FLAI_AGENT_MODEL=haiku scripts/agent-container.sh S-0505 <prompt-file>`. Claude Code (claude-haiku-5-5) loaded `mcp__flai__inbox` with ToolSearch and called it over HTTP. The server answered as `agent-S-0505` with `awaiting_you: 1` and the thread `TH-0676`. Claude replied `1` and `TH-0676`, which matches the inbox, in 3 turns, 8.4 s, and $0.0047. The server's log records the sessions under the agent name `agent-S-0505`.

## Results

T-2166 and T-2173 fill this section.

## Recommendation

T-2177 fills this section with adopt, adapt, or drop.
