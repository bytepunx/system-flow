---
id: S-0294
type: story
nature: remediation
title: Claude Code ends a headless agent ten minutes after its turn ends, even while its background sub-agent is still working, and flai serve leaves the story in progress with no agent
status: backlog
owner: alex
created: 2026-10-06T11:44:49Z
updated: 2026-10-06T23:40:32Z
transitions: []
tags: [cli, serve]
topics: [automation]
touches: [flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, flai/internal/serve/restart.go, flai/internal/config/config.go, flai/internal/config/config_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, design/adrs, design/system/flai-cli.md, design/system/workflow.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/index.md, docs/operators/settings.md, design/issues/I-0084-claude-code-ends-a-headless-agent-ten-minutes-after-its-turn-ends-even-while-its-background-sub-agent-is-still-working-and-flai-serve-leaves-the-story-in-progress-with-no-agent.md, design/issues/summary.md]
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
    - kind: planner
      seconds: 614
      estimated: true
      models:
        - model: claude-haiku-4-5-20251001
          input: 507
          output: 4242
          cache_read: 3640961
          cache_write: 173470
          cost: 0.7745
        - model: claude-opus-5-5
          input: 996
          output: 69209
          cache_read: 50650550
          cache_write: 711112
          cost: 27.9391
cost_of_delay:
  inputs:
    time_lost_per_cycle: 20m
    by: flai
    at: 2026-10-06T11:44:49Z
  value: 50
  by: planner-S-0294
  at: 2026-10-06T22:53:45Z
forecast:
  duration: 40m
  delivery: 2026-10-07T10:45:00Z
  basis: "Its own forecast of 40m; 40th in the pull order with an in-progress limit of 3, behind S-0302, S-0303, S-0273, S-0228, S-0269, S-0270, S-0271, S-0212, S-0213, S-0214, S-0215, S-0216, S-0232, S-0233, S-0234, S-0235, S-0236, S-0237, S-0238, S-0239, S-0241, S-0245, S-0246, S-0251, S-0254, S-0265, S-0272, S-0274, S-0275, S-0277, S-0279, S-0280, S-0281, S-0286, S-0287, S-0288, S-0289, S-0290, S-0291 and S-0293."
  by: flai
  at: 2026-10-06T23:40:32Z
finalized:
  by: alex
  at: 2026-10-06T22:50:34Z
---
# S-0294 Claude Code ends a headless agent ten minutes after its turn ends, even while its background sub-agent is still working, and flai serve leaves the story in progress with no agent

## Goal

This story remediates [I-0084](../../../design/issues/I-0084-claude-code-ends-a-headless-agent-ten-minutes-after-its-turn-ends-even-while-its-background-sub-agent-is-still-working-and-flai-serve-leaves-the-story-in-progress-with-no-agent.md), "Claude Code ends a headless agent ten minutes after its turn ends, even while its background sub-agent is still working, and flai serve leaves the story in progress with no agent". The issue recommends this solution:

Two separate things went wrong, and each has its own fix.

1. **The agent has no safe way to wait.** The prompt and `delegation.md` name one that keeps the session alive for a sub-agent of any length: a launch with `run_in_background: false`, once tested against a sub-agent that runs longer than ten minutes. This is the first change I-0083 recommends, and one story can make it for both issues.
2. **A story whose agent ended without finishing waits for a person.** flai serve could restart such an agent itself: once, or a small number of times, when the run ended with the story still in progress, no thread awaiting the operator, and no block on the story; and tell the operator when the limit is reached. Whether flai serve restarts on its own is the operator's decision, since ADR-0043 gave the restart to a person.

S-0285 made the first remediation ([ADR-0092](../adrs/0092-a-story-s-agent-waits-for-a-sub-agent-by-launching-it-in-the-foreground-and.md)). It measured the ten minutes again on 2.1.290: a turn ended at 10:34:48Z with a background sub-agent out, and the process exited at 10:44:49Z with the sub-agent cut off. It also measured the remedy: a launch with `run_in_background` false returned an 11-minute sub-agent's result as the tool's result, and three launched in one message ran together. The start prompt and `delegation.md` now name that way and say never to end the turn on a background sub-agent. The second remediation is still open.

