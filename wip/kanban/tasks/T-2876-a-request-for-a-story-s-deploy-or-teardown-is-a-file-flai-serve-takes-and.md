---
id: T-2876
type: task
nature: feature
title: A request for a story's deploy or teardown is a file flai serve takes and answers, and the namespace's record names who asked
status: in-progress
parent: S-0660
owner: alex
created: 2026-11-10T08:26:49Z
updated: 2026-11-13T07:33:53Z
transitions:
  - to: ready
    at: 2026-11-13T07:33:53Z
    by: agent-S-0660
  - to: in-progress
    at: 2026-11-13T07:33:53Z
    by: agent-S-0660
stream: S-0660
tags: [cli]
touches: [flai/internal/kluster/request.go, flai/internal/kluster/request_test.go, flai/internal/kluster/deploy.go, flai/internal/kluster/deploy_test.go, flai/internal/inbox/events.go, flai/internal/inbox/events_test.go, flai/internal/inbox/inbox.go, flai/internal/inbox/inbox_test.go]
---
# T-2876 A request for a story's deploy or teardown is a file flai serve takes and answers, and the namespace's record names who asked

## Work

ADR-0292 §1, §5, and §6 name the request, its files, its wait, and the `kluster` event. Add them in a new `flai/internal/kluster/request.go`, modelled on `flai/internal/closeout/closeout.go`, which `flai verify --job` and flai serve's verify pass share:

- A request is `S-nnnn.<kind>.request.json` under `.flai-cache/kluster/requests/`, with the story, the kind (`deploy` or `teardown`; `checks` is S-0661's), who asked, and when. flai serve takes it by renaming it to `.taken.json` and answers in `.answer.json`.
- `Ask` writes it, waits at most ten seconds for flai serve to take it, withdrawing it and wrapping an `ErrServeNotTaking` otherwise, then at most ten seconds for the answer. The answer is the job's id, its kind, the number of jobs ahead of it, whether this request queued it or found the story's job already waiting or running, or, for a story in review, the review's job or its passing record; or a refusal with its reason and what to do.
- The job's state is `S-nnnn.job.json` in the same folder, one per story: its id, kind, story, who asked, queued, started, and ended times, and once ended its commit, outcome (`passed`, `failed`, `refused`, or `interrupted`), failed step, exit code, each step's name, outcome, exit code, and log, and the log's path. `WriteJob` and `ReadJob` keep it.
- `Await` waits up to a given duration for the job of an answer to end, and says whether it ended, so the CLI and the MCP tool can wait for the outcome (criterion 2) and otherwise say it is still queued or running and how to read it later.

In `flai/internal/kluster/deploy.go`, `Record` gains `asked_by` and the lock's `Job` gains `by`, so that the namespace's record, `flai kluster deployments`, and the dashboard's running job say who asked. A record written before reads with none.

In `flai/internal/inbox`, the `kluster` event (ADR-0292 §5): `KlusterEvent` with the story as its ID and the fields `job`, `kind`, `commit`, `outcome`, `failed_step`, `exit_code`, and `log`, described in one sentence, and reported, as `verified` is, to the agent of its story and to readers with no story of their own.

It waits for nothing: the first layer.

## Done when

- Tests of `Ask` cover a request taken and answered, one not taken within its wait, a refusal, and a job found already queued; tests of `Await` one that ends within its wait and one that outlasts it.
- A record and a lock job written with who asked read back with it, and ones written before read with none.
- The inbox reports a `kluster` event to its own story's agent and a reader with no story, not to another story's agent, with its summary.
- `flai test` on the changed paths passes.

## Notes
