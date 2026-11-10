---
id: T-2841
type: task
nature: remediation
title: The dev server's 404 test allows the fake vite time to start, and a start delay reproduces the missing 404
status: done
parent: S-0652
owner: alex
created: 2026-11-10T07:00:59Z
updated: 2026-11-10T08:09:54Z
transitions:
  - to: ready
    at: 2026-11-10T08:05:40Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T08:05:41Z
    by: system-flow
  - to: done
    at: 2026-11-10T08:09:54Z
    by: agent-S-0652
stream: S-0652
tags: [testing]
touches: [flai/internal/components/devserver_test.go]
usage:
  source: log
  seconds: 253
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 7
      output: 748
      cache_read: 1001537
      cache_write: 7533
      cost: 0.2756
---
# T-2841 The dev server's 404 test allows the fake vite time to start, and a start delay reproduces the missing 404

## Work

`TestStartDevServerSaysWhatAViteThatDoesNotAnswer200Answered` gives `StartDevServer` a 600 ms `Timeout`. Under load the fake vite, the test binary started again by a shell script, is not serving by then, so the deadline error has no "(it answered 404, and the harness answers 200)".

- Reproduce first: run the test with the fake vite's start delay, which S-0647 adds, set longer than 600 ms, and see it fail on the missing 404.
- Give the test a deadline that outlasts the fake's start under load, with a comment saying why. Keep it as short as that allows, since the test waits the whole deadline.
- Keep the pid check S-0647 leaves, the pid of the process flai started.
- Run the test again with the delay, and many times with `-count`, and see it pass.

It waits for no task: it touches only the dev server test file, which no other task of this story touches.

## Done when

- The test fails with the start delay before the change and passes with it after.
- `flai test flai/internal/components/devserver_test.go` passes.

## Notes
