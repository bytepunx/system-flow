---
id: T-1275
type: task
nature: improvement
title: permission_prompt bounds its wait, keeps an unanswered thread open, and takes its answer on the retry
status: in-progress
parent: S-0309
owner: alex
created: 2026-10-07T23:27:25Z
updated: 2026-10-08T05:53:19Z
transitions:
  - to: ready
    at: 2026-10-08T05:53:18Z
    by: agent-S-0309
  - to: in-progress
    at: 2026-10-08T05:53:19Z
    by: agent-S-0309
stream: S-0309
tags: [mcp, permission-prompt]
touches: [flai/internal/mcpserver/permission.go, flai/internal/mcpserver/permission_test.go]
after: [T-1274]
---
# T-1275 permission_prompt bounds its wait, keeps an unanswered thread open, and takes its answer on the retry

## Work

Build in `flai/internal/mcpserver/permission.go` what T-1274's ADR decides. It waits for T-1274 because the ADR settles the bound and the retry's rules.

- `awaitAnswer` stops after the bound as well as on the context or `closing`. The bound is a package variable beside `permissionPoll`, so tests can shorten it.
- `askOperator` leaves the thread open when it stops unanswered, and refuses with a reason naming the thread, that it stays open, that the owner's answer is taken when the agent makes the same write again, and that the agent goes on meanwhile with work that does not need the write.
- Before opening a thread, `askOperator` looks for an open thread on the story that the agent opened for the same tool, path, and input. An answer already on it decides at once and settles the thread; no answer yet holds the call again, bounded as above.
- `permissionPromptDescription` says the same.

Write tests in `flai/internal/mcpserver/permission_test.go` that reproduce I-0103: with a short bound and no answer, the call refuses within the bound, names the thread, and the thread stays open. Then an allow by the owner lets the same request through at once and settles the thread, a refusal refuses the retry, and a request with other input opens a new thread.

## Done when

- The tests above pass with `flai test flai/internal/mcpserver`, and the existing permission tests still pass.
- No request holds a call longer than the bound.

## Notes
