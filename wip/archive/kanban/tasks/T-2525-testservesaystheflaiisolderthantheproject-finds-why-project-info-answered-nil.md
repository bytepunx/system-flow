---
id: T-2525
type: task
nature: remediation
title: TestServeSaysTheFlaiIsOlderThanTheProject finds why project.info answered nil on a loaded host and waits for an answer that holds the project
status: done
parent: S-0569
owner: alex
created: 2026-11-07T03:47:37Z
updated: 2026-11-07T04:40:37Z
transitions:
  - to: ready
    at: 2026-11-07T04:17:12Z
    by: agent-S-0569
  - to: in-progress
    at: 2026-11-07T04:17:15Z
    by: agent-S-0569
  - to: done
    at: 2026-11-07T04:40:37Z
    by: agent-S-0569
stream: S-0569
tags: [flai, tests]
touches: [flai/internal/serve/outdated_test.go]
usage:
  source: log
  seconds: 1402
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 44
      output: 11630
      cache_read: 5493613
      cache_write: 184767
      cost: 2.3613
---
# T-2525 TestServeSaysTheFlaiIsOlderThanTheProject finds why project.info answered nil on a loaded host and waits for an answer that holds the project

## Work

I-0247's first instance: `TestServeSaysTheFlaiIsOlderThanTheProject` (`flai/internal/serve/outdated_test.go`) failed at line 175, `project.info: <nil>`. The test asks `project.info` once (`conn.Ask`, 5-second context, `flai/internal/channel/channeltest/dashboard.go`) right after the dashboard's first proven connection, with no retry. `project.info` works out `FlaiOutdated` at request time (`flai/internal/hostapi/hostapi.go`), so the value is not stale. An answer with no body points to the project not being registered with the channel yet when the first ask lands.

First reproduce it: run the test with `GOMAXPROCS=2 go test -race -count=20 -run TestServeSaysTheFlaiIsOlderThanTheProject ./internal/serve` beside a CPU load, and find which step answers an empty body. Then make the test wait for what it asserts: ask again until the answer holds the project, within a bound a loaded host does not reach. If the empty answer comes from `hostapi` or `serve`, answering before the project is ready, that is a defect in flai, not in the test. Fix it there, widen the task's touches to the file, and say so in the narrative. The 3-second loops at lines 106, 126, and 177 end silently on timeout: make each fail with what it waited for. Waits for no other task.

## Done when

- [ ] The narrative names the step that answered an empty `project.info`
- [ ] `GOMAXPROCS=2 go test -race -count=10 -run TestServeSaysTheFlaiIsOlderThanTheProject ./internal/serve` passes beside a CPU load, recorded in the narrative
- [ ] The test asserts no less than before

## Notes
