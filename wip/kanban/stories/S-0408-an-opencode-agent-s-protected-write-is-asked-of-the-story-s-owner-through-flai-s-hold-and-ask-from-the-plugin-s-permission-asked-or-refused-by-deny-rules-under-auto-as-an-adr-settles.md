---
id: S-0408
type: story
nature: feature
title: An OpenCode agent's protected write is asked of the story's owner through flai's hold-and-ask from the plugin's permission.asked, or refused by deny rules under --auto, as an ADR settles
status: cancelled
parent: E-0023
owner: alex
created: 2026-10-16T05:16:14Z
updated: 2026-11-07T21:53:43Z
transitions:
  - to: ready
    at: 2026-11-07T20:43:08Z
    by: alex
  - to: in-progress
    at: 2026-11-07T21:46:09Z
    by: agent-S-0408
  - to: cancelled
    at: 2026-11-07T21:53:43Z
    by: alex
tags: [cli, template]
topics: [agents, cli, template]
touches: [flai/cmd/guard.go, flai/cmd/guard_test.go, flai/internal/guard/call.go, flai/internal/guard/call_test.go, flai/internal/guard/opencode.go, flai/internal/guard/opencode_test.go, flai/internal/harness/opencode.go, flai/internal/harness/opencode_test.go, flai/internal/harness/opencodeconfig.go, flai/internal/harness/opencodeconfig_test.go, ".opencode/plugins/flai-guard.js", template/root/.opencode/plugins/flai-guard.js, flai/tests/integration/opencode_test.go, design/adrs, design/system/agent-adapters.md, design/system/flai-cli.md, docs/operators/index.md, docs/users/flai.md, template/CHANGELOG.md, flai/internal/harness/harness_test.go, design/issues/I-0306-adr-new-over-mcp-refuses-an-adr-for-wip-overlap-findings-between-other-stories.md, design/issues/summary.md]
after: [S-0407]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1701
  priced_by: harness
  turns:
    - day: 2026-11-07
      ceremony: 4
      hand_edits: 1
      work: 36
  models:
    - model: claude-opus-5-5
      input: 192
      output: 69054
      cache_read: 27556070
      cache_write: 1132273
      cost: 13.5799
  strategic:
    - kind: orchestrator
      seconds: 463
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 83
          output: 1136
          cache_read: 13232302
          cache_write: 91199
          cost: 3.3412
cost_of_delay:
  value: 76.92
  by: planner-E-0023
  at: 2026-10-20T21:13:40Z
forecast:
  duration: 50m
  delivery: 2026-11-07T22:45:00Z
  basis: "Its own forecast of 50m; 1st in the pull order with an in-progress limit of 5, with nothing ahead of it."
  by: flai
  at: 2026-11-07T21:45:57Z
finalized:
  by: orchestrator
  at: 2026-10-20T21:14:05Z
---
# S-0408 An OpenCode agent's protected write is asked of the story's owner through flai's hold-and-ask from the plugin's permission.asked, or refused by deny rules under --auto, as an ADR settles

## Goal

