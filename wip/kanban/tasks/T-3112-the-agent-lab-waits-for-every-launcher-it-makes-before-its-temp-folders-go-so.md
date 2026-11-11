---
id: T-3112
type: task
nature: remediation
title: The agent lab waits for every launcher it makes before its temp folders go, so the restart test's fresh launcher no longer writes into a folder being removed
status: backlog
parent: S-0715
owner: alex
created: 2026-11-11T09:04:38Z
updated: 2026-11-11T09:04:38Z
transitions: []
stream: S-0715
tags: [flai, serve, test]
touches: [flai/internal/serve/agents_test.go]
---
# T-3112 The agent lab waits for every launcher it makes before its temp folders go, so the restart test's fresh launcher no longer writes into a folder being removed

## Work

Confirm the cause the story's `### Planning` proposes, then fix it in `flai/internal/serve/agents_test.go`. This task waits for no other: it is the first layer.

1. Reproduce the failure before the fix with repeated runs: `GOMAXPROCS=2 go test -race -count=20 -run TestARestartStartsWhatIsReadyAndNothingTwice ./internal/serve` from `flai/`. Write in the task's `## Notes` how many runs failed. If none fail, say so and go on: the race is narrow, and the code shows it.
2. Give the lab a way to make a launcher it keeps, such as `lab.launcher()`, and have `TestARestartStartsWhatIsReadyAndNothingTwice` take `fresh` from it in place of `newLauncher(lab.o, lab.l.entry)`.
3. Make `settle` wait until no launcher the lab keeps has a PID in `waiting`, `lab.l` among them. Lock each launcher's `mu` to read it, as `settle` does for `lab.l` today.
4. When `settle`'s deadline passes, log with `t.Logf` the launcher and the PIDs it still waited for.
5. Add a test that fits: a second launcher from the lab starts an agent, and once `settle` returns, that launcher waits for no PID.

## Done when

- `TestARestartStartsWhatIsReadyAndNothingTwice` takes its second launcher from the lab, and `settle` waits for every launcher the lab keeps.
- The new test passes, and `GOMAXPROCS=2 go test -race -count=20 -run 'TestARestartStartsWhatIsReadyAndNothingTwice|TestNothingIsStartedWhenItShouldNotBe' ./internal/serve` passes 20 of 20. The counts before and after are written in `## Notes`.
- `flai test flai/internal/serve/agents_test.go` passes.

## Notes
