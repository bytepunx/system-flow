---
id: TH-1019
title: "Plan for S-0695: a 30 s timeout on repo.test.ts's monorepo describe, two tasks in two layers"
anchor:
  path: wip/kanban/stories/S-0695-repo-test-ts-s-whole-monorepo-read-times-out-at-vitest-s-5s-default-under-load.md
  item: S-0695
status: open
participants: [planner-S-0695]
created: 2026-11-10T10:29:36Z
updated: 2026-11-10T10:29:36Z
---

# TH-1019 Plan for S-0695: a 30 s timeout on repo.test.ts's monorepo describe, two tasks in two layers

On wip/kanban/stories/S-0695-repo-test-ts-s-whole-monorepo-read-times-out-at-vitest-s-5s-default-under-load.md.

## Entries

### 2026-11-10T10:29:36Z planner-S-0695
S-0695 is planned with two tasks in two layers. No input is missing.

## Proposed fix

Give the `Repo on the monorepo` describe in `flaiover/src/lib/server/repo.test.ts` `{ timeout: 30000 }`. `flaiover/src/lib/server/search.test.ts`'s `adrs on the monorepo` describe already has the same timeout. All three of I-0269's instances failed on a loaded host. The read, `flai hostapi items.list` with the archive, took 0.75 to 0.90 s here at a load average of 130 on 24 cores, so 5 s leaves too little room when vitest's workers share the host.

## Tasks

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-2969 Set the 30 s timeout, and check D-0524's concurrent-write cause | `flaiover/src/lib/server/repo.test.ts` | none |
| 2 | T-2972 Close I-0269 with the fix as its reason | the I-0269 file, `design/issues/summary.md` | T-2969 |

## Figures

- Touches: the three files above. No folder touch, and no co-change entry was relevant.
- Forecast: 14m, flai's figure, kept.
- Cost of delay: 45 USD a week, from the operator's 18m `time_lost_per_cycle` input.
- No `after`: S-0695 reads no other story's output.

## Assumptions

- A timeout that depends on load has no test that reproduces it reliably, so criterion 1's "where one fits" applies. The story's Notes record why.
- D-0524 also names concurrent writes to the main checkout. If T-2969 finds that cause real, it gets its own issue and is not fixed here.
- Moving the test out of the `flaiover-flai` tier was ruled out: that tier is where it checks `flai hostapi` against real data.
