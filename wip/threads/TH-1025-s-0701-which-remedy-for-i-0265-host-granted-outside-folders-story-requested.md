---
id: TH-1025
title: "S-0701: which remedy for I-0265 — host-granted outside folders, story-requested within a host list, or a prompt warning only?"
anchor:
  path: wip/kanban/stories/S-0701-a-story-s-sub-agents-write-and-monitor-are-refused-every-path-outside-its.md
  item: S-0701
status: resolved
participants: [planner-S-0701, orchestrator, alex]
created: 2026-11-10T10:30:43Z
updated: 2026-11-11T09:03:13Z
---

# TH-1025 S-0701: which remedy for I-0265 — host-granted outside folders, story-requested within a host list, or a prompt warning only?

On wip/kanban/stories/S-0701-a-story-s-sub-agents-write-and-monitor-are-refused-every-path-outside-its.md.

## Entries

### 2026-11-10T10:30:43Z planner-S-0701
The fix decides what a story's agent may reach on your host, which ADR-0038 makes yours alone ("a story names only the tunables its adapter checks"). I need your choice before I draft S-0701's tasks.

## What causes I-0265

- flai serve starts `claude -p` in the main checkout with `--permission-prompt-tool mcp__flai__permission_prompt --permission-mode acceptEdits --allowedTools Bash,mcp__flai`.
- A Read, Grep, Glob, or Write of a path outside the project, and every Monitor call, falls to `permission_prompt`. In `flai/internal/mcpserver/permission.go` (lines 177–185) it refuses every tool but a write, and every path outside a story's worktree.
- Bash passes because `--allowedTools` names it.
- In a container, the outside folders are not mounted at all, so even Bash cannot reach them there.

## Recommended: A. The host grants outside folders per project

- You name the folders on the host, such as `flai serve agent outside ~/git/bytepunx/kluster /tmp`, stored in your flai configuration under `agent.projects`.
- flai serve passes each as `--add-dir` to every claude-code story agent of the project. The session's sub-agents inherit them. The container harness mounts each at the same path; it still refuses home itself and any folder holding the project.
- `Monitor` joins the default `--allowedTools`, since it runs a shell command as Bash does.
- The start prompt names the granted folders and says any other outside path is refused.
- ADR-0038 holds: only you widen what an agent reaches. Today you could already add `--add-dir` to the harness arguments by hand, but containers would not mount it and no prompt would say so.

## Alternatives

- **B. A story requests, the host list grants.** A story names the folders in `agent.config.add_dirs`. flai passes only those inside your host list, and refuses the start, naming the folder, for any other. This is narrower per story, but it adds a story field and a check on top of A, about 1.5× the work.
- **C. A prompt warning only.** The start prompt says outside paths are refused to Read, Write, Monitor, and sub-agents, and to use Bash and the worktree's git-ignored `bin/<story>/`, as D-0480 records. It is the smallest change, but the friction stays.

Reply A, B, or C, or name another remedy. Meanwhile I have set the cost of delay value: 65 USD a week, which `flai cod` works out from flai's input of 26m lost per 168h cycle at 150 USD an hour. Once you answer I will write the touches, the forecast, and the tasks.

### 2026-11-10T10:36:41Z planner-S-0701
The plan for S-0701 is drafted on remedy A. It is easy to redraw if you choose B or C. Seven tasks in six layers:

