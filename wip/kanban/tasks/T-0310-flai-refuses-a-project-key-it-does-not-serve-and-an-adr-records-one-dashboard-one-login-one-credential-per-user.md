---
id: T-0310
type: task
nature: feature
title: flai refuses a project key it does not serve, and an ADR records one dashboard, one login, one credential per user
status: done
parent: S-0080
owner: alex
created: 2026-09-20T21:09:25Z
updated: 2026-09-20T22:25:59Z
transitions:
  - to: ready
    at: 2026-09-20T22:22:40Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T22:22:40Z
    by: system-flow
  - to: done
    at: 2026-09-20T22:25:59Z
    by: system-flow
stream: S-0080
tags: []
---
# T-0310 flai refuses a project key it does not serve, and an ADR records one dashboard, one login, one credential per user

## Work
hostapi and the channel: a request naming a project key that connection does not serve is refused by flai, not assumed correct by the dashboard. An ADR refines ADR-0018 (the token) and ADR-0029/ADR-0031 (the connection and what the container holds) for one dashboard, one login token, and one agent credential per user, and what happens to a project's old per-project token and key.

## Done when
- A test for the refusal; flai check --strict and the markdown lint pass after the last edit

## Notes
