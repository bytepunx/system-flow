---
id: T-3123
type: task
nature: remediation
title: Record the per-login budget key and the guarded clear in an ADR, linked from flai-cli.md
status: done
parent: S-0720
owner: alex
created: 2026-11-11T23:11:37Z
updated: 2026-11-11T23:12:27Z
transitions:
  - to: ready
    at: 2026-11-11T23:11:41Z
    by: agent-S-0720
  - to: in-progress
    at: 2026-11-11T23:11:41Z
    by: agent-S-0720
  - to: done
    at: 2026-11-11T23:12:27Z
    by: agent-S-0720
stream: S-0720
tags: []
touches: [design/system/flai-cli.md]
usage:
  source: log
  seconds: 46
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 28
      output: 7917
      cache_read: 2168701
      cache_write: 36668
      cost: 0.8146
---
# T-3123 Record the per-login budget key and the guarded clear in an ADR, linked from flai-cli.md

## Work

- Record in an ADR, with `adr_new` and commit, that the host's budget state is kept per provider, else per harness program and login, that a state of the old shape holds its whole harness until it clears, that a change of a harness's program, arguments, or login clears the old identity's state, and that `flai serve budget clear` clears one by hand and `flai guard` refuses it to served agents and sub-agents. It refines ADR-0209 §5, which set the host-wide hold.
- Link the ADR from `design/system/flai-cli.md` § The host's budget state.
- Waits for the six other tasks: it records what they built. The close-out passed before it, so it runs alone.

## Done when

- The ADR is in `design/adrs/` with its index row, and flai-cli.md links it.
- The markdown lint passes on the changed files.

## Notes
