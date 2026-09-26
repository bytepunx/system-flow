---
title: "Runbook: back up"
updated: 2026-09-26
status: active
---

# Back up

What there is to keep, and how. A project's work is in git; flai keeps a little on the host beside it; the dashboard keeps nothing.

## What exists, and what is worth keeping

| Where | What | Keep? |
|-------|------|-------|
| The project's git repository | Every committed work item, narrative, decision, and document, and each story's branch `story/<story>` | Yes: it is the project |
| Uncommitted files in the main checkout | flai writes `wip/` there and commits it at acceptance, so work items and narratives of stories in progress are often uncommitted | Yes |
| `.flai-cache/worktrees/<story>` in a project | Each story in progress's worktree, with what its agent has not committed | Yes, what is uncommitted |
| `~/.flai/config.json` | Your settings, host actions, agent and check commands, import folders ([settings](../settings.md#flai-configuration)) | Yes |
| `~/.flai/serve/dashboard.token`, `dashboard.agent-key` | The dashboard's login token and the credential `flai serve` proves itself with; secrets, mode 0600 | Yes, to keep logins working; otherwise they are made again |
| `~/.flai/serve/journal.jsonl` | Every host action asked for, and what became of it | Yes, as a record |
| `~/.flai/serve/agents.json`, `agents/`, `checks/` | Each story's newest agent run, agents' logs, check results and logs | Yes, on the same machine; see [restore](restore.md) |
| `~/.flai/serve/projects.json`, `removed.json`, `dashboards.json` | Which projects `flai serve` serves, which it does not serve from below a folder because you removed them, and where their dashboards are, by absolute path | On the same machine only |
| `~/.flai/serve/state.json`, `requests.json`, `serve.log`, `~/.flai/host/` | What runs now, recent requests, logs, and the host's token, remade at every start | No |
| `~/.flai/cache`, `~/.flai/bin` | Template clones and the binary | No: fetched again |
| A project's `.flai-cache/` apart from worktrees | The MCP server's token, address, state, and log, agents' read markers, edit notices | No: made again; agents over HTTP are given a new token |
| The flaiover container and image | Nothing of yours: no volume, no project file | No |

With `FLAI_CONFIG` set, the configuration and the `serve` and `host` folders are beside the file it names instead of in `~/.flai`.

## flai

1. The host's state. It can be copied while flai runs; for a copy in which the journal and the agents' record agree to the second, `flai host stop` first and `flai host start` after.

   ```bash
   umask 077
   tar -C ~/.flai -czf ~/flai-host-$(date +%F).tar.gz config.json serve
   ```

   The archive holds the dashboard's token and credential. Keep it where only you can read it, as you would an SSH key; if it is lost, rotate the token ([restore](restore.md#if-a-backup-with-secrets-is-lost)).
2. Each project's history, every branch included, without publishing anything:

   ```bash
   cd my-project
   git bundle create ~/my-project-$(date +%F).bundle --all
   ```

   Pushing to your remote also keeps it; a story branch is pushed only when you choose to (`git push origin story/S-0042`).
3. What is not committed, in the main checkout and in each story's worktree:

   ```bash
   git ls-files -z --modified --others --exclude-standard | tar --null -T - -czf ~/my-project-uncommitted-$(date +%F).tar.gz
   git -C .flai-cache/worktrees/S-0042 ls-files -z --modified --others --exclude-standard \
     | tar -C .flai-cache/worktrees/S-0042 --null -T - -czf ~/my-project-S-0042-uncommitted-$(date +%F).tar.gz
   ```

   `git worktree list` names the worktrees. A deleted file is not in these archives; `git status` before and after says whether one matters.
4. Check each archive lists what you expect: `tar -tzf <archive>`, `git bundle verify <bundle>`.

## flaiover

Nothing to back up. The container has no volume and holds no file of any project: everything it shows it asks of `flai serve`. Its settings are in `~/.flai/config.json` and the project's `system-flow.yaml`, and its two secrets are the files in `~/.flai/serve` above.
