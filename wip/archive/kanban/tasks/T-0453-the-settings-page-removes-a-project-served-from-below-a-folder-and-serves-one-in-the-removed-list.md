---
id: T-0453
type: task
nature: improvement
title: The settings page removes a project served from below a folder and serves one in the removed list
status: done
parent: S-0123
owner: alex
created: 2026-09-26T07:06:22Z
updated: 2026-09-26T07:16:41Z
transitions:
  - to: ready
    at: 2026-09-26T07:06:29Z
    by: agent-S-0123
  - to: in-progress
    at: 2026-09-26T07:11:24Z
    by: agent-S-0123
  - to: done
    at: 2026-09-26T07:16:41Z
    by: agent-S-0123
stream: S-0123
tags: [dashboard]
touches: [flaiover, flai/cmd]
---
# T-0453 The settings page removes a project served from below a folder and serves one in the removed list

## Work

- `settings.get`'s projects view says `removed: true` on an unserved project in the list.
- The settings page offers Remove, behind the same confirmation as a registered project, on a project served from below a folder, and says Serve brings it back; on an unserved project in the list it says it was removed and offers Serve, which waits until the switcher has it.
- Component tests for both.

## Done when

- flaiover's tests, lint, and check pass; flai's `make test` passes.

## Notes
