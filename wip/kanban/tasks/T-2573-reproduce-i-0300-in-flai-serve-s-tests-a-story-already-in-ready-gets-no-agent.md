---
id: T-2573
type: task
nature: remediation
title: "Reproduce I-0300 in flai serve's tests: a story already in ready gets no agent after the agent action is turned on"
status: backlog
parent: S-0589
owner: alex
created: 2026-11-07T10:05:50Z
updated: 2026-11-07T10:05:50Z
transitions: []
stream: S-0589
tags: [flai, serve, test]
touches: [flai/internal/serve/agents_test.go, flai/internal/serve/serve_test.go]
---
# T-2573 Reproduce I-0300 in flai serve's tests: a story already in ready gets no agent after the agent action is turned on

## Work

Write the tests that reproduce I-0300 and name its cause before anything is fixed, as the story's goal asks. The issue's own guess, that flai serve starts an agent only when a story enters ready, does not match the code: `(*launcher).look` in `flai/internal/serve/agents.go` runs on every change under `kanban/` and `threads/` and every minute, and reads the `agent` action on each look. So the cause is elsewhere. Test each candidate:

- In `flai/internal/serve/agents_test.go`, beside `TestNothingIsStartedWhenItShouldNotBe`, add `TestAStoryAlreadyReadyGetsItsAgentWhenTheActionIsTurnedOn`: with `newAgentLab`, a story in ready and `lab.cfg.Enabled` false, look; set `Enabled` true with a command set, look once, and assert one run. If it passes, the launcher alone is not the cause; say so in the test's comment.
- In `flai/internal/serve/serve_test.go`, run serve on a project with a story already in ready and the `agent` action off, then turn the action on as `flai serve enable agent` does, through the host configuration, and assert the agent starts within the next look. Then the same with the orchestrator's or a role's wake firing more often than once a minute: the project loop's `case <-time.After(time.Minute)` in `flai/internal/serve/serve.go` is made again on every pass, so frequent wakes may starve the minute look, and turning the action on writes no path the watcher sees.
- A project the board's Import adds while serve runs, as S-0564's flai-scratch-import was, is served from then on: check that its first look, and the minute look after it, run with its launcher's configuration.

Record in the story's Notes, under a `### Cause` heading, which test fails on the unfixed code and why. The tests stay failing until the next task fixes the cause; run them with `flai test` on these two files only.

## Done when

- At least one test in `agents_test.go` or `serve_test.go` fails on the unfixed code for the cause I-0300 describes.
- The story's Notes name that cause, with the file and function it lies in.

## Notes
