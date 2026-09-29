---
id: T-0587
type: task
nature: feature
title: The design and the user guide say what the review's diff shows, every tier passes, and the panel is looked at in light and dark
status: done
parent: S-0164
owner: alex
created: 2026-09-29T21:12:38Z
updated: 2026-09-29T21:22:42Z
transitions:
  - to: ready
    at: 2026-09-29T21:13:08Z
    by: agent-S-0164
  - to: in-progress
    at: 2026-09-29T21:19:43Z
    by: agent-S-0164
  - to: done
    at: 2026-09-29T21:22:42Z
    by: agent-S-0164
stream: S-0164
tags: []
touches: [design/system/flaiover-dashboard.md, docs/users/flaiover.md]
usage:
  source: log
  seconds: 179
  estimated: true
  models:
    - model: claude-fable-5-1
      input: 29
      output: 11644
      cache_read: 1639850
      cache_write: 51043
      cost: 2.0133
---

# T-0587 The design and the user guide say what the review's diff shows, every tier passes, and the panel is looked at in light and dark

## Work

Say in `design/system/flaiover-dashboard.md` and `docs/users/flaiover.md` what the review's diff shows. Run the three tiers and every lint. Render the panel from a real patch in light and dark and look at it. Tick the criteria that were verified and say in the story's notes how each was.

## Done when

- The design and the user guide describe the toggle, the margin, and the outlined runs.
- `make flai-test`, the dashboard's tests and lints, the markdown lint, and `flai check --strict` pass, or what fails is named with its cause.
- Each acceptance criterion is ticked with what verified it, or left unticked with why.

## Notes
