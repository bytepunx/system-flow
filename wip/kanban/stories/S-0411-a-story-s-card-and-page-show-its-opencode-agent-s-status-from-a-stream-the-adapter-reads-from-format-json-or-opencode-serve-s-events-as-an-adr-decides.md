---
id: S-0411
type: story
nature: feature
title: A story's card and page show its OpenCode agent's status from a stream the adapter reads, from --format json or opencode serve's events as an ADR decides
status: backlog
parent: E-0023
owner: alex
created: 2026-10-16T05:19:30Z
updated: 2026-10-21T12:51:45Z
transitions: []
tags: [cli]
topics: [agents, cli]
touches: [flai/internal/serve/stream.go, flai/internal/serve/stream_test.go, flai/internal/serve/stream_opencode.go, flai/internal/serve/stream_opencode_test.go, flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, design/adrs, design/system/agent-adapters.md, design/system/flai-cli.md, docs/operators/index.md]
after: [S-0410]
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
    - kind: orchestrator
      seconds: 148
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 7
          output: 135
          cache_read: 2226442
          cache_write: 1761
          cost: 0.5586
cost_of_delay:
  value: 53.85
  by: planner-E-0023
  at: 2026-10-20T21:13:44Z
forecast:
  duration: 35m
  delivery: 2026-10-21T18:21:00Z
  basis: "Its own forecast of 35m; 7th in the pull order with an in-progress limit of 5, behind S-0461, S-0471, S-0473, S-0479, S-0241, S-0406, S-0407, S-0408, S-0409 and S-0410."
  by: flai
  at: 2026-10-21T12:51:45Z
finalized:
  by: orchestrator
  at: 2026-10-20T21:14:09Z
---
# S-0411 A story's card and page show its OpenCode agent's status from a stream the adapter reads, from --format json or opencode serve's events as an ADR decides

## Goal

The operator watching the board sees what a story's agent on OpenCode is doing, as they see one on Claude Code: its last step, the tool it is in, its sub-agents, and whether it is waiting. `flai serve` builds the card's and the page's stream through the adapter of the run's agent instead of parsing Claude Code's stream-json outright, and the OpenCode stream comes from `--format json` or from `opencode serve` and its event stream, whichever S-0406's finding on `--attach` and this story's ADR choose.

## Acceptance criteria

- [ ] An ADR, refining ADR-0130, records where an OpenCode agent's status comes from, `opencode run --format json` or `opencode serve`'s event stream with `--attach`, and why, including what each costs when several stories run on one host.
- [ ] `flai serve`'s stream (`Stream`, `PlanStream`, `RoleStream`) reads entries through the run's adapter, and a Claude Code run's entries are the same as before.
- [ ] A story on `opencode` shows on its card and page its last step, the tool it is in, its sub-agents, and its waits, in the same entry shape, so the dashboard needs no change.
- [ ] The stream is tested on recorded OpenCode output, and, where the ADR chooses `opencode serve`, on a recorded event stream.
- [ ] `design/system/flai-cli.md`, `agent-adapters.md`, and `docs/operators/index.md` § Harnesses say where the status of a story on OpenCode comes from.

## Tasks
- T-1712 An ADR records whether an OpenCode agent's status comes from opencode run --format json or from opencode serve's event stream, and what each costs on a busy host
- T-1713 flai serve's stream reads its entries through the run's adapter, and a Claude Code run's entries stay the same
- T-1714 An OpenCode stream reader gives a story's card and page its agent's last step, tool, sub-agents, and waits from the source the ADR chose
- T-1715 The design and operator docs say where the status of a story on OpenCode comes from

## Notes

- Drafted by the planner for E-0023, after S-0410, whose recorded output the stream tests reuse.
- Should the ADR choose `opencode serve` with `--attach`, `flai serve` gains a long-running OpenCode server per host, and this story grows: say so on the plan thread before building it, and split the server's life cycle into a story of its own.

### Planning

- **Touches.** None were declared. Every touch is the planner's:
  - **Layout:** `flai/internal/serve/stream.go`, which the explorer found parses stream-json outright and says so in its header, with its test; `stream_opencode.go`, new, for OpenCode's entries; `agents.go`, whose `StoryActivity` and `Activity` feed the card, with its test. The dashboard reads the same entry shape, so no `flaiover/` path is touched.
  - **Design:** `agent-adapters.md`, `flai-cli.md`, and `docs/operators/index.md` for criterion 5.
  - **Folder kept:** `design/adrs`, because the ADR's file name is not known until `flai adr new`; it is a shared path, so it holds nothing.
  - **Co-change:** `flai touches suggest` lists `docs/users/flai.md` in half the seed commits, from `agents.go`'s long history; this story changes no user-facing behaviour, so it is not added. The rest are hub documents.
- **Forecast.** 35m against flai's 24m. Moving `stream.go`'s parsing behind a per-harness reader with Claude Code's entries unchanged, plus an ADR, adds about half. If the ADR chooses `opencode serve`, the story grows and the Notes say to split it. The delivery, 2026-10-16T16:51Z, is flai's play-out moved by the raise times the cycle factor.
- **Cost of delay.** 53.85 USD a week. E-0023 carries `penalty_per_week: 500`, the operator's answer to TH-0467 question 1, confirmed on TH-0470 on 2026-10-20 (as on E-0019). `flai cod` shares it among the epic's seven open stories by forecast, 35m of 5h25m. Kept as flai gives it.
