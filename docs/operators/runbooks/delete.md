---
title: "Runbook: delete"
updated: 2026-09-24
status: active
---

# Delete

Taking flai and the dashboard off a machine. Your projects are git repositories and stay as they are: this removes the tools, their state, and their secrets, not the standard's files in a project. Deleting `~/.flai` cannot be undone, so take a [backup](backup.md) first if you may come back.

## flaiover

Removing the dashboard alone leaves flai working in a shell.

1. See which projects it serves: `flai dashboard status`.
2. In each of them, or in the folder you started it from, take it out:

   ```bash
   flai dashboard stop
   ```

   Each unregisters its project from `flai serve`. The shared `flaiover` container stops when the last one is gone, and Docker removes it, since it was started with `--rm`. To stop it at once whatever is registered: `docker stop flaiover`.
3. Remove the image, and a local build if you made one:

   ```bash
   docker image ls ghcr.io/bytepunx/flaiover
   docker image rm ghcr.io/bytepunx/flaiover:latest   # and each tag listed
   docker image rm flaiover:local
   ```

4. If flai logged Docker in to pull it, and nothing else of yours uses GHCR: `docker logout ghcr.io`.
5. The login token and the agent credential are files in `~/.flai/serve/` (`dashboard.token`, `dashboard.agent-key`). With no container they open nothing; they go with `~/.flai` below. A project's `.flai-cache/dashboard.token` and `.flai-cache/dashboard.agent-key` are left from releases before one container served every project; delete them.

Check: `docker ps --filter name=flaiover` lists nothing.

## flai

1. Stop agents from being started, then see whether any still run:

   ```bash
   flai serve disable agent --all-projects
   flai serve journal       # each agent start, with its PID
   ```

   An agent `flai serve` started outlives it. Let it finish its story, or stop it by its PID. Its story stays where it left it.
2. Stop everything flai runs in the background:

   ```bash
   flai host stop           # flai serve and every MCP server it runs
   flai mcp stop            # in each project where you started one by hand
   ```

   Check: `flai host status` says it is not running.
3. Save the work only a project's `.flai-cache` holds. Each story in progress has a worktree at `.flai-cache/worktrees/<story>` on the branch `story/<story>`:

   ```bash
   git worktree list
   git -C .flai-cache/worktrees/S-0042 status    # for each: commit or discard what is not committed
   git push origin story/S-0042                  # if the branch should outlive this clone
   git worktree remove .flai-cache/worktrees/S-0042
   ```

4. Delete each project's `.flai-cache/`. It is git-ignored and holds only what flai made: the MCP server's token, address, state, and log, agents' read markers, edit notices, and the worktrees you just removed.
5. Delete flai's own folder and the binary:

   ```bash
   rm -rf ~/.flai          # the configuration, cache, serve and host state, secrets, logs, and ~/.flai/bin/flai
   ```

   With `FLAI_CONFIG` set, the configuration and the `serve` and `host` folders are beside the file it names instead; with `FLAI_INSTALL_DIR` or another install, the binary is there (`command -v flai`).
6. Remove the `PATH` line `install.sh` had you add to your shell profile, and `FLAI_*` variables you set there.

A project keeps its `system-flow.yaml`, `design/`, `docs/`, `wip/`, and `.mcp.json`. An agent in it that starts `flai mcp` from `.mcp.json` now fails to start that server, and says so.
