---
id: S-0335
type: story
nature: improvement
title: A story's agent that waits only on another agent's reply ends, and flai serve starts it again when the reply or a new message to its story comes
status: backlog
parent: E-0018
owner: alex
created: 2026-10-07T20:11:14Z
updated: 2026-10-07T20:24:18Z
transitions: []
tags: [flai, flaiover, template]
topics: [cli, dashboard, conventions, template]
touches: [flai/internal/mcpserver/server.go, flai/internal/mcpserver/server_test.go, flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, flaiover/src/lib/components/StoryAgent.svelte, flaiover/src/lib/components/StoryAgent.svelte.test.ts, design/system/workflow.md, design/system/agent-narrative.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flaiover.md]
after: [S-0331]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 300
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 45
          output: 818
          cache_read: 3661515
          cache_write: 10355
          cost: 0.9048
cost_of_delay:
  value: 93.02
  by: planner-E-0018
  at: 2026-10-07T20:22:04Z
forecast:
  duration: 36m
  delivery: 2026-10-08T10:04:00Z
  basis: "flai forecast's size 19 at the feature rate of 114 s per unit instead of the improvement rate of 78 s, because it adds a restart trigger to flai serve, whose asked runs S-0317 shows are fragile; delivery shifted by the same 11m."
  by: planner-E-0018
  at: 2026-10-07T20:19:22Z
finalized:
  by: orchestrator
  at: 2026-10-07T20:24:18Z
---
# S-0335 A story's agent that waits only on another agent's reply ends, and flai serve starts it again when the reply or a new message to its story comes

## Goal

Make a message to another agent cost no more than a question to the operator. Since S-0272 a story's agent that has nothing left but the operator's answer ends, and flai serve starts it again in its session when the answer comes. A story's agent that waits on another agent's reply holds `wait_for_events` instead, paying a turn every wake, and an agent that ended on a question never sees a message sent to its story until the operator answers. Treat a message awaiting a reply as a question: the agent ends on it, and flai serve starts it again when the reply, or any new message to its story, arrives.

## Acceptance criteria

- [ ] `wait_for_events` answers `end: true` with a `why` naming the conversations when flai serve started the agent for its story, a conversation of its story awaits the other story's reply, and no task of the story is in progress, as it does for a question to the operator.
- [ ] flai serve records such a run `asked`, and starts it again in its session when the other story replies, as it does on a thread's answer.
- [ ] flai serve starts again, in its session, an agent that ended `asked` when a new message to its story arrives, even while its question to the operator is unanswered.
- [ ] `agent.status` and the dashboard's agent state say the agent waits on a story's agent, naming the story, rather than on the operator.
- [ ] `design/system/workflow.md` § Branches and collisions, `agent-narrative.md`, and `work-management.md` in both copies describe it.

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
