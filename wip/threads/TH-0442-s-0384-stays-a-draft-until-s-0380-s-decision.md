---
id: TH-0442
title: S-0384 stays a draft until S-0380's decision
anchor:
  path: wip/kanban/stories/S-0384-an-agent-s-proposal-for-the-design-system-is-a-story-whose-branch-changes-it-and-the-design-page-previews-that-branch-beside-the-main-branch-s-before-the-operator-accepts-it.md
  item: S-0384
status: resolved
participants: [orchestrator]
created: 2026-10-13T17:53:31Z
updated: 2026-10-15T00:40:26Z
---

# TH-0442 S-0384 stays a draft until S-0380's decision

On wip/kanban/stories/S-0384-an-agent-s-proposal-for-the-design-system-is-a-story-whose-branch-changes-it-and-the-design-page-previews-that-branch-beside-the-main-branch-s-before-the-operator-accepts-it.md.

## Entries

### 2026-10-13T17:53:31Z orchestrator
**Recommendation:** keep S-0384 a draft until S-0380 is decided, then finalize it as its planner replans it.

flai lists it complete (1h5m, 250.00 USD/week), but its criteria rest on assumptions S-0380 confirms or corrects (TH-0438 § Assumptions): a proposal as a story whose branch changes the document, previewed beside main's, with no second store of drafts. If I finalize it now, its planner can no longer rewrite those words without asking. S-0380's T-1531 names on TH-0438 what each draft must change. It is held behind S-0383 by its `after` meanwhile, so nothing waits on it.

### 2026-10-15T00:40:26Z orchestrator
Resolved: The operator finalized S-0384 on 2026-10-15T00:36Z. S-0380's decision is in (ADR-0146 to ADR-0152). S-0384 waits for S-0383 by its after, and its agent reviews its tasks against those ADRs and TH-0438's list when it pulls it.
