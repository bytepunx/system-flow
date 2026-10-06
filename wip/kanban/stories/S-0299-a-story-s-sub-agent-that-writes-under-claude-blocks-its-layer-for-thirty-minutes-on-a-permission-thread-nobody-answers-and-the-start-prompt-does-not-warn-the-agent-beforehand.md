---
id: S-0299
type: story
nature: remediation
title: A story's sub-agent that writes under .claude/ blocks its layer for thirty minutes on a permission thread nobody answers, and the start prompt does not warn the agent beforehand
status: ready
owner: alex
created: 2026-10-06T21:00:33Z
updated: 2026-10-06T22:12:16Z
transitions:
  - to: ready
    at: 2026-10-06T21:39:22Z
    by: alex
tags: [flai, guard, harness]
topics: [cli, conventions]
touches: [flai/internal/guard/guard.go, flai/internal/guard/guard_test.go, flai/cmd/guard.go, flai/cmd/guard_test.go, flai/internal/hostapi/writes.go, docs/users/flai-reference.md, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, design/conventions/delegation.md, design/system/flai-cli.md, docs/users/flai.md, docs/operators/settings.md, design/adrs, design/adrs/README.md, design/issues/I-0093-a-story-s-sub-agent-that-writes-under-claude-blocks-its-layer-for-thirty-minutes-on-a-permission-thread-nobody-answers-and-the-start-prompt-does-not-warn-the-agent-beforehand.md, design/issues/summary.md]
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
      seconds: 443
      estimated: true
      models:
        - model: claude-haiku-4-5-20251001
          input: 154
          output: 8534
          cache_read: 1111223
          cache_write: 97672
          cost: 0.276
        - model: claude-opus-5-5
          input: 68
          output: 22649
          cache_read: 2370457
          cache_write: 102839
          cost: 1.7501
cost_of_delay:
  inputs:
    time_lost_per_cycle: 6m
    by: flai
    at: 2026-10-06T21:00:33Z
  value: 15
  by: planner-S-0299
  at: 2026-10-06T21:07:32Z
forecast:
  duration: 40m
  delivery: 2026-10-06T23:00:00Z
  basis: "Its own forecast of 40m; 1st in the pull order with an in-progress limit of 3, behind S-0229."
  by: flai
  at: 2026-10-06T22:12:16Z
finalized:
  by: alex
  at: 2026-10-06T21:01:56Z
---
# S-0299 A story's sub-agent that writes under .claude/ blocks its layer for thirty minutes on a permission thread nobody answers, and the start prompt does not warn the agent beforehand

## Goal

This story remediates [I-0093](../../../design/issues/I-0093-a-story-s-sub-agent-that-writes-under-claude-blocks-its-layer-for-thirty-minutes-on-a-permission-thread-nobody-answers-and-the-start-prompt-does-not-warn-the-agent-beforehand.md), "A story's sub-agent that writes under .claude/ blocks its layer for thirty minutes on a permission thread nobody answers, and the start prompt does not warn the agent beforehand". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0093 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0093 is closed with `flai issue close I-0093 --reason` saying what fixed it

## Tasks
- T-1035 flai guard refuses a sub-agent's write under .claude/ at once while auto-approve is off, so its layer never waits on a permission thread
- T-1036 The story agent's start prompt says, before it launches a sub-agent, that a write under .claude/ is its own and waits on the operator
- T-1037 An ADR, the design, the docs, and delegation.md say a .claude/ write is the story agent's, and I-0093 is closed

## Notes

Cost of delay inputs set by flai from I-0093. time_lost_per_cycle 6m: 6m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-06T20:53:14Z, 0 days before this story; under one cycle counts as one).

### Planning

The proposed fix, from I-0093's one instance (S-0223's T-0957), has two parts:

- `flai guard` refuses a sub-agent's write to a file in a `.claude/` folder at once while `auto-approve` is off (T-1035). It already sees `Edit|Write|NotebookEdit` and knows a sub-agent by its `agent_id`, so no `.claude/settings.json` change is needed. A refusal costs the sub-agent one turn instead of holding its layer on a thread.
- The start prompt warns the story's agent before it launches anything (T-1036). Such writes are its own, made after the layer, or handed to the owner as `cp` commands on a thread when they may be away.

T-1037 then records the decision in an ADR, the design, and the docs, and closes the issue.

Touches, by where each came from:

- Declared: none; the story declared no touches, so `flai touches suggest` needed a starting set.
- Design and layout:
  - `flai/internal/guard/guard.go` and `guard_test.go`: `Guard.Decide` decides a sub-agent's call.
  - `flai/cmd/guard.go` and `guard_test.go`: wire auto-approve into the hook command and test it end to end.
  - `flai/internal/hostapi/writes.go`: `ActionAutoApprove`'s description.
  - `docs/users/flai-reference.md`: generated from `flai guard`'s help.
  - `flai/internal/harness/harness.go` and `harness_test.go`: `delegation()` is the start prompt's sub-agent text.
  - `design/conventions/delegation.md`: its project addition tells agents to write `.claude/` files with `Edit` or `Write`, which is how T-0957 came to.
  - `design/system/flai-cli.md`, `docs/users/flai.md` § Writes under .claude/, and `docs/operators/settings.md`: describe the guard and `auto-approve`.
  - `design/adrs/README.md`: indexes the new ADR.
  - The I-0093 file and `design/issues/summary.md`: `flai issue close` writes both.
- Co-change: `flai touches suggest` listed only paths changed with nearly everything (the dashboard design, the operators' index), none bearing on this fix, so none were taken. `docs/operators/settings.md` was among them and is kept for the design reason above.
- Folder touch kept: `design/adrs`. T-1037 adds a new ADR whose number and file name `flai adr new` gives only when it runs.

Figures:

- Forecast: 40m, against `flai forecast`'s 18m (64 s per unit of size × size 16). The four closest stories each changed `permission_prompt`, the guard, or the start prompt: S-0283 2010 s, S-0257 2168 s, S-0285 2506 s, S-0266 2559 s, a median of about 39m. Delivery moved 22m later than flai's 2026-10-07T08:12Z to match.
- Cost of delay: 15 USD a week, as `flai cod` gives it from flai's input of 6m lost per 168h cycle at 150 USD an hour. Left unchanged. The input likely understates the cost: the watcher's restart cut the hold to six minutes where it could have run thirty, and the board watch's notes on TH-0194 say S-0222 also had a session stopped for it. The input is the operator's to raise, with `flai issue bump I-0093` or by setting it.
