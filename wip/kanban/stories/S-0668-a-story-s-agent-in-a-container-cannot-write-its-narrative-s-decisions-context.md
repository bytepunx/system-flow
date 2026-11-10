---
id: S-0668
type: story
nature: remediation
title: A story's agent in a container cannot write its narrative's Decisions, Context, or Open questions, since wip is read-only there and flai stream state writes only Current state and Next steps
status: backlog
owner: alex
created: 2026-11-10T08:34:42Z
updated: 2026-11-10T08:54:23Z
transitions: []
tags: [flai, mcp, narrative]
topics: [narrative, mcp, harness]
touches: [flai/internal/workitem/narrative.go, flai/internal/workitem/narrative_test.go, flai/cmd/stream.go, flai/cmd/stream_state_test.go, flai/internal/mcpserver/stream.go, flai/internal/mcpserver/stream_test.go, flai/internal/hostapi/writes_streams.go, flai/internal/hostapi/writes_test.go, design/system/agent-narrative.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/settings.md, design/conventions/tooling.md, template/root/design/conventions/tooling.md, template/CHANGELOG.md, design/issues/I-0277-a-story-s-agent-in-a-container-cannot-write-its-narrative-s-decisions-context.md, design/issues/summary.md]
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
      seconds: 21
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 12
          output: 44
          cache_read: 1035481
          cache_write: 28453
          cost: 0.0712
        - model: claude-opus-5-5
          input: 94
          output: 1581
          cache_read: 6510298
          cache_write: 273653
          cost: 1.7397
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-11-10T08:34:42Z
  value: 12.5
  by: planner-S-0668
  at: 2026-11-10T08:50:21Z
forecast:
  duration: 45m
  delivery: 2026-11-10T15:33:00Z
  basis: "Its own forecast of 45m; 14th in the pull order with an in-progress limit of 5, behind S-0680, S-0682, S-0684, S-0685, S-0687, S-0674, S-0673, S-0672, S-0671, S-0241, S-0659, S-0660, S-0661, S-0662, S-0663, S-0665, S-0666 and S-0667."
  by: flai
  at: 2026-11-10T08:54:23Z
finalized:
  by: alex
  at: 2026-11-10T08:51:18Z
---
# S-0668 A story's agent in a container cannot write its narrative's Decisions, Context, or Open questions, since wip is read-only there and flai stream state writes only Current state and Next steps

## Goal

This story remediates [I-0277](../../../design/issues/I-0277-a-story-s-agent-in-a-container-cannot-write-its-narrative-s-decisions-context.md), "A story's agent in a container cannot write its narrative's Decisions, Context, or Open questions, since wip is read-only there and flai stream state writes only Current state and Next steps". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0277 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0277 is closed with `flai issue close I-0277 --reason` saying what fixed it

## Tasks
- T-2944 SetStreamState replaces Context and appends Decisions and Open questions bullets, keeping the threads block
- T-2947 flai stream state takes --context, --decision, and --question, and the reference and flag index list them
- T-2949 The MCP tool stream_state takes context, decisions, and questions, with a test of a container's session writing them
- T-2950 The host channel's stream.state passes context, decisions, and questions to flai stream state
- T-2951 The design, the user guide, and tooling.md in both copies say stream state writes Context, Decisions, and Open questions, and I-0277 is closed

## Notes

Cost of delay inputs set by flai from I-0277. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-07T06:20:58Z, 3.1 days before this story; under one cycle counts as one).

### Planning

Proposed fix, from I-0277's one instance (S-0556 put its decisions in T-2459's Notes). The MCP server runs on the host, and a container's agent reaches it over HTTP (`flai/cmd/mcp.go`). So `stream_state` already writes the main checkout's `wip/` from a container. The fix widens what it writes, rather than mounting `wip/` read-write, which ADR-0253 refused.

- `stream_state`, `flai stream state`, and the host channel's `stream.state` gain three texts:
  - `context` replaces what is under `## Context`.
  - `decisions` appends each text given as a bullet under `## Decisions`, in the shape the digest reads (`flai/internal/digest/sources.go`).
  - `questions` appends each text as a dated bullet under `## Open questions`, above the `<!-- threads:start -->` block, which stays as it is.
- Appending, not replacing, keeps what an earlier session recorded. `flai stream answer` still moves an answered question to Decisions. A question written this way holds the move to review as a hand-written one does (`flai/internal/workitem/rules.go`).
- The narrative's other sections, its stamps, the lint, the heading refusal, and the refusals by story state stay as `SetStreamState` has them.

Touches, all files, no folder:

- Declared: none before this plan.
- Design: `flai/internal/workitem/narrative.go` and its test hold `SetStreamState`. `flai/cmd/stream.go` and `flai/cmd/stream_state_test.go` hold the CLI. `flai/internal/mcpserver/stream.go` and its test hold the tool. `flai/internal/hostapi/writes_streams.go` and `writes_test.go` hold the host channel's `stream.state`. `design/system/agent-narrative.md` § Narrative structure says the three sections are edited by hand. `docs/users/flai.md` § Agent narratives and its MCP row say the same.
- Layout: `docs/users/flai-reference.md` and `docs/operators/settings.md` hold the flag rows `make flai-reference` writes. `design/conventions/tooling.md` and `template/root/design/conventions/tooling.md` (line 17) say Context, Decisions, and Open questions are edited by hand, and `template/CHANGELOG.md` records a template change, as S-0271 did. The I-0277 issue file and `design/issues/summary.md` change when the issue is closed.
- Co-change: `design/system/flai-cli.md`, 23% of the commits that changed these paths, whose command table (`flai stream state`) and `flai mcp` tool list describe the call. The rest of the list (`docs/operators/index.md`, `flai/internal/hostapi/writes.go`, `flai/internal/harness/harness.go`, and others) changes nothing this story needs.

Figures:

- Forecast 45m, up from `flai forecast`'s 34m (106 s per unit of size over 63 large remediation stories, size 19). S-0271, which built `stream state` across the same CLI, MCP tool, and host channel, took 2,477 s, about 41m. This story adds the care the Open questions threads block needs. Delivery is flai's, moved by the 11 minutes added.
- Cost of delay 12.50 USD a week, as `flai cod` works it from the inputs flai set from I-0277: 5m lost per 168h cycle at 150 USD an hour. It stands. The inputs stay as flai set them.
- No `after`: the story reads no other open story's output.
