---
id: T-0882
type: task
nature: feature
title: The orchestrator's prompt, its claude-code definition, and settings that run flai guard on its file edits
status: backlog
parent: S-0218
owner: alex
created: 2026-10-05T04:45:56Z
updated: 2026-10-05T04:47:09Z
transitions: []
stream: S-0218
tags: [flai, template]
touches: [flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/internal/harness/adapters.go, ".claude/agents/orchestrator.md", ".claude/settings.json", template/root/.claude/agents/orchestrator.md, template/root/.claude/settings.json, template/CHANGELOG.md, template/template.yaml]
---
# T-0882 The orchestrator's prompt, its claude-code definition, and settings that run flai guard on its file edits

## Work

Give the orchestrator a prompt and a definition, as S-0208 gave the planner (ADR-0082, `design/system/strategic-agents.md` § The planner › What it is told):

- `flai/internal/harness/harness.go`: `Prompt` dispatches `Role` `orchestrate` to an `orchestratePrompt`. It says: prime with `--role orchestrate`; call `inbox` and read the board; act only within `orchestration.permissions`, using flai's commands for the arithmetic (`flai order --by`, `flai promote --candidates`, `flai release --evaluate`, from S-0217) and never doing it yourself; log each decision with `activity_log`, kind `orchestrator`, saying what, why, and the policy figure behind it; then hold `wait_for_events` and repeat, without ending. Never edit code or documents, never work a story, and never work around a refusal of `flai guard`. `roleEnv` gives `FLAI_ROLE=orchestrate` and no `FLAI_ITEM` or `FLAI_STORY`.
- `flai/internal/harness/adapters.go`: on `claude-code`, a request with role `orchestrate` runs with `--agent orchestrator` over the project's `.claude/agents/orchestrator.md`, passed in `--agents` beside the explorer and the verifier, named `<key> orchestrate`. A project with no `orchestrator.md` is refused, naming `flai upgrade`.
- `.claude/agents/orchestrator.md` and the template's copy: the same rules in short, and its tools: `Read`, `Grep`, `Glob`, `Bash`, `Agent`, flai's MCP read tools, `inbox`, `board`, `item_edit`, `item_move`, `thread_open`, `thread_reply`, `activity_log`, `wait_for_events`, the MCP tool `plan`, and S-0217's MCP tools. No `Edit`, `Write`, or `NotebookEdit`.
- `.claude/settings.json` and the template's: the `Edit|Write|NotebookEdit` `PreToolUse` entry runs the guard when `FLAI_ROLE` is `plan` or `orchestrate`, and exits at once otherwise.
- A template release: `template/CHANGELOG.md` and `template/template.yaml`.

A story agent cannot write under `.claude/` without the operator: `flai serve`'s permission prompt asks on a thread (S-0257). Ask there for `orchestrator.md` and `settings.json`, with the whole contents.

This task waits for nothing and runs with T-0881, whose paths it does not share. T-0883 waits for it, since it starts the run this prompt describes.

## Done when

- `Prompt` for role `orchestrate` holds each instruction above, and a test pins them
- the `claude-code` adapter builds `--agent orchestrator` with the definition in `--agents`, refuses a project without it, and a test pins both
- both settings files run the guard on file edits for `plan` and `orchestrate` only, and both copies of `orchestrator.md` match
- `go test ./internal/harness/` passes, and the template smoke passes

## Notes
