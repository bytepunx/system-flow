---
id: T-2912
type: task
nature: improvement
title: I-0278 is closed with the cause and the fix
status: done
parent: S-0669
owner: alex
created: 2026-11-10T08:47:51Z
updated: 2026-11-10T09:54:50Z
transitions:
  - to: ready
    at: 2026-11-10T09:54:36Z
    by: agent-S-0669
  - to: in-progress
    at: 2026-11-10T09:54:37Z
    by: agent-S-0669
  - to: done
    at: 2026-11-10T09:54:50Z
    by: agent-S-0669
stream: S-0669
tags: [issues]
touches: [design/issues/I-0278-after-a-stop-from-the-dashboard-the-story-s-agent-page-shows-the-agent-stopped.md]
after: [T-2909]
usage:
  source: log
  seconds: 13
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 3
      output: 255
      cache_read: 596853
      cache_write: 1816
      cost: 0.1387
---
# T-2912 I-0278 is closed with the cause and the fix

## Work

This task waits for T-2909, because the reason it writes names the fix T-2909 made. Close the issue:

```bash
flai issue close I-0278 --reason "<cause>; fixed by <change>, covered by <test>"
```

The reason names the cause from the narrative's `## Decisions`, the change T-2909 made, and the test that covers it. `flai issue close` rewrites `design/issues/summary.md` as well.

## Done when

- [ ] I-0278 is closed, and its reason names the cause, the fix, and the test that covers it.
- [ ] `flai check --strict` reports nothing on I-0278.

## Notes
