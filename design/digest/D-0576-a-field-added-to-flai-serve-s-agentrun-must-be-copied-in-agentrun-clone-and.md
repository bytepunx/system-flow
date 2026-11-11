---
id: D-0576
title: A field added to flai serve's AgentRun must be copied in AgentRun.clone, and TestAgentRunCloneCoversEveryField lists the fields clone copies, so the change touches agents_state.go and agents_test.go too
kind: decision
sources: [S-0665]
topics: [go, code, testing]
paths: [flai/internal/serve/agents_state.go, flai/internal/serve/agents_test.go]
date: 2026-11-10
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0576 A field added to flai serve's AgentRun must be copied in AgentRun.clone, and TestAgentRunCloneCoversEveryField lists the fields clone copies, so the change touches agents_state.go and agents_test.go too

S-0665 added the `Image`, `Container`, and `Mounts` fields to a run's record. Its task T-2943 had to change two files its touches did not name, from S-0665's Decisions (`wip/archive/agents/S-0665.md`):

```text
T-2943 also changed `flai/internal/serve/agents_state.go` and `agents_test.go`, outside its touches: `AgentRun.clone` must copy the new `Mounts` slice, and `TestAgentRunCloneCoversEveryField` lists the fields it copies.
```

So a story or task that adds a field to `AgentRun` claims `flai/internal/serve/agents_state.go` and `flai/internal/serve/agents_test.go` from the start: `clone` copies the new field, a slice as a copy of its own, and the field joins the list in `TestAgentRunCloneCoversEveryField`.

Source: S-0665 (T-2943).
