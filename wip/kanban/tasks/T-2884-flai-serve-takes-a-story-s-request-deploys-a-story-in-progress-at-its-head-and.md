---
id: T-2884
type: task
nature: feature
title: flai serve takes a story's request, deploys a story in progress at its head, and keeps its namespace until acceptance, cancellation, or a teardown asked for
status: backlog
parent: S-0660
owner: alex
created: 2026-11-10T08:27:42Z
updated: 2026-11-10T08:27:42Z
transitions: []
stream: S-0660
tags: [cli]
touches: [flai/internal/serve/kluster.go, flai/internal/serve/kluster_test.go]
after: [T-2876, T-2879]
---
# T-2884 flai serve takes a story's request, deploys a story in progress at its head, and keeps its namespace until acceptance, cancellation, or a teardown asked for

## Work

In `flai/internal/serve/kluster.go`, the kluster pass takes each request T-2876 defines and answers it:

- It refuses, with the reason and what to do, a story not in progress or in review, a worktree with uncommitted changes, the `kluster` host action off, no `deploy` in the block, and a deploy that would make more story namespaces live than T-2879's limit allows. A story whose namespace already lives does not count against itself.
- It queues the deploy behind the other jobs, in the story's worktree at its branch's head, as `lookReview` does for a story in review. `SystemKluster.Deploy` stops hard-coding `workitem.Review`, and `run` stops refusing a story in progress.
- A teardown asked for runs `teardown` and deletes the namespace.
- `lookLeft` no longer tears down a story in progress. A namespace goes only when the story is accepted or cancelled, or when its agent or the operator asks.
- `lookReview` does not deploy a story whose record holds its head already, and deploys it in place at a new head.
- A failed job someone asked for skips `failed` and `report`, so it opens no thread. The answer gives the step, the exit code, and the log.
- The journal entry and the namespace's record name who asked.

It waits for T-2876, whose request and answer it reads, and for T-2879, whose limit it enforces.

## Done when

- `flai/internal/serve/kluster_test.go` covers a story in progress deployed on request at its head; each refusal; a namespace kept while the story is in progress and torn down on acceptance, on cancellation, and on a teardown asked for; review at the deployed head with no new deploy, and at a new head deployed in place; a failed asked-for deploy with no thread; and the journal and the record naming who asked.
- `go test ./internal/serve/...` passes.

## Notes
