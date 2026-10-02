---
id: TH-0062
title: S-0176 and S-0194 conflict when merged
anchor:
  path: wip/kanban/stories/S-0194-an-experiment-story-is-accepted-with-its-results-document-under-design-experiments.md
  item: S-0194
status: resolved
participants: [flai, agent-S-0194, alex]
created: 2026-10-02T10:21:22Z
updated: 2026-10-02T10:23:41Z
---

# TH-0062 S-0176 and S-0194 conflict when merged

On wip/kanban/stories/S-0194-an-experiment-story-is-accepted-with-its-results-document-under-design-experiments.md.

## Entries

### 2026-10-02T10:21:22Z flai
A trial merge of story/S-0176 with story/S-0194 at flai stream sync conflicts in:

- `design/conventions/delegation.md`
- `template/root/design/conventions/delegation.md`

Whichever of S-0176 and S-0194 is accepted second will stop on these paths when it rebases. Settle between the two stories who changes what: one narrows its change, or names the other in `after:` and waits for it. Ask the designer when it is not clear. The next sync that finds the two merging cleanly resolves this thread.

### 2026-10-02T10:21:29Z agent-S-0194
Same as TH-0061: S-0194 changes neither `design/conventions/delegation.md` nor `template/root/design/conventions/delegation.md`. The designer's commits on `main` (2d4e090, f45f959) change both, and S-0176's branch has not been synced onto them. My branch has those changes only through its rebase onto `main`, so this is S-0176's to resolve on its next `flai stream sync`.

### 2026-10-02T10:23:41Z alex
Resolved.
