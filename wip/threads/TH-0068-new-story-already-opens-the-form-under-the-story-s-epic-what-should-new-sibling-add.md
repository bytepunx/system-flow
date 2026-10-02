---
id: TH-0068
title: "New story already opens the form under the story's epic: what should New sibling add?"
anchor:
  path: wip/kanban/stories/S-0192-add-a-way-to-create-a-sibling-story.md
  item: S-0192
status: resolved
participants: [agent-S-0192, alex]
created: 2026-10-02T16:27:59Z
updated: 2026-10-02T16:33:57Z
---

# TH-0068 New story already opens the form under the story's epic: what should New sibling add?

On wip/kanban/stories/S-0192-add-a-way-to-create-a-sibling-story.md.

## Entries

### 2026-10-02T16:27:59Z agent-S-0192
On a story's page, **New story** already opens `/new?type=story&parent=<the story's epic>` (S-0171), so the form starts under the same epic. A **New sibling** that does the same would be two buttons with one meaning.

My recommendation: make them differ.

- **New story** opens `/new?type=story` with no epic, a story on its own.
- **New sibling** comes after it and opens `/new?type=story&parent=<epic>`, a story under the same epic. It shows only when the story has an epic that is open (not done, cancelled, or archived), since the form drops a parent that is not open, as Create story does on an epic (S-0191).

Alternatively: keep New story as it is and add New sibling with the same link, as the criterion reads literally.

Which do you want? I'll update the design and user guide to match.

### 2026-10-02T16:30:32Z alex
go with recommendation

### 2026-10-02T16:33:57Z agent-S-0192
Resolved: Done as recommended in T-0703 (67f6a6b): New story opens the form for a story with no epic, and New sibling opens it under the story's epic while that epic is open.
