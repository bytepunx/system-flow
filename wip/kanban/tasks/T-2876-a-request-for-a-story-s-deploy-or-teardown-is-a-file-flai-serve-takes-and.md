---
id: T-2876
type: task
nature: feature
title: A request for a story's deploy or teardown is a file flai serve takes and answers, and the namespace's record names who asked
status: backlog
parent: S-0660
owner: alex
created: 2026-11-10T08:26:49Z
updated: 2026-11-10T08:26:49Z
transitions: []
stream: S-0660
tags: [cli]
touches: [flai/internal/kluster/request.go, flai/internal/kluster/request_test.go, flai/internal/kluster/deploy.go, flai/internal/kluster/deploy_test.go]
---
# T-2876 A request for a story's deploy or teardown is a file flai serve takes and answers, and the namespace's record names who asked

## Work

Add the request a story's agent makes, in a new `flai/internal/kluster/request.go`, modelled on `flai/internal/closeout/closeout.go`'s `Ask`, the hand-off `flai verify --job` makes. A request is a file under `.flai-cache/kluster/` naming the story, the kind (deploy or teardown), the commit, who asked, and when. `Ask` writes it and waits, first for flai serve to take it, then for the answer, for as long as S-0659's ADR sets. The answer holds the outcome, the commit, each step's result with its exit code, and the log's path, or a refusal with its reason and what to do. A request that outlasts its wait answers that the job is still queued or running, and how to read it later.

In `flai/internal/kluster/deploy.go`, `Record` and the lock's `Job` gain who asked, so that the namespace's record and `flai kluster status` say it.

It waits for nothing: the first layer.

## Done when

- Tests of `Ask` cover a request taken and answered, one not taken within its wait, one taken and still running past its wait, and a refusal, each with what the caller is told.
- A record written for a job someone asked for reads back with who asked, and a record written before this change reads with none.
- `go test ./internal/kluster/...` passes.

## Notes
