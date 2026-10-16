---
id: T-1852
type: task
nature: feature
title: The CLI design, the user docs, and the design system tool's design say file.get and file.save are built, what they reach, and the size limit
status: done
parent: S-0436
owner: alex
created: 2026-10-16T21:31:53Z
updated: 2026-10-16T22:37:40Z
transitions:
  - to: ready
    at: 2026-10-16T22:33:36Z
    by: agent-S-0436
  - to: in-progress
    at: 2026-10-16T22:33:37Z
    by: agent-S-0436
  - to: done
    at: 2026-10-16T22:37:40Z
    by: agent-S-0436
stream: S-0436
tags: [cli, dashboard]
touches: [design/system/flai-cli.md, docs/users/flai.md, docs/users/flaiover.md, design/system/design-system-tool.md]
after: [T-1851]
usage:
  source: log
  seconds: 238
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 70
      output: 334
      cache_read: 2991018
      cache_write: 116389
      cost: 1.4434
---
# T-1852 The CLI design, the user docs, and the design system tool's design say file.get and file.save are built, what they reach, and the size limit

## Work

- `design/system/flai-cli.md` and `docs/users/flai.md`: `flai file show` and `flai file save`, what they reach and refuse, the size limit, the exit codes, and that a token save writes the stylesheet with it.
- `docs/users/flaiover.md` § Design system: the token save now lands on the host, with the stylesheet in the same commit.
- `design/system/design-system-tool.md` § file.get and file.save: built, with the limit.

Waits for T-1851, so that the docs describe what the code does.

## Done when

The four files say the same as the code and pass `flai test` on them.

## Notes

None.
