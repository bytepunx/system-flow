---
id: T-3287
type: task
nature: feature
title: A docker step of a pull or a build streams its output, is stopped after the stall limit, and is tried again on a network failure
status: in-progress
parent: S-0740
owner: alex
created: 2026-11-13T11:43:30Z
updated: 2026-11-13T11:43:55Z
transitions:
  - to: ready
    at: 2026-11-13T11:43:55Z
    by: agent-S-0740
  - to: in-progress
    at: 2026-11-13T11:43:55Z
    by: agent-S-0740
stream: S-0740
tags: []
touches: [flai/internal/agentimage/run.go, flai/internal/agentimage/run_test.go]
usage:
  source: log
  seconds: 36
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 3
      output: 36
      cache_read: 416404
      cache_write: 1911
      cost: 0.1561
---
# T-3287 A docker step of a pull or a build streams its output, is stopped after the stall limit, and is tried again on a network failure

## Work

`flai/internal/agentimage/run.go`: a `Runner` that runs one docker step (a pull, a tag, a build) through an `Exec` adapter, writing its output to a writer as it comes, stopping it when it writes nothing for the stall limit (ADR-0295 §5, 10 minutes by default), and running it again after 30 seconds and then 2 minutes when its last output names a network failure (a connection reset, a timeout, a DNS failure, a TLS handshake failure, a registry 5xx). `SystemExec` is the `os/exec` adapter. Waits for nothing: T-3169 and T-3170 run their docker steps through it.

## Done when

- A test over a fake `Exec` shows a silent step stopped and reported with its name and the limit, and not tried again.
- A test shows a step failing on each network failure tried again up to twice more, with the delays, and a step failing otherwise failing at once.
- A test shows the output reaching the writer as the step writes it.

## Notes
