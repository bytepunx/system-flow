---
id: T-0435
type: task
nature: remediation
title: Every project's Repo listens for flai's changes
status: done
parent: S-0117
owner: arobson
created: 2026-09-26T03:29:51Z
updated: 2026-09-26T03:31:37Z
transitions:
  - to: ready
    at: 2026-09-26T03:29:54Z
    by: agent-S-0117
  - to: in-progress
    at: 2026-09-26T03:29:54Z
    by: agent-S-0117
  - to: done
    at: 2026-09-26T03:31:37Z
    by: agent-S-0117
stream: S-0117
tags: [dashboard]
touches: [flaiover/src, design/system/flaiover-dashboard.md]
---
# T-0435 Every project's Repo listens for flai's changes

## Work

- In `repo()` (`flaiover/src/lib/server/repo.ts`), start the new Repo's watcher when it is made for a key, so a Repo for a named project listens for flai's change, connected, and gone events like the default one.
- Test in `project-events.test.ts`, with real connections: a named project's board is asked, flai reports a removed story file, and the board is asked again of flai rather than answered from the cache.
- Test that a Repo given its own source does not listen, and that `watch()` twice registers one set of listeners.
- Say so in `design/system/flaiover-dashboard.md` beside S-0095's entry.

## Done when

- The new tests fail without the change and pass with it, and the flaiover tests and checks pass.

## Notes
