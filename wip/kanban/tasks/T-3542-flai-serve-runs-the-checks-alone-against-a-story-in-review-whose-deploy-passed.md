---
id: T-3542
type: task
nature: remediation
title: flai serve runs the checks alone against a story in review whose deploy passed at its head and whose checks have not, and tries each head's checks once
status: in-progress
parent: S-0835
owner: alex
created: 2026-11-17T08:30:57Z
updated: 2026-11-28T00:35:36Z
transitions:
  - to: ready
    at: 2026-11-28T00:35:36Z
    by: agent-S-0835
  - to: in-progress
    at: 2026-11-28T00:35:36Z
    by: agent-S-0835
stream: S-0835
tags: [flai]
touches: [flai/internal/serve/kluster.go, flai/internal/serve/kluster_test.go, flai/internal/serve/checks_test.go, flai/internal/serve/kluster_requests_test.go, flai/internal/serve/statejob.go, flai/internal/serve/statejob_test.go]
after: [T-3540]
---
# T-3542 flai serve runs the checks alone against a story in review whose deploy passed at its head and whose checks have not, and tries each head's checks once

## Work

Waits for T-3540, whose judgement this reads in place of `AcceptGate`.

Where the review's deploy is queued: `lookReview` in `flai/internal/serve/kluster.go` today; the state job's `kluster:` steps in `flai/internal/serve/statejob.go` once S-0810 has moved them there. Make the change wherever S-0810 left it.

1. Read T-3540's judgement for each story entering review or at a new head. Nothing: queue nothing, as today. A deploy: queue `klusterDeploy`, then its checks, as today. The checks alone: queue `klusterChecks` for the story at its head, or, in the state job, skip the `deploy` step with the reason that it passed at the head and run `checks`.
2. The checks run the review asks for takes the deploy record at the head as `kluster.AtHead` judges it. Today `klusterEnvWhy` refuses a record whose commit is not the exact head. Keep the exact head for a checks request a story's agent makes, as `design/system/kluster.md` § At the head says, and widen only the review's.
3. Record a checks-alone try per head, as `writeKlusterTried` does for a review deploy, so that a checks run that is refused or interrupted at a head is not queued again every minute. A new head is tried again.
4. A failed checks run alone sends the story back with the checks failure thread, as a failed checks run after a review deploy does today (`runChecks`, `sendBack`).
5. Tests that reproduce I-0362: a story whose asked-for deploy passed at its head and whose checks failed enters review and gets one checks run, no deploy; the same with checks never run; a record followed only by `design`, `docs`, and `wip` commits gets the checks alone; a failed checks run alone sends the story back. Change the assertion of `TestAnAskedDeploysFailingChecksOpenNoThreadAndReviewDeploysAgain` in `kluster_requests_test.go`, which pins today's deploy, to the checks run alone, and rename it to say so.

## Done when

- `flai test` passes on the changed paths.
- A test shows a story in review whose deploy passed at its head and whose checks failed gets one checks run and no deploy, and the existing review-pass tests in `kluster_test.go` and `checks_test.go` still pass.

## Notes
