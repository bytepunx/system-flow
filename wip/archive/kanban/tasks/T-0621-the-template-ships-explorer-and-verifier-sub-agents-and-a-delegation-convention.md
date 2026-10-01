---
id: T-0621
type: task
nature: improvement
title: The template ships explorer and verifier sub-agents and a delegation convention
status: done
parent: S-0175
owner: arobson
created: 2026-10-01T07:53:01Z
updated: 2026-10-01T08:08:28Z
transitions:
  - to: ready
    at: 2026-10-01T07:53:25Z
    by: agent-S-0175
  - to: in-progress
    at: 2026-10-01T08:06:56Z
    by: agent-S-0175
  - to: done
    at: 2026-10-01T08:08:28Z
    by: agent-S-0175
stream: S-0175
tags: []
touches: [template/, design/conventions/, ".claude/agents"]
usage:
  source: log
  seconds: 92
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 36
      output: 13993
      cache_read: 3336995
      cache_write: 42384
      cost: 1.2393
---
# T-0621 The template ships explorer and verifier sub-agents and a delegation convention

## Work

- `template/root/.claude/agents/explorer.md` (read-only) and `verifier.md` (reads, runs tests and lint, no edits), with `tools` that leave out `item_move`, `item_edit`, `item_new`, `inbox`, `wait_for_work`, `wait_for_events`, `thread_open`, `thread_reply`, and `thread_resolve`; each primes with `prime --role`.
- A baseline convention on delegation, in the template and here, with `roles` on every convention the roles read.
- The same agent definitions in this repository's `.claude/agents/`; check `.gitignore` keeps them.
- Template version and changelog; `make smoke` renders it.

## Done when

- A rendered project has both definitions and the convention, `flai check --strict` passes on it and here, and a headless session started in the project lists both agents.

## Notes
