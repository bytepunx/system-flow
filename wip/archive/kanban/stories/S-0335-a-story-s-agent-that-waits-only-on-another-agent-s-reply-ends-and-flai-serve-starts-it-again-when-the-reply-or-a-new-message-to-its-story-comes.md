---
id: S-0335
type: story
nature: improvement
title: A story's agent that waits only on another agent's reply ends, and flai serve starts it again when the reply or a new message to its story comes
status: done
parent: E-0018
owner: alex
created: 2026-10-07T20:11:14Z
updated: 2026-10-07T22:10:05Z
transitions:
  - to: ready
    at: 2026-10-07T21:43:11Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-07T21:43:16Z
    by: agent-S-0335
  - to: review
    at: 2026-10-07T22:09:19Z
    by: agent-S-0335
  - to: done
    at: 2026-10-07T22:10:05Z
    by: orchestrator
tags: [flai, flaiover, template]
topics: [cli, dashboard, conventions, template]
touches: [flai/internal/mcpserver/server.go, flai/internal/mcpserver/server_test.go, flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, flaiover/src/lib/components/StoryAgent.svelte, flaiover/src/lib/components/StoryAgent.svelte.test.ts, design/system/workflow.md, design/system/agent-narrative.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flaiover.md, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/internal/serve/commit.go, flai/internal/serve/commit_test.go, flai/internal/serve/restart.go, flai/internal/serve/start.go, flai/internal/serve/stop.go, flai/internal/messages/messages.go, flai/internal/messages/messages_test.go, flaiover/src/lib/activity.ts]
after: [S-0331]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1576
  turns:
    - day: 2026-10-07
      ceremony: 4
      hand_edits: 2
      work: 56
  models:
    - model: claude-opus-5-5
      input: 296
      output: 119700
      cache_read: 16406659
      cache_write: 529866
      cost: 8.9359
  strategic:
    - kind: orchestrator
      seconds: 1911
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 99
          output: 1715
          cache_read: 12389311
          cache_write: 25232
          cost: 3.059
        - model: claude-sonnet-5-5
          input: 14
          output: 73
          cache_read: 210210
          cache_write: 52287
          cost: 0.2263
cost_of_delay:
  value: 93.02
  by: planner-E-0018
  at: 2026-10-07T20:22:04Z
forecast:
  duration: 36m
  delivery: 2026-10-08T09:05:00Z
  basis: "Its own forecast of 36m; 41st in the pull order with an in-progress limit of 3, behind S-0232, S-0233, S-0234, S-0235, S-0236, S-0237, S-0238, S-0239, S-0241, S-0280, S-0287, S-0288, S-0289, S-0290, S-0291, S-0297, S-0304, S-0305, S-0306, S-0309, S-0310, S-0312, S-0313, S-0314, S-0315, S-0316, S-0317, S-0318, S-0319, S-0320, S-0321, S-0322, S-0323, S-0324, S-0325, S-0326, S-0327, S-0332, S-0333 and S-0334."
  by: flai
  at: 2026-10-07T21:42:43Z
finalized:
  by: orchestrator
  at: 2026-10-07T20:24:18Z
---
# S-0335 A story's agent that waits only on another agent's reply ends, and flai serve starts it again when the reply or a new message to its story comes

## Goal

Make a message to another agent cost no more than a question to the operator. Since S-0272 a story's agent that has nothing left but the operator's answer ends, and flai serve starts it again in its session when the answer comes. A story's agent that waits on another agent's reply holds `wait_for_events` instead, paying a turn every wake, and an agent that ended on a question never sees a message sent to its story until the operator answers. Treat a message awaiting a reply as a question: the agent ends on it, and flai serve starts it again when the reply, or any new message to its story, arrives.

## Acceptance criteria

