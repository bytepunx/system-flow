---
id: T-1768
type: task
nature: feature
title: "An ADR records the formula file: where it lives, its steps and inputs, its body, its marker, and the instance's nature and formula field"
status: done
parent: S-0421
owner: alex
created: 2026-10-16T05:47:09Z
updated: 2026-11-02T02:20:04Z
transitions:
  - to: ready
    at: 2026-11-02T02:09:46Z
    by: system-flow
  - to: in-progress
    at: 2026-11-02T02:09:47Z
    by: system-flow
  - to: done
    at: 2026-11-02T02:20:04Z
    by: agent-S-0421
stream: S-0421
tags: [cli]
touches: [design/adrs, design/adrs/README.md]
usage:
  source: log
  seconds: 617
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 15
      output: 5690
      cache_read: 855314
      cache_write: 26790
      cost: 0.4447
---
# T-1768 An ADR records the formula file: where it lives, its steps and inputs, its body, its marker, and the instance's nature and formula field

## Work

Write the ADR with `flai adr new`, on TH-0471's answers to questions 2, 3, and 5. It records:

- the folder;
- the front matter: `name`, `description`, `nature`, `inputs`, a `schedule` when it has one, and `steps` with `id`, `title`, `run` (a flai command, a script under `scripts/`, or a role with its prompt), and `after`;
- how an input is named in a step;
- the body as the formula's description, and the marker;
- the instance's nature and its `formula` field;
- who works the instance.

First layer: every other task of the story reads the format.

## Done when

The ADR is accepted and listed in `design/adrs/README.md`. `flai check --strict` reports nothing on it.

## Notes

Touches the shared `design/adrs` folder, since the ADR's file name is known only once it is made.
