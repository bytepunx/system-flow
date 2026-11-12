---
id: T-3177
type: task
nature: feature
title: flai serve runs the image ensure as a host job, one at a time, that held starts wait on and that the journal records
status: backlog
parent: S-0751
owner: alex
created: 2026-11-12T16:01:38Z
updated: 2026-11-12T16:01:38Z
transitions: []
stream: S-0751
tags: [cli]
touches: [flai/internal/serve/imagejob.go, flai/internal/serve/imagejob_test.go]
---
# T-3177 flai serve runs the image ensure as a host job, one at a time, that held starts wait on and that the journal records

## Work

Write `flai/internal/serve/imagejob.go`. An image job runs `agentimage.Ensure` (S-0740, S-0745) for one project on the host. Its rules:

- At most one job runs per host. A request for images a running job covers joins it. A request for others queues behind it.
- Each job has an id, a start and end time, its way (pull or build) per image, its log file, and on failure the failing step and cause.
- The serve journal logs the start and the end with those fields, in the words of the close-out job's entries (`verifyjob.go`, `journal`, line 502).
- Callers wait on a job through a channel or callback that answers its outcome. A cancelled `flai serve` stops the job's docker as a stop does (ADR-0058).
- A failed job is kept, so the next request starts a new one rather than looping on it.

Model it on the close-out job (`verifyjob.go`: `startJob` at line 281, `runJob` at line 356). Test over a fake ensure:

- one job at a time;
- a joined request;
- a queued request;
- a pass;
- a failure;
- a stop.

It is the first layer.

## Done when

- `go test ./internal/serve/ -run ImageJob` passes, under `-race` too.
- A test reads the journal's two entries for a job and finds every field.

## Notes
