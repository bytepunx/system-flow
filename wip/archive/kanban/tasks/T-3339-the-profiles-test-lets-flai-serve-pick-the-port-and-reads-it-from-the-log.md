---
id: T-3339
type: task
nature: remediation
title: The profiles test lets flai serve pick the port and reads it from the log
status: done
parent: S-0779
owner: alex
created: 2026-11-13T21:45:58Z
updated: 2026-11-13T21:51:16Z
transitions:
  - to: ready
    at: 2026-11-13T21:46:19Z
    by: agent-S-0779
  - to: in-progress
    at: 2026-11-13T21:46:20Z
    by: agent-S-0779
  - to: done
    at: 2026-11-13T21:51:16Z
    by: agent-S-0779
stream: S-0779
tags: []
touches: [flai/cmd/serve_pprof_test.go]
usage:
  source: log
  seconds: 296
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 8
      output: 2341
      cache_read: 1013449
      cache_write: 12441
      cost: 0.3491
---

# T-3339 The profiles test lets flai serve pick the port and reads it from the log

## Work

`TestProfilesAreOfferedOnLoopbackOnly` (`flai/cmd/serve_pprof_test.go`) asks the kernel for a port, closes it, and gives it to `servePprof` in `FLAI_PPROF_ADDR`, counting on it staying free until `servePprof` listens: the pattern I-0372 records. Criterion 1 covers it. Give `servePprof` `127.0.0.1:0` instead, and read the port it picked from its `profiles offered` log line.

Waits for nothing in the story: no other task touches the file.

## Done when

- The test gives `servePprof` no port it closed, and still shows the profiles answered on loopback.
- `flai test` on `flai/cmd/serve_pprof_test.go` passes.

## Notes
