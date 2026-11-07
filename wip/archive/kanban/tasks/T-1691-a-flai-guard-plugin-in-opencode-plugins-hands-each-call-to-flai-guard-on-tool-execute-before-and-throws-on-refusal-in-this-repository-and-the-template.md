---
id: T-1691
type: task
nature: feature
title: A flai-guard plugin in .opencode/plugins hands each call to flai guard on tool.execute.before and throws on refusal, in this repository and the template
status: done
parent: S-0407
owner: alex
created: 2026-10-16T05:15:34Z
updated: 2026-11-07T21:13:06Z
transitions:
  - to: ready
    at: 2026-11-07T21:03:56Z
    by: agent-S-0407
  - to: in-progress
    at: 2026-11-07T21:03:57Z
    by: agent-S-0407
  - to: done
    at: 2026-11-07T21:13:06Z
    by: agent-S-0407
stream: S-0407
tags: [cli, template]
touches: [".opencode/plugins/flai-guard.js", template/root/.opencode/plugins/flai-guard.js, template/CHANGELOG.md]
after: [T-1690]
usage:
  source: log
  seconds: 549
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 7
      output: 2737
      cache_read: 1126795
      cache_write: 14020
      cost: 0.3815
---
# T-1691 A flai-guard plugin in .opencode/plugins hands each call to flai guard on tool.execute.before and throws on refusal, in this repository and the template

## Work

- Write `flai-guard.js`, an OpenCode plugin of a few lines, in `.opencode/plugins/` and `template/root/.opencode/plugins/`, the two copies the same. On `tool.execute.before` it runs `flai guard --harness opencode` with the input T-1690 defined on standard input, and throws an error carrying the guard's stderr when it exits 2. It passes the session and sub-agent events T-1690 records on to the same command.
- When `flai` is not on `PATH`, it lets every call through and logs once that the guard is off.
- Name the plugin in `template/CHANGELOG.md`.
- Waits for T-1690: the plugin sends the input that reader defines.

## Done when

Both copies are the same, load in the known `opencode` version without error, and refuse a call `flai guard` refuses.

## Notes
