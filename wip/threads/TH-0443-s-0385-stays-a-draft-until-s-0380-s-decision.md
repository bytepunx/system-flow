---
id: TH-0443
title: S-0385 stays a draft until S-0380's decision
anchor:
  path: wip/kanban/stories/S-0385-flai-writes-the-design-system-s-tokens-into-the-stylesheet-the-manifest-names-and-flai-check-reports-the-stylesheet-when-it-falls-behind-them.md
  item: S-0385
status: resolved
participants: [orchestrator]
created: 2026-10-13T17:53:34Z
updated: 2026-10-15T00:40:27Z
---

# TH-0443 S-0385 stays a draft until S-0380's decision

On wip/kanban/stories/S-0385-flai-writes-the-design-system-s-tokens-into-the-stylesheet-the-manifest-names-and-flai-check-reports-the-stylesheet-when-it-falls-behind-them.md.

## Entries

### 2026-10-13T17:53:34Z orchestrator
**Recommendation:** keep S-0385 a draft until S-0380 is decided, then finalize it as its planner replans it.

flai lists it complete (27m, 103.85 USD/week), but its criteria rest on assumptions S-0380 confirms or corrects (TH-0438 § Assumptions): tokens reaching code through one generated stylesheet the manifest names, checked by `design.stale`. If I finalize it now, its planner can no longer rewrite those words without asking. S-0380's T-1531 names on TH-0438 what each draft must change. It is held behind S-0381 by its `after` meanwhile, so nothing waits on it.

### 2026-10-15T00:40:27Z orchestrator
Resolved: The operator finalized S-0385 on 2026-10-15T00:37Z. S-0380's decision is in (ADR-0146 to ADR-0152). S-0385 waits for S-0381 by its after, and its agent reviews its tasks against those ADRs and TH-0438's list when it pulls it.
