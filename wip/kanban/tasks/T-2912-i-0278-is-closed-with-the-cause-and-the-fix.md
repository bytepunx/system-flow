---
id: T-2912
type: task
nature: improvement
title: I-0278 is closed with the cause and the fix
status: backlog
parent: S-0669
owner: alex
created: 2026-11-10T08:47:51Z
updated: 2026-11-10T08:47:51Z
transitions: []
stream: S-0669
tags: [issues]
touches: [design/issues/I-0278-after-a-stop-from-the-dashboard-the-story-s-agent-page-shows-the-agent-stopped.md]
after: [T-2909]
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
