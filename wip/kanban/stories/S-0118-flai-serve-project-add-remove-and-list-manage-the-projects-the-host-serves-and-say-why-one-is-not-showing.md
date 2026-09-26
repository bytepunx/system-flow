---
id: S-0118
type: story
nature: improvement
title: flai serve project add, remove, and list manage the projects the host serves and say why one is not showing
status: in-progress
owner: alex
created: 2026-09-26T05:18:06Z
updated: 2026-09-26T05:35:26Z
transitions:
  - to: ready
    at: 2026-09-26T05:24:21Z
    by: alex
  - to: in-progress
    at: 2026-09-26T05:24:40Z
    by: agent-S-0118
tags: [cli]
touches: [flai/cmd/serve_project.go, flai/cmd/serve.go, flai/cmd/dashboard_agent.go, flai/internal/serve/serve.go, flai/internal/serve/projects.go, flai/internal/channel/channeltest, flaiover/src/lib/server/agent.ts, docs/users/flai.md, docs/users/flai-reference.md, docs/operators, design/system/flai-cli.md, design/system/dashboard-host-channel.md, design/system/flaiover-dashboard.md]
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
- [x] `flai serve project add [dir]` registers the project in `dir` (default: the current folder) with the host's dashboard, and refuses a folder with no `system-flow.yaml`, a manifest with no key, or a key another project already has, naming the conflict
- [x] `flai serve project remove <key|dir>` unregisters it and touches none of its files, and the dashboard's switcher drops it without a restart
- [x] `flai serve project list` (and `--json`) shows every served project with its key, root, dashboard address, and connected state or last error, plus the import candidates and the unserved projects under the import roots
- [x] A registered project whose root no longer exists or no longer has a manifest is reported by `list` and `flai serve status`, and is not retried forever in the log
- [x] `flai dashboard` and `flai dashboard stop` keep working as before and use the same code to register and unregister
- [x] Documented in `docs/users/flai.md`, `design/system/flai-cli.md`, and the operators' runbooks

## Tasks
- T-0433 The registry refuses a key another project has, and flai serve reports a registered root that is gone or has no manifest instead of retrying it
- T-0434 flai serve tells the dashboard a project was removed, and the switcher drops it without a restart
- T-0435 flai serve project add, remove, and list, registering through the same code as flai dashboard
- T-0436 Document flai serve project in the user guide, the CLI design, the channel design, and the operators' runbooks

## Notes

Today the registry (`~/.flai/serve/projects.json`) can only be changed by `flai dashboard` or the board's import, and editing `~/.flai/config.json` has no effect on it. `flai serve import add|remove|list` manages the import roots, not the served projects, which is easy to confuse. Related to [I-0045](../../../design/issues/I-0045-a-project-imported-with-flai-import-on-the-command-line-is-neither-served-nor-offered-so-the-dashboard-never-shows-it.md).
