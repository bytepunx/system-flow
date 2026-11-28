---
id: T-3657
type: task
nature: remediation
title: flai serve's kluster request pass takes requests from each story's agent cache folder and answers there
status: backlog
parent: S-0849
owner: alex
created: 2026-11-28T00:04:55Z
updated: 2026-11-28T00:04:55Z
transitions: []
stream: S-0849
tags: [flai, kluster, serve]
touches: [flai/internal/serve/kluster_requests.go, flai/internal/serve/kluster_requests_test.go]
after: [T-3656]
---
# T-3657 flai serve's kluster request pass takes requests from each story's agent cache folder and answers there

## Work

In `flai/internal/serve/kluster_requests.go`, make `takeRequests` look in each story's container requests folder as well as the main checkout's.

- List the `S-nnnn` folders under `<root>/.flai-cache/agents/` and take from each one's `AgentRequestsDir`. Pass over role folders, such as `orchestrator`, and anything that is not a story's ID.
- Refuse a request whose story is not the folder's story, and queue nothing. Write the refusal in that folder.
- Write each answer in the folder its request came from.
- Keep the job state in the main checkout's `requests/` as now, since the review pass and `flai kluster deployments` read it there. Also write it, in `startAsked`, `endAsked`, and the queueing in `answerRequest`, to the story's agent folder when the job was asked from there. Record on the job where it was asked from, or write it to both whenever the story's agent folder exists.
- Make `sweepAsked` end interrupted the jobs an earlier flai serve left in either place, with one event per job.
- Tests in `kluster_requests_test.go`: a request in a story's agent folder is taken, answered there, and queues the job; a request for S-0002 in S-0001's folder is refused there; a role's folder is passed over; the job's state reaches the agent folder at its start and its end.

It waits for T-3656, whose folder-based request functions it calls. It shares no path with T-3658, so the two run together.

## Done when

- A request written under `.flai-cache/agents/S-nnnn/kluster-requests/` is taken and answered there by the pass, and its job's state is written there as well as in the main checkout's `requests/`.
- A request for another story in a story's folder is refused, and nothing is queued.
- The existing tests pass, the new ones pass, and `flai test` on the two paths passes.

## Notes
