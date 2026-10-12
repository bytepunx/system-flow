---
id: ADR-0143
title: "flai guard refuses a story's agent, the planner, and the analyzer a call that names an author or session other than its own with --by, --placed-by, FLAI_AGENT, or FLAI_SESSION"
status: accepted
date: 2026-10-12
supersedes: []
superseded_by: []
refines: [ADR-0060, ADR-0097]
---

# ADR-0143 flai guard refuses a story's agent, the planner, and the analyzer a call that names an author or session other than its own with --by, --placed-by, FLAI_AGENT, or FLAI_SESSION

## Context

[ADR-0097](0097-permission-prompt-takes-an-answer-from-the-story-s-owner-or-the-project-s-owner.md) has `permission_prompt` take an answer on its thread from the story's owner or the project's owner, read from the entry's author. A thread entry's author is whatever the command that wrote it says: `--by`, else `FLAI_AGENT`, else the config's author ([I-0091](../issues/I-0091-a-thread-entry-s-author-is-whatever-by-says-so-an-agent-with-a-shell-could-answer-a-permission-thread-as-the-owner.md)). A story's agent with a shell could run `flai thread reply TH-nnnn "allow" --by alex`, or `FLAI_AGENT=alex flai thread reply …`, and answer its own permission thread as the owner. The same goes for every flai command that records who acted: `move`, `edit`, `criteria`, `accept`, `stream answer`, `order --placed-by`, and the `message` commands.

[ADR-0060](0060-a-claude-code-pretooluse-hook-flai-guard-refuses-any-sub-agent-s-call-that.md) has `flai guard` hold sub-agents, and later stories held the planner, the orchestrator, and the analyzer to their roles; a story's own agent's calls passed untouched but for `wait_for_events`. The orchestrator's rule already refuses it a thread call whose `--by` names another (S-0220), but none looks at `FLAI_AGENT`.

A check inside flai that compares `--by` with `FLAI_AGENT` is defeated by setting both, and would refuse the dashboard, which runs flai with `FLAI_AGENT=flaiover` and `--by=<owner>`. The guard runs in the hook's environment, which the agent's command line cannot change, so it can tell the agent's own name from the one a command names.

The orchestrator recommended, and the operator confirmed on TH-0418, refusing only a name that is not the agent's own, so that an agent that sets `FLAI_AGENT` to its own name, as `tooling.md` asks, still passes.

## Decision

`flai guard` refuses a served agent's own `Bash` call that names an author or a session other than its own. A served agent is a story's agent (no `FLAI_ROLE`, and `FLAI_STORY` or `FLAI_STARTED_BY=flai-serve` set), the planner (`FLAI_ROLE=plan`), or the analyzer (`FLAI_ROLE=analyze`).

- **Its own name** is the `FLAI_AGENT` the hook sees, as flai serve set it or the project's Claude Code settings override it, or the name flai serve gives the run: `agent-S-nnnn` for a story's agent, `planner-<item>` (from `FLAI_ITEM`) for the planner, `analyzer` for the analyzer. **Its own session** is the hook's `FLAI_SESSION`.
- **Flags.** A flai command whose `--by` or `--placed-by` names anyone else is refused. `--by` on `flai order`, a policy, and on `flai stats`, a grouping, names no author and is left alone. An empty `--by` falls back to `FLAI_AGENT` and passes.
- **Environment.** An assignment of `FLAI_AGENT` or `FLAI_SESSION` to anything but its own, the empty value included, is refused: before a command, after `export`, `declare`, `typeset`, `local`, `readonly`, or `env`, anywhere on the command line, whether or not flai follows, since a script or make target may run flai out of the guard's sight. So are `unset` and `env -u` of either, and `env -i`, which leave flai with no `FLAI_AGENT` and so with the config's author, the operator. A value the guard cannot read, such as `$VAR` or a command substitution, is not its own.
- The guard reads `sh -c`, `bash -c`, and `eval` scripts and heredocs fed to a shell as it reads the rest of the command line.
- The refusal says to leave out `--by`, or give the agent's own name, and to leave `FLAI_AGENT` and `FLAI_SESSION` as they are set.
- Not held: the orchestrator, which keeps its own rule; sub-agents, which may run no flai write at all; the operator's own session, which no `FLAI_STORY`, `FLAI_ROLE`, or `FLAI_STARTED_BY` marks; and the dashboard, which runs flai on the host outside any guard.

## Consequences

- A thread entry, a move, a tick, or a message a served agent writes through the command line carries its own name, so `permission_prompt` takes no answer it wrote as an owner's.
- The guard reads `FLAI_AGENT`, `FLAI_SESSION`, and `FLAI_ITEM` from its own environment, besides `FLAI_ROLE`, `FLAI_STORY`, and `FLAI_STARTED_BY`.
- An agent that follows a project rule to set `FLAI_AGENT` to a name of its choosing, such as its model's, is refused when it is served: flai serve has set it already, and the refusal says to leave it.
- `FLAI_AGENT=$FLAI_AGENT` is refused though it changes nothing: the guard is not a shell and cannot read the value.
- Residuals the guard does not close: a story's own agent's shell write into a thread file (S-0368's refusal covers sub-agents only); a call to the HTTP MCP server with another `X-Flai-Agent`; the orchestrator's `FLAI_AGENT` on the command line; and a name hidden from the guard's reading, in a variable, a script file, or a make target, as ADR-0060 says of every rule it keeps.
- `design/system/agent-adapters.md`, `design/system/strategic-agents.md`, `docs/users/flai.md`, and `flai guard`'s help say what it refuses.

## Alternatives considered

- Refuse `--by`, `--placed-by`, `FLAI_AGENT`, and `FLAI_SESSION` as such, as the planner drafted: it refuses agents that set `FLAI_AGENT` to their own name as `tooling.md` asks.
- Check inside flai that `--by` matches `FLAI_AGENT` when it is set (I-0091's first direction): `FLAI_AGENT=alex --by alex` passes it, and it refuses the dashboard's writes for the owner.
- Refuse only `thread reply`, `resolve`, and `confirm` (I-0091's second direction): it leaves `move`, `accept`, `criteria`, and the messages open to the same forgery.
- Hold the orchestrator to the environment rule too: the operator confirmed the three agents only; its `--by` rule stays as S-0220 made it.
