---
id: TH-0655
title: S-0501 is verified but cannot enter the merge queue while TH-0654 is open
anchor:
  path: wip/kanban/stories/S-0501-search-test-ts-finds-a-body-phrase-across-wip-and-design-times-out-at-5-s-on-the-monorepo-while-other-stories-close-outs-run.md
  item: S-0501
status: open
participants: [orchestrator]
created: 2026-10-26T04:40:16Z
updated: 2026-10-26T04:40:16Z
---

# TH-0655 S-0501 is verified but cannot enter the merge queue while TH-0654 is open

On wip/kanban/stories/S-0501-search-test-ts-finds-a-body-phrase-across-wip-and-design-times-out-at-5-s-on-the-monorepo-while-other-stories-close-outs-run.md.

## Entries

### 2026-10-26T04:40:16Z orchestrator
**Recommended:** answer TH-0654 ("no story now" is my recommendation there), and resolve it. I then enqueue S-0501 at once.

## What blocks it

queue_add refused S-0501:

```text
not eligible for the merge queue: thread_open: threads TH-0653, TH-0654 are not resolved on S-0501 or its tasks
```

I resolved TH-0653, the planner's plan thread, since only notes were left on it. TH-0654 asks you whether to make the one-shot search index cheaper in a story of its own. That is a scope question, so it is yours. Either answer leaves S-0501 as it is.

## What is ready

- flai verify passed at the head, 042d7aea, including the flaiover tier.
- The verifier matched both criteria. `search.test.ts` gives the monorepo block a 30 s timeout with the measured cause. I-0185 is closed with its reason, and the story's notes say why no reproducing test fits.
- Every changed file is within touches, and nothing is under `.claude/`.
- The close-out recorded I-0208: main's "lists ADRs with the supersession chain" test fails since ADR-0202 became `rejected`. S-0501 adds `rejected` to that test's list so its own tier passes. Whether `rejected` is a named ADR status is left open on I-0208.

`design/issues/summary.md` changed: I-0185 left the open table and I-0208 joined it.
