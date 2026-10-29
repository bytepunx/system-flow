---
id: S-0410
type: story
nature: feature
title: A story on OpenCode has its usage read from OpenCode's JSON events, priced by the harness or its gateway's spend log, and reaches its provider through an entry derived from the host's
status: backlog
parent: E-0023
owner: alex
created: 2026-10-16T05:18:34Z
updated: 2026-10-29T00:05:49Z
transitions: []
tags: [cli]
topics: [agents, cli]
touches: [flai/internal/usage/opencode.go, flai/internal/usage/opencode_test.go, flai/internal/usage/testdata/opencode/run.jsonl, flai/internal/usage/testdata/opencode/resumed.jsonl, flai/internal/harness/opencode.go, flai/internal/harness/opencode_test.go, flai/internal/harness/opencodeconfig.go, flai/internal/harness/opencodeconfig_test.go, flai/internal/serve/usage.go, flai/internal/serve/usage_test.go, flai/internal/serve/activity.go, flai/internal/serve/activity_test.go, flai/internal/serve/gateway_test.go, design/system/metrics.md, design/system/flai-cli.md, docs/users/flai.md, docs/operators/index.md, flai/internal/harness/harness_test.go]
after: [S-0409]
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
      seconds: 149
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 7
          output: 135
          cache_read: 2226442
          cache_write: 1761
          cost: 0.5586
cost_of_delay:
  value: 76.92
  by: planner-E-0023
  at: 2026-10-20T21:13:43Z
forecast:
  duration: 50m
  delivery: 2026-10-29T08:18:00Z
  basis: "Its own forecast of 50m; 10th in the pull order with an in-progress limit of 5, behind S-0535, S-0532, S-0537, S-0534, S-0536, S-0241, S-0406, S-0407, S-0408 and S-0409."
  by: flai
  at: 2026-10-29T00:05:49Z
finalized:
  by: orchestrator
  at: 2026-10-20T21:14:08Z
---
# S-0410 A story on OpenCode has its usage read from OpenCode's JSON events, priced by the harness or its gateway's spend log, and reaches its provider through an entry derived from the host's

## Goal

A story worked on OpenCode carries the same `usage` as one worked on Claude Code: tokens per model, seconds, turns, and a cost with its `priced_by`. The adapter's `Reader` turns OpenCode's `--format json` events into flai's neutral usage events, so nothing downstream reads OpenCode's format. A story with a provider reaches it through an OpenCode provider entry flai derives from the host's `providers` map (ADR-0129), and is priced from the gateway's spend log, as a Claude Code story is (ADR-0132).

## Acceptance criteria

- [ ] `flai/internal/usage` has an OpenCode reader that turns `step_start`, `tool_use`, `text`, and `step_finish` into neutral events: session start, a model call with its tokens (input, output, reasoning, cache read and write), its model, and its cost when reported, the sub-agent call that started it, and the run's end. It is tested on recorded output from the known version.
- [ ] The adapter's `Reader` returns it and its `Capabilities` say it reports cost. `flai serve` reads a story's, and a strategic agent's, usage and activity through the adapter of the agent that ran, in place of the Claude Code reader it names today in `serve/usage.go` and `serve/activity.go`.
- [ ] A story on `opencode` with no provider is `priced_by: harness`, from `step_finish`'s cost. One with a provider is priced from the gateway's spend log and is `priced_by: gateway`, with OpenCode's own figure kept as tokens only.
- [ ] For a story with a provider, the configuration carries an OpenCode provider entry derived from the host's entry: the package for its `api` (`openai-chat`, `openai-responses`, or `anthropic-messages`), its `base_url`, the key as `{env:<key_env>}`, and its `models`. No key is written to the file, a log, or a work item. A provider whose `api` OpenCode cannot speak is refused at start, naming it.
- [ ] A gateway test runs a story's agent on OpenCode through OpenRouter only where `OPENROUTER_API_KEY` and the known `opencode` are both present, asserts `priced_by: gateway` and a cost above zero, and is skipped, saying which is missing, otherwise.
- [ ] `design/system/metrics.md` names OpenCode's events where it names the harness's, and `design/system/flai-cli.md`, `docs/users/flai.md`, and `docs/operators/index.md` § Harnesses say how a story on OpenCode is measured and priced.

## Tasks
- T-1706 flai/internal/usage reads OpenCode's JSON events into neutral usage events, tested on recorded output from the known version
- T-1707 The opencode configuration carries a provider entry derived from the host's providers map, with the key named and never written, and refuses an api OpenCode cannot speak
- T-1708 The opencode adapter returns the OpenCode usage reader and states that it reports cost
- T-1709 flai serve measures usage and activity through the adapter of the agent that ran, and prices an OpenCode run by the harness or by its gateway's spend log
- T-1710 A gateway test runs a story's agent on OpenCode through OpenRouter only where the key and the known opencode are present, and asserts a gateway-priced cost
- T-1711 metrics.md names OpenCode's events, and the design, user, and operator docs say how a story on OpenCode is measured and priced

## Notes

- Drafted by the planner for E-0023, after S-0409, in the epic's order. It needs only S-0406's adapter; the holds on the shared adapter files would run it after the roles story in any case.
- S-0406's finding on `step_finish` decides how a run joins the gateway's spend log: a generation ID per call where OpenCode gives one, else the key's spend before and after the run.

### Planning

- **Touches.** None were declared. Every touch is the planner's:
  - **Layout:** `flai/internal/usage/opencode.go`, new, beside `claudecode.go`, with recorded output under `testdata/opencode/`, two files named for a first run and a resumed one; `flai/internal/harness/opencode.go` and `opencodeconfig.go` for the reader, the capability, and the provider entry, and `harness_test.go` for the capabilities test; `flai/internal/serve/usage.go` and `activity.go`, where the explorer found the Claude Code reader named outright; `flai/internal/serve/gateway_test.go`, where E-0019's OpenRouter tests and their skip live.
  - **Design:** `design/system/metrics.md`, which ADR-0130 says changes where the events are named; the docs for criterion 6.
  - **Co-change:** `flai touches suggest` lists only hub documents, none of which this story's criteria reach; none was added.
- **Forecast.** 50m against flai's 37m. Recorded fixtures, replacing the Claude Code reader named outright in two serve files without changing Claude Code runs, and a gateway test, as S-0357's pricing work had, add a third. The delivery, 2026-10-16T16:12Z, is flai's play-out moved by the raise times the cycle factor.
- **Cost of delay.** 76.92 USD a week. E-0023 carries `penalty_per_week: 500`, the operator's answer to TH-0467 question 1, confirmed on TH-0470 on 2026-10-20 (as on E-0019). `flai cod` shares it among the epic's seven open stories by forecast, 50m of 5h25m. Kept as flai gives it.
