---
id: TH-1208
title: "S-0744 cannot be accepted: uncommitted digest entry D-0627 on main fails markdown lint (MD028), as on TH-1206"
anchor:
  path: wip/kanban/stories/S-0744-flai-move-flai-accept-the-formula-commands-and-the-acceptance-preview-act-on-a.md
  item: S-0744
status: resolved
participants: [orchestrator]
created: 2026-11-15T01:48:46Z
updated: 2026-11-15T18:43:16Z
---

# TH-1208 S-0744 cannot be accepted: uncommitted digest entry D-0627 on main fails markdown lint (MD028), as on TH-1206

On wip/kanban/stories/S-0744-flai-move-flai-accept-the-formula-commands-and-the-acceptance-preview-act-on-a.md.

## Entries

### 2026-11-15T01:48:46Z orchestrator
**Recommended:** fix line 20 of `design/digest/D-0627-a-workitem-repo-whose-workflow-is-nil-runs-the-one-baseline-workflow-every-such.md` in the main checkout. Put `>` on the blank line between its two blockquotes, or have the curator rewrite it with `flai digest edit`. I accept S-0744, S-0738, and S-0739 on my next run. The cause is the same as on TH-1206 and TH-1207.

**What holds it:**

- Apart from this, S-0744 is ready to accept. The close-out passed at the branch head `18bbb55`. The verifier found all seven criteria met. The `accept` dry run listed no blocker.
- The acceptance itself refused:

```text
design/digest/D-0627-a-workitem-repo-whose-workflow-is-nil-runs-the-one-baseline-workflow-every-such.md:20 error MD028/no-blanks-blockquote Blank line inside blockquote
correct them in the main checkout and run flai accept S-0744 again
```

**Alternative:** fix the file, then accept it yourself with `flai accept S-0744 --by alex`.

### 2026-11-15T18:43:16Z orchestrator
Resolved: D-0627 is committed on main (19b4036a5), so the acceptance no longer brings it in uncommitted. The other uncommitted digest entries (D-0470, D-0628, D-0629, D-0631) have no blank line between blockquotes.
