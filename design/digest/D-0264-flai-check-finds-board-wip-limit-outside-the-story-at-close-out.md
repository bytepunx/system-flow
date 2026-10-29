---
id: D-0264
title: "flai check finds `board.wip-limit` outside the story at close-out"
kind: cause
sources: [I-0123]
topics: [cli, conventions, template, flai, go, code, all]
paths: [design/adrs, flai/internal/check/scope.go, flai/internal/check/scope_test.go, flai/cmd/check.go, flai/cmd/check_test.go, docs/users/flai-reference.md, docs/users/flai.md, design/system/flai-cli.md, design/system/continuous-improvement.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, design/issues/I-0123-flai-check-finds-board-wip-limit-outside-the-story-at-close-out.md, design/issues/summary.md]
date: 2026-10-08
status: absorbed
absorbed_by: "design/system/flai-cli.md#commands"
asked_again: []
cited: []
---

# D-0264 flai check finds `board.wip-limit` outside the story at close-out

I-0123 was closed on 2026-10-08.

The cause:

```text
flai check finds `board.wip-limit` outside the story at close-out
```

The fix:

```text
S-0348 (ADR-0133): a check scoped to a story leaves out every board.wip-limit outside it, so a close-out neither notes nor records a column of the main checkout's board over its limit; the main checkout's own flai check still reports it. TestCheckRecordIssuesLeavesOutABoardWIPLimitOutsideTheStory reproduces S-0321's and S-0339's close-outs.
```

Source: I-0123.
