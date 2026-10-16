---
id: S-0408
type: story
nature: feature
title: An OpenCode agent's protected write is asked of the story's owner through flai's hold-and-ask from the plugin's permission.asked, or refused by deny rules under --auto, as an ADR settles
status: backlog
parent: E-0023
owner: alex
created: 2026-10-16T05:16:14Z
updated: 2026-10-16T07:30:22Z
transitions: []
tags: [cli, template]
topics: [agents, cli, template]
touches: [flai/cmd/guard.go, flai/cmd/guard_test.go, flai/internal/guard/call.go, flai/internal/guard/call_test.go, flai/internal/guard/opencode.go, flai/internal/guard/opencode_test.go, flai/internal/harness/opencode.go, flai/internal/harness/opencode_test.go, flai/internal/harness/opencodeconfig.go, flai/internal/harness/opencodeconfig_test.go, ".opencode/plugins/flai-guard.js", template/root/.opencode/plugins/flai-guard.js, flai/tests/integration/opencode_test.go, design/adrs, design/system/agent-adapters.md, design/system/flai-cli.md, docs/operators/index.md, docs/users/flai.md, template/CHANGELOG.md, flai/internal/harness/harness_test.go]
after: [S-0407]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
forecast:
  duration: 50m
  delivery: 2026-10-16T16:40:00Z
  basis: "Its own forecast of 50m; 18th in the pull order with an in-progress limit of 5, behind S-0426, S-0241, S-0392, S-0393, S-0394, S-0395, S-0396, S-0397, S-0398, S-0399, S-0400, S-0401, S-0402, S-0403, S-0404, S-0405, S-0406 and S-0407."
  by: flai
  at: 2026-10-16T07:30:22Z
---
# S-0408 An OpenCode agent's protected write is asked of the story's owner through flai's hold-and-ask from the plugin's permission.asked, or refused by deny rules under --auto, as an ADR settles

## Goal

A story's agent on OpenCode asks before it writes a protected path, as one on Claude Code does through `permission_prompt`. Either the plugin answers OpenCode's `permission.asked` with the owner's answer from flai's hold-and-ask (`ask.Asker.Hold`, ADR-0124's thread and retry), or the adapter runs `--auto` behind `deny` rules for the protected list and the agent is told to ask on a thread. An ADR settles which from what S-0406 found, so that a story on OpenCode no longer needs `auto-approve` or `--deny-protected` to start.

## Acceptance criteria

- [ ] An ADR, refining ADR-0124 and ADR-0130, records which way OpenCode asks, OpenCode's bound for a held call, and why the other way was not taken.
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

- **Touches.** None were declared. Every touch is the planner's:
  - **Layout:** `flai/cmd/guard.go` and `flai/internal/guard/call.go`, where `Hold` and the `Ask` verdict live, with their tests; `flai/internal/guard/opencode.go`, S-0407's reader, for the `permission.asked` input; `flai/internal/harness/opencode.go` and `opencodeconfig.go` for the capability and the `ask` or `deny` rules, and `harness_test.go` for the capabilities test; the two plugin copies; `flai/tests/integration/opencode_test.go`. `flai/internal/ask/ask.go` is reused unchanged.
  - **Design:** the docs, `agent-adapters.md`, and `template/CHANGELOG.md` for criterion 5.
  - **Folder kept:** `design/adrs`, because the ADR's file name is not known until `flai adr new`; it is a shared path, so it holds nothing.
  - **Co-change:** `flai touches suggest` lists hub documents and `harness_test.go`; only `harness_test.go` is added.
- **Forecast.** 50m against flai's 39m. The ADR and the first caller of `guard.Hold`, whose bounded waits and thread reuse need timed tests, add a quarter. The delivery, 2026-10-16T14:32Z, is flai's play-out moved by the raise times the cycle factor.
- **Cost of delay.** Not set: E-0023 has no inputs, and `flai cod` refuses. Asked on TH-0467.
