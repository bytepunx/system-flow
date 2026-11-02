---
id: T-2436
type: task
nature: feature
title: flai accept sends a story whose rebase it cannot settle back to its agent, whoever accepts
status: done
parent: S-0566
owner: alex
created: 2026-11-02T18:35:08Z
updated: 2026-11-02T19:01:48Z
transitions:
  - to: ready
    at: 2026-11-02T18:36:51Z
    by: system-flow
  - to: in-progress
    at: 2026-11-02T18:36:52Z
    by: system-flow
  - to: done
    at: 2026-11-02T19:01:48Z
    by: agent-S-0566
stream: S-0566
tags: []
touches: [flai/cmd/accept.go, flai/cmd/accept_conflict_test.go, flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, flai/internal/storygit/sync.go, flai/internal/preview/accept.go, flai/internal/preview/accept_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go]
usage:
  source: log
  seconds: 1496
  estimated: true
  priced_by: harness
  models:
    - model: claude-fable-5-1
      input: 421
      output: 12159
      cache_read: 4609183
      cache_write: 56026
      cost: 2.627
---
# T-2436 flai accept sends a story whose rebase it cannot settle back to its agent, whoever accepts

## Work

ADR-0248 (`design/adrs/0248-*.md`), the send-back part. Today `flai accept` aborts a rebase that stops on conflicts after the merge of the issue files, the generated files, and the shared markdown files, moves the story back to in-progress, opens a thread, and answers that it was sent back only under `--by orchestrator` (`sendBackUnsettled` in `flai/cmd/accept.go`, ADR-0189); every other acceptance returns the `*syncStopped` error, which tells the operator to resolve each path and `git rebase --continue`. Make the send back the answer of every acceptance that lands here: the operator's by hand, the dashboard's (which runs `flai accept` through flai serve), and the orchestrator's; keep `--batch`, the merge queue's landing, as it is. Generalize the words of the reason, the thread's text, the `sentBack.Reason`, and the help text, which say "the orchestrator's acceptance": name the acceptance by whoever accepted (`res.By`, or `the acceptance of S-nnnn`), and keep the JSON refusal shape `{"refused": {"id", "reason", "sent_back", "branch", "base", "conflicts", "thread"}}` with exit 4. In `flai/internal/serve/agents.go`, `startAcceptSentBack` tells the agent "the orchestrator's acceptance found that …": say "the acceptance of S-nnnn found that …" (the thread's title carries the paths; its opener is still `storygit.SentBackAuthor`), and fix the comments in `flai/internal/storygit/sync.go` that say the orchestrator's acceptance alone sends back. A rebase found already in progress in the worktree (`stop.Found`) is reached only past the preview's uncommitted-worktree blocker, so leave its handling as it is. In `flai/internal/preview/accept.go`, the conflict-markers blocker (`ConflictMarkers`, "resolve each conflict … commit, and accept again") is read by the operator on the review page: say instead that the story is sent back to in-progress so that its agent resolves each conflict, keeping what both sides meant, and comes back to review; keep its path and line list. Update the `flai accept` help text in `accept.go` (the paragraph "When its rebase of the story branch stops on conflicts …" and "Any other acceptance leaves the stopped rebase in the worktree") to the new rule, citing ADR-0248. Tests: in `flai/cmd/accept_conflict_test.go`, add a test that the operator's acceptance by hand (no `--by orchestrator`) of a story whose rebase stops aborts the rebase, moves the story to in-progress with the paths in the reason, opens the thread by `flai`, leaves the worktree with no rebase in progress, and prints the `sent_back` refusal with `--json`; keep `TestTheOrchestratorSendsBackAStoryWhoseRebaseItCannotSettle` passing (its words may change). In `flai/internal/serve/agents_test.go`, adjust the test of `startAcceptSentBack`'s words if one pins them. Run `flai test` on the paths you changed.

This task waits for no other: it is layer 1.

## Done when

- `flai accept S-nnnn --by alex` on a story whose rebase stops on an unsettled conflict aborts the rebase, moves the story to in-progress as `flai`, opens (or writes on) the `Acceptance sent back:` thread, and exits saying the story was sent back; with `--json` it exits 4 with the `sent_back` refusal.
- No text `flai accept` prints, no thread it writes, and no prompt `flai serve` gives the agent tells the operator to rebase, resolve a conflict, or `git rebase --continue`; the conflict-markers blocker says to send the story back.
- `flai serve` starts the story's agent on the thread as it did for the orchestrator's send back, with words that name the acceptance, not the orchestrator.
- The tests named above pass with `flai test`.

## Notes

Acceptance criteria 1 and 2 of S-0566.
