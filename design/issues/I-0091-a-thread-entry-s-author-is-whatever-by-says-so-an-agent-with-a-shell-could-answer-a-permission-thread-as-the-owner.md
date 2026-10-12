---
id: I-0091
title: A thread entry's author is whatever --by says, so an agent with a shell could answer a permission thread as the owner
class: impression
status: closed
count: 1
first_reported: 2026-10-06T19:45:52Z
last_reported: 2026-10-06T19:45:52Z
updated: 2026-10-12T18:41:24Z
---

# I-0091 A thread entry's author is whatever --by says, so an agent with a shell could answer a permission thread as the owner

## Description

The operator approves an agent's write under `.claude/` by answering `allow` on a thread, and flai knows the operator by the name on the entry. That name is not checked. `flai thread reply --by <name>` writes whatever name it is given, and a story's agent runs flai through its shell without `flai guard` holding its own calls. An agent that replied `allow --by alex` on its own permission thread would be let through.

No agent has done this. The conventions tell an agent never to work around a refusal, and the agents seen on 2026-10-06 waited or asked. The same is true of every decision flai reads from a thread's author: an answer that restarts an agent, and a recommendation's confirmation.

## Instances

### 2026-10-06T19:45:52Z
Story: S-0284.
Found reviewing S-0284, not seen happening. permission_prompt lets a write under .claude/ through when an entry by the story's owner or the project's owner starts with allow (ADR-0086, ADR-0097). flai thread reply takes its author from --by as given (threadAuthor in flai/cmd/thread.go), and flai guard does not hold a story's own agent's shell calls. So the approval rests on a name any session can type; only the conventions keep an agent from typing it.

## Remediation

Directions to weigh: flai refuses a `--by` that is not the session's own name when `FLAI_AGENT` is set, so that an agent's entries always carry its name; and `flai guard` refuses a story's agent `flai thread reply`, `resolve`, and `confirm` with a `--by` at all. Either leaves the operator's own session, where `FLAI_AGENT` is unset, free to name itself. The dashboard already writes the operator's name itself.

Story S-0369 remediates this issue, created from it at 2026-10-12T16:52:42Z.
Closed 2026-10-12T18:41:24Z: Fixed by S-0369 (ADR-0143): flai guard refuses a story's agent, the planner, and the analyzer a Bash call whose flai --by or --placed-by names anyone but the agent, or that sets, empties, or unsets FLAI_AGENT or FLAI_SESSION to anything but its own, in env, export and the other declarers, shell -c scripts, eval, and heredocs a shell reads, so none can answer its own permission thread as the owner. flai/internal/guard/authors_test.go and flai/cmd/guard_test.go (TestGuardHoldsAServedAgentToItsOwnName) reproduce the cause; threadAuthor now falls back to the config's author, so the operator's own shell writes as the owner. Left open, as ADR-0143 says: a story's own agent's shell write into a thread file, an HTTP MCP call with another X-Flai-Agent, the orchestrator's FLAI_AGENT on the command line, and a name hidden in a variable, a script file, or a make target.
