---
id: S-0730
type: story
nature: improvement
title: A story's agent can overwrite a work item file whole with Write, front matter included, and nothing refuses it
status: review
owner: alex
created: 2026-11-12T03:43:39Z
updated: 2026-11-13T20:44:41Z
transitions:
  - to: ready
    at: 2026-11-13T19:57:17Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-13T19:57:25Z
    by: agent-S-0730
  - to: review
    at: 2026-11-13T20:44:41Z
    by: agent-S-0730
tags: [flai, template]
topics: [guard]
touches: [flai/internal/guard/call.go, flai/internal/guard/claudecode.go, flai/internal/guard/claudecode_test.go, flai/internal/guard/opencode.go, flai/internal/guard/opencode_test.go, flai/internal/guard/items.go, flai/internal/guard/items_test.go, design/system/agent-adapters.md, docs/users/flai.md, design/conventions/tooling.md, template/root/design/conventions/tooling.md, template/CHANGELOG.md, design/issues/I-0352-a-story-s-agent-can-overwrite-a-work-item-file-whole-with-write-front-matter.md, design/issues/summary.md, design/issues/I-0365-flai-test-reports-a-go-test-failure-in-flai-cmd-that-names-no-test-and-the-same.md, design/adrs/0298-flai-guard-refuses-a-write-that-replaces-a-work-item-file-whole-with-content.md, design/adrs/README.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 2859
  priced_by: harness
  turns:
    - day: 2026-11-13
      ceremony: 1
      work: 133
  models:
    - model: claude-opus-5-5
      input: 326
      output: 46955
      cache_read: 39400527
      cache_write: 377275
      cost: 11.6426
  strategic:
    - kind: planner
      seconds: 278
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 36
          output: 11831
          cache_read: 971377
          cache_write: 87110
          cost: 0.0265
        - model: claude-opus-5-5
          input: 54
          output: 18202
          cache_read: 2043999
          cache_write: 110414
          cost: 1.6564
    - kind: orchestrator
      seconds: 246
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 66
          output: 466
          cache_read: 2997883
          cache_write: 201155
          cost: 0.9347
        - model: claude-sonnet-5-5
          input: 4
          output: 860
          cache_read: 73863
          cache_write: 11380
          cost: 0.067
cost_of_delay:
  inputs:
    time_lost_per_cycle: 1m
    by: flai
    at: 2026-11-12T03:43:39Z
  value: 2.5
  by: planner-S-0730
  at: 2026-11-13T19:54:55Z
forecast:
  duration: 45m
  delivery: 2026-11-13T20:51:00Z
  basis: "flai forecast's 107 s per unit of size over 179 done large-band improvement stories gives about 27m; raised to 45m for two harness readers, a new guard rule with its tests, and a template convention change over four tasks."
  by: planner-S-0730
  at: 2026-11-13T19:54:55Z
finalized:
  by: orchestrator
  at: 2026-11-13T19:57:00Z
---
# S-0730 A story's agent can overwrite a work item file whole with Write, front matter included, and nothing refuses it

## Goal

This story remediates [I-0352](../../../design/issues/I-0352-a-story-s-agent-can-overwrite-a-work-item-file-whole-with-write-front-matter.md), "A story's agent can overwrite a work item file whole with Write, front matter included, and nothing refuses it". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0352 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0352 is closed with `flai issue close I-0352 --reason` saying what fixed it

## Tasks
- T-3319 flai guard carries the content a Write or an OpenCode write gives into its Call
- T-3320 flai guard refuses a Write that replaces a work item file with content that is not that item
- T-3321 The guard's design, the user guide, and the tooling convention say a work item file is never replaced whole by a Write
- T-3322 I-0352 is closed with the guard rule that fixed it

## Notes

Cost of delay inputs set by flai from I-0352. time_lost_per_cycle 1m: 1m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-11T22:53:06Z, 0.2 days before this story; under one cycle counts as one).

### Planning

The proposed fix follows the issue's own instance. `flai guard` refuses a whole-file Write to a work item file when the new content is not that item: it has no front matter, its front matter does not parse, or its id or type differs from the file on disk. The refusal points the agent to `Edit` or `flai edit`. The guard drops a Write's content today (`Input` in `flai/internal/guard/claudecode.go`, `openCodeArgs` in `flai/internal/guard/opencode.go`), so the content has to reach `Call` first.

Each touch names a file. There is no folder touch.

| Touch | From |
|-------|------|
| `flai/internal/guard/call.go`, `claudecode.go`, `opencode.go` | layout: where the hook input is read into `Call` and decided (`Decide`) |
| `flai/internal/guard/claudecode_test.go`, `opencode_test.go` | layout: the reader tests of each harness |
| `flai/internal/guard/items.go`, `items_test.go` | layout: a new file for the rule and its test that reproduces I-0352 |
| `design/system/agent-adapters.md` | design: `### The guard` lists what the hook reads and refuses |
| `docs/users/flai.md` | design: the user guide's list of guard refusals |
| `design/conventions/tooling.md`, `template/root/design/conventions/tooling.md` | design: the baseline line on hand-editing item bodies |
| `template/CHANGELOG.md` | co-change: `flai touches suggest` lists it in 94 of the matching commits; a template convention change takes an entry |
| `design/issues/I-0352-…md`, `design/issues/summary.md` | goal and criterion 2: `flai issue close` writes both |

`flai touches suggest` also lists `design/system/flai-cli.md` and `docs/users/flai-reference.md`. I left them out: `flai-cli.md` lists the guard by package, not by file, and the guard command's flags do not change. `task_done` widens the touches if the story's agent needs them.

The story names no `after`. It reads no other story's output.

Forecast 45m, delivery 2026-11-13T20:51:00Z. `flai forecast` gives 27m: a median of 107 s per unit of size over 179 done large-band improvement stories, times size 15. I raised it to 45m. The work spans two harness readers, a new guard rule with its tests, and a template convention change, over four tasks that each close on their own. The delivery is flai's, moved by the same 18 minutes. It assumes the operator finalizes the draft at once, since the story is 2nd in the pull order behind S-0760.

Cost of delay 2.50 USD a week, as `flai cod` gives it from flai's input: 1m of time lost per 168h cycle at 150 USD an hour. It stands as is: I-0352 has one instance, costing 1m.
