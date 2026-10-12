---
id: T-1465
type: task
nature: improvement
title: An ADR refining ADR-0060 and ADR-0097 decides that flai guard refuses a served agent a flai command that names who wrote it, and agent-adapters.md and strategic-agents.md record it
status: backlog
parent: S-0369
owner: alex
created: 2026-10-12T17:05:34Z
updated: 2026-10-12T17:05:34Z
transitions: []
stream: S-0369
tags: [flai, guard, adr]
touches: [design/adrs/README.md, design/system/agent-adapters.md, design/system/strategic-agents.md]
---
# T-1465 An ADR refining ADR-0060 and ADR-0097 decides that flai guard refuses a served agent a flai command that names who wrote it, and agent-adapters.md and strategic-agents.md record it

## Work

Write an ADR with `flai adr new` that refines ADR-0060 (the guard) and ADR-0097 (whose answer `permission_prompt` takes). It decides how I-0091's cause is closed. Weigh the issue's two directions against what the code does today:

- `flai guard` never looks at a story's own agent's `Bash` calls: `Decide` in `flai/internal/guard/guard.go` runs only `wait` for it, and `TestTheStorysAgentIsNeverRefused` asserts that.
- The planner's rules let `flai thread new` and `reply` pass with any `--by`. The analyzer's are to be checked. The orchestrator's `thread` rule already refuses a `--by` naming another (`writesAsSelf`).
- A check inside flai that compares `--by` to `FLAI_AGENT` is defeated by `FLAI_AGENT=alex flai thread reply …`. It would also refuse the dashboard, which runs flai with `FLAI_AGENT=flaiover` and `--by=<owner>` (`flai/internal/hostapi/writes.go`).
- With no `--by`, `threadAuthor` takes `FLAI_AGENT` from the command's own environment. So an env prefix forges the author as well as `--by` does.
- In this project `.claude/settings.local.json` sets `FLAI_AGENT=system-flow`, so the hook's own `FLAI_AGENT` may not be the run's name. `FLAI_STORY` and `FLAI_ROLE`, which flai serve sets, are what the guard can trust.

The planned direction, to be confirmed or replaced by the ADR: the guard refuses a served agent (a story's agent, the planner, the analyzer) any `flai` command that carries `--by` or `--placed-by`, or that sets `FLAI_AGENT` or `FLAI_SESSION` for flai. The commands are `thread new`, `reply`, `confirm`, `resolve`, `message`, `stream answer`, `accept`, `move`, `edit`, `criteria`, and `order`. The orchestrator keeps its rule. The operator's own session, which no `FLAI_STORY` or `FLAI_ROLE` marks, is not held. The ADR names the residuals the guard does not close: a shell redirection into a thread file (S-0368 covers sub-agents only) and a call to the HTTP MCP server with another `X-Flai-Agent`.

Record the decision in `design/system/agent-adapters.md` under `### Permissions` and `### Asking the owner before a protected write`. Record it in `design/system/strategic-agents.md` under the planner's `### The guard` and `### The analyzer's guard`. Add the ADR to `design/adrs/README.md`. The ADR's own file cannot be named before `flai adr new` numbers it; `flai task done` widens this task's touches with it.

## Done when

- An accepted ADR says which agents the guard holds to their own name, which commands and environment it refuses, what it leaves to the operator's session and the dashboard, and the residuals.
- `design/system/agent-adapters.md` and `design/system/strategic-agents.md` say the same and link the ADR.
- `design/adrs/README.md` lists it, and `flai check --strict` is clean.

## Notes

First layer, with the `threadAuthor` task: the guard task waits for this decision.
