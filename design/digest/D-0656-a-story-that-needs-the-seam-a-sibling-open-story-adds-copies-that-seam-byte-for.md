---
id: D-0656
title: A story that needs the seam a sibling open story adds copies that seam byte for byte into its own branch rather than writing one of its own, so the two branches merge cleanly whichever lands first
kind: decision
sources: [S-0742, S-0743]
topics: [planning, flai, go, code]
paths: [flai/internal/serve/flow.go, flai/internal/serve/kinds_test.go]
date: 2026-11-17
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0656 A story that needs the seam a sibling open story adds copies that seam byte for byte into its own branch rather than writing one of its own, so the two branches merge cleanly whichever lands first

S-0742 and S-0743 were open at once, and both needed flai serve to open its repositories with the workflow its options give. S-0743 added that seam: `Options.Workflow` and `openRepo` in `flai/internal/serve/flow.go`. S-0742's narrative's Decisions (`wip/archive/agents/S-0742.md`) record how it took the seam:

> S-0742 takes S-0743's seam (`Options.Workflow`, `openRepo` in `flow.go`) byte for byte rather than a seam of its own, so that the two branches merge cleanly whichever lands first; told on MS-0567. S-0743's `flow_test.go` is not copied, since it calls S-0743's `workflow.AwaitsAcceptance`; its fixture is copied as `kindsFlow` in `kinds_test.go`.

S-0743 was accepted first, and S-0742 landed after it.

So a story that needs a helper or seam another open story adds copies it exactly, the same lines in the same file, and tells the other story in their conversation (ADR-0121). It copies nothing that depends on code only the other branch has, such as a test calling a function the other story adds. A fixture it needs goes into a file of its own under its own name, not into the other story's test file.

Sources: S-0742's narrative, `wip/archive/agents/S-0742.md`, § Decisions; S-0743.
