---
id: D-0662
title: Since main's 4052fc603 this repository has no design/workflows/formulas/publish.md, so the orchestrator publishes with release_publish, and a test that reads the publish formula reads the template's copy until S-0837 retires both
kind: decision
sources: [S-0827, I-0401, S-0837]
topics: [release, orchestration, testing, go]
paths: [design/workflows/formulas/publish.md, template/root/design/workflows/formulas/publish.md, flai/cmd/publish_judgement_test.go]
date: 2026-11-17
status: active
absorbed_by: ""
asked_again: []
cited: [S-0811, T-3443, TH-1250]
---

# D-0662 Since main's 4052fc603 this repository has no design/workflows/formulas/publish.md, so the orchestrator publishes with release_publish, and a test that reads the publish formula reads the template's copy until S-0837 retires both

The operator's commit 4052fc603 on main, "chore: remove the publish formula so the orchestrator publishes with release_publish (stopgap for S-0837)", deleted `design/workflows/formulas/publish.md` from this repository. Its message says: "With no publish formula the orchestrator's role file has it call release_publish, which makes no story. ADR-0316 and S-0837 retire the formula in flai and the template; five documents still link the file until then."

So, until S-0837 lands, the template still ships `template/root/design/workflows/formulas/publish.md`, and this repository's own copy is gone. A test that read this repository's copy failed: `TestPublishFormulaUnderJudgement` (`flai/cmd/publish_judgement_test.go`) failed main's go-test tier and every close-out whose change selects `./cmd` (I-0401). S-0827 fixed it in T-3557 by pointing the test at the template's copy. From its Decisions (`wip/archive/agents/S-0827.md`):

```text
Added T-3557 rather than waiting for S-0837. Main's 4052fc603 removed design/workflows/formulas/publish.md, so TestPublishFormulaUnderJudgement fails on main, and S-0827's change selects ./cmd, so its close-out would fail too (I-0401). The test now reads the template's copy, the formula the template ships, as its comment always said. This keeps the operator's stopgap in place, and S-0837 retires both the test and the copy.
```

Do not restore this repository's copy. A story that touches `design/workflows/formulas/publish.md`, or a test that reads it, works against the template's copy or leaves it to S-0837. I-0401 also names `flai/cmd/orchestrate_permissions_test.go` and `flai/internal/mcpserver/orchestrate_test.go`, which name the file. The older entries on publish formula instances (D-0603, D-0610, D-0613) describe runs from before the stopgap.

Sources: S-0827 (`wip/archive/agents/S-0827.md`, Decisions); I-0401; S-0837; ADR-0316.
