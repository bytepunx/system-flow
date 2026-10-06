---
id: T-0968
type: task
nature: feature
title: flai serve starts the analyzer behind the analyze host action, records its run, and logs the report it wrote with the run's cost
status: done
parent: S-0223
owner: alex
created: 2026-10-05T05:47:21Z
updated: 2026-10-06T20:27:49Z
transitions:
  - to: ready
    at: 2026-10-06T20:15:10Z
    by: agent-S-0223
  - to: in-progress
    at: 2026-10-06T20:15:10Z
    by: agent-S-0223
  - to: review
    at: 2026-10-06T20:27:49Z
    by: agent-S-0223
  - to: done
    at: 2026-10-06T20:27:49Z
    by: agent-S-0223
stream: S-0223
tags: [flai]
touches: [flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/internal/serve/analyze.go, flai/internal/serve/analyze_test.go, flai/internal/serve/agents.go, flai/internal/serve/activity.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go]
after: [T-0949, T-0957]
usage:
  source: log
  seconds: 758
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 159
      output: 815
      cache_read: 10644843
      cache_write: 226274
      cost: 4.7359
---
# T-0968 flai serve starts the analyzer behind the analyze host action, records its run, and logs the report it wrote with the run's cost

## Work

Add the `analyze` host action, off by default, beside `plan` (`ActionAnalyze` in `flai/internal/hostapi/writes.go`, with its description), and the hostapi method `analyze.run`, which takes an optional `focus` (`bottlenecks`, `intent`, or `risk`), refuses another, and starts a run only while the action is on.

In `flai/internal/serve/analyze.go`, as `plan.go` does for the planner: start one analyzer run per project at a time in the main checkout, with `analysis.agent` over the project's agent (`AnalysisAgent()`), the harness request with `Role` `analyze` and the focus, and `FLAI_ROLE=analyze`; refuse a second while one runs, naming it. Record it as other runs are (`flai/internal/serve/agents.go`: the newest analyzer run per project beside `Plans`, with its trigger, start, end, outcome, and log), log the host entry as the planner's run does, and show it in `flai serve`'s status (`flai/cmd/serve_actions.go`).

When the run ends, its entry in `wip/agents/analyzer.md` names the report it wrote, the newest file under `design/analysis/` changed since its start, beside the run's seconds and cost (`flai/internal/serve/activity.go`); a run that wrote none says so.

It waits for T-0949, whose `AnalysisAgent()` it reads, and T-0957, whose prompt it runs.

## Done when

- a test starts a run through `analyze.run` with the action on, and one is refused with it off and with a bad focus
- a second run while one runs is refused, naming the first
- the run is in the serve status, and its activity entry names the report, the seconds, and the cost
- `go test ./internal/hostapi/ ./internal/serve/ ./cmd/` passes

## Notes
