---
title: "Runbook: restore"
updated: 2026-09-24
status: active
---

# Restore

From what [backup.md](backup.md) made. On another machine, or where the projects now lie at other paths, follow [migrate.md](migrate.md#to-another-machine) instead: the host's registry names projects by absolute path.

## flai

1. Install flai if it is not there ([install.md](install.md#flai)).
2. Stop what runs, so nothing rewrites the files as you put them back:

   ```bash
   flai host stop
   ```

3. Put the host's state back:

   ```bash
   tar -C ~/.flai -xzf ~/flai-host-2026-09-24.tar.gz
   chmod 600 ~/.flai/serve/dashboard.token ~/.flai/serve/dashboard.agent-key
   ```

4. A project whose clone was lost: clone it again from the bundle, point it at its remote, and put back each story's worktree where flai keeps it and what was not committed:

   ```bash
   git clone ~/my-project-2026-09-24.bundle my-project && cd my-project
   git remote set-url origin <the project's remote>
   git worktree add .flai-cache/worktrees/S-0042 story/S-0042
   tar -xzf ~/my-project-uncommitted-2026-09-24.tar.gz
   tar -C .flai-cache/worktrees/S-0042 -xzf ~/my-project-S-0042-uncommitted-2026-09-24.tar.gz
   ```

   Clone it at the path it had, so that the registry and the host actions, which name it by path, still find it. A clone that survived needs none of this.
5. Start again and check:

   ```bash
   flai dashboard            # in a project: starts the host, flai serve, and the dashboard
   flai serve status         # the projects you expect, each connected
   flai serve actions        # the host actions on where you had them
   flai check --strict       # in each project
   ```

6. Agents that reach flai's MCP server over HTTP need its token again. The project's `.flai-cache/mcp.token` is made anew when the server first needs it; `flai mcp token` prints it, and `flai mcp status` says where the server listens now.

`agents.json` says which ready stories have had an agent since they entered ready, so restoring it keeps `flai serve` from starting those again. Without it, with the `agent` host action on, every ready story may be started once more.

## flaiover

Nothing of the dashboard's needs restoring. With `~/.flai/serve` back, `flai dashboard` starts it with the same login token: the old login link works, and so do browser sessions.

Without the secrets, `flai dashboard` makes new ones. A dashboard container still running from before holds the old ones, so `docker stop flaiover` first. Then log in again with the link `flai dashboard` prints.

### If a backup with secrets is lost

Make the lost token worthless:

```bash
flai dashboard token --rotate   # a new login token; every other session ends
flai mcp token --rotate         # in each project whose .flai-cache went with it
```

The agent credential has no rotate of its own. `flai dashboard` makes one when the file is missing, and the container reads it only when it starts, so replace it with the container stopped, then have `flai serve` read the new one. This follows from how flai starts the container and has not been tried:

```bash
docker stop flaiover
rm ~/.flai/serve/dashboard.agent-key
flai dashboard              # in a project: a new credential, and the container started with it
flai host restart serve
flai serve status           # every project connected again
```
