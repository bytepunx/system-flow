---
id: TH-1212
title: "S-0747 cannot be accepted: uncommitted digest entry D-0627 on main fails markdown lint (MD028), as on TH-1206"
anchor:
  path: wip/kanban/stories/S-0747-the-host-api-s-writes-kluster-components-and-people-methods-act-on-a-state-s.md
  item: S-0747
status: resolved
participants: [orchestrator]
created: 2026-11-15T02:44:56Z
updated: 2026-11-15T18:43:20Z
---

# TH-1212 S-0747 cannot be accepted: uncommitted digest entry D-0627 on main fails markdown lint (MD028), as on TH-1206

On wip/kanban/stories/S-0747-the-host-api-s-writes-kluster-components-and-people-methods-act-on-a-state-s.md.

## Entries

### 2026-11-15T02:44:56Z orchestrator
**Recommended:** fix line 20 of `design/digest/D-0627-a-workitem-repo-whose-workflow-is-nil-runs-the-one-baseline-workflow-every-such.md` in the main checkout. The line holds only two spaces between two blockquotes. Put `>` on it, or drop it, or have the curator rewrite the entry with `flai digest edit`. Then I accept S-0747 and the six stories held for the same reason (TH-1206 to TH-1211) on my next run.

**What holds it:**

- Apart from this, S-0747 is ready to accept. Its close-out passed at the branch head `be8fafb`. The verifier found all five criteria met. The `accept` dry run listed no blocker.
- The acceptance itself was refused:

```text
design/digest/D-0627-a-workitem-repo-whose-workflow-is-nil-runs-the-one-baseline-workflow-every-such.md:20 error MD028/no-blanks-blockquote Blank line inside blockquote
correct them in the main checkout and run flai accept S-0747 again
```

- Review is over its limit (7 of 5), so no ready story can be pulled until stories are accepted.

**Alternative:** fix the file, then accept it yourself with `flai accept S-0747 --by alex`.

### 2026-11-15T18:43:20Z orchestrator
Resolved: D-0627 is committed on main (19b4036a5), so the acceptance no longer brings it in uncommitted. The other uncommitted digest entries (D-0470, D-0628, D-0629, D-0631) have no blank line between blockquotes.
