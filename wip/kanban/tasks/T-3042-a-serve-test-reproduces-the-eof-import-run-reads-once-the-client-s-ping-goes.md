---
id: T-3042
type: task
nature: remediation
title: A serve test reproduces the EOF import.run reads once the client's ping goes unanswered by the test dashboard
status: done
parent: S-0703
owner: alex
created: 2026-11-10T10:34:32Z
updated: 2026-11-11T23:17:58Z
transitions:
  - to: ready
    at: 2026-11-11T23:17:53Z
    by: agent-S-0703
  - to: in-progress
    at: 2026-11-11T23:17:53Z
    by: agent-S-0703
  - to: done
    at: 2026-11-11T23:17:58Z
    by: agent-S-0703
stream: S-0703
tags: [flai, tests]
touches: [flai/internal/serve/candidates_test.go]
after: [T-3044]
usage:
  source: log
  seconds: 5
  estimated: true
  priced_by: estimate
  models: []
---
# T-3042 A serve test reproduces the EOF import.run reads once the client's ping goes unanswered by the test dashboard

## Work

Reproduce I-0267 without host load. In `flai/internal/serve/candidates_test.go`, beside `TestARepositoryIsOfferedImportedAndThenServed`, add a test whose `NewClient` sets a short `PingEvery`, such as 50 ms, and that waits longer than two pings after the candidate's connection is proven before it asks `import.run`. Today it should read "failed to read frame header: EOF", because `channeltest`'s dashboard reads only inside `Ask` and `Next`, so it answers no ping.

Also confirm that S-0560's drain (6920a061c) closed the halt path: the existing test passes on main. Say in the narrative's `## Decisions` which path caused the 13346d1 instance, the ping or the halt. If the ping does not reproduce the EOF, ask on a thread on S-0703 before the next task starts: its fix assumes the ping does.

Waits on nothing: it is the first layer.

## Done when

- The new test fails on main with the EOF, and the narrative's `## Decisions` names the cause it shows.

## Notes
