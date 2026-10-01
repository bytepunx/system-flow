---
id: TH-0047
title: S-0184 and S-0185 conflict when merged
anchor:
  path: wip/kanban/stories/S-0184-the-dashboard-recovers-from-a-stray-stop-and-is-not-taken-over-by-a-second-flai-serve.md
  item: S-0184
status: resolved
participants: [flai, agent-S-0184]
created: 2026-10-01T09:02:51Z
updated: 2026-10-01T09:04:47Z
---

# TH-0047 S-0184 and S-0185 conflict when merged

On wip/kanban/stories/S-0184-the-dashboard-recovers-from-a-stray-stop-and-is-not-taken-over-by-a-second-flai-serve.md.

## Entries

### 2026-10-01T09:02:51Z flai
A trial merge of story/S-0184 with story/S-0185 at flai stream sync conflicts in:

- `design/issues/summary.md`

Whichever of S-0184 and S-0185 is accepted second will stop on these paths when it rebases. Settle between the two stories who changes what: one narrows its change, or names the other in `after:` and waits for it. Ask the designer when it is not clear. The next sync that finds the two merging cleanly resolves this thread.

### 2026-10-01T09:02:58Z agent-S-0184
S-0184 has not changed `design/issues/summary.md` yet. This sync rebased story/S-0184 onto main, which now carries S-0179's change to that file, so the conflict is between main and S-0185's branch. It goes away when S-0185 syncs onto main. In its last task S-0184 will close I-0025 and I-0029 with `flai issue close`, which rewrites summary.md. If S-0185 is accepted first, I will sync and let `flai issue` regenerate the file rather than merge it by hand. Neither story needs to name the other in `after:`.

### 2026-10-01T09:04:47Z flai
Resolved: story/S-0184 and story/S-0185 merge cleanly at the sync of S-0185
