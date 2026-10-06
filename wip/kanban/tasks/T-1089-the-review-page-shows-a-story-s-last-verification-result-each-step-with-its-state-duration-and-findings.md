---
id: T-1089
type: task
nature: improvement
title: The review page shows a story's last verification result, each step with its state, duration, and findings
status: backlog
parent: S-0270
owner: alex
created: 2026-10-06T22:52:56Z
updated: 2026-10-06T22:52:56Z
transitions: []
stream: S-0270
tags: [flaiover, dashboard]
touches: ["flaiover/src/routes/api/items/[id]/verify/+server.ts", "flaiover/src/routes/api/items/[id]/verify/verify.test.ts", flaiover/src/lib/server/agent.ts, flaiover/src/lib/review.ts, flaiover/src/lib/review.test.ts, flaiover/src/lib/components/Review.svelte, flaiover/src/lib/components/Review.svelte.test.ts]
after: [T-1076]
---
# T-1089 The review page shows a story's last verification result, each step with its state, duration, and findings

## Work

Criterion 2: the dashboard's story page can show the last result.

- Add `flaiover/src/routes/api/items/[id]/verify/+server.ts`, beside the `checks` route. Its GET calls the host method `verify.status` through `flaiover/src/lib/server/agent.ts`, as the `checks` route calls `checks.status`, and answers nothing when there is no result or no host.
- In `flaiover/src/lib/review.ts`, add the result's types and a function that turns it into rows: step, state, duration, and the findings under a failed step, with the notes apart.
- In `flaiover/src/lib/components/Review.svelte`, show the last result, with the commit it verified and when, beside the checks the panel shows already. Show nothing when there is none. No button runs it from here.
- Tests beside each file.
- Waits for T-1076, whose `verify.status` it reads.

## Done when

- [ ] A story in review with a stored result shows each step's state and duration, and the findings of the step that failed.
- [ ] A story with no result, or with no host connected, shows nothing and no error.
- [ ] `scripts/flaiover-test.sh` passes.

## Notes

Drafted by the planner. It assumes "the story's page" in criterion 2 is the review page, `/review/[id]`, where `Review.svelte` already shows a story's checks. This assumption is on the plan thread.
