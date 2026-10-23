---
id: D-0277
title: "flai check finds `threads.archived` outside the story at close-out"
kind: cause
sources: [I-0135]
topics: [flai, cli, go, template, conventions, code, all]
paths: [design/adrs, design/adrs/README.md, flai/internal/check/scope.go, flai/internal/check/scope_test.go, flai/cmd/check.go, flai/cmd/check_test.go, docs/users/flai-reference.md, docs/users/flai.md, design/system/flai-cli.md, design/system/continuous-improvement.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, design/issues/I-0135-flai-check-finds-threads-archived-outside-the-story-at-close-out.md, design/issues/summary.md]
date: 2026-10-12
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0277 flai check finds `threads.archived` outside the story at close-out

I-0135 was closed on 2026-10-12.

The cause:

```text
flai check finds `threads.archived` outside the story at close-out
```

The fix:

```text
S-0366: ADR-0137 decides that a check scoped to a story leaves out every threads.archived outside it, and ScopeToStory (flai/internal/check/scope.go) now drops it beside item.archive, board.wip-limit, and issues.no-story, so a close-out neither notes nor records an open thread on an archived item it did not touch, as TH-0389 on S-0338 was at the close-outs of S-0343, S-0289, and S-0342. TestScopedCheckLeavesOutAThreadsArchivedOutsideTheStory and TestCheckRecordIssuesLeavesOutAThreadsArchivedOutsideTheStory reproduce it. One on a thread anchored on the story or its tasks stays inside, and flai check in the main checkout still warns on every open thread on an archived item.
```

Source: I-0135.
