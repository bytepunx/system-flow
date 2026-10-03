---
id: TH-0073
title: "Conflicts in a publish are between accepted stories, which take no new tasks: where do the conflict tasks go?"
anchor:
  path: wip/kanban/stories/S-0242-a-publish-refused-because-the-remote-moved-drops-its-tags-takes-the-remote-by-rebase-or-merge-and-works-conflicts-through-tasks-and-threads.md
  item: S-0242
status: resolved
participants: [agent-S-0242, alex]
created: 2026-10-03T01:25:35Z
updated: 2026-10-03T01:51:01Z
---

# TH-0073 Conflicts in a publish are between accepted stories, which take no new tasks: where do the conflict tasks go?

On wip/kanban/stories/S-0242-a-publish-refused-because-the-remote-moved-drops-its-tags-takes-the-remote-by-rebase-or-merge-and-works-conflicts-through-tasks-and-threads.md.

## Entries

### 2026-10-03T01:25:35Z agent-S-0242
TH-0070 says merge conflicts in a publish are worked "through tasks on the stories in question, in threads with the operator". But the stories whose changes conflict when a clone takes its remote are accepted: done and archived, on both sides. `flai task new` refuses a task under a done or archived story (workitem.Create: "S-nnnn is done; reopen or pick another parent"), and nothing moves a story back out of done.

Options:
1. (Recommended) The agent publishing opens one remediation story for the publish's conflicts ("Resolve the conflicts taking <remote>/<branch> into the publish"), with one task per accepted story whose changes conflict. Each task names that story and the conflicting paths, and has a thread with the operator on it. Accepting that story is what publishes again. This needs no change to flai.
2. Let `flai task new` add a task to a done story when the task's purpose is a publish conflict. That reopens accepted work, and flai has no such state today. It's a larger change than S-0242's criteria.
3. Use threads only, on each conflicting story (threads on archived items work), and no tasks.

Until you answer, I'll write the convention with option 1 and change it if you choose otherwise.

### 2026-10-03T01:46:29Z agent-S-0242
An update while this waits. S-0242's work is committed, and the close-out passes the lint and every test tier. It stops only at `flai check --strict`, on `wip.overlap` warnings. S-0198, which is in progress, and its T-0731 also claim `docs/users/flai.md`, `docs/users/flai-reference.md`, and `design/system/flai-cli.md`, and S-0242 edits all three. `flai stream sync` says the two branches merge cleanly. Unless you say otherwise, I'll wait for S-0198 to be accepted, sync onto it, run the close-out again, and then move S-0242 to review. That way I also have your answer here before the convention's conflict sentence goes to review.

### 2026-10-03T01:50:36Z alex
I would like a different convention that 1 - I don't want new stories created for a story stuck in review since the overlaps would likely cause unwanted blocking behavior. Instead I prefer option 3 - use threads to resolve conflicts.

### 2026-10-03T01:51:01Z agent-S-0242
Resolved: Option 3. Publish conflicts are worked on a thread with the operator on each conflicting accepted story, with no new story or tasks. The work-management baseline and its template copy, pushing-from-the-board.md, and docs/users/flai.md say so (d8eaf47), and S-0242's notes record it.
