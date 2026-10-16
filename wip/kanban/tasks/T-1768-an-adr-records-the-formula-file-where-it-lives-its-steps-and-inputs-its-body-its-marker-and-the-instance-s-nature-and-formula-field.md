---
id: T-1768
type: task
nature: feature
title: "An ADR records the formula file: where it lives, its steps and inputs, its body, its marker, and the instance's nature and formula field"
status: backlog
parent: S-0421
owner: alex
created: 2026-10-16T05:47:09Z
updated: 2026-10-16T05:47:09Z
transitions: []
stream: S-0421
tags: [cli]
touches: [design/adrs, design/adrs/README.md]
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
