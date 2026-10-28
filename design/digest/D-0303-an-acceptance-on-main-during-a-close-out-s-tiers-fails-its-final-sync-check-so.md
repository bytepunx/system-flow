---
id: D-0303
title: An acceptance on main during a close-out's tiers fails its final sync check, so the whole close-out runs again
kind: cause
sources: [I-0136]
topics: [cli, conventions, template, testing, flai, go, code, all]
paths: [design/adrs/README.md, flai/internal/verify/story.go, flai/internal/verify/paths.go, flai/internal/verify/story_test.go, flai/internal/verify/paths_test.go, scripts/close-out.sh, template/root/scripts/close-out.sh, template/CHANGELOG.md, design/system/workflow.md, design/system/flai-cli.md, docs/users/flai.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, design/issues/I-0136-an-acceptance-on-main-during-a-close-out-s-tiers-fails-its-final-sync-check-so-the-whole-close-out-runs-again.md, design/issues/summary.md, design/adrs/0156-a-close-out-that-main-outran-syncs-and-verifies-again-at-most-three-rounds-and.md, docs/users/flai-reference.md, flai/cmd/verify.go]
date: 2026-10-14
status: absorbed
absorbed_by: "design/adrs/0156-a-close-out-that-main-outran-syncs-and-verifies-again-at-most-three-rounds-and.md#decision"
asked_again: []
cited: []
---

# D-0303 An acceptance on main during a close-out's tiers fails its final sync check, so the whole close-out runs again

I-0136 was closed on 2026-10-14.

The cause:

```text
An acceptance on main during a close-out's tiers fails its final sync check, so the whole close-out runs again
```

The fix:

```text
Fixed by S-0374 (ADR-0156). When the close-out's last check finds that main moved during its tiers, scripts/close-out.sh, here and in the template, now runs flai stream sync itself and, when the sync ends clean, verifies, commits, and checks again, at most three rounds. Across that clean rebase flai verify reuses the tiers the last run passed while the branch's own change keeps the same patch id, and runs again only those that main's new paths select, leaving out wip paths the branch does not change. Tests reproduce I-0136 with fake git and in a real repository.
```

Source: I-0136.
