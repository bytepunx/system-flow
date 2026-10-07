---
id: T-1181
type: task
nature: improvement
title: Item, thread, inbox, review, activity, and host views show their times in local time
status: in-progress
parent: S-0329
owner: alex
created: 2026-10-07T19:49:23Z
updated: 2026-10-07T20:26:26Z
transitions:
  - to: ready
    at: 2026-10-07T20:26:25Z
    by: agent-S-0329
  - to: in-progress
    at: 2026-10-07T20:26:26Z
    by: agent-S-0329
stream: S-0329
tags: [dashboard]
touches: ["flaiover/src/routes/items/[id]/+page.svelte", "flaiover/src/routes/items/[id]/item.svelte.test.ts", flaiover/src/lib/components/Threads.svelte, flaiover/src/lib/components/Threads.svelte.test.ts, flaiover/src/routes/threads/threads.svelte.test.ts, flaiover/src/lib/components/InboxView.svelte, flaiover/src/lib/components/Inbox.svelte.test.ts, flaiover/src/lib/components/OpenQuestions.svelte, flaiover/src/lib/components/OpenQuestions.svelte.test.ts, flaiover/src/lib/components/Review.svelte, flaiover/src/lib/components/Review.svelte.test.ts, flaiover/src/lib/components/ActivityView.svelte, flaiover/src/lib/components/ActivityView.svelte.test.ts, flaiover/src/routes/activity/+page.svelte, flaiover/src/routes/activity/activity.svelte.test.ts, flaiover/src/lib/components/AgentStopConfirm.svelte, flaiover/src/lib/components/HostFlai.svelte, flaiover/src/lib/components/HostFlai.svelte.test.ts, flaiover/src/lib/components/HostProcesses.svelte, flaiover/src/lib/components/HostProcesses.svelte.test.ts]
after: [T-1180]
---
# T-1181 Item, thread, inbox, review, activity, and host views show their times in local time

## Work

These views print flai's UTC strings as they are. Replace each with `localTime` or `localDate` from `flaiover/src/lib/localtime.ts`, and update the tests that assert the UTC text to the zone the tests run in.

| File | What it shows today |
|------|---------------------|
| `flaiover/src/routes/items/[id]/+page.svelte` | `item.created`, each transition's `t.at`, each blocked interval's `from` and `until` raw; the acceptance date as `t.at.slice(0, 10)` |
| `flaiover/src/lib/components/Threads.svelte` | each entry's `e.at` raw |
| `flaiover/src/lib/components/InboxView.svelte` | each change's `e.at` raw |
| `flaiover/src/lib/components/OpenQuestions.svelte` | each question's `q.at` raw |
| `flaiover/src/lib/components/Review.svelte` | the checks run's `started` and `ended` raw |
| `flaiover/src/lib/components/ActivityView.svelte` | `last_log.at` raw, and `updated` raw in a `title` |
| `flaiover/src/routes/activity/+page.svelte` | the orchestrator's `started` and `ended` raw |
| `flaiover/src/lib/components/AgentStopConfirm.svelte` | the run's `started` raw |
| `flaiover/src/lib/components/HostFlai.svelte` | `connected since` raw in a `title` |
| `flaiover/src/lib/components/HostProcesses.svelte` | a release's `published` as `.slice(0, 10)` |

A relative age (`age()`, `3h ago`) stays as it is. A `title` counts as a display: it shows local time too.

It waits for T-1180 because it calls its formatter. It shares no file with the other two view tasks, so the three run together.

## Done when

- None of the files above shows a time or a date in UTC; each shows it through `localTime` or `localDate`.
- Their tests assert the local text in the tests' pinned zone, and `flai test` on the files above passes.

## Notes

- `HostProcesses.svelte` is the Host › Updates list; `AgentStopConfirm.svelte` has no test file of its own, so its time is checked through the test of the view that opens it, or a new test beside it.
