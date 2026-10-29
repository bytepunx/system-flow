---
id: D-0270
title: "flai check finds `board.wip-limit` outside the story at close-out"
kind: cause
sources: [I-0122]
topics: [cli, conventions, template, flai, go, code, all]
paths: [design/adrs, flai/internal/check/scope.go, flai/internal/check/scope_test.go, flai/cmd/check.go, flai/cmd/check_test.go, docs/users/flai-reference.md, docs/users/flai.md, design/system/flai-cli.md, design/system/continuous-improvement.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, design/issues/I-0123-flai-check-finds-board-wip-limit-outside-the-story-at-close-out.md, design/issues/summary.md]
date: 2026-10-12
status: absorbed
absorbed_by: "design/system/flai-cli.md#commands"
asked_again: []
cited: []
---

# D-0270 flai check finds `board.wip-limit` outside the story at close-out

I-0122 was closed on 2026-10-12.

The cause:

```text
flai check finds `board.wip-limit` outside the story at close-out
```

The fix:

```text
Fixed by S-0348 (ADR-0133), released in flai 1.39.15: a check scoped to a story, as the close-out runs it, leaves out every board.wip-limit outside the story. I-0122's one instance, S-0291's close-out reporting 6 stories in in-progress over a limit of 5, came before that fix and is the "in progress over its limit" case of TestScopedCheckLeavesOutABoardWIPLimitOutsideTheStory in flai/internal/check/scope_test.go, which S-0361 names I-0122 in. Duplicate of I-0123.
```

Source: I-0122.
