---
id: T-2863
type: task
nature: feature
title: flai serve's kluster pass runs the checks a story's agent asks for against its namespace at the head, after the job queue, records who asked, and opens no thread on a failure
status: backlog
parent: S-0661
owner: alex
created: 2026-11-10T08:23:22Z
updated: 2026-11-10T08:25:55Z
transitions: []
stream: S-0661
tags: [cli]
touches: [flai/internal/serve/kluster.go, flai/internal/serve/kluster_test.go, flai/internal/serve/checks.go, flai/internal/serve/checks_test.go, flai/internal/kluster/request.go, flai/internal/kluster/request_test.go]
---
# T-2863 flai serve's kluster pass runs the checks a story's agent asks for against its namespace at the head, after the job queue, records who asked, and opens no thread on a failure

## Work

Read first the ADR S-0659 accepted, on the checks a story's agent asks for and what is recorded of them, and the request S-0660 built for a deploy the agent asks for, in `flai/internal/kluster/request.go`. The checks request is a second kind of that request: same hand-off, same refusals, same answer.

In `flai/internal/kluster/request.go`, add the checks kind to the request and its answer: each check's name, outcome, exit, and log, the run's log, the commit, and the namespace.

In `flai/internal/serve/kluster.go`:

- Give `KlusterEnv` a reason when it gives no variables: the action off, no `kluster:` block, no record, a record that did not pass, no branch, or a record at another commit than the branch's head. The checks request is refused on each one, saying to deploy first at the head, with the command S-0660 built. The review pass keeps passing over them in silence.
- Take a checks request from a story's agent for its own story, in progress or in review, as a job on the kluster pass's queue (`enqueue`, `next`, `drain`). It runs one at a time per project, after the jobs queued before it, as the review checks do.
- Run it with `RunChecksIn`, the namespace's variables, and the story's check container, as `checks` does for a story in review.
- On a failure, do not call `p.report`. The answer holds the failed check, its exit, and the log. A review checks run still opens its thread as before.

In `flai/internal/serve/checks.go`, add to `ChecksRun` who asked (the agent's name, or `flai serve` for the review pass), the commit it ran at, and the namespace. `checks/<story>.json` then holds the run where S-0662 reads it. Journal each run with who asked.

Tests in `kluster_test.go`, `checks_test.go`, and `request_test.go` fake `kluster` and the checks, as S-0579's do. They cover a passing run, a failing run with no thread, each refusal, and a request queued behind a deploy. It waits for no task: it is the first layer.

## Done when

- A checks request for a story deployed at its head runs on the kluster pass after the jobs ahead of it, with the namespace's variables, and records its commit, outcome, and who asked in `checks/<story>.json` and the host journal.
- A request for a story not deployed at its head is refused, saying to deploy first; a failed run asked for by the agent opens no thread.
- `go test ./internal/serve/... ./internal/kluster/...` passes.

## Notes
