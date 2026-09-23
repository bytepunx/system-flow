---
id: T-0371
type: task
nature: feature
title: A settings host action, enabled only in a shell, lets a dashboard change the host's settings through flai
status: done
parent: S-0105
owner: alex
created: 2026-09-23T20:36:54Z
updated: 2026-09-23T20:44:05Z
transitions:
  - to: ready
    at: 2026-09-23T20:37:12Z
    by: system-flow
  - to: in-progress
    at: 2026-09-23T20:37:12Z
    by: system-flow
  - to: done
    at: 2026-09-23T20:44:05Z
    by: system-flow
stream: S-0105
tags: []
---

# T-0371 A settings host action, enabled only in a shell, lets a dashboard change the host's settings through flai

## Work
- A new host action, `settings`, is off until `flai serve enable settings` in a shell. No method turns it on or off.
- `settings.get` (read) returns, for the project: every host action and where it is on; the project's default agent; the agent's name, attended minutes, and command; each harness's program and arguments; the checks commands and timeout; the import folders; and the project's MCP server. It also says whether settings may be changed here and host-wide.
- Writes, each a validated `flai` command run like the dashboard's other writes, and gated by `settings`:
  - host actions on and off, and the default agent: for the project;
  - the agent's name, attended minutes, and command, each harness, the checks commands and timeout, the import folders, and a rotated dashboard token: host-wide, so they need `settings` on for every project;
  - rotate the project's MCP token.
- `flai serve` passes its config path on to the flai it runs.

## Done when
- hostapi tests cover each method's command line, the refusals when `settings` is off here or host-wide, and values that are not values.
- CLI tests cover the new flags.

## Notes
