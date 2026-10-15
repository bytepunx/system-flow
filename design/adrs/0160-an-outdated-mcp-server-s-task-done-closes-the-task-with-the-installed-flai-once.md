---
id: ADR-0160
title: "An outdated MCP server's task_done closes the task with the installed flai once it has caught up"
status: proposed
date: 2026-10-15
supersedes: []
superseded_by: []
refines: [ADR-0142, ADR-0158]
---

# ADR-0160 An outdated MCP server's task_done closes the task with the installed flai once it has caught up

## Context

[ADR-0142](0142-flai-task-done-refuses-to-close-a-task-while-the-running-flai-s-major-or-minor.md) refuses to close a task while the running flai's major or minor version is behind the newest flai release in the project's history. [ADR-0158](0158-flai-serve-checks-each-served-project-for-a-newer-flai-release-every-minute.md) has `flai serve` upgrade the host, or ask the operator, within a minute of such a release.

Neither reaches a story session that is already running. Its MCP server is a stdio child of the agent's session, started from the binary path `flai serve` gave it, and it keeps that process until the session ends. Once the host is upgraded, the file at that path is the new flai, but the server still runs the old one, so `task_done` keeps refusing. [I-0147](../issues/I-0147-the-mcp-server-s-flai-lags-a-fresh-publish-so-task-done-refuses-and-a-story-s-agent-falls-back-to-the-tree-s-cli.md) records four stories whose agents started before a release and closed every task with the tree's `scripts/flai.sh task done` instead. A project made from the template has no tree CLI to fall back to.

## Decision

When ADR-0142 would refuse, the MCP tool `task_done` closes the task with the flai installed at the server's own path, once that flai has caught up.

1. **Which flai.** The server records the path of its own binary when it starts (`os.Executable`; on Linux Go drops the ` (deleted)` that `/proc/self/exe` gains once `flai self-upgrade` has replaced the file). That path is where `flai self-upgrade` and `flai host upgrade` install.
2. **Caught up.** The server runs `<path> version --json` and reads its `version`. The installed flai has caught up when ADR-0142 would not refuse it: its major and minor versions are not below the newest release's. A patch behind, or a build that is not a release, has caught up, as ADR-0142 holds for the running flai.
3. **The close.** It runs `<path> task done <task> --json`, with `-m` and `--log` as given, in the project's main checkout, with `FLAI_AGENT` and `FLAI_SESSION` set to the server's agent and session, and answers the result it prints, field for field, with `closed_by` naming the version and path that closed it. A close that stops at a step exits non-zero and still prints its result: that is the answer, as for the running flai. Output that is not a result is the tool's error, with what the command printed.
4. **Otherwise it refuses.** When the installed flai is behind too, or its version cannot be read, `task_done` refuses as ADR-0142 does, changing nothing, and adds the installed flai's version, or why it could not be read, and to call `task_done` again once the host is upgraded.
5. **Only `task_done`.** `flai task done` in a shell already runs the installed binary, and the host channel's `task.done` runs the host's flai, which ADR-0158 restarts on the upgrade. Every other MCP tool keeps ADR-0117's warn-only rule.

## Consequences

- A story's agent whose session started before a release closes its tasks through MCP as soon as the host is upgraded, with no fallback to a tree CLI.
- A task is still never closed by a flai a minor version behind the newest release: the close runs in the installed flai, by its rules.
- The rest of an outdated server's tools still run the old flai until the session ends; the inbox's `flai_outdated` keeps saying so.
- The close costs one more process start and one `flai version` call, only on the path ADR-0142 used to refuse.

## Alternatives considered

- **Re-execute the MCP server on the new binary.** The stdio server would lose its MCP session with the agent.
- **Restart the agent's session after an upgrade.** It cuts the agent's work mid-task.
- **Delegate every write tool to the installed flai.** Wider than the issue; only `task_done` refuses.
- **Look the installed flai up on `PATH`.** The server's own path is the binary `flai serve` chose for it; `PATH` may name another.
