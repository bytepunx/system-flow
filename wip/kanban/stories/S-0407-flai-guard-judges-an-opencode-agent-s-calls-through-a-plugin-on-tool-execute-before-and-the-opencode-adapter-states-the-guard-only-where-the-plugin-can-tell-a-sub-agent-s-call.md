---
id: S-0407
type: story
nature: feature
title: flai guard judges an OpenCode agent's calls through a plugin on tool.execute.before, and the opencode adapter states the guard only where the plugin can tell a sub-agent's call
status: backlog
parent: E-0023
owner: alex
created: 2026-10-16T05:15:14Z
updated: 2026-10-16T05:24:55Z
transitions: []
tags: [cli, template]
topics: [agents, cli, template]
touches: [flai/internal/guard/opencode.go, flai/internal/guard/opencode_test.go, flai/internal/guard/guard.go, flai/internal/guard/subagents.go, flai/internal/guard/subagents_test.go, flai/cmd/guard_test.go, flai/internal/harness/opencode.go, flai/internal/harness/opencode_test.go, ".opencode/plugins/flai-guard.js", template/root/.opencode/plugins/flai-guard.js, flai/tests/integration/opencode_test.go, design/system/agent-adapters.md, design/system/flai-cli.md, docs/users/flai.md, docs/operators/index.md, template/CHANGELOG.md, flai/internal/harness/harness_test.go]
after: [S-0406]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
forecast:
  duration: 45m
  delivery: 2026-10-16T13:39:00Z
  basis: "Its own forecast of 45m; 17th in the pull order with an in-progress limit of 5, behind S-0241, S-0392, S-0393, S-0394, S-0395, S-0396, S-0397, S-0398, S-0399, S-0400, S-0401, S-0402, S-0403, S-0404, S-0405 and S-0406."
  by: flai
  at: 2026-10-16T05:24:55Z
---
# S-0407 flai guard judges an OpenCode agent's calls through a plugin on tool.execute.before, and the opencode adapter states the guard only where the plugin can tell a sub-agent's call

## Goal

A story's agent on OpenCode is held by the same guard as one on Claude Code. A plugin flai ships into `.opencode/plugins/` hands each call to `flai guard` before it runs, and throws when the guard refuses, so a sub-agent cannot change a work item, a thread, or the history, and a planner on OpenCode stays a planner. The adapter's `Capabilities` say `Guard: true` only where S-0406 found that the plugin can tell a sub-agent's call from the primary's; where it cannot, the roles stay off under ADR-0131 and `agent-adapters.md` says why.

## Acceptance criteria

- [ ] `flai/internal/guard` has an `opencode` input reader that maps the plugin's input onto the neutral `Call`: OpenCode's tool names (`bash`, `edit`, `write`, `patch`, and flai's MCP tools under their OpenCode names) onto `Kind`, its sub-agent sign onto `Subagent`, and a read-only sub-agent onto `ReadOnly`. `flai guard --harness opencode` decides with the one `Decide` and answers in the form the plugin reads.
- [ ] `.opencode/plugins/flai-guard.js`, in this repository and in the template, runs `flai guard --harness opencode` on `tool.execute.before` with the call on standard input and throws with the guard's reason when it refuses; it passes every call when `flai` is not on `PATH`, saying so once, so a project opened outside `flai serve` still works.
- [ ] Where the plugin can tell a sub-agent's call, the adapter states `Guard: true`, and a story on `opencode` with roles is no longer refused for the guard. Where it cannot, `Guard` stays false and `agent-adapters.md` § OpenCode records why and what OpenCode would need to change.
- [ ] Where `Guard` is true, the integration test shows a sub-agent refused a `flai` write and the story's own agent let through, on the known `opencode` version, and is skipped, saying so, without it.
- [ ] `docs/operators/index.md` § Harnesses, `docs/users/flai.md`, and `design/system/flai-cli.md` say what the guard holds on OpenCode, and `template/CHANGELOG.md` names the plugin.

## Tasks
- T-1690 flai guard reads an OpenCode plugin's call into the neutral Call and answers in the form the plugin reads
- T-1691 A flai-guard plugin in .opencode/plugins hands each call to flai guard on tool.execute.before and throws on refusal, in this repository and the template
- T-1692 The opencode adapter states Guard true where the plugin can tell a sub-agent's call, so a story on it with roles is no longer refused for the guard
- T-1693 The OpenCode integration test shows the plugin refusing a sub-agent's flai write and letting the story's own agent through
- T-1694 The operator, user, and design docs say what flai guard holds on OpenCode, and agent-adapters.md records the sub-agent finding's consequence

## Notes

- Drafted by the planner for E-0023, after S-0406, whose finding on the sub-agent sign decides criterion 3.
- `.opencode/` is on the protected list for `opencode` only (S-0406), so this repository's own agent on `claude-code` writes the plugin without a permission thread.
- Claude Code records sub-agents through `SubagentStart` and `SubagentStop` for ADR-0092; the reader records OpenCode's sub-agent sessions the same way where the plugin sees them start and end.

### Planning

- **Touches.** None were declared. Every touch is the planner's:
  - **Layout:** `flai/internal/guard/opencode.go`, new, beside `claudecode.go`, and `guard.go`, whose `readers` map registers it; `subagents.go`, which records Claude Code's sub-agents for ADR-0092 and gains OpenCode's; `flai/cmd/guard_test.go` for `--harness opencode`, while `flai/cmd/guard.go` already looks the reader up by name; `flai/internal/harness/opencode.go` for the capability, and `harness_test.go`, where `TestEachAdapterStatesItsCapabilities` lives; the plugin in `.opencode/plugins/` and `template/root/.opencode/plugins/`; `flai/tests/integration/opencode_test.go`, S-0406's.
  - **Design:** `design/system/agent-adapters.md` for criterion 3; the docs and `template/CHANGELOG.md` for criterion 5.
  - **Co-change:** `flai touches suggest` lists hub documents, `guard_test.go`, and `harness_test.go`. `harness_test.go` is added; `guard_test.go` is not, since the reader's tests go in `opencode_test.go`.
- **Forecast.** 45m against flai's 34m. A JavaScript plugin has no test pattern in this repository, and the sub-agent tracking is checked against the real binary: a third again. The delivery, 2026-10-16T13:36Z, is flai's play-out moved by the raise times the cycle factor.
- **Cost of delay.** Not set: E-0023 has no inputs, and `flai cod` refuses. Asked on TH-0467.
