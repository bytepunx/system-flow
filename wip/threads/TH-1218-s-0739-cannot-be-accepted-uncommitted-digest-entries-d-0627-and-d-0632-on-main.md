---
id: TH-1218
title: "S-0739 cannot be accepted: uncommitted digest entries D-0627 and D-0632 on main fail markdown lint (MD028)"
anchor:
  path: wip/kanban/stories/S-0739-workitem-s-epic-following-board-view-task-plan-holds-order-and-promotion-read-a.md
  item: S-0739
status: resolved
participants: [orchestrator]
created: 2026-11-15T07:45:19Z
updated: 2026-11-15T18:43:21Z
---

# TH-1218 S-0739 cannot be accepted: uncommitted digest entries D-0627 and D-0632 on main fail markdown lint (MD028)

On wip/kanban/stories/S-0739-workitem-s-epic-following-board-view-task-plan-holds-order-and-promotion-read-a.md.

## Entries

### 2026-11-15T07:45:19Z orchestrator
**Recommended:** fix line 20 of two digest entries in the main checkout. Each has a blank line between two blockquotes. Put `>` on that line, or drop it, or have the curator rewrite the entries with `flai digest edit`. Then commit them. Once they are fixed, I accept S-0739 on my next run.

- `design/digest/D-0627-a-workitem-repo-whose-workflow-is-nil-runs-the-one-baseline-workflow-every-such.md`: changed again, uncommitted. Its `cited:` list grew, and the edit brings it back into the acceptance's lint. Line 20 still fails, as it did on TH-1207.
- `design/digest/D-0632-a-workflow-rule-s-check-gives-one-refusal-text-and-one-warning-text-the-same.md`: a new entry, not yet committed, with the same blank line between its two blockquotes.

**What holds it:**

- Apart from this, S-0739 is ready to accept. Its close-out passed every step at the branch head `da74987`, over main `9fb9d24`. The verifier found all five criteria met. The `accept` dry run listed no blocker.
- The acceptance itself was refused:

```text
design/digest/D-0627-a-workitem-repo-whose-workflow-is-nil-runs-the-one-baseline-workflow-every-such.md:20 error MD028/no-blanks-blockquote Blank line inside blockquote
design/digest/D-0632-a-workflow-rule-s-check-gives-one-refusal-text-and-one-warning-text-the-same.md:20 error MD028/no-blanks-blockquote Blank line inside blockquote
correct them in the main checkout and run flai accept S-0739 again
```

**Recurring:** this is the third time the curator's digest entries have held an acceptance (TH-1206 to TH-1212). `flai digest new` and `flai digest edit` should refuse an entry that MD028 fails, as flai already does for its own `wip` writes. That would be a story.

**Alternative:** fix the files, then accept it yourself with `flai accept S-0739 --by alex`.

### 2026-11-15T18:43:21Z orchestrator
Resolved: D-0627 and D-0632 are committed on main (19b4036a5, 7ecfcb8c2), so the acceptance no longer brings them in uncommitted. The other uncommitted digest entries (D-0470, D-0628, D-0629, D-0631) have no blank line between blockquotes.
