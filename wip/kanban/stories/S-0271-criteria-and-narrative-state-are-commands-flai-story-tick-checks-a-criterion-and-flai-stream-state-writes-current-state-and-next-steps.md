---
id: S-0271
type: story
nature: improvement
title: "Criteria and narrative state are commands: flai story tick checks a criterion and flai stream state writes Current state and Next steps"
status: ready
parent: E-0017
owner: alex
created: 2026-10-05T01:35:30Z
updated: 2026-10-06T23:33:22Z
transitions:
  - to: ready
    at: 2026-10-06T22:48:17Z
    by: alex
tags: [cli, mcp]
topics: [automation, mcp, hostapi, conventions]
touches: [flai/cmd/stream.go, flai/cmd/stream_state_test.go, flai/internal/workitem/narrative.go, flai/internal/workitem/narrative_test.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/stream.go, flai/internal/mcpserver/stream_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/internal/check/check.go, flai/internal/check/check_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, design/conventions/tooling.md, design/conventions/work-management.md, design/conventions/session-start.md, template/root/design/conventions/tooling.md, template/root/design/conventions/work-management.md, template/root/design/conventions/session-start.md, template/CHANGELOG.md, design/system/agent-narrative.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/settings.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  value: 16
  by: planner-E-0017
  at: 2026-10-06T11:36:20Z
forecast:
  duration: 39m
  delivery: 2026-10-07T01:35:00Z
  basis: "Its own forecast of 39m; 7th in the pull order with an in-progress limit of 3, behind S-0300, S-0302, S-0273, S-0303, S-0228, S-0269 and S-0270."
  by: flai
  at: 2026-10-06T23:33:22Z
finalized:
  by: alex
  at: 2026-10-06T22:48:13Z
---
# S-0271 Criteria and narrative state are commands: flai story tick checks a criterion and flai stream state writes Current state and Next steps

## Goal

Agents tick acceptance criteria with `sed -i` on the story file (119 turns across 85 stories) and rewrite a narrative's `## Current state` and `## Next steps` with Edit or Write (111 turns across 44 stories), each a model turn that reads the file first. `flai story tick S-nnnn <n>` (and `--untick`) checks the nth criterion and refuses one that does not exist; `flai stream state S-nnnn --current "<text>" --next "<text>"` replaces the two sections, leaving the rest of the narrative alone, and appends nothing. Both are `item_tick` and `stream_state` over MCP and `item.tick` and `stream.state` on the host channel, so the dashboard's story page can tick a criterion too. `flai check` then knows the sections' shape.

## Acceptance criteria
- [ ] `flai story tick` and `flai stream state` exist with the behaviour above, as text and `--json`, and refuse a story that is not in progress or review
- [ ] The same operations exist over MCP and on the host channel, and the dashboard's story page ticks a criterion through them
- [ ] The conventions, the template's copies, the harness prompt, `design/system/flai-cli.md`, `agent-narrative.md`, and the user guide send the agent to them instead of editing the files

## Tasks
- T-1054 The workitem package replaces a narrative's Current state and Next steps and leaves the rest of it alone
- T-1055 flai stream state writes a narrative's Current state and Next steps, as text and --json
- T-1056 The MCP tool stream_state writes a narrative's Current state and Next steps
- T-1058 The host channel method stream.state writes a narrative's Current state and Next steps
- T-1059 flai check knows the shape of a narrative's Current state and Next steps
- T-1061 The conventions, their template copies, and the harness prompt send the agent to flai stream state
- T-1063 The design and the user guide describe flai stream state, stream_state, and stream.state

## Notes

From the epic's log classification. `flai task done` (its sibling story) may take `--tick <n>` once this exists.

### Planning

The ticking half of the story already exists. ADR-0089 (S-0282) added:

- the command `flai criteria tick` and `untick` (`flai/cmd/criteria.go`);
- the MCP tool `criteria_tick` (`flai/internal/mcpserver/criteria.go`);
- the host method `item.criteria` (`flai/internal/hostapi/writes.go`), through which the dashboard's item page ticks.

alex answered TH-0207 on 2026-10-06 with option (a): that ticking meets the tick half of criteria 1 and 2. No `flai story tick`, `item_tick`, or `item.tick` is added. When the story's agent ticks those two criteria, it records in these Notes that the existing names deliver them. So the new work here is `flai stream state` and the check rule. The story's words are left as the operator wrote them.

Tasks, in four layers:

1. T-1054, the shared writer.
2. T-1055 (the command), T-1056 (MCP), and T-1059 (the check rule).
3. T-1058 (the host method, which runs the command) and T-1061 (the conventions and the harness prompt).
4. T-1063 (the design and the docs).

All the touches are files, none is a folder. `flai touches suggest S-0271` lists co-changes from 409 of 1050 commits, but none of its suggestions is predicted here: they are the board's general churn files. Each touch is listed below with where it came from.

- `flai/cmd/stream.go` and `flai/cmd/stream_state_test.go`: declared (layout). The stream command and a new test.
- `flai/internal/workitem/narrative.go`: declared (co-change). Narrative sections are parsed here.
- `flai/internal/workitem/narrative_test.go`: layout, added. The test for T-1054.
- `flai/internal/mcpserver/folder.go`: declared (layout). The tool is registered here.
- `flai/internal/mcpserver/stream.go` and `stream_test.go`: layout, added. New files modelled on `criteria.go`.
- `flai/internal/hostapi/writes.go` and `writes_test.go`: declared (layout). They hold `stream.log` and `stream.answer`.
- `flai/internal/check/check.go` and `check_test.go`: layout, added. This is the goal's "`flai check` then knows the sections' shape". `narrativeSections` lives here.
- `flai/internal/harness/harness.go` and `harness_test.go`: declared (design, criterion 3).
- The conventions:
  - `design/conventions/tooling.md`, `work-management.md`, and `session-start.md`: declared (design).
  - Their `template/root` copies and `template/CHANGELOG.md`: declared (design).
- `design/system/agent-narrative.md`, `flai-cli.md`, `docs/users/flai.md`, and `flai-reference.md`: declared (co-change and criterion 3).
- `docs/operators/settings.md`: layout, added. `scripts/flai-reference.sh` regenerates its flag index, which gains `--current` and `--next`.
- Left out: `flai/cmd/criteria.go` and the flaiover item page. Ticking already exists there, and TH-0207 adds no alias.

Forecast: flai's 39m stands. It is the median of 83 s per unit over 25 done large improvement stories, times size 28 (3 criteria and 25 touches, up from 22 with the six touches added). Seven small tasks, each modelled on an existing sibling (`stream log`, `criteria_tick`, `stream.log`), fit that rate. The delivery is flai's, and it replans it as the board moves. It was 2026-10-07T02:01Z after S-0264 was accepted.

Cost of delay: 16 USD a week stands, against flai's 77.75. flai shares E-0017's 900 USD a week by forecast duration. The epic's planner shared it by the turns each story removes, and the operator resolved that plan (TH-0176). This story removes the 111 narrative edits. The 119 `sed` ticks came before `criteria_tick` existed, so they are not counted.
