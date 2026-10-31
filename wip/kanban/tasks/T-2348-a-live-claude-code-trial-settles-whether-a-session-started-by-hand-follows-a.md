---
id: T-2348
type: task
nature: research
title: "A live Claude Code trial settles whether a session started by hand follows a symbolic link under .claude/agents, and an ADR refining ADR-0207 § 4 records where the generated definition lives"
status: backlog
parent: S-0553
owner: alex
created: 2026-10-30T20:47:05Z
updated: 2026-10-30T20:47:05Z
transitions: []
stream: S-0553
tags: [cli]
touches: [design/adrs/README.md, design/system/agent-adapters.md]
---
# T-2348 A live Claude Code trial settles whether a session started by hand follows a symbolic link under .claude/agents, and an ADR refining ADR-0207 § 4 records where the generated definition lives

## Work

In a scratch repository outside this one, give `claude` a `.claude/agents/<name>.md` that is a symbolic link to a definition file elsewhere in the repository. Start a session by hand and see whether it lists and runs the agent, with `--agent <name>` and as a sub-agent. Record the version of Claude Code, the commands, and what it did.

Write an ADR refining ADR-0207 § 4 from the result. It says where `flai role definition --write` puts a generated definition, outside `.claude/`. It says what the operator does once: create the link if Claude Code follows one, else start a hand session with `--agents`. It also says what `flai check` accepts for a built-in role's `.claude/agents/<name>.md`. Add it to `design/adrs/README.md`. Rewrite the role-definition section of `design/system/agent-adapters.md` to the decision.

First layer: it waits for nothing and shares no path with the generation task. The command, the `.claude/agents` pair, and the docs wait for it.

## Done when

The ADR is accepted with the trial's observations in its context. `design/adrs/README.md` lists it. `agent-adapters.md` names one way, not two. `flai check --strict` reports nothing.

## Notes

Drafted by the planner of E-0024 on 2026-10-30. ADR-0207 § 4 left the question to S-0419, which did not generate the definitions (TH-0760). The ADR file's name is not known until it is written, so the story keeps `design/adrs` as a folder touch and `flai task done` widens this task's.
