---
id: D-0645
title: "S-0739 cannot be enqueued: its branch is 37 commits behind main, and the close-out stops at sync: done"
kind: answer
sources: [TH-1226, TH-1227, TH-1228, S-0798]
topics: [orchestration, acceptance, verify, git, flai, close-out, merge-queue, all]
paths: [flai/internal/preview/accept.go, flai/cmd/accept.go, scripts/close-out.sh, template/root/scripts/close-out.sh, flai/internal/serve/verifyjob_container.go, design/adrs/0310-a-close-out-whose-first-flai-verify-stops-at-sync-syncs-and-verifies-again.md]
date: 2026-11-16
status: absorbed
absorbed_by: design/adrs/0310-a-close-out-whose-first-flai-verify-stops-at-sync-syncs-and-verifies-again.md
asked_again: []
cited: []
---

# D-0645 S-0739 cannot be enqueued: its branch is 37 commits behind main, and the close-out stops at sync: done

Since S-0798 (ADR-0310, 2026-11-16), a close-out asked for a story whose branch is behind main syncs the branch itself. When the round's first `flai verify` stops at sync, the close-out commits what is outstanding, runs `flai stream sync`, and verifies again, within ADR-0156's three rounds. So the orchestrator clears a `verify_behind_main` blocker by asking for the close-out again, as ADR-0291 §4 says. The operator no longer syncs the branch by hand. A stop at sync now lasts only while the sync conflicts, which needs the story's agent.

S-0798's agent chose this fix over the issue's two fallbacks: "Fix chosen: the close-out syncs when its round's first flai verify stops at sync, committing what is outstanding first, within ADR-0156's three rounds. The issue's fallback alternatives (keep a passing report, or a verify_failed that tells a stop at sync apart) are not built: once the close-out syncs, a stop at sync lasts only while the sync conflicts, which needs the agent anyway" (wip/archive/agents/S-0798.md, Decisions).

Before that fix, the answer was the operator's. On 2026-11-15 three stories in review were behind main: S-0739 (37 commits behind, TH-1226), S-0742 (51 behind, TH-1227), and S-0747 (51 behind, TH-1228). Their agents had ended, and each accept dry run had one blocker, `verify_behind_main` ([ADR-0291](../adrs/0291-every-acceptance-but-the-merge-queue-s-landing-waits-for-a-flai-verify-over.md)). The close-out the orchestrator asked for stopped at `sync`, and `queue_add` then refused the story as `verify_failed`. The orchestrator recommended running `flai stream sync` in each worktree and said "Rebasing the story branch is the story's work, not mine." The operator, alex, synced the three branches and answered each thread "done" between 2026-11-16T00:27:35Z and 00:28:30Z. The orchestrator then asked for each close-out again and enqueued each story once it passed.

Sources: TH-1226, TH-1227, TH-1228, S-0798 ([ADR-0310](../adrs/0310-a-close-out-whose-first-flai-verify-stops-at-sync-syncs-and-verifies-again.md), I-0383). D-0646 and D-0647 are merged here. D-0633 covers a different send-back: a story whose verify tier main's new paths sent back, which its own agent syncs and closes out again.
