---
id: ADR-0192
title: "flai verify's patch id leaves out the files a sync's rebase writes again, so a round after a rebase that rewrote design/issues/summary.md still reuses the tiers main's new paths do not select"
status: accepted
date: 2026-10-23
supersedes: []
superseded_by: []
refines: [ADR-0156]
---

# ADR-0192 flai verify's patch id leaves out the files a sync's rebase writes again, so a round after a rebase that rewrote design/issues/summary.md still reuses the tiers main's new paths do not select

## Context

[ADR-0156](0156-a-close-out-that-main-outran-syncs-and-verifies-again-at-most-three-rounds-and.md) lets `flai verify` reuse, across a clean rebase onto a newer main branch, the tiers that passed while the branch's own change stays the same: its patch id, `git patch-id --stable` of its diff from the merge base. A tier runs again only when the paths main brought in select it.

A story that records or closes an issue changes `design/issues/summary.md`. When main changed the issues too, the rebase stops on the summary alone and the sync writes it again from the issue files ([ADR-0098](0098-flai-stream-sync-and-flai-accept-regenerate-design-issues-summary-md-when-a.md)). The branch's summary then differs from before the rebase, so its patch id differs, and the close-out's next round reused nothing: go-test ran again, about six minutes, though main had changed only `design/issues` ([I-0196](../issues/I-0196-a-close-out-s-second-round-runs-go-test-again-when-the-rebase-regenerates-design-issues-summary-md-which-changes-the-branch-s-patch-id-so-a-close-out-main-outruns-still-takes-about-11-minutes-with-the-queue-on.md)).

## Decision

The patch id `flai verify` records and compares leaves out the files flai writes again at a sync's rebase, the paths `issues.Generated` names: `design/issues/summary.md` alone today. `verify.BranchPatch` takes them and diffs with an exclude pathspec for each. The patch is still recorded only from a worktree with nothing uncommitted, the left-out files included, and is empty when the branch changes nothing else.

ADR-0156's gate is unchanged: one patch id, and the tiers main's new paths select run again with every tier after them.

## Consequences

- A generated file is a function of others. The branch's own issue files stay in the patch, and main's issue files are among main's new paths, which select the tiers they would. Leaving the summary out loses nothing a tier could read.
- A round after such a rebase reuses the tiers main's paths do not select, go-test among them when main changed only `design/issues`.
- A report written before this change has a patch id that covered the summary, so the first rebase after an upgrade reuses nothing, once.
- A branch whose only change is a generated file records no patch, and reuses nothing across a rebase.
- A file added to `issues.Generated` later is left out of the patch with no change here.

## Alternatives considered

- Leave the whole `design/issues` folder out of the patch. Rejected: the branch's own issue files are its change, and a tier may one day read them.
- Compare the patches with the summary's diff stripped of its timestamps. Rejected: the regenerated summary also differs by main's issues, which is main's change, not the branch's, and parsing it ties verify to its format.
- Run go-test again and accept the cost. Rejected: it is the cost ADR-0156 set out to remove.
