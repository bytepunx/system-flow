---
id: D-0555
title: A story's agent tells each task sub-agent never to print the process environment, since under flai serve it holds the host's API keys and a scratch test that printed it put them in the session's tool output
kind: decision
sources: [S-0628, T-2777]
topics: [delegation, testing, code]
paths: [design/conventions/safety.md, design/conventions/delegation.md]
date: 2026-11-10
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0555 A story's agent tells each task sub-agent never to print the process environment, since under flai serve it holds the host's API keys and a scratch test that printed it put them in the session's tool output

S-0628's narrative records, in its `## Decisions` (`wip/archive/agents/S-0628.md`):

```text
While working T-2777, a sub-agent's scratch test printed the host's process environment, API keys among them, into this session's tool output. Nothing reached the repository. Later sub-agents were told never to print the environment. Reported to the operator.
```

A story's agent and its sub-agents run with the host's environment, which holds the provider keys. `safety.md` keeps secrets out of the repository, narratives, and reports, but a tool's output is the session's log, which the harness keeps. So a story's agent whose task sub-agent may write scratch tests or debug output, as T-2777's did for the container close-out, says in the brief never to print the environment or a variable that holds a key, and tests read only the variables they name. When it happens anyway, the agent tells the operator, as S-0628 did; this entry records no value.

Sources: S-0628 (`wip/archive/agents/S-0628.md`, Decisions); T-2777.
