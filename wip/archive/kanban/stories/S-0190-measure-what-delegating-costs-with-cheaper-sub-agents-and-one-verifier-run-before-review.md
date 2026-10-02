---
id: S-0190
type: story
nature: research
title: Measure what delegating costs with cheaper sub-agents and one verifier run before review
status: done
owner: arobson
created: 2026-10-01T10:52:53Z
updated: 2026-10-02T17:01:34Z
transitions:
  - to: ready
    at: 2026-10-01T11:46:27Z
    by: alex
  - to: in-progress
    at: 2026-10-02T12:09:34Z
    by: agent-S-0190
  - to: ready
    at: 2026-10-02T12:41:31Z
    by: agent-S-0190
  - to: backlog
    at: 2026-10-02T12:41:31Z
    by: agent-S-0190
  - to: ready
    at: 2026-10-02T16:53:13Z
    by: alex
  - to: in-progress
    at: 2026-10-02T16:53:34Z
    by: agent-S-0190
  - to: review
    at: 2026-10-02T16:59:13Z
    by: agent-S-0190
  - to: done
    at: 2026-10-02T17:01:34Z
    by: alex
blocked:
  - from: 2026-10-02T12:29:38Z
    until: 2026-10-02T16:53:11Z
    reason: "Waiting for the next story run on claude-opus-5-5 that changes code (S-0231, S-0191, S-0192, or S-0178) to end; its log is the second run (TH-0064)"
tags: [cli]
topics: [conventions]
touches: [design/system/agent-context.md]
after: [S-0189]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1731
  models:
    - model: claude-opus-5-5
      input: 234
      output: 67705
      cache_read: 10920140
      cache_write: 230847
      cost: 5.3858
    - model: claude-sonnet-5-5
      input: 18
      output: 7129
      cache_read: 202347
      cache_write: 70969
      cost: 0.2892
---
# S-0190 Measure what delegating costs with cheaper sub-agents and one verifier run before review

## Goal

S-0189 ran the explorer on `haiku` and the verifier on `sonnet`, let roles carry their own model (`agent.roles`), and had the verifier's run replace the story's agent's own full-suite runs. Measure whether delegation now costs less than not delegating, the way S-0188 measured S-0184 and S-0185 (`design/system/agent-context.md § Sub-agents › Measured`).

## Acceptance criteria

- [x] Two stories run by `flai serve` with a flai that includes S-0189 are measured as S-0188 measured S-0184 and S-0185, with the same reading of the logs and against the same non-delegating comparables
- [x] The table in `agent-context.md § Sub-agents › Measured` gains their rows and shows whether the cost of delegation fell below the cost of not delegating, with each sub-agent's model and share of the cost
- [x] The section says what to change next, if anything

## Tasks
- T-0688 Measure S-0194 and the next Opus code story, the runs on a flai with S-0189, as S-0188 measured
- T-0689 Add the runs to agent-context.md § Sub-agents › Measured with each sub-agent's model and share
- T-0690 Say what to change next in § Measured

## Notes

- Split from S-0189 criterion 4 on the designer's word (TH-0057): it can be measured only after S-0189 is accepted, flai is released and installed, and two more stories have run.
- Pull it once two stories have run with the release that includes S-0189.
