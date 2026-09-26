---
id: S-0118
type: story
nature: improvement
title: flai serve project add, remove, and list manage the projects the host serves and say why one is not showing
status: backlog
owner: alex
created: 2026-09-26T05:18:06Z
updated: 2026-09-26T05:18:06Z
transitions: []
tags: [cli]
touches: [flai/cmd, flai/internal/serve]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0118 flai serve project add, remove, and list manage the projects the host serves and say why one is not showing

## Goal

The operator can add a project to the host flai, remove it, and see why it is or is not showing on the dashboard with one command each, without knowing that `flai dashboard` is what registers a project.

## Acceptance criteria
- [ ] `flai serve project add [dir]` registers the project in `dir` (default: the current folder) with the host's dashboard, and refuses a folder with no `system-flow.yaml`, a manifest with no key, or a key another project already has, naming the conflict
- [ ] `flai serve project remove <key|dir>` unregisters it and touches none of its files, and the dashboard's switcher drops it without a restart
- [ ] `flai serve project list` (and `--json`) shows every served project with its key, root, dashboard address, and connected state or last error, plus the import candidates and the unserved projects under the import roots
- [ ] A registered project whose root no longer exists or no longer has a manifest is reported by `list` and `flai serve status`, and is not retried forever in the log
- [ ] `flai dashboard` and `flai dashboard stop` keep working as before and use the same code to register and unregister
- [ ] Documented in `docs/users/flai.md`, `design/system/flai-cli.md`, and the operators' runbooks

## Tasks

## Notes

Today the registry (`~/.flai/serve/projects.json`) can only be changed by `flai dashboard` or the board's import, and editing `~/.flai/config.json` has no effect on it. `flai serve import add|remove|list` manages the import roots, not the served projects, which is easy to confuse. Related to [I-0045](../../../design/issues/I-0045-a-project-imported-with-flai-import-on-the-command-line-is-neither-served-nor-offered-so-the-dashboard-never-shows-it.md).
