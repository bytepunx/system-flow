---
id: D-0274
title: "S-0337's branch ref is an empty file left by the host's crash at 05:04; please delete it so the worktree can be opened: Recommendation: run the agent's `rm`, but first complete S-0338's acceptance…"
kind: answer
sources: [TH-0390]
topics: [cli, dashboard, git, flai, go, flaiover, sveltekit, code, all]
paths: [design/adrs, design/system/metrics.md, flai/internal/storygit/conflicts.go, flai/internal/storygit/conflicts_test.go, flai/internal/storygit/sync.go, flai/internal/storygit/sync_test.go, flai/cmd/accept.go, flai/cmd/accept_conflict_test.go, flai/internal/metrics/coordination.go, flai/internal/metrics/coordination_test.go, flai/internal/metrics/metrics.go, flai/internal/statsread/statsread.go, flai/internal/statsread/statsread_test.go, flai/cmd/stats.go, flai/cmd/check_stats_test.go, flaiover/src/lib/viz/charts.ts, flaiover/src/lib/viz/charts.test.ts, "flaiover/src/routes/charts/[kind]/+page.svelte", "flaiover/src/routes/charts/[kind]/charts.svelte.test.ts", design/system/flaiover-dashboard.md, docs/users/flai.md, docs/users/flaiover.md, docs/users/flai-reference.md, design/system/flai-cli.md, design/issues/I-0125-a-settled-conversation-between-two-stories-agents-keeps-waking-each-in-turn-because-every-reply-awaits-the-other.md, design/issues/summary.md, design/issues/I-0137-a-host-crash-during-flai-story-start-leaves-the-story-s-branch-ref-worktree-and-narrative-as-empty-files-and-no-agent-can-repair-the-ref.md, design/issues/I-0140-testactivitywriterstaketurns-fails-about-one-run-in-three-two-writers-that-both-find-the-activity-lock-stale-can-both-take-it-and-an-entry-is-lost.md, design/issues/I-0139-flai-check-finds-markdown-md012-outside-the-story-at-close-out.md, design/issues/I-0136-an-acceptance-on-main-during-a-close-out-s-tiers-fails-its-final-sync-check-so-the-whole-close-out-runs-again.md]
date: 2026-10-12
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0274 S-0337's branch ref is an empty file left by the host's crash at 05:04; please delete it so the worktree can be opened: Recommendation: run the agent's `rm`, but first complete S-0338's acceptance…

The operator, alex, answered TH-0390 at 2026-10-12T15:29:06Z.

What was asked:

```text
S-0337 was started at 11:04:51Z (05:04 MDT), and the host went down moments later: it has been up only since about 14:32 MDT. Every file the start wrote was left at zero bytes: the narrative, all 4,147 files of the worktree, its .git/worktrees metadata, and the branch ref `.git/refs/heads/story/S-0337` with its reflog. No commit was ever made on the branch, so there is no work to lose. `git fsck --connectivity-only` finds nothing wrong but that ref, and no other file changed since 04:30 is empty.
```

The answer:

````text
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
````

Source: TH-0390.
