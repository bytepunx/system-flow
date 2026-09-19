---
id: T-0205
type: task
nature: improvement
title: Verify in a browser at three widths in both themes, and update the docs if what a card shows changed
status: done
parent: S-0054
owner: alex
created: 2026-09-19T05:44:38Z
updated: 2026-09-19T05:55:05Z
transitions:
  - to: ready
    at: 2026-09-19T05:55:04Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T05:55:04Z
    by: system-flow
  - to: done
    at: 2026-09-19T05:55:05Z
    by: system-flow
stream: S-0054
tags: []
---

# T-0205 Verify in a browser at three widths in both themes, and update the docs if what a card shows changed

## Work
Look at the board in a browser at desktop, tablet, and phone widths, in the light and the dark theme, with "epics and tasks too" on and off, measuring as S-0048 did: the parent bottom right, the title ending above the divider, no card overflowing. Update `docs/users/flaiover.md` only if what a card shows changed. Run `make flaiover-test` and `flai check --strict`. Tick the story's criteria for what was seen.

## Done when
- The measurements at three widths in both themes are in the narrative
- Every criterion on S-0054 is checked, or unchecked with the reason in the story notes

## Notes
