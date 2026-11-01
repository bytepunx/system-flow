---
id: D-0445
title: flai guard refuses a story's sub-agent flai adr index, so the story's agent runs it after a task changes the ADR files
kind: decision
sources: [S-0550]
topics: [flai, cli, conventions, code]
paths: [design/adrs/README.md, flai/internal/guard/guard.go, design/conventions/delegation.md]
date: 2026-11-01
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0445 flai guard refuses a story's sub-agent flai adr index, so the story's agent runs it after a task changes the ADR files

S-0550 made `design/adrs/README.md` an index generated from the ADR files, written by `flai adr index` (ADR-0228 §1). Its narrative records, for T-2405:

```text
T-2405: `flai guard` refuses a sub-agent `flai adr index`, so the story's agent ran it; the regeneration changed exactly the 11 cells ADR-0228 lists.
```

ADR-0127 lets a story's sub-agent run `flai adr new`, `topics`, and `accept`, which write the index themselves, but not `flai adr index`. So when a task sub-agent's work leaves the index out of step, the story's agent runs `flai adr index` itself, or `flai check --fix`, rather than hand it to the sub-agent.

Source: S-0550's narrative, `wip/archive/agents/S-0550.md` § Decisions.
