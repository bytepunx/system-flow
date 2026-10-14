---
id: T-1578
type: task
nature: feature
title: "POST /api/items/[id]/accept forwards keep_removals to accept.run and passes a refusal's undoes_main through"
status: backlog
parent: S-0388
owner: alex
created: 2026-10-14T21:08:50Z
updated: 2026-10-14T21:08:50Z
transitions: []
stream: S-0388
tags: [dashboard]
touches: ["flaiover/src/routes/api/items/[id]/accept/+server.ts", "flaiover/src/routes/api/items/[id]/accept/accept.test.ts"]
---
# T-1578 POST /api/items/[id]/accept forwards keep_removals to accept.run and passes a refusal's undoes_main through

## Work

- In `flaiover/src/routes/api/items/[id]/accept/+server.ts`, read `keep_removals` from the body: a list of `{ path, reason }`.
- Pass it to `accept.run` beside `include_uncommitted`, and only when it is not empty.
- Answer 400 for an entry with no path or an empty reason, before calling flai.
- Take the field names from `accept.run` as S-0371's T-1575 gives them in `flai/internal/hostapi/writes.go`.
- The refusal's `data`, which S-0378 passes through, now carries `undoes_main` as flai gives it. Test that it is passed through whole. Change the pass-through only if the test shows it drops the field.
- It waits for no task of this story.

## Done when

- [ ] `accept.test.ts` shows `keep_removals` forwarded, left out when empty, and a bad entry answered 400
- [ ] `accept.test.ts` shows a refusal whose data carries `undoes_main` passed through whole
- [ ] `flai test` on the two files passes

## Notes
