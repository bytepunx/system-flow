---
id: S-0272
type: story
nature: improvement
title: "An agent with an open question ends instead of waiting: flai serve restarts it on the answer, and wait_for_events keeps a timeout only for an agent with work in hand"
status: backlog
parent: E-0017
owner: alex
created: 2026-10-05T01:35:31Z
updated: 2026-10-06T23:33:30Z
transitions: []
tags: [cli, mcp]
topics: [automation, mcp, conventions, metrics, template]
touches: [flai/internal/mcpserver/server.go, flai/internal/mcpserver/server_test.go, flai/internal/serve/restart.go, flai/internal/serve/restart_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/internal/usage/log.go, flai/internal/usage/log_test.go, flai/internal/metrics/waiting.go, flai/internal/metrics/waiting_test.go, flai/cmd/stats.go, flai/cmd/check_stats_test.go, design/system/metrics.md, design/adrs, design/conventions/work-management.md, design/conventions/delegation.md, template/root/design/conventions/work-management.md, template/root/design/conventions/delegation.md, template/CHANGELOG.md, design/system/flai-cli.md, design/system/workflow.md, design/system/agent-narrative.md, docs/users/flai.md, docs/users/flai-reference.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  value: 96
  by: planner-E-0017
  at: 2026-10-06T11:36:16Z
forecast:
  duration: 50m
  delivery: 2026-10-07T08:57:00Z
  basis: "Its own forecast of 50m; 27th in the pull order with an in-progress limit of 3, behind S-0300, S-0302, S-0303, S-0273, S-0228, S-0269, S-0270, S-0271, S-0212, S-0213, S-0214, S-0215, S-0216, S-0232, S-0233, S-0234, S-0235, S-0236, S-0237, S-0238, S-0239, S-0241, S-0245, S-0246, S-0251, S-0254 and S-0265."
  by: flai
  at: 2026-10-06T23:33:30Z
finalized:
  by: alex
  at: 2026-10-06T22:48:30Z
---
# S-0272 An agent with an open question ends instead of waiting: flai serve restarts it on the answer, and wait_for_events keeps a timeout only for an agent with work in hand

## Goal

650 turns across 38 story runs woke from `wait_for_events` to find nothing, each a full-context model call. flai serve already restarts a story agent when its thread is answered, so an agent whose only pending work is a question has nothing to wait for. The convention and the harness prompt say to end the session on an open question after writing the narrative's state; `wait_for_events` answers at once with `end: true` when the caller's story has an open thread and no task in progress, and keeps its timeout for an agent that has work to go on with. The MCP server's instructions say the same.

## Acceptance criteria
- [ ] `wait_for_events` answers `end: true` and why when the calling agent's story has an open thread to the designer and no task in progress, and the harness prompt and conventions tell the agent to end then
- [ ] flai serve's restart on the answer is tested end to end: an agent that ended on a question is started again when the thread is answered, with the answer in its first inbox
- [ ] `design/system/flai-serve.md` (or where serve is described), the conventions, the template's copies, and the user guide describe the loop
- [ ] `flai stats` counts the empty wakes so that the saving is measured

## Tasks
- T-1087 wait_for_events answers end: true with why when the caller's story has an open question to the designer and no task in progress
- T-1090 The story agent's start prompt says to write the narrative's state and end on an open question, and to end when wait_for_events answers end
- T-1094 work-management.md and delegation.md, and their template copies, say a story's agent ends on an open question and is started again on the answer
- T-1097 An ADR and metrics.md define an empty wake and where flai stats reports it
- T-1101 flai stats counts each story's empty wakes from its agents' run logs, as text and --json
- T-1104 An end-to-end test: an agent that wait_for_events ends on a question is started again by flai serve on the answer, and its first inbox holds the answer
- T-1112 The design and the user guide describe the loop: end on a question, wait_for_events' end answer, the restart on the answer, and the empty-wake count

## Notes

From the epic's log classification: MCP wait_for_events 650 turns, 1,235 minutes, 38 stories.

### Planning

Planned by planner-S-0272 on 2026-10-06. It revisits planner-E-0017's enrichment and adds the tasks.

**Tasks, in three layers:**

- **Layer 1, waiting for nothing:**
  - T-1087: `wait_for_events` answers `end`.
  - T-1090: the start prompt.
  - T-1094: the conventions and their template copies.
  - T-1097: the ADR and `metrics.md`.
- **Layer 2:**
  - T-1101: the empty-wake count in `flai stats`. It waits for T-1097's definition.
  - T-1104: the end-to-end test. It waits for T-1087's `end`.
- **Layer 3:** T-1112, the design and the user guide. It waits for all the rest, because it describes what they built.

**Touches.** All 21 that planner-E-0017 declared are kept. `flai touches suggest S-0272` lists 425 of 1,050 commits co-changing them. None of its top entries are reached by the goal (the dashboard design, the operator docs, `folder.go`, `hostapi/writes.go`), so none was added.

- `flai/internal/mcpserver/server.go`, `server_test.go`: declared, layout. `waitForEvents`, `WaitOut`, the tool's description and the server's instructions are here (T-1087).
- `flai/internal/serve/restart.go`: declared. It is the operator's restart from the dashboard. The restart on an answer is in `serve/agents.go`, and `TestAnAgentThatEndedAskingIsStartedAgainWhenAnswered` already tests serve's half. No task changes `restart.go`. It is kept because it was declared.
- `flai/internal/serve/restart_test.go`: declared, layout. This is where T-1104's new test goes. If the test finds a gap, T-1104 adds `agents.go` to its touches.
- `flai/internal/harness/harness.go`, `harness_test.go`: declared, layout. The story-run rules and the commit-run prompt both say to hold `wait_for_events` for an answer (T-1090).
- `flai/internal/usage/log.go`, `flai/internal/metrics/waiting.go`, `waiting_test.go`, `flai/cmd/stats.go`, `flai/cmd/check_stats_test.go`: declared, design (T-1101).
- `flai/internal/usage/log_test.go`: added, layout. The tests of the log parsing that T-1101 extends.
- `design/system/metrics.md`: declared, design (T-1097).
- `design/adrs`: declared, design. **Kept as a folder touch** because the ADR's file name is not known until `flai adr new` numbers it. It is a shared path in the manifest.
- `design/conventions/work-management.md`, `delegation.md`, their `template/root` copies, `template/CHANGELOG.md`: declared, layout (T-1094).
- `design/system/flai-cli.md`, `docs/users/flai.md`, `docs/users/flai-reference.md`: declared, co-change (T-1112).
- `design/system/workflow.md`: added, design. Its **An answered agent** paragraph is where serve's restart on an answer is described. There is no `flai-serve.md` (criterion 3).
- `design/system/agent-narrative.md`: added, design. Its lines on `wait_for_events` say how long it holds.
- Left out:
  - `flai/internal/guard`: S-0285's refusal of `wait_for_events` while a sub-agent runs is unchanged.
  - `serve/agents.go`: serve already restarts on an answer.
  - `mcpserver/issues.go`: the story is read from `FLAI_STORY`, without changing its resolver.

**Topics.** Added `metrics` (T-1097, T-1101) and `template` (T-1094).

**Forecast: 50m**, delivery 2026-10-07T08:45Z. flai gives 39m: 83 s per unit over 25 done large improvement stories, times size 28 (4 criteria, 24 touches), delivered 08:34Z. I kept the earlier 50m and moved the delivery by the 11m added. The reasons:

- seven tasks in three layers
- an ADR
- a new end-to-end test across serve and the MCP server

S-0285, on the same tool and prompt, took 42m with less. Serve's half of the restart is already tested, which offsets some of this.

**Cost of delay: 96 USD a week**, unchanged. flai gives 115.98, E-0017's 900 USD a week shared by forecast duration. The 96 is the epic's share by the turns each story removes, as planner-E-0017 set it: this story removes the 650 empty wakes. It is kept so that the siblings' values still add up to the epic's 900.

**Assumptions:**

- The caller's story is `FLAI_STORY`, which flai serve sets. An agent run by hand, which serve would not restart, keeps its timeout.
- The planner, the orchestrator and the analyzer keep holding for answers.
- An empty wake is a timed-out call with no events and no changed paths.
- Criterion 4 overlaps S-0293. TH-0176's third proposal, to move it there, was resolved without a change to the criteria, so T-1097 and T-1101 deliver it here.
