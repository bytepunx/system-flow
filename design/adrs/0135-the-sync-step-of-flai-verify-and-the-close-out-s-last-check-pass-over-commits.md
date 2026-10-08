---
id: ADR-0135
title: "The sync step of flai verify and the close-out's last check pass over commits on main that change only wip paths the branch does not change"
status: accepted
date: 2026-10-08
supersedes: []
superseded_by: []
refines: [ADR-0110]
---

# ADR-0135 The sync step of flai verify and the close-out's last check pass over commits on main that change only wip paths the branch does not change

## Context

`flai verify`'s `sync` step fails when the story's branch lacks any commit of the main branch (`git rev-list --count HEAD..<main>` above zero). The close-out runs `flai verify --record-issues`, commits, and then checks again with `git merge-base --is-ancestor <main> HEAD` that the branch contains the main branch ([ADR-0110](0110-a-story-s-agent-runs-the-close-out-which-runs-flai-verify-itself-before-review.md)).

A close-out takes four to thirteen minutes, most of it in the integration and smoke tiers. While it runs, `flai` commits `wip/` on the main branch: forecast replans, edits of drafts, criteria, and touches, and the moves of other stories. I-0119 has six instances, from the close-outs of S-0324, S-0316, S-0320, S-0321, S-0322, and S-0297 on 2026-10-08, where every step of `flai verify` passed and the last check failed. In most of them every commit the branch lacked changed only `wip/`; S-0321's second run lacked 17 such commits. Each failure meant another `flai stream sync` and another full run, which the next replan could race again.

Such a commit cannot change what the run verified:

1. **`wip/` lives on the main branch** ([ADR-0019](0019-story-branches-and-touches.md)). A story branch holds `wip/` as the main branch last committed it and does not change it; flai writes work items, narratives, and threads in the main checkout.
2. **The tiers that read the whole repository scope themselves to the story** under `CLOSE_OUT_STORY` ([ADR-0085](0085-a-close-out-s-flai-check-reports-findings-outside-the-story-as-notes-and.md)), and leave out the `wip/` markdown the story does not change (S-0345).
3. **Acceptance merges the branch**, so the merge brings the `wip/` commits as it brings any other.

## Decision

The `sync` step of `flai verify` passes when the branch contains the main branch but for commits that change only paths under the manifest's `wip` folder, none of which the branch changes.

1. **What the branch lacks is what the main branch changed since the branch left it**: the paths `git diff --no-renames --name-only HEAD...<main>` lists, from the merge base to the main branch, a rename as both of its paths. When none is outside the `wip` folder and none is among the paths the branch changed (`verify.ChangedPaths`, its commits and what it has not committed), the step passes. Its result carries a note naming how many commits it passed over and their short hashes, which the text prints under the step. The note is the step's own, not a check note: `--record-issues` records nothing for it.
2. **Any other commit fails the step as now**, naming `flai stream sync`: an acceptance, a release, or a `wip/` commit that changes a path the branch changes too.
3. **`flai verify <story> --sync-only`** runs the `rebase` and `sync` steps alone, prints them as a full run does, exits as a full run does, and stores no report, so the stored report stays the last full run's. It is refused with `--last` and with `--record-issues`.
4. **The close-out's last check is `flai verify <story> --sync-only`**, in this repository's `scripts/close-out.sh` and the template's, in place of `git merge-base --is-ancestor`. The rule lives in one place, and the close-out and `flai verify` keep agreeing.
5. **The definition of done says it**: the story's branch contains the main branch but for commits that change only `wip/` paths it does not change.

This refines ADR-0110, whose close-out checks again after its commit that the branch contains the main branch.

## Consequences

- A close-out whose run overlaps flai's own `wip/` commits on the main branch passes, as five of I-0119's six instances would have. I-0119 can be closed.
- A story still goes to review with its branch behind the main branch by those commits. Acceptance merges them; the merge cannot conflict, since the branch changes none of their paths.
- An acceptance or a release during the run still stops the close-out at its last check, as in I-0119's third instance, since its commits change paths outside `wip/` that the tiers may read.
- The template's `scripts/close-out.sh` needs a flai that has `--sync-only`. The template's changelog names the version.
- The note on the `sync` step is printed by `flai verify` and kept in the stored report; the dashboard's review page does not show it yet.

## Alternatives considered

- **Sync automatically and run everything again.** It costs the same full run each time, and the next replan can race it again.
- **Resume at the failed step, as S-0341 makes `flai verify` do.** A sync rebases the branch, which changes its head, and the resume requires the same head, so every tier runs again.
- **Check each commit the branch lacks, rather than the net change.** A commit outside `wip/` that a later commit undoes would fail the step though the main branch's content outside `wip/` is what the run verified. The net change from the merge base is what acceptance merges.
- **Hold flai's `wip/` commits while a close-out runs.** It would delay the operator's and other agents' writes for minutes, and couple every writer of `wip/` to every close-out.
