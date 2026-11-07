---
id: T-2483
type: task
nature: feature
title: The review checks of a deployed story run after its deploy ends, with the namespace and the kube context in their environment
status: done
parent: S-0579
owner: alex
created: 2026-11-06T21:44:28Z
updated: 2026-11-07T03:13:50Z
transitions:
  - to: ready
    at: 2026-11-07T03:05:32Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T03:05:33Z
    by: system-flow
  - to: done
    at: 2026-11-07T03:13:50Z
    by: agent-S-0579
stream: S-0579
tags: [kluster, serve, checks]
touches: [flai/internal/serve/checks.go, flai/internal/serve/checks_test.go, flai/cmd/checks.go, flai/internal/serve/kluster.go, flai/internal/serve/kluster_test.go, flai/cmd/serve.go, flai/internal/serve/serve.go]
after: [T-2482]
usage:
  source: log
  seconds: 497
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 48
      output: 295
      cache_read: 8152369
      cache_write: 61069
      cost: 2.9309
---
# T-2483 The review checks of a deployed story run after its deploy ends, with the namespace and the kube context in their environment

## Work

It waits for T-2482: it starts the checks from the kluster pass in `flai/internal/serve/kluster.go`, which T-2481 writes and T-2482 extends, and reads the deployment records T-2481's deploys write.

flai serve starts no checks run on its own today: a run is started by `flai checks run` or the dashboard's `checks.run` (`flai/cmd/checks.go`), and runs through `Dir.RunChecks` in `flai/internal/serve/checks.go`. So:

- `RunChecks` takes the extra environment a run gets, and adds it to each check's process environment. With none, it runs as before.
- A checks run for a story whose namespace record (`kluster.ReadRecord`) has passed at the story branch's head, while the `kluster` host action is on and the manifest has a `kluster:` block, gets ADR-0254 §4's variables for that namespace (`Cluster.Env`), whoever starts it. Otherwise it gets none.
- With `kluster` and `checks` both on, the kluster pass starts the story's checks run itself once its deploy has passed at the head, in the same job queue, so it runs after the deploy and before the next job. A failed deploy starts none. A failed checks run against the namespace opens one thread on the story as a failed deploy does, and moves nothing.

Test in `checks_test.go` and `kluster_test.go`: the variables present after a passing deploy at the head, absent after a failing one, absent at an older head, absent with the action off; the pass starting the checks after a passing deploy with `checks` on, none with it off, and none after a failure.

## Done when

- A deployed story's checks see the namespace and the context, and the pass runs them after its deploy ends.
- Every case above is tested, and `flai test` on the four paths passes.

## Notes
