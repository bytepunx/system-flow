---
id: T-0208
type: task
nature: feature
title: Verify blocked and dragging on every tint in a browser in both themes, and describe the coding in the design and the user docs
status: done
parent: S-0055
owner: alex
created: 2026-09-19T06:05:58Z
updated: 2026-09-19T06:22:49Z
transitions:
  - to: ready
    at: 2026-09-19T06:19:55Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T06:19:55Z
    by: system-flow
  - to: done
    at: 2026-09-19T06:22:49Z
    by: system-flow
stream: S-0055
tags: []
---

# T-0208 Verify blocked and dragging on every tint in a browser in both themes, and describe the coding in the design and the user docs

## Work
Run the branch's dev server against this repository with a throwaway token and look at the real board with epics and tasks shown, in the light and the dark theme, at desktop and phone widths: every nature present is tinted, types can be told apart, a blocked card and a card being dragged stand out on each tint, the legend reads correctly, nothing overflows. Measure the computed colours against the tokens. Describe the coding, the legend, and the colour blindness waiver with its reason in `design/system/flaiover-dashboard.md` and `docs/users/flaiover.md`. Tick the story's criteria only for what was seen.

## Done when
- The browser check is recorded in the narrative with what was measured
- The design document and the user docs describe the coding and record the waiver
- The story's criteria are ticked and `flai check --strict` is clean

## Notes
