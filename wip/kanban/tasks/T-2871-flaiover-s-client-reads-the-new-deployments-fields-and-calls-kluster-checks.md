---
id: T-2871
type: task
nature: feature
title: flaiover's client reads the new deployments fields and calls kluster.checks through a route of its own
status: backlog
parent: S-0663
owner: alex
created: 2026-11-10T08:25:19Z
updated: 2026-11-10T08:25:19Z
transitions: []
stream: S-0663
tags: [dashboard]
touches: [flaiover/src/lib/server/repo.ts, flaiover/src/lib/server/agent.ts, flaiover/src/routes/api/kluster/checks/+server.ts, flaiover/src/routes/api/kluster/checks/checks.test.ts, flaiover/src/routes/api/kluster/deployments.test.ts, flaiover/src/routes/api/kluster/forward/forward.test.ts]
after: [T-2869]
---
# T-2871 flaiover's client reads the new deployments fields and calls kluster.checks through a route of its own

## Work

Carry the new answer and the new write through flaiover's server side.

- In `flaiover/src/lib/server/repo.ts`, the types `KlusterDeployment` and `KlusterDeployments` gain the fields T-2869 adds: the story's column, whether the commit is the head, the last deploy and checks run with who asked, the stories in progress with no record, and the count against the limit with the refused requests. A field an older flai does not send is optional. Add a client method `klusterChecks(story)` beside `klusterDeploy`.
- In `flaiover/src/lib/server/agent.ts`, add `kluster.checks` to the `kluster.*` allow-list.
- Add `flaiover/src/routes/api/kluster/checks/+server.ts`, a POST as `deploy/+server.ts` is, that calls `klusterChecks` and shows flai's refusal word for word, with its test `checks/checks.test.ts`.
- In `deployments.test.ts`, show the route passing a story in progress's row and the count against the limit through. In `forward/forward.test.ts`, show a forward into a story in progress's namespace.

It waits for T-2869, whose answer the types follow. It runs beside T-2870: the method's name, `kluster.checks`, and its `{"story"}` parameter are fixed by this plan, and the two share no file.

## Done when

- The route tests pass with a fake host that answers the new fields and `kluster.checks`.
- `flai test` passes on the changed paths, `flaiover-lint`, `vitest`, and `svelte-check` among them.

## Notes
