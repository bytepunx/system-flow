---
id: T-1105
type: task
nature: improvement
title: The design and the user guide describe --commit, the new MCP tools, and the host channel operations
status: in-progress
parent: S-0275
owner: alex
created: 2026-10-06T22:53:23Z
updated: 2026-10-07T08:45:12Z
transitions:
  - to: ready
    at: 2026-10-07T08:45:11Z
    by: agent-S-0275
  - to: in-progress
    at: 2026-10-07T08:45:12Z
    by: agent-S-0275
stream: S-0275
tags: [docs]
touches: [design/system/flai-cli.md, design/system/continuous-improvement.md, design/system/dashboard-host-channel.md, docs/users/flai.md, docs/users/flai-reference.md]
after: [T-1077, T-1083, T-1093]
---
# T-1105 The design and the user guide describe --commit, the new MCP tools, and the host channel operations

## Work

Criterion 3, for the design and the user guide. It waits for T-1077, T-1083, and T-1093, so that it describes the CLI, the MCP tools, and the host channel operations as they were built.

- `design/system/flai-cli.md`: the `flai issue` and `flai adr` sections give `--commit`, what it commits, where, with which prefix, and the touches it widens.
- `design/system/continuous-improvement.md`: recording an occurrence from a story is one call.
- `design/system/dashboard-host-channel.md`: add `issue.new`, `issue.bump`, and `issue.close` beside the existing operations.
- `docs/users/flai-reference.md`: `--commit` under `flai issue bump`, `close`, and `new` and `flai adr new`, and `--autocommit` and `--trailer` where T-1093 added them. If the reference is generated from the commands' help, regenerate it instead of editing by hand.
- `docs/users/flai.md`: the guide's walk-through of recording an issue or an ADR from a story uses the one call, and names the MCP tools.

## Done when

- Each of the five documents names the one call, and none tells a story's agent to commit or claim an issue or ADR file by hand.
- `flai check --strict` and the markdown lint are clean on the changed files.

## Notes

Written by the planner.
