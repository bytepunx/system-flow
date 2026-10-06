---
id: T-1104
type: task
nature: improvement
title: "An end-to-end test: an agent that wait_for_events ends on a question is started again by flai serve on the answer, and its first inbox holds the answer"
status: backlog
parent: S-0272
owner: alex
created: 2026-10-06T22:53:22Z
updated: 2026-10-06T22:53:30Z
transitions: []
stream: S-0272
tags: [cli, mcp]
touches: [flai/internal/serve/restart_test.go]
after: [T-1087]
---
# T-1104 An end-to-end test: an agent that wait_for_events ends on a question is started again by flai serve on the answer, and its first inbox holds the answer

## Work

Criterion 2. `TestAnAgentThatEndedAskingIsStartedAgainWhenAnswered`, in `flai/internal/serve/agents_test.go` (about line 739), already proves serve's half, with a fake agent command. A run that ends with its thread open is `OutcomeAsked`, and it is started again in its session with `answered` set once the thread is answered. Its MCP half is not tested.

Write a new test in `flai/internal/serve/restart_test.go`, on the same agent lab (`newAgentLab`), that drives the agent's side through the MCP server in process. The `serve` tests may import `mcpserver`, which does not import `serve`. It uses the run's agent name and `FLAI_STORY`:

1. The agent opens a thread on its story with `thread_open`.
2. `wait_for_events` answers `end: true`.
3. The fake agent ends, and the run is `OutcomeAsked`.
4. Nothing is started again before the answer.
5. The operator replies.
6. The agent is started again, as the same agent in its session.
7. `inbox` for that agent name lists the thread, awaiting it, with the answer as its last entry.

If the test shows a gap in serve, fix it in `agents.go`. Add that file to this task's touches with `flai touches` first, and say so in the narrative.

It waits for T-1087, which gives `wait_for_events` its `end` answer.

## Done when

- The test passes with `scripts/flai-test.sh` and fails if `wait_for_events` does not end, or if the agent is not started again with the answer in its first `inbox`.

## Notes

Drafted by the planner. `flai/internal/serve/restart.go` is the operator's restart from the dashboard, which this story leaves alone. The restart on an answer is in `agents.go`.
