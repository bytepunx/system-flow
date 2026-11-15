---
id: TH-1206
title: "S-0738 cannot be accepted: uncommitted digest entry D-0627 on main fails markdown lint (MD028)"
anchor:
  path: wip/kanban/stories/S-0738-the-rules-repo-move-enforces-are-named-rules-from-a-vocabulary-flai-implements.md
  item: S-0738
status: resolved
participants: [orchestrator]
created: 2026-11-15T00:52:33Z
updated: 2026-11-15T06:52:20Z
---

# TH-1206 S-0738 cannot be accepted: uncommitted digest entry D-0627 on main fails markdown lint (MD028)

On wip/kanban/stories/S-0738-the-rules-repo-move-enforces-are-named-rules-from-a-vocabulary-flai-implements.md.

## Entries

### 2026-11-15T00:52:33Z orchestrator
**Recommended:** in the main checkout, join the two blockquotes at line 19–21 of `design/digest/D-0627-a-workitem-repo-whose-workflow-is-nil-runs-the-one-baseline-workflow-every-such.md` (put `>` on the blank line between them, or drop the blank line), or have the curator rewrite the entry with `flai digest edit`. Then I accept S-0738 on my next run.

**What holds it:**

- S-0738 is ready to accept. The close-out passed at the branch head `66323bb`. The verifier found all six criteria met. The `accept` dry run listed no blocker.
- The acceptance itself refused, because the main checkout's uncommitted files hold lint the markdown fix could not fix:

```text
design/digest/D-0627-a-workitem-repo-whose-workflow-is-nil-runs-the-one-baseline-workflow-every-such.md:20 error MD028/no-blanks-blockquote Blank line inside blockquote
correct them in the main checkout and run flai accept S-0738 again
```

- D-0627, D-0628, and D-0629 are untracked on main, written by the curator. The file is not S-0738's, and I edit no files.

**Alternative:** accept it yourself with `flai accept S-0738 --by alex` once the file is fixed.

### 2026-11-15T06:52:20Z orchestrator
Resolved: Its concern is gone: the operator committed the MD028 fix to D-0627 on main (f01a053e7), and S-0738 is back in progress on TH-1213, whose dry run lists only the verify-coverage blocker. A fresh close-out decides S-0738's next acceptance.
