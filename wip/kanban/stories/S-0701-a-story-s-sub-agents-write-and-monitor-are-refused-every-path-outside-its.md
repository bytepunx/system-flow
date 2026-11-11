---
id: S-0701
type: story
nature: improvement
title: A story's sub-agents, Write, and Monitor are refused every path outside its worktree, though an experiment story may be written to work outside the repository
status: backlog
owner: alex
created: 2026-11-10T08:55:34Z
updated: 2026-11-11T09:01:05Z
transitions: []
tags: [harness, serve, docker]
touches: [design/adrs, design/adrs/README.md, design/system/agent-adapters.md, flai/internal/config/config.go, flai/internal/config/config_test.go, flai/cmd/serve_harness.go, flai/cmd/serve_agent.go, flai/cmd/serve_actions_test.go, docs/operators/settings.md, docs/users/flai-reference.md, flai/internal/harness/harness.go, flai/internal/harness/adapters.go, flai/internal/harness/harness_test.go, flai/internal/serve/agents_start.go, flai/internal/serve/agents_test.go, flai/internal/harness/container.go, flai/internal/harness/container_test.go, flai/internal/pilot/request.go, flai/internal/pilot/request_test.go, docs/operators/index.md, docs/users/flai.md, docs/operators/runbooks/agent-container.md, design/issues/I-0265-a-story-s-sub-agents-write-and-monitor-are-refused-every-path-outside-its.md, design/issues/summary.md]
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
      seconds: 42
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 12
          output: 76
          cache_read: 372455
          cache_write: 69627
          cost: 0.1235
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 26m
    by: flai
    at: 2026-11-10T08:55:34Z
  value: 65
  by: planner-S-0701
  at: 2026-11-10T10:30:33Z
forecast:
  duration: 1h30m
  delivery: 2026-11-11T22:04:00Z
  basis: "Its own forecast of 1h30m; 18th in the pull order with an in-progress limit of 3, behind S-0694, S-0698, S-0702, S-0706, S-0710, S-0719, S-0724, S-0716, S-0689, S-0699, S-0705, S-0681, S-0718, S-0241, S-0659, S-0660, S-0661, S-0662, S-0663 and S-0683."
  by: flai
  at: 2026-11-11T09:01:05Z
---
# S-0701 A story's sub-agents, Write, and Monitor are refused every path outside its worktree, though an experiment story may be written to work outside the repository

## Goal

This story remediates [I-0265](../../../design/issues/I-0265-a-story-s-sub-agents-write-and-monitor-are-refused-every-path-outside-its.md), "A story's sub-agents, Write, and Monitor are refused every path outside its worktree, though an experiment story may be written to work outside the repository". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0265 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0265 is closed with `flai issue close I-0265 --reason` saying what fixed it

## Tasks
- T-3048 Record in an ADR that the operator's outside folders reach a story's agent and its sub-agents as --add-dir and container mounts, and that Monitor joins the default tools
- T-3054 flai serve agent harness takes --outside folders, checks and stores them with the harness's settings, and shows them
- T-3059 The claude-code adapter passes the start request's outside folders as --add-dir, and Monitor joins the default --allowedTools
- T-3063 The container harness mounts each outside folder at its own path and carries them to pilot, refusing what mountable refuses
- T-3064 A story's start prompt names the outside folders its agent and sub-agents may reach, and says every other outside path is refused
- T-3067 The operators guide, the container runbook, and the user guide say how to give a story's agent outside folders, and that Monitor is allowed
- T-3068 I-0265 is closed with the reason that the operator's outside folders reach a story's agent and its sub-agents and Monitor is allowed

## Notes

Cost of delay inputs set by flai from I-0265. time_lost_per_cycle 26m: 13m per occurrence × 2 occurrences ÷ 1 cycle of 168h (first reported 2026-11-06T23:16:17Z, 3.4 days before this story; under one cycle counts as one).

### Planning

