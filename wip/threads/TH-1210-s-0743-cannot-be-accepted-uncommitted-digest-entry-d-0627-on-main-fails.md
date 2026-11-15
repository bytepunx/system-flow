---
id: TH-1210
title: "S-0743 cannot be accepted: uncommitted digest entry D-0627 on main fails markdown lint (MD028), as on TH-1206"
anchor:
  path: wip/kanban/stories/S-0743-flai-serve-s-kluster-deploys-forge-pass-and-merge-queue-act-on-a-state-s-kind.md
  item: S-0743
status: resolved
participants: [orchestrator]
created: 2026-11-15T02:22:04Z
updated: 2026-11-15T18:43:18Z
---

# TH-1210 S-0743 cannot be accepted: uncommitted digest entry D-0627 on main fails markdown lint (MD028), as on TH-1206

On wip/kanban/stories/S-0743-flai-serve-s-kluster-deploys-forge-pass-and-merge-queue-act-on-a-state-s-kind.md.

## Entries

### 2026-11-15T02:22:04Z orchestrator
**Recommended:** fix line 20 of `design/digest/D-0627-a-workitem-repo-whose-workflow-is-nil-runs-the-one-baseline-workflow-every-such.md` in the main checkout. The line holds only two spaces between two blockquotes. Put `>` on it, or drop it, or have the curator rewrite the entry with `flai digest edit`. Then I accept S-0743, S-0738, S-0739, S-0744, and S-0789 on my next run. The cause is the same as on TH-1206 to TH-1209.

**What holds it:**

- Apart from this, S-0743 is ready to accept. Its close-out passed at the branch head `8dd5578`. The verifier found all five criteria met. The `accept` dry run listed no blocker.
- The acceptance itself was refused:

```text
design/digest/D-0627-a-workitem-repo-whose-workflow-is-nil-runs-the-one-baseline-workflow-every-such.md:20 error MD028/no-blanks-blockquote Blank line inside blockquote
correct them in the main checkout and run flai accept S-0743 again
```

**Alternative:** fix the file, then accept it yourself with `flai accept S-0743 --by alex`.

### 2026-11-15T18:43:18Z orchestrator
Resolved: D-0627 is committed on main (19b4036a5), so the acceptance no longer brings it in uncommitted. The other uncommitted digest entries (D-0470, D-0628, D-0629, D-0631) have no blank line between blockquotes.
