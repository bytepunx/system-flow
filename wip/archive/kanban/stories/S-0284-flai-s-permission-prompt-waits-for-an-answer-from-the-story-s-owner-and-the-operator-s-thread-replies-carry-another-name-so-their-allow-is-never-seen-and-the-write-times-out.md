---
id: S-0284
type: story
nature: remediation
title: flai's permission_prompt waits for an answer from the story's owner, and the operator's thread replies carry another name, so their allow is never seen and the write times out
status: done
owner: alex
created: 2026-10-06T03:45:19Z
updated: 2026-10-06T19:45:24Z
transitions:
  - to: ready
    at: 2026-10-06T06:20:07Z
    by: alex
  - to: in-progress
    at: 2026-10-06T19:35:46Z
    by: agent-S-0284
  - to: review
    at: 2026-10-06T19:44:34Z
    by: agent-S-0284
  - to: done
    at: 2026-10-06T19:45:24Z
    by: alex
tags: []
topics: [cli]
touches: [flai/internal/mcpserver/permission.go, flai/internal/mcpserver/permission_test.go, design/adrs, design/system/flai-cli.md, docs/users/flai.md, design/issues/I-0081-flai-s-permission-prompt-waits-for-an-answer-from-the-story-s-owner-and-the-operator-s-thread-replies-carry-another-name-so-their-allow-is-never-seen-and-the-write-times-out.md, design/issues/summary.md, design/issues/I-0078-flai-check-finds-item-archive-outside-the-story-at-close-out.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 561
  models:
    - model: claude-haiku-4-5-20251001
      input: 250
      output: 6971
      cache_read: 1794209
      cache_write: 79620
      cost: 0.3141
    - model: claude-opus-5-5
      input: 116
      output: 34363
      cache_read: 3776979
      cache_write: 203340
      cost: 2.748
    - model: claude-sonnet-5-5
      input: 14
      output: 2575
      cache_read: 187443
      cache_write: 42217
      cost: 0.1688
cost_of_delay:
  inputs:
    time_lost_per_cycle: 30m
    by: flai
    at: 2026-10-06T03:45:19Z
  value: 75
  by: planner-S-0284
  at: 2026-10-06T06:23:46Z
forecast:
  duration: 30m
  delivery: 2026-10-06T20:30:00Z
  basis: "Its own forecast of 30m; 2nd in the pull order with an in-progress limit of 3, behind S-0296."
  by: flai
  at: 2026-10-06T19:35:23Z
finalized:
  by: alex
  at: 2026-10-06T06:18:57Z
---
# S-0284 flai's permission_prompt waits for an answer from the story's owner, and the operator's thread replies carry another name, so their allow is never seen and the write times out

## Goal

This story remediates [I-0081](../../../design/issues/I-0081-flai-s-permission-prompt-waits-for-an-answer-from-the-story-s-owner-and-the-operator-s-thread-replies-carry-another-name-so-their-allow-is-never-seen-and-the-write-times-out.md), "flai's permission_prompt waits for an answer from the story's owner, and the operator's thread replies carry another name, so their allow is never seen and the write times out". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0081 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0081 is closed with `flai issue close I-0081 --reason` saying what fixed it

## Tasks
- T-0997 An ADR refining ADR-0086 lets the project's owner answer a permission thread as well as the story's owner
- T-0998 permission_prompt takes an answer from the story's owner or the project's owner, with a test reproducing I-0081
- T-0999 The design and the user guide say the story's owner or the project's owner answers a permission thread
- T-1000 Close I-0081 saying what fixed it

## Notes

Cost of delay inputs set by flai from I-0081. time_lost_per_cycle 30m: 30m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-05T07:45:13Z, 0.8 days before this story; under one cycle counts as one).

### Planning

Planned by planner-S-0284 on 2026-10-06.

Touches, none declared before planning:

- `flai/internal/mcpserver/permission.go`: design. I-0081's instance names `awaitAnswer` there, which takes an answer from the story's owner only.
- `flai/internal/mcpserver/permission_test.go`: co-change. `flai touches suggest` from `permission.go` gives it in 2 of 2 commits, and the first criterion asks for a test.
- `design/adrs`: design. ADR-0086's decision 2 says entries by anyone but the story's owner are not answers; changing that needs a new ADR that refines it.
- `design/system/flai-cli.md` and `docs/users/flai.md`: design. Both say only the story's owner answers a permission thread.
- The I-0081 file and `design/issues/summary.md`: layout. The second criterion closes I-0081, which rewrites the issue and the summary.

Topic `cli` added: ADR-0086 and the `flai mcp` tool are filed under it.

Forecast 30m, delivery 2026-10-06T10:45Z. `flai forecast` gave 18m (median 116 s per unit over 10 done medium-band remediation stories on claude-opus-5-5, times size 9, 2 criteria and 7 touches), delivering at 10:33Z, 9th in the pull order. Raised by 12m because the story also writes an ADR and proposes the fix before building it, and the criteria count neither. Delivery shifted by the same 12m.

Cost of delay 75 USD a week, as `flai cod` gives from the operator's input (30m lost per 168h cycle at 150 USD an hour). It stands: the issue has one instance. But each instance holds a story's agent blocked in a tool call until the MCP idle timeout, so a recurrence costs more than the input counts.

Overlap: S-0283 (I-0082) also changes `permission.go`, and S-0220, in progress, touches `flai/internal/mcpserver`. flai's overlap hold runs them one after another, so no `after` is set.
