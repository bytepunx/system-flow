---
id: T-2897
type: task
nature: improvement
title: pilot starts a probe without the story's FLAI_STORY and refuses a probe that names a role
status: done
parent: S-0682
owner: alex
created: 2026-11-10T08:46:08Z
updated: 2026-11-10T09:07:50Z
transitions:
  - to: ready
    at: 2026-11-10T08:59:41Z
    by: agent-S-0682
  - to: in-progress
    at: 2026-11-10T08:59:42Z
    by: agent-S-0682
  - to: done
    at: 2026-11-10T09:07:50Z
    by: agent-S-0682
stream: S-0682
tags: [pilot, containers]
touches: [flai/internal/pilot/request.go, flai/internal/pilot/request_test.go, flai/internal/pilot/run.go, flai/internal/pilot/run_test.go]
after: [T-2896]
usage:
  source: log
  seconds: 488
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 11
      output: 2059
      cache_read: 1241514
      cache_write: 22099
      cost: 0.4382
---
# T-2897 pilot starts a probe without the story's FLAI_STORY and refuses a probe that names a role

## Work

Waits for T-2896, whose `Probe` field and probe prompt this task's request and tests read.

- In `flai/internal/pilot/request.go`, have `Check` refuse a request with `probe` set that also names a `Role` or an `Item`: a probe is a trial of the image, never a strategic agent's run. The refusal is one line naming the fields, as `Check`'s others are.
- In `flai/internal/pilot/run.go`, have `childEnv` leave `FLAI_STORY` unset for a probe, so the inner harness's flai and the guard hook do not take the session for the story's agent; `FLAI_AGENT` and `FLAI_SESSION` stay the request's.
- The test that reproduces I-0224, in `flai/internal/pilot/run_test.go`: the request `scripts/agent-container.sh S-0507` wrote before this story, with `Name` `agent-S-0507` and no probe, starts claude with the story prompt, which names `story_start`; the same request with `probe` true and `Name` `probe-S-0507` starts it with the probe prompt, with `FLAI_AGENT=probe-S-0507` and no `FLAI_STORY` in the child's environment.
- In `flai/internal/pilot/request_test.go`: `Read` takes `probe`, and `Check` refuses a probe with a `Role` or an `Item`.

## Done when

- A probe start request reaches claude with the probe prompt and without `FLAI_STORY`, and a probe naming a role or an item is refused.
- `flai test flai/internal/pilot/request.go flai/internal/pilot/run.go flai/internal/pilot/request_test.go flai/internal/pilot/run_test.go` passes.

## Notes
