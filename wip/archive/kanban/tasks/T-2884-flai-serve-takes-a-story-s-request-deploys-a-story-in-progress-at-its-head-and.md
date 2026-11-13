---
id: T-2884
type: task
nature: feature
title: flai serve takes a story's request, deploys a story in progress at its head, and keeps its namespace until acceptance, cancellation, or a teardown asked for
status: done
parent: S-0660
owner: alex
created: 2026-11-10T08:27:42Z
updated: 2026-11-13T08:28:26Z
transitions:
  - to: ready
    at: 2026-11-13T08:04:06Z
    by: agent-S-0660
  - to: in-progress
    at: 2026-11-13T08:04:06Z
    by: agent-S-0660
  - to: done
    at: 2026-11-13T08:28:26Z
    by: agent-S-0660
stream: S-0660
tags: [cli]
touches: [flai/internal/serve/kluster.go, flai/internal/serve/kluster_test.go, flai/internal/serve/kluster_requests.go, flai/internal/serve/kluster_requests_test.go]
after: [T-2876, T-2879]
usage:
  source: log
  seconds: 1460
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 116
      output: 55515
      cache_read: 9065523
      cache_write: 284375
      cost: 4.5667
---
# T-2884 flai serve takes a story's request, deploys a story in progress at its head, and keeps its namespace until acceptance, cancellation, or a teardown asked for

## Work

ADR-0292 §1, §3, §4, and §5, in flai serve's kluster pass. Put the request handling in a new `flai/internal/serve/kluster_requests.go` and change the looks in `kluster.go`:

- **Take requests every second**, beside the pass's minute looks, while the pass runs for the project: each `S-nnnn.<kind>.request.json` T-2876 defines is taken and answered at once.
- **Refuse**, with the reason and what to do, and queue nothing: the `kluster` host action off or no `kluster:` block (naming `flai serve enable kluster` as the operator's); no `deploy`, for a deploy (deploys not set up yet); a story not in progress or in review; a story with no worktree; uncommitted changes to tracked files in the worktree, naming them; a teardown of a story with no record; and a deploy for a story with no namespace while `kluster.agent_namespaces` stories in progress hold one, naming the limit and each holder with its last deploy's commit and time. A story whose namespace is there never counts against itself, and stories in review never count.
- **A story in review** is answered with the review's job when one waits or runs, or with its record when it passed at its head, and queues nothing more.
- **Queue** a deploy at the worktree's head, or a teardown, in the pass's one queue, answering the job's id and the jobs ahead. While a job of the story waits or runs, answer that job; a deploy at a newer head takes the waiting one's place. An asked-for job always runs: the tried-once rule does not hold it back. Write the job's state at queue, start, and end, and at its end append the `kluster` event with `inbox.AppendEvent`. A job waiting or running when flai serve stopped ends `interrupted` at the next start, with its event.
- **Deploy a story in progress**: `run` and `deploy` stop requiring review for an asked-for job; the deploy builds every component and records who asked.
- **No thread for an asked-for failure**: an asked-for deploy or teardown that fails opens no thread, moves nothing, and starts no agent; its job state and event carry the step, the exit code, and the log.
- **Keep the namespace**: `lookLeft` tears down no namespace whose story is in progress, whoever moved it there; ready, backlog, cancelled, and accepted still tear it down. `lookReview` queues no deploy for a story whose record passed at its branch's head (with `checks` on, its checks too), and deploys in place at a new head.
- **The journal** entries of an asked-for job name who asked.

It waits for T-2876, whose files and event it writes, and T-2879, whose limit it reads.

## Done when

- Tests cover a story in progress deployed on request at its head; each refusal; a request while a job waits answered with it; the limit refusing a third story and never a story already holding one; a namespace kept while its story is in progress and torn down on a move to ready, on cancellation, on acceptance, and on a teardown asked for; review at the deployed head with no new deploy, and at a new head deployed in place; a failed asked-for deploy with no thread and its event; and the journal and the record naming who asked.
- `flai test` on the changed paths passes.

## Notes
