---
id: T-1817
type: task
nature: feature
title: The seed is run once on this repository, its entries committed and counted, and the docs describe harvest and seed
status: done
parent: S-0429
owner: alex
created: 2026-10-16T06:04:22Z
updated: 2026-10-23T18:22:17Z
transitions:
  - to: ready
    at: 2026-10-23T17:37:52Z
    by: agent-S-0429
  - to: in-progress
    at: 2026-10-23T18:20:20Z
    by: agent-S-0429
  - to: done
    at: 2026-10-23T18:22:17Z
    by: agent-S-0429
stream: S-0429
tags: [cli]
touches: [design/digest, design/digest/README.md, design/system/flai-cli.md, docs/users/flai.md]
after: [T-1816]
usage:
  source: log
  seconds: 117
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 35
      output: 16112
      cache_read: 2641120
      cache_write: 85532
      cost: 1.3646
---
# T-1817 The seed is run once on this repository, its entries committed and counted, and the docs describe harvest and seed

## Work

- In the story's worktree, run `scripts/flai.sh digest harvest --json` and record the candidates by kind, then `scripts/flai.sh digest seed --json` and record what it wrote and skipped. Run `flai check --strict` and fix through `flai digest edit` whatever it reports; an entry with a secret's shape is refused, which the notes record by thread ID, never by its text.
- Record the counts in the story's notes.
- Describe `flai digest harvest` and `flai digest seed` in `design/system/flai-cli.md` and `docs/users/flai.md`, and say in `design/digest/README.md` that the seed's entries are a starting point the curator merges.

Waits for T-1816: the seed runs through the command.

## Done when

The entries are committed on the story branch, a second `digest seed` writes nothing, `flai check --strict` reports nothing on them, and the notes give the counts.

## Notes

The entries' file names are not known until the seed writes them, so the task touches the `design/digest` folder; no other story writes there until S-0433.
