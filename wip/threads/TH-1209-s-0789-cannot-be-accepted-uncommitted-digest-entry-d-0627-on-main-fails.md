---
id: TH-1209
title: "S-0789 cannot be accepted: uncommitted digest entry D-0627 on main fails markdown lint (MD028), as on TH-1206"
anchor:
  path: wip/kanban/stories/S-0789-publish-what-is-accepted-and-not-yet-released.md
  item: S-0789
status: open
participants: [orchestrator]
created: 2026-11-15T02:00:35Z
updated: 2026-11-15T02:00:35Z
---

# TH-1209 S-0789 cannot be accepted: uncommitted digest entry D-0627 on main fails markdown lint (MD028), as on TH-1206

On wip/kanban/stories/S-0789-publish-what-is-accepted-and-not-yet-released.md.

## Entries

### 2026-11-15T02:00:35Z orchestrator
**Recommended:** fix line 20 of `design/digest/D-0627-a-workitem-repo-whose-workflow-is-nil-runs-the-one-baseline-workflow-every-such.md` in the main checkout. Put `>` on the blank line between its two blockquotes, or have the curator rewrite it with `flai digest edit`. Then I accept S-0789, S-0738, S-0739, and S-0744 on my next run. The cause is the same as on TH-1206, TH-1207, and TH-1208.

**What holds it:**

- Apart from this, S-0789 is ready to accept. Its close-out passed at the branch head `721f679`, and the `accept` dry run listed no blocker.
- The formula ran cleanly. T-3366 (fetch) is done. T-3367 (evaluate) stopped the run with exit 2: the threshold was not met, with 0 of 1 stories and 0 of 1000 USD/week. T-3368 (publish) was skipped.
- The acceptance itself was refused:

```text
design/digest/D-0627-a-workitem-repo-whose-workflow-is-nil-runs-the-one-baseline-workflow-every-such.md:20 error MD028/no-blanks-blockquote Blank line inside blockquote
correct them in the main checkout and run flai accept S-0789 again
```

**Alternative:** fix the file, then accept it yourself with `flai accept S-0789 --by alex`.
