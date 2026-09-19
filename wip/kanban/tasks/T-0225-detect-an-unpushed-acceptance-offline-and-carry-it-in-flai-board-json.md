---
id: T-0225
type: task
nature: feature
title: Detect an unpushed acceptance offline, and carry it in flai board --json
status: done
parent: S-0063
owner: alex
created: 2026-09-19T08:34:59Z
updated: 2026-09-19T08:37:33Z
transitions:
  - to: ready
    at: 2026-09-19T08:34:59Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T08:35:00Z
    by: system-flow
  - to: done
    at: 2026-09-19T08:37:33Z
    by: system-flow
stream: S-0063
tags: []
---

# T-0225 Detect an unpushed acceptance offline, and carry it in flai board --json

## Work
A package `flai/internal/pending` that asks git, offline and with no credential, how the main checkout's branch stands against its remote-tracking branch: commits ahead, which of them are acceptances (subject `chore: [ID] accept and archive`) and of which items, the tags reachable from the branch and not from the remote-tracking branch, and whether the remote-tracking branch has commits the branch lacks. Nothing when there is no git, no upstream, or nothing ahead. `flai board --json` carries it as `unpushed`, and the plain board prints one line. Tests with real git and a bare remote: nothing, ahead without an acceptance, ahead with one, tags, no remote, diverged.

## Done when
- The detection is tested against real git for each case
- `flai board --json` has `unpushed` only when there is something

## Notes
