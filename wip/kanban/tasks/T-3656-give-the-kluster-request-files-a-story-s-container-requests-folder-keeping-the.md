---
id: T-3656
type: task
nature: remediation
title: Give the kluster request files a story's container requests folder, keeping the main checkout's for the host
status: backlog
parent: S-0849
owner: alex
created: 2026-11-28T00:04:46Z
updated: 2026-11-28T00:04:46Z
transitions: []
stream: S-0849
tags: [flai, kluster]
touches: [flai/internal/kluster/request.go, flai/internal/kluster/request_test.go]
after: [T-3655]
---
# T-3656 Give the kluster request files a story's container requests folder, keeping the main checkout's for the host

## Work

In `flai/internal/kluster/request.go`, let each request function work on a requests folder it is given, so the asker and flai serve can use the story's container folder as well as the main checkout's.

- Add `AgentRequestsDir(root, story)`, `<root>/.flai-cache/agents/S-nnnn/kluster-requests`, beside `AgentCopyDir` and outside it.
- Give `Ask`, `Pending`, `Take`, `WriteAnswer`, `ReadAskedJob`, `WriteAskedJob`, `Await`, and `HasAskedJob` a form that takes the requests folder in place of root. Keep the root-based names as thin wrappers over `RequestsDir(root)`, so `flai/internal/mcpserver/kluster.go` and the other callers compile unchanged.
- The folder's mode follows ADR-0319's copy: 0700 for folders and 0600 for files in the agent folder, since only that story's container reads it. The main checkout's folder keeps its modes.
- Tests in `request_test.go`: an `Ask` against an agent folder is taken and answered by `Take` and `WriteAnswer` on that folder; `Pending` on it lists only its requests; `AgentRequestsDir` lies outside `AgentCopyDir`.

It waits for T-3655, whose ADR fixes the folder's name and place.

## Done when

- `AgentRequestsDir` exists, and every request function has a form that takes a requests folder.
- The root-based functions behave as before, and their existing tests pass unchanged.
- The new tests pass, and `flai test` on the two paths passes.

## Notes
