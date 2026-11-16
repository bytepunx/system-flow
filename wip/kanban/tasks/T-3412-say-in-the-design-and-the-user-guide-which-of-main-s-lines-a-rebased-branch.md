---
id: T-3412
type: task
nature: remediation
title: Say in the design and the user guide which of main's lines a rebased branch undoes
status: done
parent: S-0800
owner: alex
created: 2026-11-16T01:45:14Z
updated: 2026-11-16T02:09:58Z
transitions:
  - to: ready
    at: 2026-11-16T01:47:32Z
    by: agent-S-0800
  - to: in-progress
    at: 2026-11-16T01:48:55Z
    by: agent-S-0800
  - to: done
    at: 2026-11-16T02:09:58Z
    by: agent-S-0800
stream: S-0800
tags: [docs, design]
touches: [design/system/workflow.md, design/system/flai-cli.md, docs/users/flai.md]
after: [T-3405]
usage:
  source: log
  seconds: 1263
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 28
      output: 7132
      cache_read: 1991302
      cache_write: 73644
      cost: 1.0172
---
# T-3412 Say in the design and the user guide which of main's lines a rebased branch undoes

## Work

Update the passages that state the undone check's rule, as S-0767 did for ADR-0297, to say what T-3405's ADR decides and link it:

- `design/system/workflow.md`, the paragraph under § Branches and collisions (ADR-0019) that starts "S-0371 catches a rebased branch that undoes main's change".
- `design/system/flai-cli.md`, the `flai stream sync <story-id>` row of § Commands.
- `docs/users/flai.md`, the "Main's change undone." entry under `flai stream sync`.

Change only the sentences about which lines count; the record, the blocker, restore, keep, and send back are unchanged.

It waits for T-3405, whose ADR it describes. It runs beside T-3407, with no file in common.

## Done when

- Each of the three passages states the rule and links the new ADR, and none still says every line of main's text counts in a path that is not markdown.
- `flai test` on the three files passes the markdown lint, and `flai check --strict` is clean.

## Notes
