---
id: T-1743
type: task
nature: feature
title: serve queues a role whose trigger names an event when that event happens, once per role and item, and drains the queue one run at a time per role
status: backlog
parent: S-0418
owner: alex
created: 2026-10-16T05:45:27Z
updated: 2026-10-16T05:45:42Z
transitions: []
stream: S-0418
tags: [cli]
touches: [flai/internal/serve/role_events.go, flai/internal/serve/role_events_test.go]
after: [T-1741]
---
# T-1743 serve queues a role whose trigger names an event when that event happens, once per role and item, and drains the queue one run at a time per role

## Work

Add a role event queue, modelled on the replanner's queue in `replan.go`.

- At each look of the launcher, read the moves, edit notices, and thread changes since the last look.
- Queue each role whose trigger names the event, while its action is on, once per role and item, adding a second trigger to the entry.
- Drain the queue one run at a time per role through `RunRole`, dropping an entry its checks refuse and logging `queued role dropped` with the reason.
- While the action is off, drop the queue and keep up with events, as the replanner does.
- Test each event and the coalescing.

Waits for T-1741. Runs beside T-1744 and T-1746, since they share no path.

## Done when

`go test ./internal/serve/...` passes.

## Notes
