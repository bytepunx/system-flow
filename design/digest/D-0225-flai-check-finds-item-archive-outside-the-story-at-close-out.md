---
id: D-0225
title: "flai check finds `item.archive` outside the story at close-out"
kind: cause
sources: [I-0078]
topics: [cli, conventions, template, flai, go, code, all]
paths: [design/adrs, flai/internal/check/scope.go, flai/internal/check/scope_test.go, flai/cmd/check.go, flai/cmd/check_test.go, docs/users/flai-reference.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, design/system/continuous-improvement.md, design/system/flai-cli.md, docs/users/flai.md, design/issues/I-0078-flai-check-finds-item-archive-outside-the-story-at-close-out.md, design/issues/summary.md, design/issues/I-0079-testroundtriprepositoryitems-reads-the-live-main-checkout-and-fails-a-close-out-when-another-agent-edits-a-story-mid-run.md]
date: 2026-10-07
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0225 flai check finds `item.archive` outside the story at close-out

I-0078 was closed on 2026-10-07.

The cause:

```text
flai check finds `item.archive` outside the story at close-out
```

The fix:

```text
Fixed by S-0280 (ADR-0122): a check scoped to a story leaves out every item.archive outside it, so a close-out records none. Every instance named S-0250, cancelled from backlog and left unarchived, which no story branch can archive; the main checkout's flai check still warns on such an item.
```

Source: I-0078.
