---
id: TH-1207
title: "S-0739 cannot be accepted: uncommitted digest entry D-0627 on main fails markdown lint (MD028), as on TH-1206"
anchor:
  path: wip/kanban/stories/S-0739-workitem-s-epic-following-board-view-task-plan-holds-order-and-promotion-read-a.md
  item: S-0739
status: resolved
participants: [orchestrator]
created: 2026-11-15T01:17:11Z
updated: 2026-11-15T07:31:45Z
---

# TH-1207 S-0739 cannot be accepted: uncommitted digest entry D-0627 on main fails markdown lint (MD028), as on TH-1206

On wip/kanban/stories/S-0739-workitem-s-epic-following-board-view-task-plan-holds-order-and-promotion-read-a.md.

## Entries

### 2026-11-15T01:17:11Z orchestrator
**Recommended:** fix line 20 of `design/digest/D-0627-a-workitem-repo-whose-workflow-is-nil-runs-the-one-baseline-workflow-every-such.md` in the main checkout. Join its two blockquotes by putting `>` on the blank line between them, or have the curator rewrite the entry with `flai digest edit`. Once it is fixed, I accept S-0739 and S-0738 on my next run. This is the same cause as TH-1206.

**What holds it:**

- S-0739 is otherwise ready to accept. The close-out passed at the branch head `9f905b8`. The verifier found all five criteria met. The `accept` dry run listed no blocker.
- The acceptance itself refused:

```text
design/digest/D-0627-a-workitem-repo-whose-workflow-is-nil-runs-the-one-baseline-workflow-every-such.md:20 error MD028/no-blanks-blockquote Blank line inside blockquote
correct them in the main checkout and run flai accept S-0739 again
```

**Alternative:** fix the file, then accept it yourself with `flai accept S-0739 --by alex`.

### 2026-11-15T07:31:45Z orchestrator
Resolved: Concern gone: alex fixed D-0627's MD028 on main at 18a62099c (07:04Z), and S-0738's acceptance, refused for the same cause, went through at 07:18Z. S-0739's remaining blocker, its stale verify, is TH-1216's, for its agent.
