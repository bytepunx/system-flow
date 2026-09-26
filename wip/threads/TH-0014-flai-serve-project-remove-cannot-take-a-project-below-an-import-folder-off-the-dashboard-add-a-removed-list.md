---
id: TH-0014
title: "flai serve project remove cannot take a project below an import folder off the dashboard: add a removed list?"
anchor:
  path: wip/kanban/stories/S-0121-flai-serve-project-add-remove-and-list-manage-the-projects-the-host-serves-and-say-why-one-is-not-showing.md
  item: S-0121
status: resolved
participants: [agent-S-0118, alex]
created: 2026-09-26T05:46:05Z
updated: 2026-09-26T05:58:33Z
---

# TH-0014 flai serve project remove cannot take a project below an import folder off the dashboard: add a removed list?

On wip/kanban/stories/S-0121-flai-serve-project-add-remove-and-list-manage-the-projects-the-host-serves-and-say-why-one-is-not-showing.md.

## Entries

### 2026-09-26T05:46:05Z agent-S-0118
Since S-0120, `flai serve` serves every project below a folder named for import, registered or not. So `flai serve project remove blog` for a project under `~/git` (an import folder) unregisters it, but `flai serve` goes on serving it from the folder, and the switcher keeps it. Today `remove` says exactly that and names `flai serve import remove ~/git`, which also stops offering every other repository there. This question blocks moving S-0121 to review.

A. **Recommended.** Accept S-0121 with that behaviour. Add a story, or fold it into S-0122, whose settings page also removes projects: a list of roots the operator removed, kept beside the registry. Folder and import-folder serving skip those roots until `flai serve project add` there again.

B. I build the removed list in S-0121 now: `remove` records the root, `serve.Place` skips it, and `add` clears it. That is about one more task, with tests and docs.

Which one?

### 2026-09-26T05:58:17Z alex
A

### 2026-09-26T05:58:33Z agent-S-0118
Resolved: A: S-0121 goes to review with remove saying a project below an import folder stays served. The removed-roots list is later work. Its story waits until the renumbering in TH-0012 is done, so it does not take S-0120.
