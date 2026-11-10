---
id: T-3042
type: task
nature: remediation
title: A serve test reproduces the EOF import.run reads once the client's ping goes unanswered by the test dashboard
status: backlog
parent: S-0703
owner: alex
created: 2026-11-10T10:34:32Z
updated: 2026-11-10T10:34:32Z
transitions: []
stream: S-0703
tags: [flai, tests]
touches: [flai/internal/serve/candidates_test.go]
---
# T-3042 A serve test reproduces the EOF import.run reads once the client's ping goes unanswered by the test dashboard

## Work

Reproduce I-0267 without host load. In `flai/internal/serve/candidates_test.go`, beside `TestARepositoryIsOfferedImportedAndThenServed`, add a test whose `NewClient` sets a short `PingEvery`, such as 50 ms, and that waits longer than two pings after the candidate's connection is proven before it asks `import.run`. Today it should read "failed to read frame header: EOF", because `channeltest`'s dashboard reads only inside `Ask` and `Next`, so it answers no ping.

Also confirm that S-0560's drain (6920a061c) closed the halt path: the existing test passes on main. Say in the narrative's `## Decisions` which path caused the 13346d1 instance, the ping or the halt. If the ping does not reproduce the EOF, ask on a thread on S-0703 before the next task starts: its fix assumes the ping does.

Waits on nothing: it is the first layer.

## Done when

- The new test fails on main with the EOF, and the narrative's `## Decisions` names the cause it shows.

## Notes
