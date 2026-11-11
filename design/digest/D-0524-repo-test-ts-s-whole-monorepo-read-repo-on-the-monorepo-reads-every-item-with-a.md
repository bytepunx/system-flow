---
id: D-0524
title: "repo.test.ts's whole-monorepo read, Repo on the monorepo reads every item with a consistent status, fails in a story's flai test or close-out while other agents write the main checkout, whatever the story changes; I-0269 records it and the close-out runs it again"
kind: cause
sources: [S-0605, I-0269, S-0695]
topics: [flaiover, code]
paths: [flaiover/src/lib/server/repo.test.ts]
date: 2026-11-08
status: active
absorbed_by: ""
asked_again: []
cited: [S-0695, T-2969, TH-1019]
---

# D-0524 repo.test.ts's whole-monorepo read, Repo on the monorepo reads every item with a consistent status, fails in a story's flai test or close-out while other agents write the main checkout, whatever the story changes; I-0269 records it and the close-out runs it again

`flaiover/src/lib/server/repo.test.ts` › Repo on the monorepo › reads every item with a consistent status lists every work item of the main checkout, archive included. It failed in stories that did not touch what it reads, at vitest's 5 s default timeout.

S-0605's narrative, in its `flaiover-flai` tier: "That test lists every item in the main checkout, which other agents were writing at the time, and settings.get is not involved. The close-out runs it again." The close-out passed. [I-0269](../issues/I-0269-repo-test-ts-s-whole-monorepo-read-times-out-at-vitest-s-5s-default-under-load.md) records the same test stopping S-0557's and S-0586's close-outs under load; each asked for the close-out again.

The cause is host load alone, not other agents' writes. S-0695 found, in its Decisions (`wip/archive/agents/S-0695.md`):

```text
D-0524's second cause, concurrent writes to the main checkout, does not hold: flai writes an item with atomicfile.WriteFile (flai/internal/workitem/itemlock.go), so a reader sees a whole file, before or after.
```

S-0695 (T-2969) fixed it: the `Repo on the monorepo` describe now has `{ timeout: 30000 }`, as `search.test.ts`'s adrs on the monorepo has. The read takes 0.75 to 0.90 s even at load 130, so 30 s leaves room without hiding a hang. I-0269 is closed. A later timeout of this test is a new fault, not this one.

Sources: S-0605; I-0269 (its cause and fix); S-0695 (`wip/archive/agents/S-0695.md`, Decisions).
