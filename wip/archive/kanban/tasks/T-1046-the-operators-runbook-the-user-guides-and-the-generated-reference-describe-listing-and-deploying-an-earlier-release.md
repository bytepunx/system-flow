---
id: T-1046
type: task
nature: improvement
title: The operators' runbook, the user guides, and the generated reference describe listing and deploying an earlier release
status: done
parent: S-0298
owner: alex
created: 2026-10-06T21:45:47Z
updated: 2026-10-07T15:11:29Z
transitions:
  - to: ready
    at: 2026-10-07T15:06:52Z
    by: agent-S-0298
  - to: in-progress
    at: 2026-10-07T15:06:52Z
    by: agent-S-0298
  - to: done
    at: 2026-10-07T15:11:29Z
    by: agent-S-0298
stream: S-0298
tags: [cli, dashboard]
touches: [docs/users/flai-reference.md, docs/operators/settings.md, docs/users/flai.md, docs/users/flaiover.md, docs/operators/index.md, docs/operators/runbooks/update.md]
after: [T-1041, T-1044, T-1045]
---
# T-1046 The operators' runbook, the user guides, and the generated reference describe listing and deploying an earlier release

## Work

Regenerate `docs/users/flai-reference.md` and the flag index in `docs/operators/settings.md` with `scripts/flai-reference.sh` for the new flags and commands. In `docs/operators/runbooks/update.md`, replace the "To go back" lines of both sections with listing the releases (`flai host versions`, `flai dashboard versions`) and deploying one (`flai host upgrade --version`, `flai dashboard upgrade --tag`, or the Updates page). Say that a dashboard version chosen on the page applies once and the configured tag is used again at the next start, and that pinning is `flai config set dashboard.tag`. In `docs/users/flaiover.md` (Host) describe the Versions control of each area; in `docs/users/flai.md` the `--list` flag and the MCP tool `versions`, which lists and deploys nothing; in `docs/operators/index.md` (the dashboard and host host actions) that each now allows deploying a chosen published release. Waits for T-1041, T-1044, and T-1045, the last of the user-facing changes, and through them for every other task.

## Done when

- The reference is regenerated and the runbook, guides, and operator index describe listing an earlier release over the CLI, the dashboard, and MCP, and deploying one over the CLI and the dashboard.
- `flai check --strict` and the markdown lint pass.

## Notes
