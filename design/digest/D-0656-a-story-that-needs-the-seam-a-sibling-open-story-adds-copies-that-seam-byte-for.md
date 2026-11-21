---
id: D-0656
title: A story that needs the seam a sibling open story adds copies that seam byte for byte into its own branch rather than writing one of its own, so the two branches merge cleanly whichever lands first
kind: decision
sources: [S-0742, S-0743, S-0810, TH-1318]
topics: [planning, flai, go, code]
paths: [flai/internal/serve/flow.go, flai/internal/serve/kinds_test.go, flai/internal/verify/select.go]
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

S-0810 did the same with S-0809's suite filter. From its Decisions (`wip/archive/agents/S-0810.md`):

> The filter mirrors S-0809's verify.SelectStorySuites and TierOptions.Suites line for line (commit 7b15a85f7 on story/S-0809), so the two branches merge cleanly whichever lands first.

When main had moved on by its rebase (TH-1318), main's field of the same meaning, S-0807's `StoryOptions.Suites`, replaced S-0810's own, and S-0810 kept only what it added on top.

So a story that needs a helper or seam another open story adds copies it exactly, the same lines in the same file, and tells the other story in their conversation (ADR-0121). It copies nothing that depends on code only the other branch has, such as a test calling a function the other story adds. A fixture it needs goes into a file of its own under its own name, not into the other story's test file. When the other branch lands first with a field or function of the same meaning, the rebase takes main's and drops the copy.

Sources: S-0742's and S-0810's narratives, `wip/archive/agents/S-0742.md` and `wip/archive/agents/S-0810.md`, § Decisions; S-0743; TH-1318.
