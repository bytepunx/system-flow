---
id: T-2362
type: task
nature: remediation
title: "The serve tests I-0226 and I-0233 name wait for what they assert, and pass ten times in a row under -race -count=10 with GOMAXPROCS=2"
status: done
parent: S-0548
owner: alex
created: 2026-10-31T07:24:50Z
updated: 2026-10-31T08:14:14Z
transitions:
  - to: ready
    at: 2026-10-31T07:32:08Z
    by: agent-S-0548
  - to: in-progress
    at: 2026-10-31T07:32:09Z
    by: agent-S-0548
  - to: done
    at: 2026-10-31T08:14:14Z
    by: agent-S-0548
stream: S-0548
tags: []
touches: [flai/internal/serve/agents_test.go, flai/internal/serve/serve_test.go, flai/internal/serve/forgemerge_test.go, flai/internal/serve/budget_test.go, flai/internal/serve/mcp_test.go, flai/internal/serve/projects_test.go]
usage:
  source: log
  seconds: 2494
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 107
      output: 467
      cache_read: 4858233
      cache_write: 131766
      cost: 1.7364
---
# T-2362 The serve tests I-0226 and I-0233 name wait for what they assert, and pass ten times in a row under -race -count=10 with GOMAXPROCS=2

## Work

I-0226: `TestAStorySentBackIsStartedAfreshOrBeforeARestart` (`flai/internal/serve/agents_test.go`) failed at "never happened: it is started for the send back" with `waitFor`'s 5-second deadline on a host at load 23.7, and checks `startedCount` right after the run's state changes, before the journal may hold the start; `TestRequestMergeIsAnsweredByFlaiServesMergePass` (`forgemerge_test.go`) now waits for its log line (S-0507's T-2317). I-0233: `TestAChangedFileReachesTheDashboardAsANotification` (`serve_test.go`) sleeps 80 ms for the watcher's first look and failed its temp folder's cleanup; S-0418's fe589e207 made `Run` wait for a project's loop and looks before it returns. Make each wait for the condition it asserts, with a deadline that a loaded host does not reach (a longer `waitFor` deadline costs nothing on a fast host), and replace the sleep with a wait for the watcher's first look.

Story criterion 1 says no test in the tree dials a port nothing listens on: on this WSL2 host a closed local port hangs to the connect timeout. `127.0.0.1:1` appears in `serve_test.go`, `budget_test.go`, `mcp_test.go` (three times), and `projects_test.go`. Replace each one that is dialed with the address of a listener the test opened and closed (`127.0.0.1:0`, its address taken, then closed); leave one that is only stored and never dialed, and say which in the final message. T-2365 does the same in `flai/cmd`. Waits for no other task.

## Done when

- [ ] `GOMAXPROCS=2 go test -race -count=10 -run 'TestAStorySentBackIsStartedAfreshOrBeforeARestart|TestRequestMergeIsAnsweredByFlaiServesMergePass|TestAChangedFileReachesTheDashboardAsANotification' ./internal/serve` passes, recorded in the narrative
- [ ] No serve test dials `127.0.0.1:1`
- [ ] No test asserts less than before
- [ ] Story criteria 1 (for serve) and 7

## Notes