| Layer | Task | Touches |
|-------|------|---------|
| 1 | T-3048 The ADR for remedy A, and `design/system/agent-adapters.md` § Permissions | `design/adrs/README.md`, `design/system/agent-adapters.md` |
| 2 | T-3054 `flai serve agent harness --outside`, checked, stored, and shown | `flai/internal/config/config.go` and its test, `flai/cmd/serve_harness.go`, `flai/cmd/serve_agent.go`, `flai/cmd/serve_actions_test.go`, `docs/operators/settings.md`, `docs/users/flai-reference.md` |
| 3 | T-3059 The start request carries the folders, claude-code passes them as `--add-dir`, and `Monitor` joins the default tools; the test reproduces I-0265 | `flai/internal/harness/harness.go`, `adapters.go`, `harness_test.go`, `flai/internal/serve/agents_start.go`, `agents_test.go`, plus the help text and reference |
| 4 | T-3063 The container mounts each folder and carries them to pilot | `flai/internal/harness/container.go` and its test, `flai/internal/pilot/request.go` and its test |
| 4 | T-3064 The start prompt names the folders and says other outside paths are refused | `flai/internal/harness/harness.go`, `harness_test.go` |
| 5 | T-3067 The operators guide, the container runbook, the user guide, and the settings row | `docs/operators/index.md`, `docs/operators/runbooks/agent-container.md`, `docs/users/flai.md`, `docs/operators/settings.md` |
| 6 | T-3068 I-0265 is closed | the issue, `design/issues/summary.md` |

Each layer waits for the one before it. The two tasks in layer 4 run together, since they share no file.

## Figures

- Forecast: 1h30m, delivered 2026-11-10T21:09Z. `flai forecast` gave 45m; I doubled it for seven tasks across four Go packages and an ADR, against S-0682's 66m for one script.
- Cost of delay: 65 USD a week, from flai's input.
- No `after`.

## Assumptions

- The outside folders are a per-harness host setting, stored like `--login`, and they apply to every story's agent of the project.
- Folders are mounted read-write in the container, as Bash on the host already reaches them. Home, `~/.flai`, the main checkout, and the Docker socket are refused, and so is any folder holding one of them.
- Strategic agents get the folders too, since they share the start request.
- `permission_prompt` is unchanged.
- `scripts/agent-container.sh` and the template's `.claude/agents/explorer.md` are unchanged.
- T-3059 asks the story's agent to check once, with a scratch `claude -p` on the cheapest model, that a sub-agent inherits `--add-dir`. That is a paid call of a few cents. Say if you would rather it rely on Claude Code's documentation.

Reply A, B, or C on this thread, as asked above.

### 2026-11-10T11:02:38Z orchestrator (recommendation)
**Recommended:** remedy A, the host grants outside folders per harness, as the plan is drafted.

- ADR-0038 makes what an agent may reach on your host yours alone. A keeps it so: only `flai serve agent harness --outside` widens it, on the host.
- B adds a story field that anyone who can edit a story in the dashboard writes. It still needs your host list, so it is A plus about half again the work.
- C changes nothing a story can do. I-0265's two instances in a day would recur with E-0032's kluster stories (D-0480).
- Allow T-3059's one scratch `claude -p` on the cheapest model to confirm that a sub-agent inherits `--add-dir`. It costs a few cents, and criterion 1 rests on it.

This is host configuration and a paid call, so it is yours alone.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/adrs/0038-flai-serve-starts-a-story-s-own-agent-through-an-adapter-with-what-the-operator.md

### 2026-11-11T09:00:54Z alex
Confirmed the recommendation of 2026-11-10T11:02:38Z orchestrator.

Source: design/adrs/0038-flai-serve-starts-a-story-s-own-agent-through-an-adapter-with-what-the-operator.md

### 2026-11-11T09:02:54Z orchestrator
Approved: S-0701's plan stands as drafted on remedy A, which you confirmed (D-0589). Its seven tasks in six layers (T-3048 to T-3068), its touches, the 1h30m forecast, and the 65 USD a week value fit the story. Your confirmation changes nothing in the plan, so I start no new planner run. T-3059 may make its one scratch `claude -p` on the cheapest model, as the recommendation you confirmed allows. S-0701 is still a draft.

Source: design/digest/D-0589-s-0701-which-remedy-for-i-0265-host-granted-outside-folders-story-requested.md

### 2026-11-11T09:03:13Z orchestrator
Resolved: The operator confirmed remedy A (D-0589). S-0701's plan was already drafted on A, and it is approved as it stands.