## Acceptance criteria
- [ ] The cause I-0084 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0084 is closed with `flai issue close I-0084 --reason` saying what fixed it

## Tasks
- T-1124 An ADR refining ADR-0043, flai-cli.md, and workflow.md say flai serve restarts a story's agent that ended with its story in progress, up to agent.auto_restarts times, then opens a thread
- T-1125 The host setting agent.auto_restarts, 2 when unset, is set with flai serve agent set --auto-restarts and shown with the rest of the agent settings
- T-1126 flai serve restarts a story's agent that ended with its story in progress, unblocked, and asking nothing, up to agent.auto_restarts times, then opens a thread to the operator, with a test that reproduces I-0084
- T-1127 The users' and operators' guides and the reference describe the automatic restart and agent.auto_restarts, and I-0084 is closed

## Notes

Cost of delay inputs set by flai from I-0084. time_lost_per_cycle 20m: 10m per occurrence × 2 occurrences ÷ 1 cycle of 168h (first reported 2026-10-06T07:00:22Z, 0.2 days before this story; under one cycle counts as one).

### Planning

The plan rests on TH-0208. On it I asked the operator whether flai serve may restart a story's agent on its own, which ADR-0043 gave to a person. I recommended yes: up to 2 times, under a host setting `agent.auto_restarts`, then a thread to the operator. The operator answered `yes` on 2026-10-06, and the tasks follow it.

Touches. None were declared. I seeded `flai touches suggest S-0294` with `flai/internal/serve/agents.go`, `restart.go`, and `flai/internal/config/config.go`, which 46 of 1050 commits changed.

- `flai/internal/serve/agents.go`, `agents_test.go`: layout. `judge` records the failed run, and `resume` and `start` restart a run.
- `flai/internal/serve/restart.go`: layout. The operator's restart resets the count.
- `flai/internal/config/config.go`, `config_test.go`: layout. `AgentStart` holds this host's agent settings.
- `flai/cmd/serve_actions.go`, `serve_actions_test.go`: co-change (20 and 14 of 46). `flai serve agent set`, and where `serve.AgentConfig` is built.
- `design/adrs`: design. The ADR refining ADR-0043. This is a folder touch, kept because the ADR's file name is not known until `flai adr new` numbers it. It is a shared path, so it holds no ready story.
- `design/system/flai-cli.md`: co-change (10 of 46). `design/system/workflow.md`: design, since it says a restart is "on the operator's word".
- `docs/users/flai.md`, `docs/users/flai-reference.md`: co-change (8 of 46). `docs/operators/index.md`: co-change (7), its Restart paragraph. `docs/operators/settings.md`: layout, the host config table.
- I-0084's file and `design/issues/summary.md`: from criterion 2.
- Left out:
  - `flai/internal/hostapi/writes.go`: co-changed, but the Retry write is unchanged.
  - `flai/internal/harness/harness.go`: the restarted agent's prompt is built from the ended run, as for a manual restart.
  - flaiover: a run's `why` already reaches the story page.
  - The dashboard's Settings page: it does not list host agent settings.

Forecast. flai gave 20m: 64 s per unit over 13 done large remediation stories, times size 18. I raised it to 40m. The story needs an ADR, a new host setting, and a launcher change tested across six conditions. S-0285, on the same issue, took 42m. The delivery is flai's, moved by the added 20m.

Cost of delay. The value of 50 USD a week stands as flai's figure from the inputs: 20m lost per 168h cycle at 150 USD an hour. Unattended, the gap I-0084 describes has no end, and it holds every ready story that overlaps. That argues for more, but the inputs are the operator's, and I leave them as they are.

Overlap. S-0272 also touches `flai/internal/serve/restart.go`, `flai-cli.md`, `workflow.md`, `docs/users/flai.md`, and `flai-reference.md`. Whichever of the two starts second is held until the other moves to review.
