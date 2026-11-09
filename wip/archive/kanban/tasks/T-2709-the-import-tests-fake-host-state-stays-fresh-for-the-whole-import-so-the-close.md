---
id: T-2709
type: task
nature: feature
title: The import tests' fake host state stays fresh for the whole import, so the close-out's integration tier does not fail them under load (I-0309)
status: done
parent: S-0606
owner: alex
created: 2026-11-09T08:25:03Z
updated: 2026-11-09T08:25:28Z
transitions:
  - to: ready
    at: 2026-11-09T08:25:12Z
    by: system-flow
  - to: in-progress
    at: 2026-11-09T08:25:13Z
    by: system-flow
  - to: done
    at: 2026-11-09T08:25:28Z
    by: agent-S-0606
stream: S-0606
tags: []
touches: [flai/cmd/import_serve_test.go]
usage:
  source: log
  seconds: 15
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 6
      output: 1448
      cache_read: 1277523
      cache_write: 34367
      cost: 0.4794
---

# T-2709 The import tests' fake host state stays fresh for the whole import, so the close-out's integration tier does not fail them under load (I-0309)

## Work

S-0606's close-out (job S-0606-20261109T081033Z) stopped at the integration tier on `TestImportWithAHostRegistersTheProject` and `TestImportWithServeStoppedSaysToStartIt` in `flai/cmd/import_serve_test.go`, at a load average of about 100. S-0606 does not reach that code. I-0309 found the cause: `fakeRunning` stamps the fake host's `state.json` with the current time, and `host.ReadStatus` counts a host gone once its status is older than `host.StaleAfter`, 10 s. An import that runs longer than that under load then reads no host. Stamp the fake state an hour ahead, so it stays fresh for the whole import. `ReadStatus` rejects only a status that is too old, never one from the future.

Waits for no open task. It blocks S-0606's close-out, and it is a five-minute fix, so it is a task on S-0606 (`work-management.md`).

## Done when

- [ ] `fakeRunning` stamps its state ahead, and the import tests pass under `flai test`.

## Notes
