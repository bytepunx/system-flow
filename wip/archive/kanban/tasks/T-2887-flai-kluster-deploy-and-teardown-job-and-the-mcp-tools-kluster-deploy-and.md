---
id: T-2887
type: task
nature: feature
title: flai kluster deploy and teardown --job and the MCP tools kluster_deploy and kluster_teardown ask flai serve for the session's own story and print the answer
status: done
parent: S-0660
owner: alex
created: 2026-11-10T08:27:58Z
updated: 2026-11-13T08:28:48Z
transitions:
  - to: ready
    at: 2026-11-13T08:04:06Z
    by: agent-S-0660
  - to: in-progress
    at: 2026-11-13T08:04:06Z
    by: agent-S-0660
  - to: done
    at: 2026-11-13T08:28:48Z
    by: agent-S-0660
stream: S-0660
tags: [cli]
touches: [flai/cmd/kluster.go, flai/cmd/kluster_test.go, flai/internal/mcpserver/kluster.go, flai/internal/mcpserver/kluster_test.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/folder_test.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/server_test.go, flai/internal/kluster/request.go, flai/internal/kluster/request_test.go, docs/users/flai-reference.md, docs/operators/settings.md]
after: [T-2876]
usage:
  source: log
  seconds: 1482
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 137
      output: 65699
      cache_read: 10728552
      cache_write: 336542
      cost: 5.4045
---
# T-2887 flai kluster deploy and teardown --job and the MCP tools kluster_deploy and kluster_teardown ask flai serve for the session's own story and print the answer

## Work

ADR-0292 §1 and §5, criteria 1 and 2:

- `flai kluster deploy S-nnnn --job` and `flai kluster teardown S-nnnn --job` in `flai/cmd/kluster.go` call T-2876's `Ask` in place of running the job. They print the answer: the job's id, its kind, and the jobs ahead, or the refusal with its reason, exit 1. `--wait <duration>` (default 0) then waits with `Await` and prints the outcome, the commit, each step's result with its exit code, and the log's path, exiting 1 on a failure; a job that outlasts the wait is said to be still queued or running, with how to read it later (`flai kluster deployments`, and the `kluster` event to `wait_for_events`). `--json` prints the answer and, with `--wait`, the job. A story other than the session's own (`FLAI_STORY`, else the story in the agent's name) is refused before any request is written; with neither, the story given is used, as the operator's shell has no session story.
- The MCP tools `kluster_deploy` and `kluster_teardown` in a new `flai/internal/mcpserver/kluster.go`, registered in `folder.go` as `verify` is, ask for the session's own story (an optional `story` argument other than it is refused), take an optional `wait_seconds`, at most 45, and answer as the CLI does as JSON. Their descriptions tell the agent to hold `wait_for_events` for the `kluster` event when the job has not ended.
- `wait_for_events` answers no `end` to a story's agent while a kluster job of its story waits or runs, as it does for a close-out job (`server.go`).

Add the tools to `TestToolsAreAdvertised`'s list in `server_test.go`. Run `make flai-reference`, which writes `docs/users/flai-reference.md` and the flag index in `docs/operators/settings.md`, so that `TestReferenceIsCurrent` and `TestReferenceCoversEveryCommand` pass.

It waits for T-2876, whose `Ask` and `Await` it calls. It runs beside T-2884, which touches none of its files.

## Done when

- `flai/cmd/kluster_test.go` and `flai/internal/mcpserver/kluster_test.go` cover the own story answered, another story refused, a refusal from flai serve, flai serve not taking the request, an answer waited for to its end, one that outlasts the wait, and the JSON output.
- A test shows `wait_for_events` not ending while the story's kluster job runs.
- `TestToolsAreAdvertised`, `TestReferenceIsCurrent`, and `TestReferenceCoversEveryCommand` pass under `flai test`.

## Notes
