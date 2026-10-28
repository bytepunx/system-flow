---
id: T-2184
type: task
nature: remediation
title: An ADR records that flai guard refuses every served session and sub-agent the flai serve sub-commands that write the host configuration or turn a host action on or off
status: backlog
parent: S-0511
owner: alex
created: 2026-10-28T08:19:23Z
updated: 2026-10-28T08:19:23Z
transitions: []
stream: S-0511
tags: [cli, guard, adr]
touches: [design/adrs/README.md]
---
# T-2184 An ADR records that flai guard refuses every served session and sub-agent the flai serve sub-commands that write the host configuration or turn a host action on or off

## Work

Propose the solution I-0210 asks for and record it with `flai adr new` in the story's worktree, which allocates the number, writes the file under `design/adrs/`, and adds its row to `design/adrs/README.md`; commit it with the task, not with `--commit`. The decision: `flai guard` refuses, to every sub-agent and to every session `flai serve` started (`FLAI_STARTED_BY=flai-serve`, `FLAI_ROLE`, or `FLAI_STORY` set), the commands that write the host configuration, `~/.flai/config.json` or the file `FLAI_CONFIG` names, or turn a host action on or off, whatever the agent's role or permissions, as ADR-0096 does for the shared paths and ADR-0166 for `flai forge push`. Name the forms that write, from `flai/cmd/serve_actions.go`, `serve_import.go`, `config.go`, and `template.go`: `flai serve enable` and `disable`; `agent set`, `agent clear`, `agent harness` with a name and a setting or `--reset`, `agent provider` with a flag or `--reset`; `checks set`, `checks timeout`, `checks clear`; `components idle`, `components max-servers`; `forge` with a flag or `--reset`; `notify set`, `notify clear`; `import add`, `import remove`; and `flai config set` and `flai template use`, which write the same file. Name the forms that read and stay open: `serve status`, `actions`, `journal`, `agent show`, `agent harness` and `agent provider` bare or with a name alone, `checks show`, `components show`, `forge` bare or with a forge alone, `notify show`, `import list`, `config get` and `path`, and `flai serve start`, `stop`, `project`, `agent start`, `restart`, `commit`, `stop`, `stream`, `usage`, and `orchestrate`, which are not host configuration and keep their own rules. Give the context: `agentEnv` in `flai/internal/serve/agents.go` drops `FLAI_CONFIG` from a served agent (S-0183, I-0044), so such a call lands in the operator's file; `safety.md` forbids an agent to change the operator's global tool configuration; ADR-0205 says the reminder channel is set in a shell; the `settings` host action (ADR-0039) is the dashboard's, not an agent's; and `flai serve enable auto-approve` or `notify set --command` would let an agent widen its own permissions or name a program the host runs as the operator. Say that the refusal tells a sub-agent to put what it needs in its final message and a served session to ask the operator on a thread, and that the operator's own session is not held. Set its `topics` to `[cli]` with `flai adr topics`.

## Done when

The ADR is accepted, its row is in `design/adrs/README.md`, `flai check --strict` is clean, and the file names each refused form and each form left open as listed above.

## Notes

The number is one more than the highest in any checkout or branch, so do not write it into other files before `flai adr new` prints it. The tasks that follow cite the ADR by the number it gets.
