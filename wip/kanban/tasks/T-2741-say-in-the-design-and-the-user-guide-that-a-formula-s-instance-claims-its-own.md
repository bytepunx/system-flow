---
id: T-2741
type: task
nature: remediation
title: Say in the design and the user guide that a formula's instance claims its own story file
status: backlog
parent: S-0620
owner: alex
created: 2026-11-09T17:06:05Z
updated: 2026-11-09T17:06:05Z
transitions: []
stream: S-0620
tags: [docs, formulas]
touches: [design/system/strategic-agents.md, design/system/flai-cli.md, docs/users/flai.md]
after: [T-2740]
---
# T-2741 Say in the design and the user guide that a formula's instance claims its own story file

## Work

Add the instance's touch to each place that says what `flai formula run` makes. The instance story carries its own story file as its one touch, so a story in progress never holds it `no-touches`, and it holds no ready story while it is in progress.

- `design/system/strategic-agents.md`, § Formulas (S-0421), the bullet "What run makes".
- `design/system/flai-cli.md`, § Commands, the `flai formula run` row.
- `docs/users/flai.md`, § A story from a formula.

Waits on T-2740, because it describes the behaviour T-2740 builds, as T-2740 left it.

## Done when

- Each of the three places names the instance's touch and why it has one.
- `flai test` on the three files passes, the markdown lint among its tiers.

## Notes
