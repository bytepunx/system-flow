---
id: TH-1211
title: "S-0742 cannot be accepted: uncommitted digest entry D-0627 on main fails markdown lint (MD028), as on TH-1206"
anchor:
  path: wip/kanban/stories/S-0742-flai-serve-starts-restarts-and-stops-agents-plans-again-and-finds-its-reasons.md
  item: S-0742
status: open
participants: [orchestrator]
created: 2026-11-15T02:23:45Z
updated: 2026-11-15T02:23:45Z
---

# TH-1211 S-0742 cannot be accepted: uncommitted digest entry D-0627 on main fails markdown lint (MD028), as on TH-1206

On wip/kanban/stories/S-0742-flai-serve-starts-restarts-and-stops-agents-plans-again-and-finds-its-reasons.md.

## Entries

### 2026-11-15T02:23:45Z orchestrator
**Recommended:** fix line 20 of `design/digest/D-0627-a-workitem-repo-whose-workflow-is-nil-runs-the-one-baseline-workflow-every-such.md` in the main checkout. The line holds only two spaces between two blockquotes. Put `>` on it, drop it, or have the curator rewrite the entry with `flai digest edit`. Then I accept S-0742, S-0738, S-0739, S-0743, S-0744, and S-0789 on my next run. The cause is the same as on TH-1206 to TH-1210.

**What holds it:**

- Apart from this, S-0742 is ready to accept. Its close-out passed at the branch head `7a2dfc7`. The verifier found all five criteria met. The `accept` dry run listed no blocker.
- The acceptance itself was refused:

```text
design/digest/D-0627-a-workitem-repo-whose-workflow-is-nil-runs-the-one-baseline-workflow-every-such.md:20 error MD028/no-blanks-blockquote Blank line inside blockquote
correct them in the main checkout and run flai accept S-0742 again
```

- Review is over its limit (6 of 5), so no ready story can be pulled until stories are accepted.

**Alternative:** fix the file, then accept it yourself with `flai accept S-0742 --by alex`.
