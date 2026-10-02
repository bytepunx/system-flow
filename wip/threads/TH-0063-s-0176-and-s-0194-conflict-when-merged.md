---
id: TH-0063
title: S-0176 and S-0194 conflict when merged
anchor:
  path: wip/kanban/stories/S-0194-an-experiment-story-is-accepted-with-its-results-document-under-design-experiments.md
  item: S-0194
status: answered
participants: [flai, agent-S-0194]
created: 2026-10-02T10:28:49Z
updated: 2026-10-02T10:32:10Z
---

# TH-0063 S-0176 and S-0194 conflict when merged

On wip/kanban/stories/S-0194-an-experiment-story-is-accepted-with-its-results-document-under-design-experiments.md.

## Entries

### 2026-10-02T10:28:49Z flai
A trial merge of story/S-0176 with story/S-0194 at flai stream sync conflicts in:

- `design/conventions/delegation.md`
- `template/root/design/conventions/delegation.md`

Whichever of S-0176 and S-0194 is accepted second will stop on these paths when it rebases. Settle between the two stories who changes what: one narrows its change, or names the other in `after:` and waits for it. Ask the designer when it is not clear. The next sync that finds the two merging cleanly resolves this thread.

### 2026-10-02T10:32:10Z agent-S-0194
This is the same conflict as TH-0061 and TH-0062. S-0194 changes neither `delegation.md`. The designer's commits on `main` change both, and S-0176's branch has not been synced onto them since. S-0176 resolves it on its next `flai stream sync`, or at its acceptance.
