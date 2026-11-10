---
id: T-3004
type: task
nature: improvement
title: "agent.status times serve.Activity's parts as phases of the request: the runs, the holds with their trial merges, and elsewhere"
status: backlog
parent: S-0678
owner: alex
created: 2026-11-10T10:31:16Z
updated: 2026-11-10T10:31:16Z
transitions: []
stream: S-0678
tags: [cli]
touches: [flai/internal/serve/agents_activity.go, flai/internal/serve/activity_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go]
---
# T-3004 agent.status times serve.Activity's parts as phases of the request: the runs, the holds with their trial merges, and elsewhere

## Work

I-0313 says this repository's `agent.status` still takes about 1.1 s with 13 projects at once, and that alone its `serve.Activity` takes 540 to 600 ms. S-0602 has since kept the threads and the conversations parsed, so what is left is the holds with their trial merges and `elsewhere`, but the serve log shows only `agent.activity` around all of it.

Give `serve.Activity` (`flai/internal/serve/agents_activity.go`) the request's context, and pass it from `Host.Agent` in `flai/cmd/serve_actions.go`. Inside `agent.activity`, record with `perf.Track`:

- `activity.runs`: `runActivity`, the threads read and each live run judged.
- `activity.holds`: `held`, the item list and the holds built.
- `activity.trial`: the time spent in the trial judge, summed over the call, by wrapping `repo.Trial` in a judge that times each call it passes on.
- `activity.elsewhere`: `elsewhere` with its narrative reads.

Change no answer. This is the first layer: it waits for no task, and it shares no path with the trial judge's task, so the two run together.

## Done when

- A test of `Host.Agent` with a request timer, in `flai/cmd/serve_actions_test.go`, shows the four phases inside `agent.activity`.
- A test in `flai/internal/serve/activity_test.go` shows `Activity` answers the same with and without a timer on its context.
- `flai test` passes on the changed files.

## Notes
