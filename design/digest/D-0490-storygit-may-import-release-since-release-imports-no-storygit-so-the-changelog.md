---
id: D-0490
title: storygit may import release, since release imports no storygit, so the changelog merge after a rebase is release.SettleChangelog, a pure function storygit calls
kind: decision
sources: [S-0570]
topics: [flai, go, code, git, release]
paths: [flai/internal/release/changelog.go, flai/internal/storygit/sync.go, flai/internal/mergequeue/batch.go]
date: 2026-11-07
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0490 storygit may import release, since release imports no storygit, so the changelog merge after a rebase is release.SettleChangelog, a pure function storygit calls

S-0570 fixed I-0251 as ADR-0257 says: after a clean rebase, each `CHANGELOG.md` the branch changed is written as main has it, with the branch's own lines under main's top `## Unreleased`. Its narrative records where the code sits: "The merge is `release.SettleChangelog`, a pure function beside `cutUnreleased`; the git step is one exported function in `storygit`, called from `Rebase` (sync, task done's sync, accept) and from the merge queue's batch rebase. `storygit` may import `release`: `release` imports no `storygit`."

So the import goes one way: `storygit` imports `release`, and `release` must not import `storygit`, which would make a cycle. A later change to how a changelog is merged belongs in `release` as a pure function over the file's text. The git work that reads and commits the file stays in `storygit`. This is the same kind of seam D-0447 records for `storygit` and `issues`.

Source: S-0570 (wip/archive/agents/S-0570.md, Decisions); ADR-0257; I-0251.
