---
id: D-0639
title: A kluster send-back record keeps the state it moved the story to, klusterSentBack.To, and a record from an older flai, which has none, is judged by the baseline's send-back transition
kind: decision
sources: [S-0743]
topics: [flai, go, code, kluster]
paths: [flai/internal/serve/kluster.go]
date: 2026-11-15
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0639 A kluster send-back record keeps the state it moved the story to, klusterSentBack.To, and a record from an older flai, which has none, is judged by the baseline's send-back transition

S-0743 moved flai serve's kluster pass onto a state's kind and a transition's effect (ADR-0301), so a story that fails its deploy or its checks run is sent back along the workflow's send-back transition, which need not end in `in-progress`. Its narrative's Decisions (`wip/archive/agents/S-0743.md`):

> A kluster send-back record now keeps the state it moved the story to (klusterSentBack.To); a record from an older flai, which ran only the baseline, is judged by the baseline's send-back effect.

`klusterSentBack` in `flai/internal/serve/kluster.go` gained `To` (`json:"to,omitempty"`) beside its head, thread, reason, and time. So a later story that reads the pass's record does not assume the story went back to `in-progress`: it reads `To`, and only when `To` is empty, a record an older flai wrote, falls back to where the baseline's send-back goes.

Source: S-0743's narrative, `wip/archive/agents/S-0743.md`, § Decisions.
