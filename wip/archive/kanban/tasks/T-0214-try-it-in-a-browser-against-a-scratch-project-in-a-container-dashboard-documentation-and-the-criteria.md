---
id: T-0214
type: task
nature: feature
title: Try it in a browser against a scratch project in a container, dashboard documentation, and the criteria
status: done
parent: S-0057
owner: alex
created: 2026-09-19T06:41:25Z
updated: 2026-09-19T07:04:37Z
transitions:
  - to: ready
    at: 2026-09-19T06:59:57Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T06:59:58Z
    by: system-flow
  - to: done
    at: 2026-09-19T07:04:37Z
    by: system-flow
stream: S-0057
tags: []
---

# T-0214 Try it in a browser against a scratch project in a container, dashboard documentation, and the criteria

## Work
Build the branch's image, run it against a scratch project (never this repository: the check writes) under its own container name and port with a throwaway token, and try: columns in pull order with unplaced stories last; drag within ready and within backlog, the card staying put after the refresh and `flai board` agreeing; a drag within in-progress doing nothing; a drag to another column still moving; the buttons and the keyboard shortcut; both themes. `design/system/flaiover-dashboard.md` and `docs/users/flaiover.md` describe it. Tick the criteria for what was seen.

## Done when
- The browser check is recorded in the narrative with what was done and seen
- The dashboard design and user docs describe reordering
- The criteria are ticked and `flai check --strict` is clean

## Notes
