---
id: D-0593
title: "No flai command writes a sub-project's install block, so a story that changes projects[].install edits system-flow.yaml by hand and validates it with flai check"
kind: decision
sources: [S-0706, I-0259]
topics: [cli, flai, flaiover]
paths: [system-flow.yaml, design/system/project-manifest.md, flai/cmd/subproject.go]
date: 2026-11-11
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0593 No flai command writes a sub-project's install block, so a story that changes projects[].install edits system-flow.yaml by hand and validates it with flai check

S-0706 fixed I-0259 by having `flaiover/pnpm-workspace.yaml` set `supportedArchitectures`, so a change to that file or to `.npmrc` must also make `flai accept` reinstall the main checkout's `flaiover/node_modules`. Its narrative (`wip/archive/agents/S-0706.md`, § Decisions) records how the manifest was changed: "`pnpm-workspace.yaml` and `.npmrc` were added to flaiover's accept `install` paths, and `pnpm-workspace.yaml` to the svelte-check tier's, both edited by hand in `system-flow.yaml`, since `flai subproject` does not set `install`."

`flai subproject add`, `set`, and `remove` write only an entry's `name`, `path`, `kind`, and `tags`. `flai subproject set --help` says every other key, "its components and install among them, is kept". `flai components config` writes `components`, but no command writes `install`. So a story that changes `projects[].install` edits `system-flow.yaml` by hand, keeps the other lines as they are, and runs `flai check`, which reports a bad block as `manifest.install`. That is an exception to the rule in `CLAUDE.md` that only `name`, `description`, and `dashboard` are hand-edited. D-0449 records the same exception for the `tests` tiers.

After S-0706, flaiover's `install.paths` are `package.json`, `pnpm-lock.yaml`, `pnpm-workspace.yaml`, and `.npmrc`. `design/system/project-manifest.md` § `projects[].install` still names only the first two.

Sources: S-0706's narrative, § Decisions; I-0259.
