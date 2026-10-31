---
id: T-2403
type: task
nature: feature
title: flai serve takes a close-out request, runs the close-out in the story's worktree as a job, appends a verified event with its verdict, and journals its start, end, duration, and verdict
status: backlog
parent: S-0546
owner: alex
created: 2026-10-31T07:44:01Z
updated: 2026-10-31T07:45:17Z
transitions: []
stream: S-0546
tags: [flai]
touches: [flai/internal/serve/verifyjob.go, flai/internal/serve/verifyjob_test.go, flai/internal/serve/serve.go, flai/internal/inbox/events.go, flai/internal/inbox/events_test.go, flai/internal/hostapi/writes.go, flai/cmd/serve_actions.go]
after: [T-2396]
---
# T-2403 flai serve takes a close-out request, runs the close-out in the story's worktree as a job, appends a verified event with its verdict, and journals its start, end, duration, and verdict

## Work

Add the job in a new `flai/internal/serve/verifyjob.go`, started per served project from `serve.go` as `newMergePass` is. It takes the close-out requests T-2396's ADR names, one job at a time per story, and gives each a job id. It runs the close-out in the story's worktree with the agent's commit options. When the job ends, it appends a `verified` event (a kind and its description in `flai/internal/inbox/events.go`) carrying the story, the job id, the commit, and the verdict, and answers the request. A request for a story whose job is running answers with that job's id and starts no second job. A job serve was running when it stopped is reported failed on the next start, not left pending.

Journal each job's start and its end, with the duration and the verdict, through `Host.Record`. Add the fields the ADR chose to `hostapi.Entry` (`flai/internal/hostapi/writes.go`), and have `flai serve journal` show them (`flai/cmd/serve_actions.go`). Export the function that writes a request, so that T-2407 calls it from the CLI and the MCP server.

Waits for T-2396, which settles the hand-off, the event, and the journal fields. Shares no path with T-2400, so the two run together.

## Done when

- [ ] A test in `flai/internal/serve/verifyjob_test.go` writes a request, sees the job run a stub close-out in a worktree, and reads a `verified` event with the verdict and two journal entries, the end with its duration and verdict
- [ ] A second request while the job runs gets the same id and starts nothing
- [ ] A test in `flai/internal/inbox/events_test.go` describes the `verified` event
- [ ] `go test -race ./internal/serve ./internal/inbox ./internal/hostapi` passes
- [ ] Story criteria 1 (serve's half) and 4 (the journal's half)

## Notes
