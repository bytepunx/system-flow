---
id: S-0407
type: story
nature: feature
title: flai guard judges an OpenCode agent's calls through a plugin on tool.execute.before, and the opencode adapter states the guard only where the plugin can tell a sub-agent's call
status: backlog
parent: E-0023
owner: alex
created: 2026-10-16T05:15:14Z
updated: 2026-10-23T17:35:05Z
transitions: []
tags: [cli, template]
topics: [agents, cli, template]
touches: [flai/internal/guard/opencode.go, flai/internal/guard/opencode_test.go, flai/internal/guard/guard.go, flai/internal/guard/subagents.go, flai/internal/guard/subagents_test.go, flai/cmd/guard.go, flai/cmd/guard_test.go, docs/users/flai-reference.md, flai/internal/harness/opencode.go, flai/internal/harness/opencode_test.go, flai/internal/harness/harness_test.go, flai/internal/serve/agents_test.go, ".opencode/plugins/flai-guard.js", template/root/.opencode/plugins/flai-guard.js, flai/tests/integration/opencode_test.go, design/system/agent-adapters.md, design/system/flai-cli.md, docs/users/flai.md, docs/operators/index.md, template/CHANGELOG.md]
after: [S-0406]
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
      seconds: 258
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 36
          output: 521
          cache_read: 14648995
          cache_write: 15789
          cost: 3.6763
cost_of_delay:
  value: 69.23
  by: planner-E-0023
  at: 2026-10-20T21:13:38Z
forecast:
  duration: 45m
  delivery: 2026-10-23T19:50:00Z
  basis: "Its own forecast of 45m; 3rd in the pull order with an in-progress limit of 5, behind S-0241 and S-0406."
  by: flai
  at: 2026-10-23T17:35:05Z
finalized:
  by: orchestrator
  at: 2026-10-20T21:14:03Z
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

- **Touches.** The story declared the E-0023 planner's 17 touches, all kept. This run added three. No folder touch is kept: every path is a file.
  - **Declared, from layout:**
    - `flai/internal/guard/opencode.go` and `opencode_test.go`, new, beside `claudecode.go`.
    - `guard.go`, whose `readers` map registers the reader.
    - `subagents.go` and `subagents_test.go`, kept for ADR-0092's record of OpenCode's sub-agents. A reader returns `Hook.SubAgent` and the record already takes it, so these may not change.
    - `flai/cmd/guard_test.go` for `--harness opencode`.
    - `flai/internal/harness/opencode.go`, `opencode_test.go`, and `harness_test.go` (`TestEachAdapterStatesItsCapabilities`) for the capability.
    - The plugin, `.opencode/plugins/flai-guard.js` and `template/root/.opencode/plugins/flai-guard.js`.
    - `flai/tests/integration/opencode_test.go`, which S-0406 creates.
  - **Declared, from design:** `design/system/agent-adapters.md` for criterion 3; `design/system/flai-cli.md`, `docs/users/flai.md`, `docs/operators/index.md`, and `template/CHANGELOG.md` for criterion 5.
  - **Added, from layout:**
    - `flai/cmd/guard.go`, whose help says claude-code is "the only reader there is" and whose `--harness` flag lists `guard.Harnesses()`. Also `docs/users/flai-reference.md`, which is generated from that help (`make flai-reference`) and must be regenerated in the same task. Both go on T-1690.
    - `flai/internal/serve/agents_test.go`, on T-1692, for a `startAllowed` test that shows criterion 3. It was approved on TH-0470.
  - **Co-change, not added:** `flai touches suggest` also lists hub documents (`flaiover-dashboard.md`, `docs/users/flaiover.md`, `design/adrs/README.md`, `docs/operators/settings.md`, `workflow.md`). None of the criteria reach them. `settings.md` already indexes `--harness` for `flai guard`, and no setting is added. The list also has `guard_test.go`, but the reader's tests go in `opencode_test.go`. `harness.go` is left out because S-0406 registers the adapter.
  - **Layout, not added:** `flai/internal/serve/agents.go`. `startAllowed` reads `Capabilities().Guard` for any harness, so a `Guard: true` lifts the refusal without a change there.
- **Forecast.** 45m against flai's 40m. flai's figure grew from 36m with the three added touches. Those are mechanical: a help text, a regenerated reference, and one test case. The raise for a JavaScript plugin with no test pattern here, and for sub-agent tracking checked against the real binary, therefore drops from a third to about an eighth. The delivery, 2026-10-16T13:42Z, is flai's 13:36Z play-out moved by the 5m raise times the cycle factor 1.11.
- **Cost of delay.** 69.23 USD a week. E-0023 carries `penalty_per_week: 500`, the operator's answer to TH-0467 question 1, confirmed on TH-0470 on 2026-10-20 (as on E-0019). `flai cod` shares it among the epic's seven open stories by forecast, 45m of 5h25m. Kept as flai gives it. TH-0470 asked the same input.
