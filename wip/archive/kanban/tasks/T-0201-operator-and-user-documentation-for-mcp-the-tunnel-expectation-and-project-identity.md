---
id: T-0201
type: task
nature: feature
title: Operator and user documentation for /mcp, the tunnel expectation, and project identity
status: done
parent: S-0043
owner: alex
created: 2026-09-19T05:17:08Z
updated: 2026-09-19T05:24:47Z
transitions:
  - to: ready
    at: 2026-09-19T05:24:02Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T05:24:03Z
    by: system-flow
  - to: done
    at: 2026-09-19T05:24:47Z
    by: system-flow
stream: S-0043
tags: []
touches: [docs]
---

# T-0201 Operator and user documentation for /mcp, the tunnel expectation, and project identity

## Work
`docs/operators/index.md`: `/mcp` (what it is, bearer only, one process per session, the idle timeout and the cap, that `wait_for_events` holds a request open for up to five minutes so a proxy's idle timeout must allow it), the tunnel expectation (TLS terminated by the tunnel or proxy; never expose the dashboard's plain HTTP port to the internet), and the project identity headers. `docs/users/flai.md`: how a remote agent's `.mcp.json` points at `/mcp` with the token, beside the stdio form. `docs/users/flaiover.md` where it lists what the dashboard offers.

## Done when
- Each document states the behaviour as built
- Markdown lint and `flai check --strict` pass

## Notes
