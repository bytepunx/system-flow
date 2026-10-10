---
id: TH-0390
title: "S-0337's branch ref is an empty file left by the host's crash at 05:04; please delete it so the worktree can be opened"
anchor:
  path: wip/kanban/stories/S-0337-flai-stats-reports-the-conversations-between-agents-the-conflicts-found-at-sync-and-at-acceptance-and-the-hold-time-shares-saved-per-week.md
  item: S-0337
status: open
participants: [agent-S-0337, orchestrator]
created: 2026-10-08T20:57:12Z
updated: 2026-10-08T20:57:26Z
---

# TH-0390 S-0337's branch ref is an empty file left by the host's crash at 05:04; please delete it so the worktree can be opened

On wip/kanban/stories/S-0337-flai-stats-reports-the-conversations-between-agents-the-conflicts-found-at-sync-and-at-acceptance-and-the-hold-time-shares-saved-per-week.md.

## Entries

### 2026-10-08T20:57:12Z agent-S-0337
S-0337 was started at 11:04:51Z (05:04 MDT), and the host went down moments later: it has been up only since about 14:32 MDT. Every file the start wrote was left at zero bytes: the narrative, all 4,147 files of the worktree, its .git/worktrees metadata, and the branch ref `.git/refs/heads/story/S-0337` with its reflog. No commit was ever made on the branch, so there is no work to lose. `git fsck --connectivity-only` finds nothing wrong but that ref, and no other file changed since 04:30 is empty.

I have removed the empty worktree folder, pruned its metadata, and opened the narrative again (it has content now). Git refuses to delete, replace, or recreate the ref while it is broken ("cannot lock ref ... reference broken"), and this session's permissions refuse any write under .git, so `flai stream open S-0337` stops at `git worktree add -b story/S-0337`.

Please run, in /home/alex/git/bytepunx/system-flow:

```bash
rm .git/refs/heads/story/S-0337 .git/logs/refs/heads/story/S-0337
```

and answer here. I will then run `flai stream open S-0337`, which creates story/S-0337 from main and its worktree, and work the story from the beginning: it has no tasks and no commits yet. Once the worktree exists I'll record the crash damage as an issue against S-0337 there.

### 2026-10-08T20:57:26Z orchestrator (recommendation)
Recommendation: run the agent's `rm`, but first complete S-0338's acceptance (TH-0389).

1. **First: `flai accept S-0338 --by alex`.** S-0338's acceptance merged, moved to done, and archived, then stopped before its commit on the empty `wip/agents/S-0337.md`. Its changes are still uncommitted in the main checkout. The narrative has content again, so the command can now finish: per `flai accept --help`, it commits an item that is done and archived but not committed. Do this before `flai stream open S-0337` writes to wip on main, so the acceptance is not swept into S-0337's commit.
2. **Then the agent's command.** I checked that both targets are 0 bytes: `.git/refs/heads/story/S-0337` and `.git/logs/refs/heads/story/S-0337`. The branch never had a commit, so deleting them loses nothing.

   ```bash
   rm .git/refs/heads/story/S-0337 .git/logs/refs/heads/story/S-0337
   ```

3. Reply here so agent-S-0337 runs `flai stream open S-0337`.

This needs you: deleting files under `.git` is a hard-to-reverse action that safety.md reserves for confirmation, and flai guard refuses it to agents.

This is a recommendation, not an answer: confirm it, or answer otherwise.

Source: design/conventions/safety.md § Rules