- [x] `wait_for_events` answers `end: true` with a `why` naming the conversations when flai serve started the agent for its story, a conversation of its story awaits the other story's reply, and no task of the story is in progress, as it does for a question to the operator.
- [x] flai serve records such a run `asked`, and starts it again in its session when the other story replies, as it does on a thread's answer.
- [x] flai serve starts again, in its session, an agent that ended `asked` when a new message to its story arrives, even while its question to the operator is unanswered.
- [x] `agent.status` and the dashboard's agent state say the agent waits on a story's agent, naming the story, rather than on the operator.
- [x] `design/system/workflow.md` § Branches and collisions, `agent-narrative.md`, and `work-management.md` in both copies describe it.

## Tasks

- T-1207 wait_for_events tells a story's agent to end when its story awaits only another story's reply
- T-1208 flai serve records an agent that ended on a conversation as asked, and starts it again on the reply or a new message to its story
- T-1209 The dashboard's story agent says an agent waits on another story's agent, naming the story
- T-1210 The workflow, the narrative design, the convention, the template, and the guides say an agent ends on another agent's reply

## Notes

### Planning

Planned by planner-E-0018 on 2026-10-07. It waits for S-0331 (`after`), whose inbox and events it builds on. It may run before S-0332 to S-0334; the hold orders it against them by their shared design documents.

Layers:

1. T-1207, the end of a wait; T-1208, the restart in flai serve. They share no path and run together.
2. T-1209, the dashboard. It waits for T-1208, which adds the field it shows.
3. T-1210, the docs.

Touches:

- **Declared:** none before planning.
- **Layout:**
  - `flai/internal/mcpserver/server.go` and its test: `endWhy` and `awaitsAnswer`.
  - `flai/internal/serve/agents.go` and its test: `judge`, `asked`, `answered`, and `Activity()`.
  - `flaiover/src/lib/components/StoryAgent.svelte` and its test: the agent's state on a story.
- **Co-change:** `flai touches suggest` from `agents.go`, `server.go`, and `restart.go` gave `server_test.go` (27%), `agents_test.go` (21%), and `docs/users/flai.md` and `flai-cli.md` (21 to 22%).
- **Design:** `design/system/workflow.md` § Branches and collisions (**An answered agent**) and `agent-narrative.md`; both copies of `work-management.md` and `template/CHANGELOG.md`; `docs/users/flaiover.md` for the waiting reason.
- **Not taken:** `flai/internal/mcpserver/folder.go` (41%), `flai/cmd/serve_actions.go` (23%), and `flai/internal/hostapi/writes.go` (16%): no tool, host action, or host write is added. `flai/internal/serve/restart.go` was a seed: the restart on an answer lives in `agents.go`; T-1208 widens its touches if it must change.

Forecast 36m, delivery 2026-10-08T10:04Z.

- `flai forecast` gave 25m: 78 s per unit over 44 done large-band improvement stories, times size 19 (5 criteria, 14 touches).
- Raised by 11m to the feature rate, 114 s per unit: it adds a restart trigger to flai serve, whose asked runs S-0317 shows are fragile, and a dashboard change. Delivery is shifted by the same 11m.

Cost of delay 93.02 USD a week: `flai cod` gave its share of E-0018's 1000 USD a week, 36m of 6h27m. It stands.

### Accepted by the orchestrator

- Verified: 372797d957f8ec13dbdc1f6a5c37e8c3e0f3db7e
- At: 2026-10-07T22:10:05Z

Verdict: all five criteria are met at 372797d9, the branch head; flai verify passed every step there, the branch changes no .claude/ path, and no convention is broken.
- 1: flai/internal/mcpserver/server.go, flai/internal/mcpserver/server_test.go
- 2: flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, flai/internal/messages/messages.go, flai/internal/harness/harness.go
- 3: flai/internal/serve/agents.go, flai/internal/serve/agents_test.go
- 4: flai/internal/serve/agents.go, flaiover/src/lib/activity.ts, flaiover/src/lib/components/StoryAgent.svelte, flaiover/src/lib/components/StoryAgent.svelte.test.ts
- 5: design/system/workflow.md, design/system/agent-narrative.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md
