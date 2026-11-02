---
id: T-2443
type: task
nature: improvement
title: flai serve does not start an agent paused on its own question for a conversation that holds only task-done notices, with a test reproducing I-0232
status: backlog
parent: S-0545
owner: alex
created: 2026-11-02T20:06:36Z
updated: 2026-11-02T20:06:36Z
transitions: []
stream: S-0545
tags: [flai]
touches: [flai/internal/serve/agents.go, flai/internal/serve/agents_test.go]
after: [T-2441]
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 3
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 0
          output: 3
          cache_read: 30802
          cache_write: 2233
          cost: 0.0078
---
# T-2443 flai serve does not start an agent paused on its own question for a conversation that holds only task-done notices, with a test reproducing I-0232

## Work

Build the rule of T-2440's ADR in `flai/internal/serve/agents.go`. `answered` falls through to `messaged`, which returns any conversation awaiting the story's reply to a message written after the run ended, and `resume` then starts the agent although its question to the operator is unanswered. While the run's own question (`run.Thread`) is unanswered, pass over a conversation whose entries to the story since the run ended all carry T-2441's notice mark. A message without the mark, a conversation that closed, and the answer to the question still start it as before. Keep the activity's why true for a paused story with notices waiting, if `waitWhy` or the activity reads them.

Reproduce I-0232 in `flai/internal/serve/agents_test.go`: a story whose agent ended asking on a thread, the `blockAsked` interval open, gets a task-done notice from another story; the look starts nothing and the story stays blocked. Then the thread is answered, and the agent is started once, on the thread, not the conversation. Add the cases that must still wake: a plain message to the paused story, and a notice to a story that ended waiting only on a conversation.

Waits for T-2441, whose mark it reads.

## Done when

- The I-0232 test fails on the code before the change and passes after it.
- The cases that must still wake pass, and the existing `resume`, `answered`, and `messaged` tests pass unchanged or with a stated reason.
- `flai test` passes on the changed paths.

## Notes
