---
id: T-3484
type: task
nature: feature
title: The state job runs the cluster's deploy and checks a state's runs names, skips them with the reason where no cluster can run them, and the cluster no longer deploys on review in code
status: done
parent: S-0810
owner: alex
created: 2026-11-16T23:51:27Z
updated: 2026-11-18T11:02:27Z
transitions:
  - to: ready
    at: 2026-11-18T10:37:52Z
    by: agent-S-0810
  - to: in-progress
    at: 2026-11-18T10:37:52Z
    by: agent-S-0810
  - to: done
    at: 2026-11-18T11:02:27Z
    by: agent-S-0810
stream: S-0810
tags: [flai]
touches: [flai/internal/serve/kluster.go, flai/internal/serve/kluster_test.go, flai/internal/serve/kluster_requests.go, flai/internal/serve/statejob.go, flai/internal/serve/statejob_test.go]
after: [T-3482]
usage:
  source: log
  seconds: 1475
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 151
      output: 70652
      cache_read: 17201163
      cache_write: 550832
      cost: 8.1373
---
# T-3484 The state job runs the cluster's deploy and checks a state's runs names, skips them with the reason where no cluster can run them, and the cluster no longer deploys on review in code

## Work

Waits for T-3482, whose state job and its step records this adds the `kluster:` steps to, in the same `flai/internal/serve/statejob.go`.

Today `lookReview` in `flai/internal/serve/kluster.go` queues a deploy, and `deploy` its checks, for each story in the state an acceptance leaves from (`flow.AwaitsAcceptance`), and `sendBack` sends a story back from there on a failed deploy or failed checks.

1. In the state job, after the suites, run `deploy` and then `checks` where the state's `runs:` names them, through the deploy and checks code `kluster.go` runs today, each as a step of the job's record with its log.
2. Skip a `kluster:` step, recording `skipped` and why, and go on to the job's end without holding the story, when: the project's `kluster:` block names no `deploy`; the kluster host action is off; `kluster.agent_namespaces` is 0, or every namespace it allows is held by another story; or, for `checks`, the `checks` host action is off. A skipped `deploy` skips `checks` with it.
3. A failed deploy or failed checks fails the job, which sends the story back as T-3482 does, with the kluster send-back thread as today.
4. Remove the deploy on `AwaitsAcceptance` from `lookReview` and `runIn`: the cluster deploys and runs its checks where a state's `runs:` names them and on no other state. Keep a deploy a story's agent asks for while in progress (ADR-0292) as it is, and make `kluster_requests.go`'s review branch read the states that name a `kluster:` step in place of the acceptance's state.
5. Tests: on the baseline a story entering `review` deploys once and runs its checks as before this story, each existing `kluster_test.go` case still passing; a workflow whose `integration` state names `deploy` deploys there and not on `review`; each skip reason is recorded and the story stays in its state; a failed deploy and failed checks send the story back by the state's send-back transition.

## Done when

- `flai test` passes on the changed paths.
- `kluster_test.go`'s baseline review cases pass unchanged, and a test shows a state whose `runs:` names no `kluster:` step deploys nothing.

## Notes
