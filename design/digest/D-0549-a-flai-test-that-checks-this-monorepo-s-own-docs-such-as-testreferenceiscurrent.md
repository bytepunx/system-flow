---
id: D-0549
title: A flai test that checks this monorepo's own docs, such as TestReferenceIsCurrent, runs under -short and fails, not skips, when the docs are absent, since the flai module's tests run only in this monorepo
kind: decision
sources: [S-0634, I-0308]
topics: [testing, go]
paths: [flai/cmd/reference_test.go, docs/users/flai-reference.md, docs/operators/settings.md]
date: 2026-11-10
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0549 A flai test that checks this monorepo's own docs, such as TestReferenceIsCurrent, runs under -short and fails, not skips, when the docs are absent, since the flai module's tests run only in this monorepo

S-0634 remediated I-0308: `flai test` on a change to the manifest's settings catalog did not run `TestReferenceIsCurrent`, so a stale `docs/users/flai-reference.md` first failed in the close-out's integration tier. The two reference tests in `flai/cmd/reference_test.go`, `TestReferenceIsCurrent` and `TestFlagIndexIsCurrent`, lost their `-short` gate.

The story chose to drop the gate outright rather than skip when the docs are missing. From S-0634's Decisions (`wip/archive/agents/S-0634.md`):

```text
The two reference tests lose their `-short` gate outright, rather than skipping when the docs are absent: `TestSettingsIndexHasEveryKey`, the precedent the planner names, reads `docs/operators/settings.md` under `-short` and fails, not skips, when it is absent, and the `flai` module's tests run only in this monorepo.
```

So a later test that reads this repository's own documents, and must run in `flai test`'s go-test or settings-index tier, runs under `-short` and fails when what it reads is missing. A skip there would let a stale document through to the close-out, as I-0308 records. I-0327 records the same gap for `TestRepositoryManifestLintsChangedFlaioverFiles`, which still skips under `-short`.

Sources: S-0634 (`wip/archive/agents/S-0634.md`, Decisions); I-0308.
