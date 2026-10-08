---
id: ADR-0126
title: "A rebase of a story branch merges an issue file both sides changed by its instances, and folds an open issue the branch added into the main branch's open issue of the same title"
status: accepted
date: 2026-10-08
supersedes: []
superseded_by: []
refines: [ADR-0085, ADR-0098]
---

# ADR-0126 A rebase of a story branch merges an issue file both sides changed by its instances, and folds an open issue the branch added into the main branch's open issue of the same title

## Context

A close-out records each check finding outside its story in an issue on the story's branch (ADR-0085), and an agent records friction there with `flai issue new` or `bump`. Both look for an open issue of the same title only in the story's own checkout. Two branches open in parallel can therefore record the same occurrence twice, and I-0112 counts both ways it happened:

- S-0275's close-out opened I-0112 while S-0265 had opened I-0111 under the same title. After S-0265 was accepted and S-0275 synced, both files stood on the branch, the next close-out bumped the lower number, and I-0112 was left a duplicate, removed by hand.
- S-0318's close-out bumped I-0118 on its branch while S-0316's bump of it reached the main branch. The sync stopped on I-0118 and `summary.md`, and the agent merged the two instances by hand.

ADR-0098 already resolves a stop on `summary.md` alone, by writing it again, because it is generated from the issue files. An issue file is not generated, but its instances are records appended by flai, so two sides' changes to one combine without a person's judgment.

## Decision

The rebase of a story branch onto the main branch merges an issue file both sides changed by its instances, and then folds each open issue the branch added into the main branch's open issue of the same title. It does so in `flai stream sync`, in the sync `flai task done` runs, and in `flai accept`:

1. **Merge.** When the rebase stops on an issue file under `design/issues` that both sides changed, flai merges it from the stop's three versions: the base's instances, then each instance either side added, in timestamp order. `count` adds both sides' increments over the base, `cost` is the average of all the occurrences, and `last_reported` and `updated` take the later of the two. Every other front matter field and the rest of the body take the side that changed them. When both sides changed one of them differently, flai merges nothing and the stop is left for the agent. A stop whose paths are all merged issue files and generated files is continued, with `summary.md` written again after the merge, and the trial merge leaves out the issue files as it leaves out the summary.
2. **Fold.** Once the rebase is clean, each open issue the branch added whose title an open issue on the main branch has too is folded into that issue: its instances, count, and cost move into it as the merge adds them, and the branch's file is deleted. flai writes `summary.md` again and commits the fold on the branch as `docs: [S-nnnn] fold I-x into I-y`.

The sync's report and its `--json` name each issue merged and each folded, with the issue it went into. Neither step reads or writes anything outside the story's worktree.

## Consequences

- A duplicate issue never reaches the main branch from a sync or an acceptance, and parallel bumps of one issue no longer stop a sync.
- An issue the branch added keeps its number on the branch only until the main branch's issue of its title arrives; text on the branch that names the folded number names an issue that is gone, and is left as it is.
- The fold changes the main branch's issue file on the branch, which the story's touches may not name; the sync lists it among the paths outside the claim, as for any change.
- A closed issue is never folded into: an occurrence after the main branch closed its issue stays an issue of its own.

## Alternatives considered

- **Number and match issues across every branch when recording**, as `flai adr new` numbers ADRs. It finds an issue another open branch opened, but the occurrence is then recorded in a file that branch also changes, so the two still conflict at acceptance, and the branch may be cancelled, taking the issue with it.
- **Record outside-the-story findings on the main branch.** A story's agent writes the main checkout only for `wip/`, and the issue would then reach the main branch before the story that found it is accepted.
- **Leave it to the agent and the close-out's check.** That is what I-0112 records, at about 8m each time.
