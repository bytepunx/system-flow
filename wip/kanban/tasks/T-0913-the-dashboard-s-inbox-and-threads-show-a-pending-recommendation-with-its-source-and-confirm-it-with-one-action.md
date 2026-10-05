---
id: T-0913
type: task
nature: feature
title: The dashboard's inbox and threads show a pending recommendation with its source and confirm it with one action
status: backlog
parent: S-0220
owner: alex
created: 2026-10-05T04:48:30Z
updated: 2026-10-05T04:48:30Z
transitions: []
stream: S-0220
tags: [dashboard]
touches: ["flaiover/src/routes/api/threads/[id]/confirm/+server.ts", flaiover/src/routes/api/threads/threads.test.ts, flaiover/src/lib/server/inbox.ts, flaiover/src/lib/server/inbox.test.ts, flaiover/src/lib/components/Threads.svelte, flaiover/src/lib/components/Threads.svelte.test.ts, flaiover/src/lib/components/InboxView.svelte, flaiover/src/lib/components/Inbox.svelte.test.ts, design/system/flaiover-dashboard.md, docs/users/flaiover.md]
after: [T-0908]
---
# T-0913 The dashboard's inbox and threads show a pending recommendation with its source and confirm it with one action

## Work

- Add `POST /api/threads/:id/confirm` beside `reply` and `resolve` (`flaiover/src/routes/api/threads/[id]/`). It calls the host API's `thread.confirm` (T-0902) and answers as `reply` does.
- `Threads.svelte` marks an entry that is a recommendation, names its source as a link to the document and heading, and, on a writable dashboard, shows **Confirm** under a thread whose recommendation is pending. Confirm posts to the route and the thread shows answered without a reload.
- The inbox (`flaiover/src/lib/server/inbox.ts`, `InboxView.svelte`) shows a thread with a pending recommendation as such, with its text and source from `inbox.designer` (T-0908), and **Confirm** in place on a writable dashboard. The entry links to the thread as now.
- Describe the route, the mark, and the action in `design/system/flaiover-dashboard.md` (the `/inbox` and `/threads` rows and the threads API row) and in `docs/users/flaiover.md` (`## Threads` and the inbox).

It waits for T-0908, which puts the recommendation and its source in `inbox.designer`. It runs with T-0914, whose paths it does not share.

## Done when

- a route test posts a confirm and finds `thread.confirm` asked with the thread's ID
- component tests find the mark, the source link, and **Confirm** on a pending recommendation, none on a read-only dashboard, and the inbox entry with Confirm in place
- both documents describe the confirm action
- `npm test` and `npm run check` pass in `flaiover`

## Notes
