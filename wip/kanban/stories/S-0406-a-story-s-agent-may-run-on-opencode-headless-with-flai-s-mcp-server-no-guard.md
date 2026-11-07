---
id: S-0406
type: story
nature: feature
title: A story's agent may run on OpenCode headless with flai's MCP server, no guard, and no roles, tested on fixtures, with the open OpenCode findings settled from its source at a pinned release
status: in-progress
parent: E-0023
owner: alex
created: 2026-10-16T05:13:48Z
updated: 2026-11-07T10:09:05Z
transitions:
  - to: ready
    at: 2026-11-07T09:38:00Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-07T09:38:06Z
    by: agent-S-0406
tags: [cli, template]
topics: [agents, cli, template]
touches: [flai/internal/harness/opencode.go, flai/internal/harness/opencode_test.go, flai/internal/harness/opencodeconfig.go, flai/internal/harness/opencodeconfig_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/internal/protected/protected.go, flai/internal/protected/protected_test.go, flai/internal/check/instructions.go, flai/internal/check/instructions_test.go, flai/internal/check/check.go, flai/internal/serve/agents_test.go, flai/internal/serve/testdata/opencode-run.jsonl, design/adrs, design/system/agent-adapters.md, design/system/flai-cli.md, docs/operators/index.md, docs/users/flai.md, docs/users/flai-reference.md, template/root/opencode.json, template/CHANGELOG.md, flai/internal/metrics/testdata/good/CLAUDE.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1937
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-07
      ceremony: 1
      hand_edits: 1
      work: 32
  models:
    - model: claude-opus-5-5
      input: 318
      output: 1910
      cache_read: 19265438
      cache_write: 669805
      cost: 7.1488
  strategic:
    - kind: orchestrator
      seconds: 245
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 38
          output: 334
          cache_read: 3281106
          cache_write: 130831
          cost: 0.8537
cost_of_delay:
  value: 97.01
  by: planner-E-0023
  at: 2026-11-07T09:36:35Z
forecast:
  duration: 1h05m
  delivery: 2026-11-07T10:50:00Z
  basis: "flai's 52m over 115 done large feature stories, raised a quarter for reading four findings from OpenCode's source at a pinned release and an ADR; the binary-gated integration test and smoke step are gone."
  by: planner-E-0023
  at: 2026-11-07T09:35:48Z
finalized:
  by: orchestrator
  at: 2026-11-07T09:37:29Z
---
# S-0406 A story's agent may run on OpenCode headless with flai's MCP server, no guard, and no roles, tested on fixtures, with the open OpenCode findings settled from its source at a pinned release

## Goal

An operator can set a story's `agent.harness` to `opencode` and `flai serve` runs it: `opencode run` headless in the story's worktree, with the harness prompt, the story's model as `provider/model`, JSON output, and a resumed session, reaching flai's MCP server as `flai mcp --agent <name>`. The adapter states that it has no guard and no hold-and-ask yet, so ADR-0131 applies: no roles, and a protected write refused unless `auto-approve` is on or the host denies the protected paths. The operator installs no OpenCode and configures no provider (TH-0532, TH-0533, D-0506), so the adapter is tested on fixtures alone, against a stub `opencode` program, and the four findings `agent-adapters.md` § OpenCode leaves to check are settled from OpenCode's source and documentation at one pinned release, which the fixtures follow, so that the later stories of E-0023 build on stated facts.

## Acceptance criteria

- [ ] `harness.Adapters` has `opencode`. Its default host runs `opencode run` in the story's worktree with the harness prompt, `--model <provider/model>`, `--format json`, and `--session <id>` when it resumes; `flai serve agent harness opencode` sets the program and arguments as it does for `claude-code`, and `TestEachAdapterStatesItsCapabilities` covers it.
- [ ] At start flai writes or merges an OpenCode configuration with a local MCP entry `flai mcp --agent <name>`, keeping every other key the project's own `opencode.json` has, and the story's worktree is left with no uncommitted change from it. An ADR records where that configuration lives and why.
- [ ] The adapter's `Capabilities` say resume only: no guard, no ask, no role models, no strategic definitions, no cost. A story on `opencode` whose `roles` name a sub-agent is refused at start with the setting that lifts it, and a start without `auto-approve` is refused unless the host sets `--deny-protected`, which the adapter writes as `deny` rules in OpenCode's `permission` map for every path the protected list holds for `opencode`.
- [x] `flai/internal/protected` holds an `opencode` list: Claude Code's paths plus `.opencode/`, `opencode.json`, and `AGENTS.md` (ADR-0130).
- [ ] `flai check` warns when the project has neither `AGENTS.md` nor `CLAUDE.md` at its root. The template keeps `CLAUDE.md` as the map, which OpenCode reads as its fallback, and ships an `opencode.json` skeleton with no MCP entry and no key.
- [x] `agent-adapters.md` § OpenCode names the pinned OpenCode release and records, for it, from its source and documentation: whether a plugin's `tool.execute.before` input can tell a sub-agent's call from the primary's, whether an `ask` in headless mode can be answered from outside, the fields `step_finish` carries, and whether `--attach` to a running `opencode serve` is better for many stories on one host. Each names the source file or page it was read from and says it was read, not run against a binary.
- [ ] The adapter and `flai serve`'s start, resume, and refusals for it are tested on fixtures alone: a stub `opencode` the test writes, as `claudecheck_test.go` writes its stub `claude`, records its arguments, environment, and configuration and replays a `--format json` stream written from the pinned release's event types. No test, tier, or close-out needs an `opencode` binary, a network, or a provider key.
- [ ] `docs/operators/index.md` gains a § Harnesses that says how to set the harness and what a story on it has and lacks today, and that the adapter is tested on fixtures, not against a live `opencode`; `docs/users/flai.md`, the generated reference, and `design/system/flai-cli.md` name the harness.

## Tasks
- T-1682 agent-adapters.md § OpenCode records the four open findings, read from OpenCode's source and documentation at a pinned release
- T-1683 The protected list holds an opencode entry with Claude Code's paths plus .opencode/, opencode.json, and AGENTS.md
- T-1684 flai check warns when a project has neither AGENTS.md nor CLAUDE.md at its root
- T-1685 The opencode adapter runs opencode run headless with the prompt, provider/model, JSON output, and a resumed session, and states resume as its only capability
- T-1686 An ADR records where flai writes a story's OpenCode configuration so the worktree stays clean, and how it merges the project's own opencode.json
- T-1687 At start the opencode adapter writes the story's OpenCode configuration with flai's MCP entry and, under deny-protected, deny rules for the protected list
- T-1688 flai serve starts a story's agent on a stub opencode that replays a fixture stream, and its refusals for opencode are tested
- T-1689 The operator, user, and design docs, the reference, and the template describe the opencode harness and what a story on it has and lacks

## Notes

- Drafted by the planner for E-0023. First of seven: every later story builds on this adapter and on the findings it records.
- The worktree must stay clean: a written `opencode.json` or `.opencode/` file left uncommitted refuses the move to review. `OPENCODE_CONFIG` pointing at a file under `.flai-cache`, or a git-ignored path, are the two ways the ADR weighs; the finding's own note is that `.opencode/` and `opencode.json` are files the harness reads to decide what its agent may do.
- Fixtures only (TH-0532, D-0506): the operator installs no OpenCode on the host and configures no OpenRouter key (TH-0533). OpenCode is MIT-licensed and its source is public, so the pinned release's event and plugin types are read there. Nothing in this story installs or runs a real `opencode`.

### Planning

- **Touches.** Re-planned on 2026-11-07 for fixtures only. Every touch is the planner's:
  - **Layout:** `flai/internal/harness/harness.go` and `harness_test.go`, the registry and the capabilities test; `opencode.go` and `opencodeconfig.go`, new, beside `adapters.go`, for criteria 1 to 3; `flai/internal/protected/protected.go`, whose comment already says `.opencode/` joins with this adapter, for criterion 4; `flai/internal/check/instructions.go`, new, and `check.go` to register it, for criterion 5; `flai/internal/serve/agents_test.go` for the stub start and the refusals of criteria 3 and 7, while `agents.go` reads the capabilities unchanged; `flai/internal/serve/testdata/opencode-run.jsonl`, new, the fixture stream, beside `usage-limit-result.jsonl`.
  - **Dropped with the binary:** `flai/tests/integration/opencode_test.go`, `scripts/opencode-smoke.sh`, `scripts/smoke.sh`, and `scripts/README.md`. A stub needs no integration tier or smoke step, and no test is skipped for want of a binary.
  - **Design:** `design/system/agent-adapters.md` for criterion 6; the docs, the reference, `template/root/opencode.json`, and `template/CHANGELOG.md` for criteria 5 and 8. `flai serve agent harness` derives its names from `harness.Names()`, so `flai/cmd/serve_actions.go` is not touched.
  - **Folder kept:** `design/adrs`, because the ADR's file name is not known until `flai adr new`; it is a shared path, so it holds nothing.
  - **Co-change:** `flai touches suggest` lists hub documents (`work-management.md`, `project-manifest.md`, `strategic-agents.md`, `template/template.yaml`) and `flai/internal/serve/agents.go`, `serve_actions.go`, and `mcpserver/server.go`, each under 6% of the seed commits. None of the criteria reach them: `agents.go` reads the capabilities unchanged and `serve_actions.go` takes its names from `harness.Names()`. None was added, and no setting key is added, so `settings.md` stays as it is.
- **Forecast.** 1h05m against flai's 52m, the median over 115 done large feature stories at size 29. Reading the four findings from OpenCode's source on the web at a pinned release, and the ADR on where the configuration lives, add a quarter. The old raise for checking by hand against a binary and for a binary-gated test and smoke step is gone with them. The delivery, 2026-11-07T10:50Z, is flai's 10:37Z play-out moved by the 13m raise.
- **Cost of delay.** 97.01 USD a week. E-0023 carries `penalty_per_week: 500`, the operator's answer to TH-0467 question 1, confirmed on TH-0470 on 2026-10-20. `flai cod` shares it among the epic's seven open stories by forecast, 1h05m of 5h35m. Kept as flai gives it: this story is the gate every later E-0023 story waits for, and that is already priced as the epic's penalty while it is open.
