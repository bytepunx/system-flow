---
title: "Runbook: migrate"
updated: 2026-09-26
status: active
---

# Migrate

Moving flai and the dashboard to another machine, moving a project on one, bringing a project to a newer template or ID format, and crossing a release that changed how things are kept.

## To another machine

The host's registry (`serve/projects.json`, `removed.json`, `dashboards.json`) and the host actions in the configuration name projects by absolute path and the credential by its file's path, so they are not carried over: on the new machine each project is registered again, which also rewrites them.

### flai

1. On the old machine, stop agents being started, let running ones finish their stories, and back up ([backup.md](backup.md)): the configuration, the two dashboard secrets, the journal if you want the record, and each project's bundle and uncommitted files.

   ```bash
   flai serve disable agent --all-projects
   flai host stop
   umask 077
   tar -C ~/.flai -czf ~/flai-move.tar.gz config.json serve/dashboard.token serve/dashboard.agent-key serve/journal.jsonl
   ```

2. On the new machine, install flai ([install.md](install.md#flai)), then put back the configuration and secrets:

   ```bash
   tar -C ~/.flai -xzf ~/flai-move.tar.gz
   flai serve actions         # the host actions and the paths they are on, from the old machine
   ```

3. Turn every host action off, so that none is on for a path that now means something else, then check what else names the old machine:

   ```bash
   for action in push agent checks dashboard host settings; do flai serve disable "$action" --all-projects; done
   flai serve agent show      # a harness --program or an agent command with an old path
   flai serve checks show     # check commands with an old path
   flai serve import list     # import folders: flai serve import remove the old ones, add the new
   ```

4. Clone each project, or put it back from its bundle ([restore.md](restore.md#flai), step 4), and in each, register it and turn on again the host actions you want for it (`flai serve project add` registers it without starting the dashboard):

   ```bash
   flai dashboard
   flai serve enable push     # and each other action you had on for it
   flai check --strict
   ```

5. Agents that reach flai's MCP server over HTTP are given the new address and token (`flai mcp status`, `flai mcp token`).

Ready stories have had no agent on this machine, so with the `agent` action on, each is started once, in pull order, as the in-progress limit allows.

### flaiover

Nothing moves: the first `flai dashboard` on the new machine starts a container of its own. With the two secrets carried over, the old login link opens it; otherwise it makes new ones and prints a new link. Stop the old machine's container (`flai dashboard stop` in each project, or `docker stop flaiover`) and remove its image ([delete.md](delete.md#flaiover)).

## A project to another folder on the same machine

1. In the project, before moving it: `flai serve project remove .` (which leaves the dashboard running for your other projects; `flai dashboard stop` also stops it if this was the last), and `flai serve disable <action>` for each host action on for it (`flai serve actions` lists them).
2. Move it, then reconnect its story worktrees, which git links by absolute path:

   ```bash
   mv ~/git/my-project ~/work/my-project && cd ~/work/my-project
   git worktree repair .flai-cache/worktrees/*
   git worktree list          # every worktree under the new path
   ```

   With `worktrees.relative_paths` on, worktrees made since then need no repair; it has a cost, in [the flai guide](../../users/flai.md#relative-worktree-links-opt-in).
3. Register it again, and enable its host actions again: `flai serve project add` (or `flai dashboard`), `flai serve enable <action>`. Forgot step 1? `flai serve project list` shows the old folder as not served because it is gone; `flai serve project remove <key>` takes it out, and then `add` in the new folder is not refused for its key.

## A project to a newer template

In the project, on a clean tree:

```bash
flai upgrade --dry-run     # what would be added, merged, replaced, or conflict
flai upgrade               # decide each conflict: keep, replace, or diff
flai check --strict
git diff --stat            # one reviewable change; commit it
```

In a script, `--keep-all` never overwrites a project's edit. A project assembled by hand, with no `system-flow.lock.yaml`, runs `flai upgrade --relock` first. How files are classified is in [the flai guide](../../users/flai.md#upgrade-to-a-newer-template). Updating flai itself does not upgrade a project.

## Work item IDs from three digits to four

A project made before IDs had four digits (`E-001`, `S-012`):

```bash
flai migrate ids --dry-run
flai migrate ids
flai check --strict
```

Every item and narrative is renamed with `git mv` and every reference in the project's folders and root files rewritten. Commit the result as one change.

## Across a release that changed how things are kept

| Story | What changed | What to do |
|-------|--------------|------------|
| S-0076, flaiover 0.22 | MCP moved from the dashboard's `/mcp` to `flai mcp` on the host | [MCP over HTTP](../index.md#mcp-over-http) |
| S-0077 | The dashboard no longer mounts the repository or holds a push key; everything goes through `flai serve` | [Upgrading from a release that mounted the repository](../index.md#upgrading-from-a-release-that-mounted-the-repository) |
| S-0080 | One container, `flaiover`, and one login token serve every project, instead of one per project | [Upgrading from one container per project](../index.md#upgrading-from-one-container-per-project) |

Each is done once, in this order, and only by a machine that ran a release from before the story; the changelog of each release names the stories it carries.
