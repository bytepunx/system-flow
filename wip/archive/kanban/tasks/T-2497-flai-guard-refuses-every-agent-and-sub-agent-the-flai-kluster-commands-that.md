---
id: T-2497
type: task
nature: feature
title: flai guard refuses every agent and sub-agent the flai kluster commands that write, and lets each run flai kluster status
status: done
parent: S-0577
owner: alex
created: 2026-11-07T00:54:32Z
updated: 2026-11-07T01:08:53Z
transitions:
  - to: ready
    at: 2026-11-07T00:55:44Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T00:55:46Z
    by: system-flow
  - to: done
    at: 2026-11-07T01:08:53Z
    by: agent-S-0577
stream: S-0577
tags: [kluster, guard]
touches: [flai/internal/guard/guard.go, flai/internal/guard/kluster_test.go, flai/internal/guard/call.go]
usage:
  source: log
  seconds: 787
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 21
      output: 6513
      cache_read: 2909088
      cache_write: 32522
      cost: 0.9415
---
# T-2497 flai guard refuses every agent and sub-agent the flai kluster commands that write, and lets each run flai kluster status

## Work

`design/system/kluster.md` § The `kluster` host action and ADR-0254 §7 say `flai guard` refuses every agent and sub-agent the `flai kluster` commands that write, since they change shared infrastructure on the operator's machine (`safety.md`), and that `flai kluster status` is a read every agent may run. No planned task covered it, so the story's agent added this one.

In `flai/internal/guard/guard.go`, refuse, as `flai forge push` is refused (S-0393): every session `flai serve` started (`FLAI_STARTED_BY=flai-serve`, `FLAI_ROLE`, or `FLAI_STORY` set) and every sub-agent, a shell line that runs `flai kluster up`, `down`, `deploy`, `teardown`, or `promote` in any form the existing rules find (`scripts/flai.sh`, `bin/flai`, a chain, an env prefix). Name the subcommands that write in one list, so S-0578 adds nothing to the guard. The refusal says why (the operator's shared infrastructure) and what to do instead: ask the operator on a thread, or run `flai kluster status`, which reads. Leave `flai kluster status` and `flai kluster --help` alone.

Test in a new `flai/internal/guard/kluster_test.go`: up and down refused to a served story's agent, to a sub-agent, and to the planner; status allowed to each; the operator's own shell (no flai serve variables, no sub-agent) allowed up and down.

It waits for no task: it reads command lines, not the manifest.

## Done when

- The guard refuses the write subcommands and lets status through, as tested.
- `flai test` on the two paths passes.

## Notes
