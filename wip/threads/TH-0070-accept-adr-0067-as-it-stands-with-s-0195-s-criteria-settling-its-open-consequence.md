---
id: TH-0070
title: Accept ADR-0067 as it stands, with S-0195's criteria settling its open consequence?
anchor:
  path: wip/kanban/stories/S-0195-publishing-is-the-one-way-accepted-work-reaches-the-remote-and-agents-no-longer-push-unpushed-acceptances.md
  item: S-0195
status: open
participants: [agent-S-0195, alex]
created: 2026-10-02T23:25:25Z
updated: 2026-10-02T23:48:37Z
---

# TH-0070 Accept ADR-0067 as it stands, with S-0195's criteria settling its open consequence?

On wip/kanban/stories/S-0195-publishing-is-the-one-way-accepted-work-reaches-the-remote-and-agents-no-longer-push-unpushed-acceptances.md.

## Entries

### 2026-10-02T23:25:25Z agent-S-0195
S-0195's first criterion needs you to accept ADR-0067 (design/adrs/0067-accepted-work-reaches-the-remote-only-when-it-is-published-and-agents-publish.md), or ask for changes, before anything is built.

Recommendation: accept it as written. Its one open point ("whether `flai push --pending` and `auto-publish` are removed or kept is decided by the story") is answered by your criteria: both stay as commands for the operator's own use, the dashboard drops them (no Push action, no push banner over the board), and the board and dashboard still show what is accepted and not yet published. I will record that in S-0195's narrative and in the living design (`pushing-from-the-board.md`), not in the ADR.

Reply "accept" and I will run `flai adr` to set it accepted, or name the changes and I will write them as a new ADR refining 0067.

### 2026-10-02T23:39:34Z alex
I have edited and accepted the ADR. Please proceed with the story's implementation.

### 2026-10-02T23:40:04Z agent-S-0195
Recorded: ADR-0067 accepted with your edits. I am building it as edited: `inbox` and `flai board` keep saying what is accepted and not yet published, with no instruction to push it; a publish from a clone behind its remote refuses and names the fetch, then merge or rebase.

One follow-up, which does not block me: the work-management baseline (`template/root/design/conventions/work-management.md` and its copy) still says "merge the remote branch (never rebase what was accepted)", and ADR-0067 now allows rebase or merge, with conflicts coordinated on a thread. Recommendation: change that line in a follow-up story after S-0196 is accepted, because S-0196 is editing the same template conventions now. Say if you want it in S-0195 instead.

### 2026-10-02T23:43:28Z alex
why wouldn't we want to rebase an accepted change set for a story onto the remote's changes?

### 2026-10-02T23:43:42Z agent-S-0195
Rebasing is safe in itself: accepted commits that never reached the remote are not public history. The rule guarded one hazard: `flai release --pending` tags and commits the version bump before it pushes, so a push refused because the remote moved leaves release tags on the local commits. A rebase then moves the commits and leaves the tags on the old ones; the rerun does not retag (tagging is idempotent) and pushes tags pointing outside the branch. A plain `git rebase` also flattens the acceptance merge commits.

T-0719 removes most of that: the publish will refuse before tagging when the remote branch has commits this clone lacks, so in the normal case nothing is tagged yet and a rebase is clean. Only a race (the remote moves between the check and the push) still leaves tags behind.

Recommendation: the convention allows either, as ADR-0067 says, adding "if the refused publish already tagged, delete those local tags before rebasing, or merge instead". I would still make that edit in a follow-up story after S-0196, unless you want it here.

### 2026-10-02T23:45:53Z alex
Lets go with the recommendation with refinements:

### 2026-10-02T23:46:41Z alex
1 - if the remote has changed, our tags are invalid so we need to delete them and either rebase on the remote OR merge and then we can re-tag after the verification completes.

### 2026-10-02T23:48:15Z alex
2 - make use of tasks attached to the stor(y|ies) in question to work through merge conflicts in threads with the operator

### 2026-10-02T23:48:37Z agent-S-0195
Recorded both refinements in S-0242 (backlog): when a publish is refused because the remote moved, the tags it made are deleted, the clone rebases onto the remote or merges it, verifies, and publishes again, which re-tags; merge conflicts are worked through tasks on the stories in question, in threads with you. It changes the work-management baseline and the release recovery after S-0196 is accepted. S-0195 goes on as planned: its publish now refuses before tagging when the remote branch has commits this clone lacks, so the tag cleanup is needed only when the remote moves mid-publish. Move S-0242 to ready when you want it worked.
