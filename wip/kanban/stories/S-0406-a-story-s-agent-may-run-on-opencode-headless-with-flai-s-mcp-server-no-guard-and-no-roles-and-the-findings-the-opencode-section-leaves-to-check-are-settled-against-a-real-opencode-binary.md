---
id: S-0406
type: story
nature: feature
title: A story's agent may run on OpenCode headless with flai's MCP server, no guard, and no roles, and the findings the OpenCode section leaves to check are settled against a real opencode binary
status: backlog
parent: E-0023
owner: alex
created: 2026-10-16T05:13:48Z
updated: 2026-10-19T20:33:45Z
transitions: []
tags: [cli, template]
topics: [agents, cli, template]
touches: [flai/internal/harness/opencode.go, flai/internal/harness/opencode_test.go, flai/internal/harness/opencodeconfig.go, flai/internal/harness/opencodeconfig_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/internal/protected/protected.go, flai/internal/protected/protected_test.go, flai/internal/check/instructions.go, flai/internal/check/instructions_test.go, flai/internal/check/check.go, flai/internal/serve/agents_test.go, flai/tests/integration/opencode_test.go, scripts/opencode-smoke.sh, scripts/smoke.sh, scripts/README.md, design/adrs, design/system/agent-adapters.md, design/system/flai-cli.md, docs/operators/index.md, docs/users/flai.md, docs/users/flai-reference.md, template/root/opencode.json, template/CHANGELOG.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
forecast:
  duration: 1h15m
  delivery: 2026-10-19T22:37:00Z
  basis: "Its own forecast of 1h15m; 8th in the pull order with an in-progress limit of 5, behind S-0404, S-0405, S-0403, S-0402, S-0401, S-0400 and S-0241."
  by: flai
  at: 2026-10-19T20:33:45Z
---
# S-0406 A story's agent may run on OpenCode headless with flai's MCP server, no guard, and no roles, and the findings the OpenCode section leaves to check are settled against a real opencode binary

## Goal

An operator can set a story's `agent.harness` to `opencode` and `flai serve` runs it: `opencode run` headless in the story's worktree, with the harness prompt, the story's model as `provider/model`, JSON output, and a resumed session, reaching flai's MCP server as `flai mcp --agent <name>`. The adapter states that it has no guard and no hold-and-ask yet, so ADR-0131 applies: no roles, and a protected write refused unless `auto-approve` is on or the host denies the protected paths. The story also settles, against a real `opencode` binary, the four findings `agent-adapters.md` § OpenCode leaves to check, so that the later stories of E-0023 build on facts.

## Acceptance criteria

- [ ] `harness.Adapters` has `opencode`. Its default host runs `opencode run` in the story's worktree with the harness prompt, `--model <provider/model>`, `--format json`, and `--session <id>` when it resumes; `flai serve agent harness opencode` sets the program and arguments as it does for `claude-code`, and `TestEachAdapterStatesItsCapabilities` covers it.
- [ ] At start flai writes or merges an OpenCode configuration with a local MCP entry `flai mcp --agent <name>`, keeping every other key the project's own `opencode.json` has, and the story's worktree is left with no uncommitted change from it. An ADR records where that configuration lives and why.
- [ ] The adapter's `Capabilities` say resume only: no guard, no ask, no role models, no strategic definitions, no cost. A story on `opencode` whose `roles` name a sub-agent is refused at start with the setting that lifts it, and a start without `auto-approve` is refused unless the host sets `--deny-protected`, which the adapter writes as `deny` rules in OpenCode's `permission` map for every path the protected list holds for `opencode`.
- [ ] `flai/internal/protected` holds an `opencode` list: Claude Code's paths plus `.opencode/`, `opencode.json`, and `AGENTS.md` (ADR-0130).
- [ ] `flai check` warns when the project has neither `AGENTS.md` nor `CLAUDE.md` at its root. The template keeps `CLAUDE.md` as the map, which OpenCode reads as its fallback, and ships an `opencode.json` skeleton with no MCP entry and no key.
- [ ] `agent-adapters.md` § OpenCode records, for the installed version, whether a plugin's `tool.execute.before` can tell a sub-agent's call from the primary's, whether an `ask` in headless mode can be answered from outside, the fields `step_finish` carries, and whether `--attach` to a running `opencode serve` is better for many stories on one host, each with how it was checked.
- [ ] An integration test starts a story's agent on OpenCode against a stub MCP-only prompt and reads its JSON output. It and a smoke step run only where `opencode` of the known version is on `PATH`, and are skipped, saying so, where it is not, so `make integration`, `make smoke`, and a close-out pass without it.
- [ ] `docs/operators/index.md` gains a § Harnesses that says how to install OpenCode and set the harness, and what a story on it has and lacks today; `docs/users/flai.md`, the generated reference, and `design/system/flai-cli.md` name the harness.

## Tasks
- T-1682 agent-adapters.md § OpenCode records the four open findings, checked against the installed opencode binary
- T-1683 The protected list holds an opencode entry with Claude Code's paths plus .opencode/, opencode.json, and AGENTS.md
- T-1684 flai check warns when a project has neither AGENTS.md nor CLAUDE.md at its root
- T-1685 The opencode adapter runs opencode run headless with the prompt, provider/model, JSON output, and a resumed session, and states resume as its only capability
- T-1686 An ADR records where flai writes a story's OpenCode configuration so the worktree stays clean, and how it merges the project's own opencode.json
- T-1687 At start the opencode adapter writes the story's OpenCode configuration with flai's MCP entry and, under deny-protected, deny rules for the protected list
- T-1688 An integration test and a smoke step run a story's agent on OpenCode only where opencode of the known version is installed, and flai serve's refusals for it are tested
- T-1689 The operator, user, and design docs, the reference, and the template describe the opencode harness and what a story on it has and lacks

## Notes

- Drafted by the planner for E-0023. First of seven: every later story builds on this adapter and on the findings it records.
- The worktree must stay clean: a written `opencode.json` or `.opencode/` file left uncommitted refuses the move to review. `OPENCODE_CONFIG` pointing at a file under `.flai-cache`, or a git-ignored path, are the two ways the ADR weighs; the finding's own note is that `.opencode/` and `opencode.json` are files the harness reads to decide what its agent may do.
- The binary is installed by the operator, never by flai (TH-0467, question 2).

### Planning

- **Touches.** None were declared. Every touch is the planner's:
  - **Layout:** `flai/internal/harness/harness.go` and `harness_test.go`, the registry and the capabilities test; `opencode.go` and `opencodeconfig.go`, new, beside `adapters.go`, for criteria 1 to 3; `flai/internal/protected/protected.go`, whose comment already says `.opencode/` joins with this adapter, for criterion 4; `flai/internal/check/instructions.go`, new, and `check.go` to register it, for criterion 5; `flai/internal/serve/agents_test.go` for the refusals of criterion 3, while `agents.go` reads the capabilities unchanged; `flai/tests/integration/opencode_test.go`, following the binary-on-`PATH` skip of `golangci_lock_test.go`, and `scripts/opencode-smoke.sh` with `scripts/smoke.sh` and `scripts/README.md`, following `scripts/gateway-smoke.sh`, for criterion 7.
  - **Design:** `design/system/agent-adapters.md` for criterion 6; the docs, the reference, `template/root/opencode.json`, and `template/CHANGELOG.md` for criteria 5 and 8. `flai serve agent harness` derives its names from `harness.Names()`, so `flai/cmd/serve_actions.go` is not touched.
  - **Folder kept:** `design/adrs`, because the ADR's file name is not known until `flai adr new`; it is a shared path, so it holds nothing.
  - **Co-change:** `flai touches suggest` lists hub documents (`docs/operators/settings.md`, `flaiover-dashboard.md`, `workflow.md`), none of which this story's criteria reach; none was added. No setting key is added, so `settings.md` stays as it is.
- **Forecast.** 1h15m against flai's 51m. The four findings are checked by hand against an external binary, and the story adds a configuration merge, an ADR, and a binary-gated test: half again. The delivery, 2026-10-16T12:46Z, is flai's play-out moved by the raise times the cycle factor.
- **Cost of delay.** Not set: E-0023 has no inputs, and `flai cod` refuses. Asked on TH-0467.
