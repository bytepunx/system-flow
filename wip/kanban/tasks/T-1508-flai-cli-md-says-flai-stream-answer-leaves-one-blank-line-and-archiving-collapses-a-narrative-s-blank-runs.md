---
id: T-1508
type: task
nature: improvement
title: flai-cli.md says flai stream answer leaves one blank line and archiving collapses a narrative's blank runs
status: backlog
parent: S-0376
owner: alex
created: 2026-10-13T16:44:17Z
updated: 2026-10-13T16:44:17Z
transitions: []
stream: S-0376
tags: [docs]
touches: [design/system/flai-cli.md]
after: [T-1506, T-1507]
---
# T-1508 flai-cli.md says flai stream answer leaves one blank line and archiving collapses a narrative's blank runs

## Work

- In `design/system/flai-cli.md`, the `flai stream answer` row: removing the bullet leaves one blank line where it was, and an answer the project's markdown lint would reject is refused with its rule and line (S-0376, I-0139).
- The `flai archive` row: each narrative it moves has every run of blank lines outside fenced code collapsed to one, so no archived narrative carries an MD012 finding that only a hand commit could clear. Say that `flai accept` archives the same way.
- It waits for T-1506 and T-1507, because it describes what they built, as they built it.

## Done when

- [ ] Both rows of `design/system/flai-cli.md` describe the behaviour T-1506 and T-1507 built.
- [ ] `flai test design/system/flai-cli.md` passes the markdown lint.

## Notes
