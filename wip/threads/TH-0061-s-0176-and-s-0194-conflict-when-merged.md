---
id: TH-0061
title: S-0176 and S-0194 conflict when merged
anchor:
  path: wip/kanban/stories/S-0194-an-experiment-story-is-accepted-with-its-results-document-under-design-experiments.md
  item: S-0194
status: resolved
participants: [flai, agent-S-0194, alex]
created: 2026-10-02T10:13:37Z
updated: 2026-10-02T10:19:03Z
---

# TH-0061 S-0176 and S-0194 conflict when merged

On wip/kanban/stories/S-0194-an-experiment-story-is-accepted-with-its-results-document-under-design-experiments.md.

## Entries

### 2026-10-02T10:13:37Z flai
A trial merge of story/S-0176 with story/S-0194 at flai stream sync conflicts in:

- `template/root/design/conventions/delegation.md`

Whichever of S-0176 and S-0194 is accepted second will stop on these paths when it rebases. Settle between the two stories who changes what: one narrows its change, or names the other in `after:` and waits for it. Ask the designer when it is not clear. The next sync that finds the two merging cleanly resolves this thread.

### 2026-10-02T10:13:57Z agent-S-0194
S-0194 does not change `template/root/design/conventions/delegation.md`: `git diff main...story/S-0194` has no path under `template/root/design/conventions/`. The conflict is between the designer's commit 2d4e090 on `main` (convention edits) and S-0176's branch, which has not been synced onto `main` since. My branch has it only because it was rebased onto `main`. It is S-0176's to settle on its next `flai stream sync`, or the designer's at its acceptance. S-0194 will not narrow anything for it.

### 2026-10-02T10:19:03Z alex
Resolved.
