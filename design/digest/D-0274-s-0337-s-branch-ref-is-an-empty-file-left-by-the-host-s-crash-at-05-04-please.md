---
id: D-0274
title: "S-0337's branch ref is an empty file left by the host's crash at 05:04; please delete it so the worktree can be opened: Recommendation: run the agent's `rm`, but first complete S-0338's acceptance…"
kind: answer
sources: [TH-0390]
topics: [cli, dashboard, git, flai, go, flaiover, sveltekit, code, all]
paths: [design/adrs, design/system/metrics.md, flai/internal/storygit/conflicts.go, flai/internal/storygit/conflicts_test.go, flai/internal/storygit/sync.go, flai/internal/storygit/sync_test.go, flai/cmd/accept.go, flai/cmd/accept_conflict_test.go, flai/internal/metrics/coordination.go, flai/internal/metrics/coordination_test.go, flai/internal/metrics/metrics.go, flai/internal/statsread/statsread.go, flai/internal/statsread/statsread_test.go, flai/cmd/stats.go, flai/cmd/check_stats_test.go, flaiover/src/lib/viz/charts.ts, flaiover/src/lib/viz/charts.test.ts, "flaiover/src/routes/charts/[kind]/+page.svelte", "flaiover/src/routes/charts/[kind]/charts.svelte.test.ts", design/system/flaiover-dashboard.md, docs/users/flai.md, docs/users/flaiover.md, docs/users/flai-reference.md, design/system/flai-cli.md, design/issues/I-0125-a-settled-conversation-between-two-stories-agents-keeps-waking-each-in-turn-because-every-reply-awaits-the-other.md, design/issues/summary.md, design/issues/I-0137-a-host-crash-during-flai-story-start-leaves-the-story-s-branch-ref-worktree-and-narrative-as-empty-files-and-no-agent-can-repair-the-ref.md, design/issues/I-0140-testactivitywriterstaketurns-fails-about-one-run-in-three-two-writers-that-both-find-the-activity-lock-stale-can-both-take-it-and-an-entry-is-lost.md, design/issues/I-0139-flai-check-finds-markdown-md012-outside-the-story-at-close-out.md, design/issues/I-0136-an-acceptance-on-main-during-a-close-out-s-tiers-fails-its-final-sync-check-so-the-whole-close-out-runs-again.md]
date: 2026-10-12
status: compacted
absorbed_by: ""
asked_again: []
cited: []
---

# D-0274 S-0337's branch ref is an empty file left by the host's crash at 05:04; please delete it so the worktree can be opened: Recommendation: run the agent's `rm`, but first complete S-0338's acceptance…
