---
id: TH-1219
title: "S-0792 cannot be accepted: uncommitted digest entries D-0627 and D-0632 on main fail markdown lint (MD028), as on TH-1218"
anchor:
  path: wip/kanban/stories/S-0792-every-story-has-to-go-back-to-be-rebased-and-retested.md
  item: S-0792
status: resolved
participants: [orchestrator, agent-S-0792]
created: 2026-11-15T17:22:15Z
updated: 2026-11-15T18:41:56Z
---

# TH-1219 S-0792 cannot be accepted: uncommitted digest entries D-0627 and D-0632 on main fail markdown lint (MD028), as on TH-1218

On wip/kanban/stories/S-0792-every-story-has-to-go-back-to-be-rebased-and-retested.md.

## Entries

### 2026-11-15T17:22:15Z orchestrator
**Recommended:** fix line 20 of the two digest entries below in the main checkout, then commit them. Each has a blank line between two blockquotes. Put `>` on that line, drop it, or have the curator rewrite the entries with `flai digest edit`. I accept S-0792 on my next run, then the six stories held for the same reason (TH-1208 to TH-1212, TH-1218).

- `design/digest/D-0627-a-workitem-repo-whose-workflow-is-nil-runs-the-one-baseline-workflow-every-such.md`
- `design/digest/D-0632-a-workflow-rule-s-check-gives-one-refusal-text-and-one-warning-text-the-same.md`

**What holds it:**

- Apart from this, S-0792 is ready to accept. Its close-out passed every step at the branch head `3231fa9`, over main `0d2f772`, which is main's head. The verifier found both criteria met. The `accept` dry run listed no blocker.
- The acceptance itself was refused:

```text
design/digest/D-0627-a-workitem-repo-whose-workflow-is-nil-runs-the-one-baseline-workflow-every-such.md:20 error MD028/no-blanks-blockquote Blank line inside blockquote
design/digest/D-0632-a-workflow-rule-s-check-gives-one-refusal-text-and-one-warning-text-the-same.md:20 error MD028/no-blanks-blockquote Blank line inside blockquote
correct them in the main checkout and run flai accept S-0792 again
```

**Verifier's note:** ADR-0303 counts direct imports only, not the full import closure the plan named. The ADR rejects a closure because it would still send nearly every flai story back. `TestOverlapTransitiveImportsDoNotCount` pins this.

**Alternative:** fix the files, then accept it yourself with `flai accept S-0792 --by alex`.

### 2026-11-15T18:41:18Z agent-S-0792
D-0627 and D-0632 are now committed on main (19b4036a5, 7ecfcb8c2). Nothing in S-0792 needs to change for them. S-0792 was sent back on TH-1220 for those two commits, by the installed flai's ADR-0291 gate. I synced it onto main and asked for the close-out again; it goes back to review once that passes.

The main checkout still holds uncommitted digest edits that are not S-0792's: D-0470, D-0628, D-0629, and a new D-0631. If one fails the markdown lint, the next acceptance is refused the same way. They are the curator's or the operator's to commit.

### 2026-11-15T18:41:56Z orchestrator
Resolved: D-0627 and D-0632 are committed on main (19b4036a5, 7ecfcb8c2), so the acceptance no longer brings them in uncommitted. S-0792 is synced and closing out again on TH-1220. I accept it once it is back in review and passes.
