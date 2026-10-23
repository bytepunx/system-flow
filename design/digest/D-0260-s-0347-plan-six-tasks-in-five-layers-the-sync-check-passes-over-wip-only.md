---
id: D-0260
title: "S-0347 plan: six tasks in five layers; the sync check passes over wip-only commits on main: The plan fits S-0347, and I recommend the operator approve it"
kind: answer
sources: [TH-0379]
topics: [testing, git, flai, cli, go, template, conventions, code, all]
paths: [design/adrs/0133-the-close-out-s-sync-check-passes-over-commits-on-main-that-change-only-wip.md, design/adrs/README.md, design/system/flai-cli.md, design/system/devex.md, flai/internal/verify/story.go, flai/internal/verify/story_test.go, flai/internal/verify/paths.go, flai/internal/verify/paths_test.go, flai/cmd/verify.go, flai/cmd/verify_test.go, docs/users/flai.md, docs/users/flai-reference.md, scripts/close-out.sh, template/root/scripts/close-out.sh, scripts/README.md, template/root/scripts/README.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, design/issues/I-0119-the-close-out-s-last-check-that-the-branch-contains-main-fails-when-flai-commits-wip-on-main-during-its-run.md, design/issues/summary.md, design/adrs/0135-the-sync-step-of-flai-verify-and-the-close-out-s-last-check-pass-over-commits.md, design/issues/I-0126-flai-task-done-commits-the-changed-paths-of-another-open-task-of-the-same-layer.md, docs/operators/settings.md]
date: 2026-10-08
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0260 S-0347 plan: six tasks in five layers; the sync check passes over wip-only commits on main: The plan fits S-0347, and I recommend the operator approve it

The operator, alex, answered TH-0379 at 2026-10-08T08:48:20Z.

What was asked:

```text
**Proposed fix.** The `sync` step of `flai verify` passes when every main commit the branch lacks changes only paths under the manifest's `wip` folder, and the branch changes none of those paths. The step names the commits it skipped in a note. The close-out's last check calls the same rule through a new `flai verify S-nnnn --sync-only` instead of `git merge-base --is-ancestor`.
```

The answer:

```text
The plan fits S-0347, and I recommend the operator approve it. One point needs your confirmation, because it is a decision not yet recorded: T-1375 and T-1379 change the definition of done ("the branch contains main"), in this repository and in the template's baseline `work-management.md`.

Why I recommend it:

- I hit this exact case on S-0320 today. Its verify at the head stopped at `sync`, 2 commits behind, and both commits were flai's own "replan forecasts" wip commits. Under ADR-0093 I could not accept a story whose work was finished, and you accepted it by hand.
- The rule is narrow. It passes over main commits that change only `wip/` paths the branch does not change. Anything else, such as an acceptance or a release, still forces a sync and a new run. Acceptance merges the branch anyway, so the skipped commits come in at the merge.
- ADR-0093 is unaffected. I still accept only a story verified at its branch head, and the head is unchanged by commits on main.

The tasks, layers, 1h30m forecast and 75 USD/week stand. The hold behind S-0341 (assumption 5) is right, so no `after` is needed. Assumption 2, a project on an older flai failing the last check until it upgrades, should carry a minimum flai version in the template's changelog entry. T-1379 plans one.

Source: design/adrs/0093-with-accept-reviews-on-the-orchestrator-accepts-a-story-in-review-through-flai.md
```

Source: TH-0379.
