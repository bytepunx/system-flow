---
id: D-0285
title: A thread entry's author is whatever --by says, so an agent with a shell could answer a permission thread as the owner
kind: cause
sources: [I-0091]
topics: [cli, conventions, flai, go, code, all]
paths: [design/adrs/README.md, design/system/agent-adapters.md, design/system/strategic-agents.md, flai/internal/guard/guard.go, flai/internal/guard/guard_test.go, flai/cmd/guard.go, flai/cmd/guard_test.go, flai/cmd/thread.go, flai/cmd/thread_test.go, docs/users/flai-reference.md, docs/users/flai.md, design/issues/I-0091-a-thread-entry-s-author-is-whatever-by-says-so-an-agent-with-a-shell-could-answer-a-permission-thread-as-the-owner.md, design/issues/summary.md, flai/internal/guard/authors.go, flai/internal/guard/authors_test.go, design/adrs/0143-flai-guard-refuses-a-story-s-agent-the-planner-and-the-analyzer-a-call-that.md]
date: 2026-10-12
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0285 A thread entry's author is whatever --by says, so an agent with a shell could answer a permission thread as the owner

I-0091 was closed on 2026-10-12.

The cause:

```text
The operator approves an agent's write under `.claude/` by answering `allow` on a thread, and flai knows the operator by the name on the entry. That name is not checked. `flai thread reply --by <name>` writes whatever name it is given, and a story's agent runs flai through its shell without `flai guard` holding its own calls. An agent that replied `allow --by alex` on its own permission thread would be let through.

No agent has done this. The conventions tell an agent never to work around a refusal, and the agents seen on 2026-10-06 waited or asked. The same is true of every decision flai reads from a thread's author: an answer that restarts an agent, and a recommendation's confirmation.
```

The fix:

```text
Fixed by S-0369 (ADR-0143): flai guard refuses a story's agent, the planner, and the analyzer a Bash call whose flai --by or --placed-by names anyone but the agent, or that sets, empties, or unsets FLAI_AGENT or FLAI_SESSION to anything but its own, in env, export and the other declarers, shell -c scripts, eval, and heredocs a shell reads, so none can answer its own permission thread as the owner. flai/internal/guard/authors_test.go and flai/cmd/guard_test.go (TestGuardHoldsAServedAgentToItsOwnName) reproduce the cause; threadAuthor now falls back to the config's author, so the operator's own shell writes as the owner. Left open, as ADR-0143 says: a story's own agent's shell write into a thread file, an HTTP MCP call with another X-Flai-Agent, the orchestrator's FLAI_AGENT on the command line, and a name hidden in a variable, a script file, or a make target.
```

Source: I-0091.
