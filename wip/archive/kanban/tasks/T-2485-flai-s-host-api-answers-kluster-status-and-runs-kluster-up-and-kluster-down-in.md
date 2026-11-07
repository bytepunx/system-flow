---
id: T-2485
type: task
nature: feature
title: flai's host API answers kluster.status and runs kluster.up and kluster.down in the background, refusing the writes while the kluster host action is off
status: done
parent: S-0580
owner: alex
created: 2026-11-06T21:45:01Z
updated: 2026-11-07T01:45:22Z
transitions:
  - to: ready
    at: 2026-11-07T01:33:36Z
    by: agent-S-0580
  - to: in-progress
    at: 2026-11-07T01:33:37Z
    by: agent-S-0580
  - to: done
    at: 2026-11-07T01:45:22Z
    by: agent-S-0580
stream: S-0580
tags: [kluster, hostapi]
touches: [flai/internal/hostapi/kluster.go, flai/internal/hostapi/kluster_test.go, flai/internal/hostapi/writes.go, flai/cmd/serve_actions.go, flaiover/src/lib/server/agent.ts]
usage:
  source: log
  seconds: 705
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 67
      output: 19960
      cache_read: 9251833
      cache_write: 111294
      cost: 3.0152
---
# T-2485 flai's host API answers kluster.status and runs kluster.up and kluster.down in the background, refusing the writes while the kluster host action is off

## Work

Add `flai/internal/hostapi/kluster.go` with its specs, merged into the method table in `writes.go` as `componentSpecs()` is:

- `kluster.status`, a read: the cluster's name, profile, provider, whether it runs, its age, whether an up or down is running, the last one's outcome and log path, and whether the project's `kluster` host action is on and its `kluster:` block set;
- `kluster.up` and `kluster.down`, writes: each refused, with a plain sentence, while the action is off or the block missing, and while another up or down runs. Each starts the S-0577 package's up or down in the background, since either takes minutes, and answers at once. Its outcome and log are kept for `kluster.status`, as `components.start` keeps its server's state.

The components writes run in flai serve's own process, which `flai/cmd/serve_actions.go` wires. Run the kluster jobs the same way, and wire them there.

Test in `kluster_test.go` with the package faked: the status read with the action on and off, each write started, refused while off, refused with no block, and refused while one runs.

This is the first layer: it waits for no task.

## Done when

- The read and the two writes work as above and are in the method table.
- Every case above is tested, and `flai test` on the four paths passes.

## Notes
