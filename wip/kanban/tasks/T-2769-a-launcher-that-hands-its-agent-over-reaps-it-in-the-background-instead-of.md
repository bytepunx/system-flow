---
id: T-2769
type: task
nature: remediation
title: A launcher that hands its agent over reaps it in the background instead of releasing it
status: done
parent: S-0632
owner: alex
created: 2026-11-09T21:57:24Z
updated: 2026-11-09T22:12:43Z
transitions:
  - to: ready
    at: 2026-11-09T22:02:47Z
    by: agent-S-0632
  - to: in-progress
    at: 2026-11-09T22:02:50Z
    by: agent-S-0632
  - to: done
    at: 2026-11-09T22:12:43Z
    by: agent-S-0632
stream: S-0632
tags: [flai, serve]
touches: [flai/internal/serve/agents.go, flai/internal/serve/handover_test.go]
usage:
  source: log
  seconds: 593
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 20
      output: 90
      cache_read: 1026039
      cache_write: 22983
      cost: 0.3835
---
# T-2769 A launcher that hands its agent over reaps it in the background instead of releasing it

## Work

In `launcher.await` (`flai/internal/serve/agents.go`), the `handOver` branch closes the output and calls `cmd.Process.Release()`. A long-lived parent, `flai mcp` over stdio, then never reaps the child, which stays a zombie (I-0299). Replace the release with a goroutine that calls `cmd.Wait()` and does nothing else: no `end`, no change to `waiting`, no record. The serving flai still settles the run in `settleOrphans`. A short-lived parent such as `flai analyze` exits as before, and the goroutine goes with it. Update the comments on `handOver` and `await` to say this.

Add `flai/internal/serve/handover_test.go`, a new file kept apart from `agents_test.go`. It starts a short-lived child through a launcher with `handOver` set, in this test process, which stays alive as `flai mcp` does. It then shows that the child leaves the process table once it exits: on Linux `/proc/<pid>` is gone, or its state is not `Z`, within a bounded wait. Skip it where `/proc` is missing.

Waits for no task: it shares no file with the `Owns` task, and the two run together.

## Done when

- `await` never calls `Release` for a handed-over process, and waits for it in the background.
- The new test fails on the old `Release` branch and passes with the fix.
- `flai test flai/internal/serve/agents.go flai/internal/serve/handover_test.go` passes.

## Notes
