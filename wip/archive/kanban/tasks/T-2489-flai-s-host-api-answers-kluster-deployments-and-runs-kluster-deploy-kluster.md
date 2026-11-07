---
id: T-2489
type: task
nature: feature
title: flai's host API answers kluster.deployments and runs kluster.deploy, kluster.teardown, and kluster.promote in the background, refusing the writes while the kluster host action is off
status: done
parent: S-0581
owner: alex
created: 2026-11-06T21:45:55Z
updated: 2026-11-07T02:55:06Z
transitions:
  - to: ready
    at: 2026-11-07T02:37:27Z
    by: agent-S-0581
  - to: in-progress
    at: 2026-11-07T02:37:28Z
    by: agent-S-0581
  - to: done
    at: 2026-11-07T02:55:06Z
    by: agent-S-0581
stream: S-0581
tags: [kluster, hostapi]
touches: [flai/internal/hostapi/kluster.go, flai/internal/hostapi/kluster_test.go, flaiover/src/lib/server/agent.ts, flai/internal/kluster/deploy.go]
usage:
  source: log
  seconds: 1058
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 58
      output: 30082
      cache_read: 8046186
      cache_write: 131565
      cost: 3.1347
---
# T-2489 flai's host API answers kluster.deployments and runs kluster.deploy, kluster.teardown, and kluster.promote in the background, refusing the writes while the kluster host action is off

## Work

In `flai/internal/hostapi/kluster.go`, beside S-0580's methods, add:

- `kluster.deployments`, a read: S-0578's list of deployment records, each with the story's title and status joined in, a running deploy marked, and how to reach each deployment for a demo as S-0576's ADR decides (such as the hosts its services answer on and the port-forward command);
- `kluster.deploy` with a story, `kluster.teardown` with a story, and `kluster.promote`, writes: each refused, with a plain sentence, while the action is off, the block is missing, or a deploy runs for the project. Each starts S-0578's deploy, teardown, or promote in the background and answers at once.

The flai serve pass of S-0579 and these writes share S-0578's lock, so they never run two deploys at once.

Test in `kluster_test.go` with the package faked: the list with and without a running deploy, each write started, refused while off, refused while a deploy runs, and a deploy of a story not in review refused with the package's sentence.

This is the first layer: it waits for no task.

## Done when

- The read and the three writes work as above.
- Every case above is tested, and `flai test` on the two paths passes.

## Notes
