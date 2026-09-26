---
id: TH-0012
title: "origin/main used S-0117 to S-0119 for other stories: how do we renumber before merging?"
anchor:
  path: wip/kanban/stories/S-0122-the-dashboard-s-settings-page-serves-removes-and-shows-the-health-of-every-project-the-host-knows.md
  item: S-0122
status: resolved
participants: [agent-S-0119, alex, agent-S-0122]
created: 2026-09-26T05:40:31Z
updated: 2026-09-26T06:32:16Z
---

# TH-0012 origin/main used S-0117 to S-0119 for other stories: how do we renumber before merging?

On wip/kanban/stories/S-0119-the-dashboard-s-settings-page-serves-removes-and-shows-the-health-of-every-project-the-host-knows.md.

## Entries

### 2026-09-26T05:40:31Z agent-S-0119
I have not started S-0119. The S-0117 acceptance cannot be pushed, and merging would break the board.

**What happened.** This checkout and origin/main split at 697c097. Both sides then created items with the same IDs. origin/main has 15 commits from another machine (macOS, per I-0045 there). Local main has 9 unpushed commits (S-0117's work and its acceptance, c1d3bf2).

| ID | Here | On origin/main (archived, done) |
|----|------|------|
| S-0117 | import is served by the host flai (accepted) | every project's cache forgets what flai says changed |
| S-0118 | flai serve project add/remove/list (in progress) | Retry button on a failed agent |
| S-0119 | this story: settings page serves/removes projects | flaiover image for arm64 |
| T-0429 to T-0436 | tasks of S-0117 and S-0118 | tasks of those three stories |
| ADR 0044 | a repository under an import folder is served | a retry with the in-progress limit full is queued |
| I-0045 | import neither served nor offered | three flai tests fail on macOS |
| TH-0011 | S-0118 conflicts with S-0117 | should Retry queue |

`flai push --pending` refuses because the branches diverged. The convention's fix is to merge origin/main, but the merge would give every ID above two items and fail `flai check --strict`.

**Recommended:** renumber the local side, because it has not been pushed. In one commit on main, before the merge:

- S-0117 becomes S-0120, S-0118 becomes S-0121, S-0119 becomes S-0122.
- T-0429 to T-0436 become T-0437 to T-0444, ADR 0044 becomes 0045, I-0045 becomes I-0047, and TH-0011 becomes TH-0012.
- Every link, parent, and narrative that uses the old IDs is updated, and so is the `story/S-0118` branch.

Then merge origin/main (no rebase), run `flai check --strict`, and `flai push --pending`. The old commit subjects keep saying [S-0117]; the renumbering commit maps each old ID to its new one. After that I work this story as S-0122.

The S-0118 agent must hold while this happens, because it renames that agent's story. I will tell it on S-0118.

**Alternative:** you renumber by hand, or you pick which side keeps the IDs.

Root cause: flai numbers new items from the files in the local checkout, so two machines working the same repository collide. Once this is fixed I will record it as an issue. Recording it now would take I-0046, which origin already uses.

Shall I do the renumbering as recommended?

### 2026-09-26T05:42:41Z alex
yes, renumber on this machine

### 2026-09-26T05:44:46Z agent-S-0119
Recorded. There is one more choice, and one request.

**Please do not accept S-0118 until it is renumbered.** Accepting it now adds another `chore: [S-0118] accept and archive` commit. After the merge that commit would name origin's Retry story.

**Order.** The agent working S-0118 is rebasing that story's branch now. Renaming its story, branch, and worktree while it works would break its session. So I renumber after S-0118 reaches review, all in one pass: S-0117 to S-0120, S-0118 to S-0121, S-0119 to S-0122, and the rest as listed. Until then I read code for this story and change nothing.

**The choice: commit subjects.** The release tooling works out what to publish from commit messages. It reads the item ID in each `chore: [ID] accept and archive` subject and gathers the commits whose message contains `[ID]`. After the merge, `[S-0117]` names two stories. Then `flai push --pending` would:

- release the import work under origin's S-0117 title. The flai bump is still a patch, because both are remediations.
- publish flaiover 0.27.9 again for origin's S-0117 files, which 0.27.8 already shipped.

A. **Recommended.** Before the merge, reword the 9 unpushed commits on main from `[S-0117]` to `[S-0120]`. Only the messages change; every file stays as it is. Then move `story/S-0121` onto the reworded main with `git rebase --onto`, which cannot conflict because the files are identical. Nothing that was pushed is rewritten. The convention forbids amending commits that are not mine from this session, so this needs your yes.

B. Keep the history. After the merge I run `flai release --pending --dry-run` and show you what it would publish before anything is pushed. The changelog will name the wrong story, and flaiover probably gets an extra patch release.

A or B?

### 2026-09-26T05:56:37Z alex
A

### 2026-09-26T06:04:19Z agent-S-0119
Renumbering is done, and main matches origin apart from one commit.

**Done**
- Reworded the 17 unpushed commits and renamed the items.
- Merged origin/main and pushed S-0120 as flai 1.18.1.
- You accepted S-0121 after I had rebased it onto the merge, and I pushed it as flai 1.18.2 and flaiover 0.27.9.
- a6e1ac0 renumbers the IDs on the lines S-0121 added. Those lines were written as S-0117 and S-0118, which on main now name origin's stories.

**Broken, my mistake: no flai binaries for 1.18.1 or 1.18.2.** I moved this story to in-progress before the renumbering commit and opened no narrative. So the tagged commits fail `flai check`, and goreleaser's `TestMonorepoIsClean` fails with them. The release runs for flai/v1.18.1 and flai/v1.18.2 fail on that; the system-flow check on main fails the same way. Re-running cannot help, because the tagged files are fixed. 86a1a8b on local main adds the narrative, and check is clean there. flaiover 0.27.9 is still building.

A. **Recommended.** Push 86a1a8b and tag flai/v1.18.3 on it, with a changelog line saying 1.18.1 and 1.18.2 were never built. No force push; the two dead tags stay.

B. Push 86a1a8b, move flai/v1.18.2 onto it with a forced tag push, and delete flai/v1.18.1. The history is tidier, but it rewrites pushed tags.

A or B?

### 2026-09-26T06:15:37Z alex
A

### 2026-09-26T06:17:01Z alex
A

### 2026-09-26T06:32:16Z agent-S-0122
Resolved: Done as A. 86a1a8b is pushed to main (flai check --strict clean on it), and flai/v1.18.3 is tagged on it and published with its binaries. The first release run failed on a flaky internal/serve test, not on the tagged files; a rerun of the job passed. T-0450 on S-0122 fixes the flake. flai/v1.18.1 and flai/v1.18.2 stay as tags without builds, and the v1.18.3 tag message says so.
