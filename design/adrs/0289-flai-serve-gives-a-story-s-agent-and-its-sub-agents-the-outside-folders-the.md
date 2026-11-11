---
id: ADR-0289
title: "flai serve gives a story's agent and its sub-agents the outside folders the operator names per harness on the host, as --add-dir and container mounts, and Monitor joins the default allowed tools"
status: accepted
date: 2026-11-11
supersedes: []
superseded_by: []
refines: [ADR-0086]
topics: [cli, permissions, harness, container]
---

# ADR-0289 flai serve gives a story's agent and its sub-agents the outside folders the operator names per harness on the host, as --add-dir and container mounts, and Monitor joins the default allowed tools

## Context

[I-0265](../issues/I-0265-a-story-s-sub-agents-write-and-monitor-are-refused-every-path-outside-its.md) counts two instances in a day. S-0575 is an experiment that reads three other repositories and runs their scripts. S-0586 drives kluster, whose source lives outside this repository. In both, the explorer was refused every Read, Grep, and Glob outside the story's worktree. S-0575's own agent was refused a Write to `/tmp/S-0575` and every Monitor call. Bash passed. The agents worked around it by hand: they did the sub-agents' work themselves, kept files under the worktree's git-ignored `bin/<story>/`, and waited with Bash loops ([D-0480](../digest/D-0480-a-story-whose-work-reads-or-runs-repositories-outside-its-worktree-is-worked-by.md)).

The cause is in how `flai serve` starts a claude-code agent. It runs `claude -p` in the main checkout. The adapter passes `--permission-prompt-tool mcp__flai__permission_prompt` (`permissionArgs` in `flai/internal/harness/adapters.go`), then the operator's host arguments. The default is `--permission-mode acceptEdits --allowedTools Bash,mcp__flai` (`DefaultHost`). A Read, Grep, Glob, or Write outside the project needs approval, and so does every Monitor call, since no allowed tool names Monitor. Each falls to `permission_prompt` (`flai/internal/mcpserver/permission.go`). It refuses every tool but an Edit, Write, MultiEdit, or NotebookEdit, and every path outside an in-progress story's worktree ([ADR-0086](0086-flai-serve-gives-a-claude-code-agent-flai-s-permission-prompt-as-its-permission.md), [ADR-0106](0106-a-story-whose-branch-changes-a-path-claude-code-protects-is-accepted-by-the.md)). Bash passes because `--allowedTools` names it.

In a container it is worse. The container harness mounts the worktree, the main checkout's git folders, and the cache. No outside folder is mounted at all. pilot starts the inner claude with an empty host (`harness.Host{}` in `flai/internal/pilot/run.go`), so host arguments do not reach it.

[ADR-0038](0038-flai-serve-starts-a-story-s-own-agent-through-an-adapter-with-what-the-operator.md) makes what an agent may reach the operator's alone: "a story names only the tunables its adapter checks". [ADR-0041](0041-a-story-in-ready-gets-its-agent-whether-it-entered-ready-before-or-after-flai.md) supersedes ADR-0038 in part, only its clause on restarts, so that principle stands. E-0032's kluster stories will meet the same refusals.

On TH-1025 the operator chose remedy A on 2026-11-11 ([D-0589](../digest/D-0589-s-0701-which-remedy-for-i-0265-host-granted-outside-folders-story-requested.md)).

## Decision

**The operator names outside folders per harness on the host. flai serve puts them in the start request. The claude-code adapter passes each as `--add-dir`, and the container harness mounts each at its own host path. `Monitor` joins the default allowed tools.** This refines ADR-0086. The operator decided it on TH-1025 (2026-11-11).

1. **The operator names the folders on the host.** `flai serve agent harness <name> --outside <dir>` sets them. The flag repeats, once per folder, and `--outside ""` clears them. They are stored with the harness's other host settings, such as its network and login: per project under `agent.projects`, or for every project with `--all-projects`. Each folder must exist. It is given absolute or under `~/`, and stored absolute. The setting refuses home itself, `~/.flai`, any folder that holds the project root, and any folder the project holds.
2. **flai serve puts the folders in the start request.** They go in `harness.Request`, not only in the host. pilot starts the inner claude in a container with an empty host, so a host setting alone would never reach it.
3. **The claude-code adapter passes each folder as `--add-dir <folder>`.** A sub-agent the session starts inherits the session's added folders. T-3059 checks that with one scratch `claude -p` on the cheapest model, and records the result in S-0701's narrative. This decision rests on that check.
4. **The container harness mounts each folder read-write at its own host path.** It refuses home, `~/.flai`, the main checkout, the Docker socket, and any folder that holds one of them, as `mountable` already does for the login.
5. **`Monitor` joins the default `--allowedTools`.** The default becomes `--permission-mode acceptEdits --allowedTools Bash,Monitor,mcp__flai`. Monitor runs a shell command, as Bash does, so allowing it gives the agent nothing Bash does not.
6. **The story's start prompt names the granted folders.** It says that Read, Grep, Glob, Write, and Monitor are refused on any other path outside the project. There, work goes through Bash, with files kept under the worktree's git-ignored `bin/<story>/` (D-0480).
7. **`permission_prompt` is unchanged.** It still approves only a write of a protected path inside an in-progress story's worktree, as ADR-0086 and ADR-0106 leave it. A tool call in a granted folder never reaches it, because the folder is one of the session's own.
8. **Strategic agents get the folders too.** The planner, the orchestrator, the analyzer, the curator, and each defined role share the start request, so they get the same `--add-dir` and the same mounts.

ADR-0038's principle holds. A story still names no folder. Only the operator widens what an agent reaches, on the host, through `flai serve agent harness`.

## Consequences

- A story that reads or runs another repository can hand its search to an explorer, Write its scratch files there, and wait with Monitor, once the operator names the folder.
- A granted folder is open to every agent of the project on that harness, and to their sub-agents, for reading and writing. That is the operator's choice, made in a shell.
- A folder the operator has not named is refused as before, now with a prompt that says so and says what to do instead.
- Any operator who set their own `--allowedTools` keeps it. Monitor joins only the default.
- Inheritance by sub-agents is Claude Code's behaviour. A Claude Code release that changes it would bring I-0265 back for sub-agents alone.
- Other harnesses get the folders in the start request, but only claude-code and the container harness act on them here.

## Alternatives considered

- **B: a story asks for folders, within a host list.** A story would name folders in `agent.config.add_dirs`, and flai would pass only those inside the operator's list. It adds a story field that anyone who can edit a story in the dashboard writes. It still needs A's host list, so it is A plus about half again the work.
- **C: a warning in the prompt only.** It is the smallest change, but the friction stays. I-0265 had two instances in a day, and E-0032's kluster stories would meet it again.
- **`--add-dir` added by hand to the harness's arguments.** The operator can already do it on the host. A container would not mount the folder, pilot would not pass it, and no prompt would name it.
