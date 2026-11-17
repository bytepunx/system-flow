---
id: T-3466
type: task
nature: feature
title: A serve test enters the integration state in the suites project and checks the slow suites run once, the cluster steps are skipped or run, and a failing tier sends the story back
status: backlog
parent: S-0815
owner: alex
created: 2026-11-16T23:50:14Z
updated: 2026-11-16T23:50:14Z
transitions: []
stream: S-0815
tags: [flai]
touches: [flai/internal/serve/suites_e2e_test.go]
after: [T-3463]
---
# T-3466 A serve test enters the integration state in the suites project and checks the slow suites run once, the cluster steps are skipped or run, and a failing tier sends the story back

## Work

It waits for T-3463 because it reads that task's fixture, `flai/cmd/testdata/e2e-suites-board.md`, and its stub tiers.

Write `suites_e2e_test.go` in package `serve`, skipped under `-short`. Make the project from the template with that fixture as its workflow file, and drive flai serve's pass in process with the labs the package's tests have: `newVerifyLab` in `verifyjob_test.go`, and `newKlusterLab`, `fakeKluster`, and `withChecks` in `kluster_test.go` and `checks_test.go`. Move a story from `review` to `integration` and await the state's job, the one S-0810 starts. Check:

- With no `kluster:` block, the job rebases the branch, runs the `integration` and `smoke` tiers once, records the deploy and checks steps as skipped with the reason, and leaves the story in `integration`.
- The same with a `kluster:` block but the kluster host action off, and with no agent namespace free.
- A failing `integration` tier sends the story back by the state's `send-back` transition, with the thread its agent is started on.
- With a `kluster:` block, a deploy, checks, and the host action on, entering `integration` deploys the story and runs its checks once, and no other state runs them.

## Done when

The test passes in the integration tier, and fails when a slow suite runs twice, when a skipped cluster step holds the story, when a failing tier leaves it in `integration`, or when a configured cluster is not deployed on entry.

## Notes
