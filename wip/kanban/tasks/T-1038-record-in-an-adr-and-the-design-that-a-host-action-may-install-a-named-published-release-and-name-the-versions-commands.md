---
id: T-1038
type: task
nature: improvement
title: Record in an ADR and the design that a host action may install a named published release, and name the versions commands
status: done
parent: S-0298
owner: alex
created: 2026-10-06T21:44:40Z
updated: 2026-10-07T14:32:39Z
transitions:
  - to: ready
    at: 2026-10-07T14:26:51Z
    by: agent-S-0298
  - to: in-progress
    at: 2026-10-07T14:26:51Z
    by: agent-S-0298
  - to: done
    at: 2026-10-07T14:32:39Z
    by: agent-S-0298
stream: S-0298
tags: [cli, dashboard]
touches: [design/adrs, design/adrs/README.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md]
usage:
  source: log
  seconds: 341
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 89
      output: 400
      cache_read: 3794985
      cache_write: 128543
      cost: 1.7774
---
# T-1038 Record in an ADR and the design that a host action may install a named published release, and name the versions commands

## Work

Today `dashboard.upgrade` and `host.upgrade` take no image, tag, or version from the dashboard: flai installs only what the host's own configuration names (S-0081, the comment above them in `flai/internal/hostapi/writes.go`). A rollback from the Updates page reverses that, so write a new ADR (`flai adr new`, from the main checkout's numbering) deciding:

- the host API's `host.upgrade` takes an optional `version` and `dashboard.upgrade` an optional `tag`, each accepted only when it names a release flai itself lists from the configured releases repository (`flai/vX.Y.Z`, `flaiover/vX.Y.Z`; no drafts, no prereleases); the image name is never taken from the dashboard;
- the names of the listing commands and methods: `flai self-upgrade --list`, `flai dashboard versions`, `flai dashboard upgrade --published`, `flai host versions`, `flai host upgrade --version`, the host API reads `host.versions` and `dashboard.versions`, and the MCP read tool `versions`;
- as the operator answered on TH-0202: MCP lists versions only and deploys nothing, and a dashboard version chosen on the Updates page applies once, with no pinning, so the configured `dashboard.tag` is used again at the container's next start;
- that a flai version below a served project's `flai.minimum` is marked in the list and warned about, not refused.

Add the ADR to `design/adrs/README.md`. In `design/system/flai-cli.md` (Versions: the host's flai and the tree) and `design/system/flaiover-dashboard.md` (the `/host` route), describe listing and deploying a chosen version. This task waits for nothing: the tasks after it build to the names it fixes.

## Done when

- The ADR is accepted in `design/adrs` and listed in its README.
- Both design documents describe listing versions over the CLI, the host API, and MCP, and deploying a chosen one, for flai and the dashboard, over the CLI and the host API.
- `flai check --strict` and the markdown lint pass.

## Notes

The ADR's number is not known until it is written, so the task touches the `design/adrs` folder.
