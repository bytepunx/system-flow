---
id: D-0549
title: A flai test that checks this monorepo's own docs, such as TestReferenceIsCurrent, runs under -short and fails, not skips, when the docs are absent, since the flai module's tests run only in this monorepo
kind: decision
sources: [S-0634, I-0308, I-0327, S-0827]
topics: [testing, go]
paths: [flai/cmd/reference_test.go, docs/users/flai-reference.md, docs/operators/settings.md, flai/internal/verify/select_test.go]
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

So a later test that reads this repository's own documents, and must run in `flai test`'s go-test or settings-index tier, runs under `-short` and fails when what it reads is missing. A skip there would let a stale document through to the close-out, as I-0308 records.

I-0327 recorded the same gap for `TestRepositoryManifestLintsChangedFlaioverFiles`. S-0827 closed it in T-3523 the same way. From I-0327's fix: "S-0827 (T-3523) removed the testing.Short() gate from TestRepositoryManifestLintsChangedFlaioverFiles in flai/internal/verify/select_test.go. The test reads only system-flow.yaml and the file tree, keeps its monorepo-not-present skip, and now runs in make test and flai test's go-test tier, so a tier it must list fails there rather than first in the close-out's integration tier." That test keeps its skip for when the monorepo is absent, so a test reading only the manifest and the tree may keep that one skip.

Sources: S-0634 (`wip/archive/agents/S-0634.md`, Decisions); I-0308; I-0327 (its fix); S-0827.
