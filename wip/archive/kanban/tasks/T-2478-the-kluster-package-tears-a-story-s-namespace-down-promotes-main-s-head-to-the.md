---
id: T-2478
type: task
nature: feature
title: The kluster package tears a story's namespace down, promotes main's head to the default namespace, and lists the deployments with the default namespace first
status: done
parent: S-0578
owner: alex
created: 2026-11-06T21:43:33Z
updated: 2026-11-07T01:47:31Z
transitions:
  - to: ready
    at: 2026-11-07T01:42:26Z
    by: agent-S-0578
  - to: in-progress
    at: 2026-11-07T01:42:26Z
    by: agent-S-0578
  - to: done
    at: 2026-11-07T01:47:31Z
    by: agent-S-0578
stream: S-0578
tags: [kluster, flai]
touches: [flai/internal/kluster/deploy.go, flai/internal/kluster/deploy_test.go]
after: [T-2477]
usage:
  source: log
  seconds: 305
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 23
      output: 8902
      cache_read: 3362870
      cache_write: 38348
      cost: 1.1276
---
# T-2478 The kluster package tears a story's namespace down, promotes main's head to the default namespace, and lists the deployments with the default namespace first

## Work

It waits for T-2477, whose deploy steps, lock, logs, and record it reuses, in the same files.

In `flai/internal/kluster/deploy.go`, add:

- teardown of a story: run the block's teardown command with the namespace and the story's variables, or, when the block names none, delete the namespace with `kubectl delete namespace` on the kube context; log it; remove the story's record. A story with no record is a no-op that says so.
- promote: build, load, and deploy main's head from the main checkout into the default namespace, with the same steps, logs, and lock as a story's deploy, recorded under the default namespace.
- list: every record, the default namespace's first, then the stories by ID, with a running deploy shown as running.

Test each in `deploy_test.go` with the fakes T-2477 wrote: a teardown with and without a teardown command, a teardown of a story never deployed, a promote that passes and one that fails, and a list with and without a running deploy.

## Done when

- Teardown, promote, and list work as above, each tested.
- `flai test` on the two paths passes.

## Notes
