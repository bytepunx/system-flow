---
id: S-0271
type: story
nature: improvement
title: "Criteria and narrative state are commands: flai story tick checks a criterion and flai stream state writes Current state and Next steps"
status: backlog
parent: E-0017
owner: alex
created: 2026-10-05T01:35:30Z
updated: 2026-10-06T18:11:07Z
transitions: []
tags: [cli, mcp]
topics: [automation, mcp, hostapi, conventions]
touches: [flai/cmd/stream.go, flai/cmd/stream_state_test.go, flai/internal/workitem/narrative.go, flai/internal/mcpserver/folder.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, design/conventions/tooling.md, design/conventions/work-management.md, design/conventions/session-start.md, template/root/design/conventions/tooling.md, template/root/design/conventions/work-management.md, template/root/design/conventions/session-start.md, template/CHANGELOG.md, design/system/agent-narrative.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
cost_of_delay:
  value: 16
  by: planner-E-0017
  at: 2026-10-06T11:36:20Z
forecast:
  duration: 33m
  delivery: 2026-10-07T07:31:00Z
  basis: "Its own forecast of 33m; 33rd in the pull order with an in-progress limit of 3, behind S-0226, S-0223, S-0224, S-0227, S-0284, S-0229, S-0278, S-0295, S-0296, S-0212, S-0213, S-0214, S-0215, S-0216, S-0228, S-0232, S-0233, S-0234, S-0235, S-0236, S-0237, S-0238, S-0239, S-0241, S-0245, S-0246, S-0251, S-0254, S-0261, S-0264, S-0265, S-0269 and S-0270."
  by: flai
  at: 2026-10-06T18:11:07Z
---
# S-0271 Criteria and narrative state are commands: flai story tick checks a criterion and flai stream state writes Current state and Next steps

## Goal

Agents tick acceptance criteria with `sed -i` on the story file (119 turns across 85 stories) and rewrite a narrative's `## Current state` and `## Next steps` with Edit or Write (111 turns across 44 stories), each a model turn that reads the file first. `flai story tick S-nnnn <n>` (and `--untick`) checks the nth criterion and refuses one that does not exist; `flai stream state S-nnnn --current "<text>" --next "<text>"` replaces the two sections, leaving the rest of the narrative alone, and appends nothing. Both are `item_tick` and `stream_state` over MCP and `item.tick` and `stream.state` on the host channel, so the dashboard's story page can tick a criterion too. `flai check` then knows the sections' shape.

## Acceptance criteria
- [ ] `flai story tick` and `flai stream state` exist with the behaviour above, as text and `--json`, and refuse a story that is not in progress or review
- [ ] The same operations exist over MCP and on the host channel, and the dashboard's story page ticks a criterion through them
- [ ] The conventions, the template's copies, the harness prompt, `design/system/flai-cli.md`, `agent-narrative.md`, and the user guide send the agent to them instead of editing the files

## Tasks

## Notes

From the epic's log classification. `flai task done` (its sibling story) may take `--tick <n>` once this exists.

### Planning

Ticking already exists, so this story's new work is `flai stream state`. ADR-0089 added `flai criteria tick` and `untick` (`flai/cmd/criteria.go`), the MCP tool `criteria_tick` (`flai/internal/mcpserver/criteria.go`), and the host method `item.criteria` (`flai/internal/hostapi/writes.go`). The dashboard's item page ticks through `flaiover/src/routes/api/items/[id]/criteria`. The harness prompt and `work-management.md` already send the agent to them. The plan thread proposes narrowing the goal and criteria to `flai stream state`. The words are left as the operator wrote them.

Touches, none declared before. `flai touches suggest S-0271` was seeded with `flai/cmd/criteria.go` and `flai/cmd/stream.go`, which 10 commits changed. They are predicted for `flai stream state`:

- `flai/cmd/stream.go`, `flai/cmd/stream_state_test.go`: layout. The stream command and a new test.
- `flai/internal/workitem/narrative.go`: co-change (3 of 10). Narrative sections are parsed here.
- `flai/internal/mcpserver/folder.go`, `flai/internal/hostapi/writes.go`, `writes_test.go`: layout. The `stream_state` tool and the `stream.state` method.
- `flai/internal/harness/harness.go`, `harness_test.go`: design (criterion 3).
- `design/conventions/tooling.md`, `work-management.md`, `session-start.md`, their `template/root` copies, and `template/CHANGELOG.md`: design. `tooling.md` says the summary sections are edited by hand. The other two say when to rewrite `## Current state` and `## Next steps`.
- `design/system/agent-narrative.md`, `flai-cli.md`, `docs/users/flai.md`, `flai-reference.md`: co-change (2, 5, 3, and 2 of 10) and criterion 3.
- Left out: `flai/cmd/criteria.go` and the flaiover item page, whose ticking already exists. Add them back if the operator keeps `flai story tick` as an alias.

Forecast: flai gave 33m (89 s per unit over 21 done large improvement stories, times size 22), and it stands. Narrowed to `flai stream state` it would be smaller. The delivery, 2026-10-07T00:39Z, is flai's.

Cost of delay: 16 USD a week, against flai's 86.84. This is E-0017's 900 USD a week shared by the turns each story removes. This one removes the 111 narrative edits. The 119 `sed` ticks were made before `criteria_tick` existed, so they are not counted.
