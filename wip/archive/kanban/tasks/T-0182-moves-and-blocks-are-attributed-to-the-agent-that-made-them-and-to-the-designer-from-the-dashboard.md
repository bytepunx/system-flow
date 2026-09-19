---
id: T-0182
type: task
nature: remediation
title: Moves and blocks are attributed to the agent that made them, and to the designer from the dashboard
status: done
parent: S-0058
owner: alex
created: 2026-09-19T02:46:21Z
updated: 2026-09-19T02:55:01Z
transitions:
  - to: ready
    at: 2026-09-19T02:52:49Z
    by: alex
  - to: in-progress
    at: 2026-09-19T02:52:50Z
    by: alex
  - to: done
    at: 2026-09-19T02:55:01Z
    by: system-flow
stream: S-0058
tags: []
touches: [flai/cmd, flaiover/src/routes/api]
---

# T-0182 Moves and blocks are attributed to the agent that made them, and to the designer from the dashboard

## Work
`flai move`, `flai block`, and `flai unblock` default `--by` the way `flai thread` does: `--by`, else `FLAI_AGENT`, else the config author; today they ignore `FLAI_AGENT`, so an agent's moves from the command line are recorded as the operator's and would be reported back to the agent as the designer's. `flai accept` keeps its default: acceptance is the operator's. The dashboard's move, block, and unblock endpoints pass `--by` with the designer (`designer()` in `src/lib/server/threads.ts`, the manifest's owner) when the request names nobody, so the board's moves stay attributed to the person and not to `flaiover`. Tests on both sides.

## Done when
- A move with `FLAI_AGENT` set and no `--by` records the agent; without it, the config author
- The dashboard's move records the manifest owner; flaiover tests pass

## Notes
Only moves were changed. `flai block` and `flai unblock` take no `--by` and a blocked interval has no author field (`from`, `until`, `reason`), so there is nothing to default. Adding one would make older flai binaries, which parse front matter strictly, refuse any item that carries it, including the installed flai that runs the MCP server in the operator's sessions. A block is therefore reported to every agent, including the one that set it. The dashboard's block and unblock endpoints were left alone for the same reason.
