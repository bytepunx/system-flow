---
id: ADR-0310
title: A close-out whose first flai verify stops at sync syncs and verifies again within its three rounds
status: accepted
date: 2026-11-16
supersedes: []
superseded_by: []
refines: [ADR-0156, ADR-0291]
---

# ADR-0310 A close-out whose first flai verify stops at sync syncs and verifies again within its three rounds

## Context

[ADR-0156](0156-a-close-out-that-main-outran-syncs-and-verifies-again-at-most-three-rounds-and.md) §1 has the close-out sync and verify again when its last check, `flai verify <story> --sync-only`, stops at sync: main moved while the tiers ran. It does nothing of the kind when the round's first step, `flai verify <story> --record-issues`, stops at sync. That happens when the branch was already behind main when the close-out was asked for. The close-out then ended at once.

[ADR-0291](0291-every-acceptance-but-the-merge-queue-s-landing-waits-for-a-flai-verify-over.md) §4 tells the operator and the orchestrator to clear a `verify_behind_main` blocker by asking for the story's close-out again, and says the close-out syncs the branch onto main. With the close-out as it was, it did not.

[I-0383](../issues/I-0383-a-close-out-asked-for-a-story-whose-branch-is-behind-main-stops-at-sync-without.md) records it. On 2026-11-15 the orchestrator asked for the close-outs of S-0739, S-0742, and S-0747, each 37 to 51 commits behind main. Each job ended in under a second at sync. Each left a failed report as its story's last, so the merge queue, which had listed all three as eligible on their earlier passing verify, refused them as `verify_failed` (TH-1226, TH-1227, TH-1228) until each was synced by hand and closed out again.

## Decision

A close-out whose first `flai verify` stops at sync syncs and verifies again within its three rounds.

1. **When.** A round's `flai verify <story> --record-issues` fails and its last line is `verify: <story> stopped at sync (…)`, and the round is not the last of three. At any other step, or in the last round, the close-out stops at `flai verify` as before.
2. **What it does.** It commits what is outstanding with the request's commit message, as the commit step does, since `flai stream sync` refuses uncommitted changes; it stops at the commit when there are changes and no message. It prints `close-out: the branch is behind main; flai stream sync <story> (round <n> of 3)`, runs `flai stream sync <story>`, and, when the sync ends clean, starts the next round with `flai verify` again. ADR-0156's reuse across a clean rebase applies to that round.
3. **A conflict stops it.** A sync that stops on a conflict stops the close-out at `flai stream sync`, with the hint ADR-0156's rounds give: resolve, `git add`, sync again, and ask for the close-out again.
4. **Where.** `scripts/close-out.sh` and the template's script do this, and so does a container story's close-out, whose commit and sync run on the host through flai's own code ([ADR-0274](0274-the-close-out-of-a-story-whose-agent-runs-in-a-container-runs-its-verify-in.md)). The three rounds are shared: a round started by the first verify's stop and one started by the last check's stop count alike.

ADR-0291 §4's words, and the `verify_behind_main` blocker's, then hold as written: asking for the close-out again syncs the branch onto main.

This refines ADR-0156, whose rounds now also start from the first verify's stop at sync, and ADR-0291, whose remedy now does what it says.

## Consequences

- A close-out asked for a branch behind main needs no agent or operator to sync it first. The orchestrator's close-outs after `verify_behind_main` sync and verify, and leave a passing report when the sync is clean.
- A stop at sync now lasts only while the sync conflicts, which needs the story's agent anyway.
- The commit before the sync is made without a verify first; the next round verifies it, as the last check's round verifies what the commit step committed.
- `flai/internal/closeout/script_test.go` runs both scripts against stubs, so a change to their steps is tested. A change to the script's steps still changes `containerCloseOut.run` in Go too (ADR-0274).
- The close-out rule in `work-management.md`, `docs/users/flai.md`, and the scripts' READMEs say a close-out asked for a branch behind main syncs first.

## Alternatives considered

- **Keep a passing report as the story's last when a close-out stops at sync before any tier ran.** It would keep the queue's eligibility, but the story would still not be synced, and `verify_behind_main` would still hold it. Once the close-out syncs, a stop at sync leaves a failed report only on a conflict, when the story is rightly not eligible.
- **Have the queue's `verify_failed` tell a stop at sync apart.** The same: it hides the stop rather than clearing it, and a story whose branch conflicts with main cannot land.
- **Have `flai verify` sync itself.** `flai verify` reads and judges; the close-out is where the branch is written. Keeping the sync in the close-out keeps one place that rebases.
