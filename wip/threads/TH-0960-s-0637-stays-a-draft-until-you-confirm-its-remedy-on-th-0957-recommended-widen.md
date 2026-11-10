---
id: TH-0960
title: "S-0637 stays a draft until you confirm its remedy on TH-0957 (recommended: widen the section merge to design and docs markdown)"
anchor:
  path: wip/kanban/stories/S-0637-stories-that-each-add-a-section-at-one-spot-in-a-shared-design-document.md
  item: S-0637
status: resolved
participants: [orchestrator, agent-S-0637]
created: 2026-11-09T22:19:42Z
updated: 2026-11-10T07:57:56Z
---

# TH-0960 S-0637 stays a draft until you confirm its remedy on TH-0957 (recommended: widen the section merge to design and docs markdown)

On wip/kanban/stories/S-0637-stories-that-each-add-a-section-at-one-spot-in-a-shared-design-document.md.

## Entries

### 2026-11-09T22:19:42Z orchestrator
**Recommended:** confirm the planner's fix on TH-0957: widen ADR-0175's section merge from markdown inside `claims.shared` to every markdown file under the manifest's design and docs folders, issue files, generated indexes, and wip left out. Then I finalize S-0637.

**Why it is not finalized yet.** `flai promote --drafts` lists S-0637 as complete, and its criteria, touches, 1h30m forecast, and 175 USD a week value are consistent. But its tasks build a remedy that refines an accepted ADR, and TH-0957 asks you to confirm it. That is a decision not yet recorded, so it is yours before an agent builds it.

**Why the fix.** All three I-0315 instances were settled the same way, keeping both sections, which the section merge already does for sections added at one spot. One matcher, `sharedMarkdown` in `flai/internal/issues/generated.go`, gates every rebase flai runs. A real same-lines conflict still stops the rebase.

**Alternative:** add each hub document to `claims.shared` by hand, which is your list and only known after a conflict.

### 2026-11-10T07:57:54Z agent-S-0637
The operator finalized S-0637 at 2026-11-10T07:44:15Z, so the remedy on TH-0957 is taken as confirmed. Building it.

### 2026-11-10T07:57:56Z agent-S-0637
Resolved: The operator finalized S-0637 at 2026-11-10T07:44:15Z.
