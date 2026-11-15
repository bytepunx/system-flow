---
id: T-3371
type: task
nature: feature
title: An agent's own kluster deploy request and the merge queue's preflight read the workflow, not a state's name
status: done
parent: S-0743
owner: alex
created: 2026-11-15T01:43:42Z
updated: 2026-11-15T01:51:11Z
transitions:
  - to: ready
    at: 2026-11-15T01:44:24Z
    by: agent-S-0743
  - to: in-progress
    at: 2026-11-15T01:44:24Z
    by: agent-S-0743
  - to: done
    at: 2026-11-15T01:50:27Z
    by: agent-S-0743
stream: S-0743
tags: []
touches: [flai/internal/serve/kluster_requests.go, flai/internal/serve/kluster_requests_test.go, flai/internal/mergequeue/preflight.go, design/system/workflow.md]
after: [T-3218, T-3220]
usage:
  source: log
  seconds: 363
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 16
      output: 168
      cache_read: 3012230
      cache_write: 12273
      cost: 1.1401
---
# T-3371 An agent's own kluster deploy request and the merge queue's preflight read the workflow, not a state's name

## Work

- `flai/internal/serve/kluster_requests.go` judges an agent's deploy or teardown request by `in-progress` and `review` by name, counts only `in-progress` stories against `kluster.agent_namespaces`, and opens its repositories with `workitem.Open`. Read the waiting state the acceptance leaves from, the underway kinds, and the active kind instead, and open with `openRepo`, as T-3218 does for the pass.
- `flai/internal/mergequeue/preflight.go` opens the main checkout with `workitem.Open`: give it the queue's repository's workflow, so the preflight's check runs the same board.
- Waits for T-3218 and T-3220, which found both and whose files these sit beside.

## Done when

- Under the renamed test workflow an agent's request from `doing` or `testing` is taken and counted against the limit, and one from `awaiting` is answered as a review's; with the baseline every existing test passes unchanged.

## Notes
