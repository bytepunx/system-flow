---
id: TH-0654
title: "The one-shot search index build costs about 2 s of CPU on the monorepo: make it cheaper in a story of its own?"
anchor:
  path: wip/kanban/stories/S-0501-search-test-ts-finds-a-body-phrase-across-wip-and-design-times-out-at-5-s-on-the-monorepo-while-other-stories-close-outs-run.md
  item: S-0501
status: resolved
participants: [agent-S-0501, orchestrator, alex]
created: 2026-10-26T04:35:36Z
updated: 2026-10-26T16:03:01Z
---

# TH-0654 The one-shot search index build costs about 2 s of CPU on the monorepo: make it cheaper in a story of its own?

On wip/kanban/stories/S-0501-search-test-ts-finds-a-body-phrase-across-wip-and-design-times-out-at-5-s-on-the-monorepo-while-other-stories-close-outs-run.md.

## Entries

### 2026-10-26T04:35:36Z agent-S-0501
**Finding.** One `bin/flai hostapi search.query` on the monorepo (4,845 markdown files) takes 1.8 to 2.0 s wall and 1.9 s of user CPU; flai's own start takes 17 ms, so nearly all of it is building the search index from nothing. That cost is CPU, not waiting, so it holds on an idle host too. 24 at once on this 24-core host took 4.8 to 6.1 s each, past vitest's 5 s, which is I-0185.

**What S-0501 does regardless.** The `adrs on the monorepo` block in `flaiover/src/lib/server/search.test.ts` gets a 30 s timeout with these figures in a comment, as `writes.test.ts` did in S-0116. That closes I-0185. The index code in `flai/internal/hostapi/searching.go` stays as it is; the plan said to ask rather than widen the story.

**Question.** Should the one-shot build be made cheaper, for example by keeping the index on disk under `.flai-cache` keyed by its signature, so that `flai hostapi`, `flai doc search`, and a restarted `flai serve` skip the rebuild?

Recommended: no story now. `flai serve` builds the index once and keeps it, so the dashboard does not pay this per query. Only the test and one-shot CLI calls pay the 2 s, and the timeout covers the test. Make a story only if a one-shot caller turns up that a person waits on.

Either answer leaves S-0501 as it is. If you want the story, I will record an issue for it so you can pick it at acceptance.

### 2026-10-26T04:35:44Z orchestrator (recommendation)
**Recommended for the operator:** take agent-S-0501's answer. Make no story now. S-0501 stays as planned: a 30 s timeout with the measured figures, then close I-0185.

- **Why no story.** `flai serve` builds the index once and keeps it, so neither the dashboard nor anyone waiting on a query pays the 2 s per call. Only the test and one-shot CLI calls pay it, and the timeout covers the test.
- **Why it is yours.** Whether to add a story is a change of scope. The work-management rules keep S-0501 to its goal and say not to invent scope.
- **If you want it later:** say so here, and S-0501's agent records an issue so you can choose a story from it at acceptance.

Source: design/conventions/work-management.md § Rules

### 2026-10-26T16:03:01Z alex
Resolved: S-0501 was accepted
