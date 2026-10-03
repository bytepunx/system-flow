---
id: TH-0082
title: "An epic accepted with its last story makes the next publish a major release: keep that?"
anchor:
  path: wip/kanban/stories/S-0200-an-epic-follows-its-stories-to-ready-and-in-progress-with-the-first-to-review-and-done-with-the-last.md
  item: S-0200
status: open
participants: [agent-S-0200]
created: 2026-10-03T07:17:38Z
updated: 2026-10-03T07:17:38Z
---

# TH-0082 An epic accepted with its last story makes the next publish a major release: keep that?

On wip/kanban/stories/S-0200-an-epic-follows-its-stories-to-ready-and-in-progress-with-the-first-to-review-and-done-with-the-last.md.

## Entries

### 2026-10-03T07:17:38Z agent-S-0200
S-0200 makes an epic go to done, and be archived, in the acceptance of its last open story, as its criteria say. An accepted epic contributes a major bump to the next publish (git.md; release.LevelFor returns Major for any epic). No epic has been accepted in this repository yet (wip/archive/epics is empty, flai is at 1.29.1), so once this ships, accepting the last story of E-0016, or of any epic, makes the next `flai release --pending` take flai to 2.0.0, and the dashboard and template too where the epic's stories touched them.

Options:
1. Keep it: the epic contributes a major bump, as git.md says today. Accepting an epic's last story is then a deliberate act the review page should name ("this also accepts E-nnnn: the next publish is a major release"). I will make the dry run and the acceptance output say so.
2. An epic that follows its stories contributes no bump of its own: its stories' bumps already cover what was delivered, and a major release stays a deliberate story (git.md already says a breaking change needs its own story). That changes git.md's baseline rule and release.LevelFor, which I would do in this story with the ADR.

My recommendation is 2: a major version should mean a breaking change, and with this story an epic's end is no longer a deliberate act, so it would cut a major release by accident. Until you answer I am building the rest of the story, which does not depend on it, and keep option 1's behaviour.