A story's agent on OpenCode asks before it writes a protected path, as one on Claude Code does through `permission_prompt`. Either the plugin answers OpenCode's `permission.asked` with the owner's answer from flai's hold-and-ask (`ask.Asker.Hold`, ADR-0124's thread and retry), or the adapter runs `--auto` behind `deny` rules for the protected list and the agent is told to ask on a thread. An ADR settles which from what S-0406 found, so that a story on OpenCode no longer needs `auto-approve` or `--deny-protected` to start.

## Acceptance criteria

- [x] An ADR, refining ADR-0124 and ADR-0130, records which way OpenCode asks, OpenCode's bound for a held call, and why the other way was not taken.
- [ ] Under hold-and-ask: the configuration sets `ask` for every path the protected list holds for `opencode`; the plugin hands `permission.asked` to `flai guard`, which holds it through `ask.Asker.Hold` within the bound, opening or reusing the story's permission thread as ADR-0124 says; and the plugin replies allow or deny with the owner's answer. Under `--auto` with deny rules: the adapter passes `--auto`, writes the `deny` rules, and the refusal names the thread the agent should open.
- [ ] The adapter's `Capabilities` say `Ask: true` under hold-and-ask, and `flai serve` then starts a story on `opencode` without `auto-approve` or `--deny-protected`.
- [ ] The integration test shows a protected write held and then allowed by an owner's answer, or refused by the deny rule, on the known `opencode` version, and is skipped, saying so, without it.
- [ ] `docs/operators/index.md` § Harnesses, `docs/users/flai.md`, `design/system/flai-cli.md`, and `agent-adapters.md` say how a protected write is asked on OpenCode, and `template/CHANGELOG.md` names the plugin's change.

## Tasks
- T-1695 An ADR records how an OpenCode agent asks before a protected write, OpenCode's bound for a held call, and the way not taken
- T-1696 flai guard holds an OpenCode permission request through ask.Asker.Hold within the ADR's bound and answers allow or deny with the owner's word
- T-1697 The opencode configuration sets ask or deny on the protected list as the ADR says, and the adapter states Ask true under hold-and-ask
- T-1698 The flai-guard plugin answers OpenCode's permission.asked with flai guard's allow or deny, in this repository and the template
- T-1699 The OpenCode integration test shows a protected write held and allowed by the owner's answer, or refused by the deny rule
- T-1700 The operator, user, and design docs say how an OpenCode agent's protected write is asked or refused

## Notes

- Drafted by the planner for E-0023, after S-0407, whose plugin this story extends.
- `guard.Hold` has no caller outside its tests today; it is the natural entry for a harness whose plugin can wait on `flai guard`.

### Planning

- **Touches.** None were declared. Every touch is the planner's, checked again on 2026-10-17:
  - **Layout:** `flai/cmd/guard.go` and `flai/internal/guard/call.go`, where `guard.Hold` and the `Ask` verdict live, with their tests; `flai/internal/guard/opencode.go`, S-0407's reader, for the `permission.asked` input; `flai/internal/harness/opencode.go` and `opencodeconfig.go` for the capability and the `ask` or `deny` rules, and `harness_test.go` for the capabilities test; the two plugin copies; `flai/tests/integration/opencode_test.go`. The OpenCode files do not exist yet: S-0405 to S-0407 add them. `flai/internal/ask/ask.go` is reused unchanged.
  - **Design:** the docs, `agent-adapters.md`, and `template/CHANGELOG.md` for criterion 5.
  - **Folder kept:** `design/adrs`, because the ADR's file name is not known until `flai adr new`; it also covers the index `flai adr new` updates. It is a shared path, so it holds nothing.
  - **Co-change:** `flai touches suggest` lists hub documents and `harness_test.go`; only `harness_test.go` is added. Left out with reasons: `flai/internal/harness/harness.go`, since `Capabilities` already has `Ask` and `flai/internal/serve/agents.go` already lifts the start refusal on it; `template/template.yaml`, which only a publish bumps; `docs/users/flai-reference.md`, needed only if `flai guard` gains a flag, which T-1696's Notes then add.
- **Forecast.** 50m against flai's 45m on 2026-10-20. The ADR and the first caller of `guard.Hold`, whose bounded waits and thread reuse need timed tests, add about a ninth. The delivery, 2026-10-21T00:23Z, is flai's play-out moved by the 5m raise.
- **Cost of delay.** 76.92 USD a week. E-0023 carries `penalty_per_week: 500`, the operator's answer to TH-0467 question 1, confirmed on TH-0470 on 2026-10-20 (as on E-0019). `flai cod` shares it among the epic's seven open stories by forecast, 50m of 5h25m. Kept as flai gives it. This story's plan thread, TH-0490, asked the same input.
- 2026-11-07T21:53:43Z: moved to cancelled: E-0023 cancelled: OpenCode isn't going to receive first-class support given its aversion to supporting before hooks for tool calls (amongst other things)
