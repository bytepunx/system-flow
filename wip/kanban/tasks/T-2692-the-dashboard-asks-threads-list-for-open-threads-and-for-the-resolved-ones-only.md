---
id: T-2692
type: task
nature: improvement
title: "The dashboard asks threads.list for open threads, and for the resolved ones only when a page asks with all=1"
status: backlog
parent: S-0611
owner: alex
created: 2026-11-07T23:56:42Z
updated: 2026-11-07T23:56:42Z
transitions: []
stream: S-0611
tags: [dashboard]
touches: [flaiover/src/lib/server/repo.ts, flaiover/src/routes/api/threads/+server.ts, flaiover/src/routes/api/threads/threads.test.ts, flaiover/src/lib/server/repo.test.ts, flaiover/src/lib/server/repo-channel.test.ts, flaiover/src/lib/server/repo-forget.test.ts]
---
# T-2692 The dashboard asks threads.list for open threads, and for the resolved ones only when a page asks with all=1

## Work

Criterion 1. `repo.threads()` and `repo.threadsFor(on)` (`flaiover/src/lib/server/repo.ts`, about lines 884 to 893) always ask flai with `all: true`, and `/api/threads` (`flaiover/src/routes/api/threads/+server.ts`) drops the resolved ones afterwards unless the page sent `all=1`. flai's `threads.list` already answers open threads alone when `all` is absent (`flai/internal/hostapi/hostapi.go`, about line 613), so this task changes only the dashboard.

- Give `threads()` and `threadsFor(on)` an `all` argument, false by default, sent to flai only when true. Put `all` in their cache keys, so an open list and a whole one are kept apart, and make the forget table (about line 1242) forget both on a `thread` or `item` change.
- `/api/threads` passes its `all=1` through and stops filtering after the fact.
- The pages need no change: `Threads.svelte` sends `all=1` only when its show-resolved box is ticked, and the item page's `askThreads` sends `all=1` for one item's threads, which it shows.

Waits for no task: no other task touches these files.

## Done when

- [ ] `repo-channel.test.ts` checks that `threads()` asks flai with no `all`, and `threads(true)` and `threadsFor(on, true)` with `all: true`.
- [ ] `repo-forget.test.ts` checks a thread change forgets the open and the whole lists.
- [ ] `threads.test.ts` checks `/api/threads` asks for resolved threads only with `all=1`.
- [ ] `flai test` passes on the changed paths.

## Notes

- Planned by planner-S-0611.
