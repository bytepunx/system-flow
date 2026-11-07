---
id: T-2582
type: task
nature: remediation
title: Every start falls back to the built-in default agent, and each refusal and waiting reason that remains names the missing setting, the dashboard page that sets it, and the command
status: in-progress
parent: S-0587
owner: alex
created: 2026-11-07T10:12:31Z
updated: 2026-11-07T11:13:26Z
transitions:
  - to: ready
    at: 2026-11-07T10:57:48Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T10:57:51Z
    by: system-flow
stream: S-0587
tags: [cli]
touches: [flai/internal/serve/start.go, flai/internal/serve/restart.go, flai/internal/serve/commit.go, flai/internal/serve/agents.go, flai/internal/serve/plan.go, flai/internal/serve/analyze.go, flai/internal/serve/curate.go, flai/internal/serve/orchestrate.go, flai/internal/serve/role.go, flai/internal/serve/claudecheck.go, flai/internal/serve/start_test.go, flai/internal/serve/agents_test.go, flai/internal/serve/default_agent_test.go, flai/internal/serve/analyze_test.go, flai/internal/serve/curate_test.go, flai/internal/serve/plan_test.go, flai/internal/serve/role_test.go, flai/internal/serve/role_events_test.go, flai/internal/serve/orchestrate_test.go, flai/internal/serve/orchestrate_hold_test.go, flai/internal/serve/serve_test.go, flai/internal/serve/claudecheck_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/cmd/analyze_test.go, flai/cmd/curate_test.go, flai/cmd/plan_test.go, flai/cmd/role_test.go, flai/cmd/role.go, flai/cmd/agent.go, flai/internal/mcpserver/role.go, flai/internal/hostapi/settings.go]
after: [T-2578, T-2580]
---
# T-2582 Every start falls back to the built-in default agent, and each refusal and waiting reason that remains names the missing setting, the dashboard page that sets it, and the command

## Work

It waits for T-2578, which adds `manifest.Agent.OrBuiltin` and `workitem.Repo.StartAgentOf`, and for T-2580, which makes the host's command per project in `AgentConfig`. ADR-0261 § 1 and § 4.

- Every path that starts, restarts, resumes, or commits a story's agent, and every view of the agent it will run (`start.go`, `restart.go`, `commit.go`, `agents.go`), resolves it with `repo.StartAgentOf(it, <command set for the project>)` in place of `repo.AgentOf(it)`, so a story whose agent and project name no harness, on a host with no command, starts on `claude-code`.
- The strategic starts (`plan.go`, `analyze.go`, `curate.go`, `orchestrate.go`, `role.go`) do the same with `OrBuiltin` on the role's agent when no command is set.
- The refusals `names no harness, and no command is set on the host` go, from serve and from `harness.For` (`flai/internal/harness/harness.go`), and so does the launcher's waiting reason that says it (`agents.go`).
- One helper in `flai/internal/serve` (or `harness`) words what remains, ADR-0261 § 4: what is missing, the dashboard page by its path, and the command:
  - a harness `command` with no command set for the project: `/settings/harnesses`, or `flai serve agent set -- <program> [args...]`;
  - a harness flai cannot start: `/settings/agents` for the project's default agent, or `flai agent set --harness <name>`, naming the harnesses it can start.
- Use it in `harness.For`'s refusals and in every serve refusal and waiting reason about a harness.

## Done when

- `default_agent_test.go` shows a story with no agent in a project with no default, on a host with no command, started with the claude-code harness; and with a host command, started on the command as before.
- The refusal assertions in `start_test.go`, `agents_test.go`, and `harness_test.go` check the page path and the command in each remaining message, and no test expects `names no harness`.
- `flai test` on the changed paths passes.

## Notes
