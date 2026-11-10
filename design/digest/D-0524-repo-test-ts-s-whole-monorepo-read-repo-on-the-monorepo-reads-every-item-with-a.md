---
id: D-0524
title: "repo.test.ts's whole-monorepo read, Repo on the monorepo reads every item with a consistent status, fails in a story's flai test or close-out while other agents write the main checkout, whatever the story changes; I-0269 records it and the close-out runs it again"
kind: cause
sources: [S-0605, I-0269]
topics: [flaiover, code]
paths: [flaiover/src/lib/server/repo.test.ts]
date: 2026-11-08
status: active
absorbed_by: ""
asked_again: []
cited: [S-0695, T-2969, TH-1019]
---

# D-0524 repo.test.ts's whole-monorepo read, Repo on the monorepo reads every item with a consistent status, fails in a story's flai test or close-out while other agents write the main checkout, whatever the story changes; I-0269 records it and the close-out runs it again

`flaiover/src/lib/server/repo.test.ts` › Repo on the monorepo › reads every item with a consistent status lists every work item of the main checkout, archive included. While other stories' agents write items there, or the host is loaded, it fails or times out at vitest's 5 s default in a story that does not touch what it reads.

S-0605's narrative, in its `flaiover-flai` tier: "That test lists every item in the main checkout, which other agents were writing at the time, and settings.get is not involved. The close-out runs it again." The close-out passed.

[I-0269](../issues/I-0269-repo-test-ts-s-whole-monorepo-read-times-out-at-vitest-s-5s-default-under-load.md) records the same test stopping S-0557's and S-0586's close-outs under load; each asked for the close-out again.
