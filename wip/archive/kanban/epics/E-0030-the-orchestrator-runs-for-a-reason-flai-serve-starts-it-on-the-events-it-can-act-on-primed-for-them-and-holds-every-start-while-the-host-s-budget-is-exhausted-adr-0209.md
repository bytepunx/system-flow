---
id: E-0030
type: epic
nature: improvement
title: "The orchestrator runs for a reason: flai serve starts it on the events it can act on, primed for them, and holds every start while the host's budget is exhausted (ADR-0209)"
status: done
owner: alex
created: 2026-10-28T20:59:30Z
updated: 2026-10-29T17:25:36Z
transitions:
  - to: ready
    at: 2026-10-28T21:16:58Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-28T21:17:05Z
    by: agent-S-0533
  - to: review
    at: 2026-10-29T17:21:54Z
    by: agent-S-0537
  - to: done
    at: 2026-10-29T17:25:36Z
    by: alex
tags: [flai, orchestrator, serve, budget]
topics: [orchestration, cli, dashboard]
usage:
  source: sum
  seconds: 33006
  priced_by: harness
  turns:
    - day: 2026-10-28
      ceremony: 4
      test_runs: 3
      hand_edits: 3
      work: 109
    - day: 2026-10-29
      ceremony: 23
      test_runs: 16
      hand_edits: 18
      work: 384
  models:
    - model: claude-opus-5-5
      input: 3946
      output: 1481361
      cache_read: 275425936
      cache_write: 9558117
      cost: 140.4805
  strategic:
    - kind: planner
      seconds: 1976
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 210
          output: 680
          cache_read: 12844746
          cache_write: 388530
          cost: 0.8668
        - model: claude-opus-5-5
          input: 112
          output: 4343
          cache_read: 7302628
          cache_write: 372766
          cost: 2.0487
    - kind: orchestrator
      seconds: 36009
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 1665
          output: 27572
          cache_read: 310553907
          cache_write: 704516
          cost: 68.5823
        - model: claude-sonnet-5-5
          input: 76
          output: 427
          cache_read: 1337207
          cache_write: 311559
          cost: 1.1378
---
# E-0030 The orchestrator runs for a reason: flai serve starts it on the events it can act on, primed for them, and holds every start while the host's budget is exhausted (ADR-0209)

## Outcome

[ADR-0209](../../../design/adrs/0209-the-orchestrator-runs-for-a-reason-and-ends-started-by-flai-serve-on-the-events.md), drafted on 2026-10-28 from a fortnight of orchestrator session logs, finds that the orchestrator's continuous loop pays the whole transcript, about 450k tokens, on every wake, and that 75 to 80 percent of its wakes end in another wait: a path change with no event, the replan's forecast edits, a task move. One system-flow session read 3.3B context tokens in a week; the week after read 3.6B in four active days with five projects served. When the account's usage limit refuses the API, flai serve starts the orchestrator again every minute with no backoff and tells nothing else: Oct 26 and 27 produced 8,986 refused starts.

When this epic is done, the orchestrator runs for a reason and ends, like the planner and the curator: flai serve queues and coalesces the reasons, names them in the start prompt, and primes a pack for them. Starts, ends, and failures of agents are events every agent sees. A usage-limit refusal is a host state that holds every start on that provider until the reset, circulated through the inbox, the journal, the dashboard, and the operator's channel. The orchestration page shows its settings above its status and activity, which the operator asked for on the same day. The orchestrator's context tokens per action, measured at about 7M on 2026-10-28, fall to a fraction, recorded in the close-out.

The ADR is in the design folder as proposed; the story that changes the orchestrator's trigger accepts it, and the living design, the convention, the template's role file, and the user docs change in the same stories.

## Stories
- S-0532 Starts, ends, and failures of a role's run and of a story's agent are events of kind agent that wait_for_events and inbox report, and the thread flai serve opens after auto_restarts cites them
- S-0533 The orchestration page shows its settings panel first, above the action's state, the run, the last decisions, and the activity log, so the operator reaches the permissions and the policy without scrolling past them
- S-0534 A usage-limit refusal sets a budget state on the host that holds every start on that provider until the reset or a backoff, probes once, and reaches every inbox, the journal, the dashboard, and the ADR-0205 channel
- S-0535 flai emits the events the orchestrator's trigger names: planned, ready-empty, release-due, a thread coming to await an agent, and the conditions awaits-orchestrator, draft-changed, and to-review
- S-0536 The orchestrator runs for a reason and ends: its trigger names events, flai check rejects continuous, flai serve coalesces the reasons into one run naming them, wait_for_events ends the run when they are done, and ADR-0209 is accepted
- S-0537 flai prime --role orchestrate takes the run's reasons and primes for them, the role's conventions once and per reason what it needs, instead of the general pack

## Notes

Evidence and figures: ADR-0209's context section. The session logs are under `~/.flai/serve/agents/<key>-orchestrator-<start>.log` on the host.
- 2026-10-28T21:16:58Z: moved to ready: follows S-0533, which moved to ready
- 2026-10-28T21:17:05Z: moved to in-progress: follows S-0533, which moved to in-progress
- 2026-10-29T17:21:54Z: moved to review: follows S-0537, which moved to review
- 2026-10-29T17:25:36Z: moved to done: follows S-0537, which moved to done
