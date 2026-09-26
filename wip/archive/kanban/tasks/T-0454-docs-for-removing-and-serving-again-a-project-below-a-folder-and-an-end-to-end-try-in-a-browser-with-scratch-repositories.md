---
id: T-0454
type: task
nature: improvement
title: Docs for removing and serving again a project below a folder, and an end-to-end try in a browser with scratch repositories
status: done
parent: S-0123
owner: alex
created: 2026-09-26T07:06:22Z
updated: 2026-09-26T07:23:09Z
transitions:
  - to: ready
    at: 2026-09-26T07:06:29Z
    by: agent-S-0123
  - to: in-progress
    at: 2026-09-26T07:16:41Z
    by: agent-S-0123
  - to: done
    at: 2026-09-26T07:23:09Z
    by: agent-S-0123
stream: S-0123
tags: [dashboard]
touches: [docs, design/system]
---
# T-0454 Docs for removing and serving again a project below a folder, and an end-to-end try in a browser with scratch repositories

## Work

- `docs/users` (flai.md and the dashboard guide), `docs/operators` (settings index for `removed.json`, if state files are listed), and `design/system` (flai-cli.md, flaiover-dashboard.md) say what Remove and Serve do for a project below a folder.
- Run a scratch flai serve and dashboard with its own config, an import folder with two scratch repositories, remove one from the settings page, see the switcher drop it, Serve it, see the switcher gain it without a reload.

## Done when

- The docs say it, `make lint-md` and `flai check --strict` pass, and the end-to-end run is recorded in the narrative.

## Notes