Cause, as found in the code: flai serve starts `claude -p` in the main checkout with `--permission-prompt-tool mcp__flai__permission_prompt` and the host's default `--permission-mode acceptEdits --allowedTools Bash,mcp__flai` (`flai/internal/harness/adapters.go`, `DefaultHost` and `permissionArgs`). A Read, Grep, Glob, or Write outside the project, and every Monitor call, falls to `permission_prompt`. `permissionPrompt` in `flai/internal/mcpserver/permission.go` refuses every tool but a write, and every path outside a story's worktree. Bash passes because `--allowedTools` names it. In a container, no outside folder is mounted at all. pilot starts the inner claude with an empty host (`flai/internal/pilot/run.go`), so anything the agent inside must get travels in the start request.

Proposed solution, remedy A on TH-1025, which awaits the operator: the operator names outside folders on the host, `flai serve agent harness <name> --outside <dir>...`, stored with the harness's other settings. flai serve puts them in the start request. The claude-code adapter passes each as `--add-dir`, which the session's sub-agents inherit. The container harness mounts each at its own path, and refuses home and any folder that holds the project, as `mountable` does for the login. `Monitor` joins the default `--allowedTools`, since it runs a shell command as Bash does. The start prompt names the granted folders and says any other outside path is refused. `permission_prompt` is unchanged. ADR-0038 holds: only the operator widens what an agent reaches. If the operator chooses B or C on TH-1025, the plan is redrawn.

Touches, each with where it came from:

| Touch | From |
|-------|------|
| `design/adrs` | design: the ADR the decision needs, whose number is not known until it is written (folder touch kept for that file alone) |
| `design/adrs/README.md` | layout: the ADR index |
| `design/system/agent-adapters.md` | design: § Today › Permissions and § Sub-agents and strategic agents describe what the agent may reach |
| `flai/internal/config/config.go`, `flai/internal/config/config_test.go` | layout: the harness settings, with `network` and `login`, and their validation |
| `flai/cmd/serve_harness.go`, `flai/cmd/serve_agent.go`, `flai/cmd/serve_actions_test.go` | layout: `flai serve agent harness` takes and shows the settings, and the help text names the default arguments |
| `docs/operators/settings.md`, `docs/users/flai-reference.md` | layout: `TestSettingsIndex` needs a row per setting, and `TestReferenceIsCurrent` needs the reference regenerated from the help |
| `flai/internal/harness/harness.go`, `flai/internal/harness/adapters.go`, `flai/internal/harness/harness_test.go` | design: the start request, the adapter's arguments, and the start prompt |
| `flai/internal/serve/agents_start.go`, `flai/internal/serve/agents_test.go` | layout: where flai serve builds the start request from the host's settings |
| `flai/internal/harness/container.go`, `flai/internal/harness/container_test.go` | layout: the container's mounts and `mountable` |
| `flai/internal/pilot/request.go`, `flai/internal/pilot/request_test.go` | layout: pilot reads the start request the folders travel in |
| `docs/operators/index.md`, `docs/users/flai.md`, `docs/operators/runbooks/agent-container.md` | design: the harness defaults, the sub-agents and protected-paths sections, and the container's mounts |
| `design/issues/I-0265-a-story-s-sub-agents-write-and-monitor-are-refused-every-path-outside-its.md`, `design/issues/summary.md` | criterion 2 |

Co-change of `flai/internal/mcpserver/permission.go` listed its test, `folder.go`, `server.go`, `flai/cmd/mcp_http.go`, and `flai/internal/ask`. None is a touch: remedy A leaves `permission_prompt` as it is.

No `after`: the story reads no other open story's output. Stories in the backlog touch `container.go` and pilot too, such as S-0705, but a trial merge judges a shared file.

Forecast: `flai forecast` gave 45m, at 103 s per unit of size over 155 done large improvement stories, size 26. Raised to 1h30m: seven tasks in six layers, across four Go packages (config, cmd, harness, serve, with pilot), an ADR, and a regenerated reference. S-0682, which changed one container script, took 66m. The delivery is flai's, 2026-11-10T20:24Z, moved 45m later to match.

Cost of delay: 65 USD a week, as `flai cod` works it out from flai's input, 26m lost per 168h cycle at 150 USD an hour. It stands: I-0265 has two instances in a day, and D-0480 records the workaround for every later story that works outside the repository, so the friction will come back with E-0032's kluster stories.
