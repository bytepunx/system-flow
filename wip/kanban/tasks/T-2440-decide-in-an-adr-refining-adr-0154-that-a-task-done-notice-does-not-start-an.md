---
id: T-2440
type: task
nature: improvement
title: Decide in an ADR refining ADR-0154 that a task-done notice does not start an agent paused on its own question, and say so in the design
status: backlog
parent: S-0545
owner: alex
created: 2026-11-02T20:05:46Z
updated: 2026-11-02T20:05:46Z
transitions: []
stream: S-0545
tags: [flai]
touches: [design/adrs/README.md, design/system/workflow.md, design/system/agent-coordination.md]
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
          input: 1
          output: 3
          cache_read: 30803
          cache_write: 2234
          cost: 0.0078
---
# T-2440 Decide in an ADR refining ADR-0154 that a task-done notice does not start an agent paused on its own question, and say so in the design

## Work

Propose the fix from I-0232's instances, as the story asks, and record it before building it. Write a new ADR with `flai adr new`, refining ADR-0154 and ADR-0205. Base it on what `resume`, `answered`, and `messaged` in `flai/internal/serve/agents.go` do today: `messaged` starts an agent that ended asking on any message to its story, even while its question to the operator is unanswered, and the `blockAsked` interval is not a block in ADR-0154 §3's sense.

The planner's proposal, to accept or change with a reason: a task-done notice is marked as a notice in its entry, as a settled reply is marked; while a run's own question to the operator is unanswered, flai serve does not start it again for a conversation whose entries to its story since it ended are only notices; the answer starts it, and its first `inbox` lists them. A message that asks something, a share, and flai's overlap question about a held ready story (ADR-0134) still start it. Say why ADR-0154's rejected alternative, holding a notice to a story with no commits, does not apply to a paused story, and name the alternatives considered: skipping only when the story's branch changes none of the notice's paths, and having flai settle the notice itself.

Update `design/system/workflow.md` § Transitions and who makes them ("An answered agent") and the messages part of `design/system/agent-coordination.md` to the decision. The ADR file is new, so this task's touches name only the index flai generates; `flai task done` widens them with the ADR's path.

Waits for nothing: it is the first layer, and the code tasks build what it decides.

## Done when

- An accepted ADR refining ADR-0154 states the rule, its scope, and the alternatives considered.
- `design/system/workflow.md` and `design/system/agent-coordination.md` state the rule and link the ADR.
- `design/adrs/README.md` lists the ADR, and `flai check --strict` passes.

## Notes
