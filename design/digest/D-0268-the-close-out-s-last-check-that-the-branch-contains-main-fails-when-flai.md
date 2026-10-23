---
id: D-0268
title: The close-out's last check that the branch contains main fails when flai commits wip on main during its run
kind: cause
sources: [I-0119]
topics: [testing, git, flai, cli, go, template, conventions, code, all]
paths: [design/adrs/0133-the-close-out-s-sync-check-passes-over-commits-on-main-that-change-only-wip.md, design/adrs/README.md, design/system/flai-cli.md, design/system/devex.md, flai/internal/verify/story.go, flai/internal/verify/story_test.go, flai/internal/verify/paths.go, flai/internal/verify/paths_test.go, flai/cmd/verify.go, flai/cmd/verify_test.go, docs/users/flai.md, docs/users/flai-reference.md, scripts/close-out.sh, template/root/scripts/close-out.sh, scripts/README.md, template/root/scripts/README.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, design/issues/I-0119-the-close-out-s-last-check-that-the-branch-contains-main-fails-when-flai-commits-wip-on-main-during-its-run.md, design/issues/summary.md, design/adrs/0135-the-sync-step-of-flai-verify-and-the-close-out-s-last-check-pass-over-commits.md, design/issues/I-0126-flai-task-done-commits-the-changed-paths-of-another-open-task-of-the-same-layer.md, docs/operators/settings.md]
date: 2026-10-08
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0268 The close-out's last check that the branch contains main fails when flai commits wip on main during its run

I-0119 was closed on 2026-10-08.

The cause:

```text
The close-out's last check that the branch contains main fails when flai commits wip on main during its run
```

The fix:

```text
Fixed by S-0347 (ADR-0135): the sync step of flai verify passes when the commits the branch lacks change only wip/ paths the branch does not change, naming them in a note on the step, and the close-out's last check is flai verify --sync-only in place of git merge-base --is-ancestor, in this repository's scripts/close-out.sh and the template's. A commit outside wip/, such as an acceptance or a release during the run, still stops the close-out and needs a sync and another run.
```

Source: I-0119.
