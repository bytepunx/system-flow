---
id: D-0339
title: "flai check finds `board.order` outside the story at close-out"
kind: cause
sources: [I-0143]
topics: [agents, cli, flai, go, code, all]
paths: [flai/internal/guard/call.go, flai/internal/guard/call_test.go, flai/internal/guard/claudecode.go, flai/internal/guard/claudecode_test.go, flai/internal/guard/guard.go, flai/internal/guard/guard_test.go, flai/internal/guard/shared.go, flai/internal/guard/shared_test.go, flai/internal/guard/subagents.go, flai/internal/guard/subagents_test.go, flai/cmd/guard.go, flai/cmd/guard_test.go, docs/operators/settings.md, docs/users/flai-reference.md, design/system/flai-cli.md, flai/internal/guard/authors.go, flai/internal/guard/authors_test.go]
date: 2026-10-21
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0339 flai check finds `board.order` outside the story at close-out

I-0143 was closed on 2026-10-21.

The cause:

```text
flai check finds `board.order` outside the story at close-out
```

The fix:

```text
Not seen since 2026-10-12: all five reports are one cause on one day, S-0353 left in board.order while in progress. Every flai move out of ready or backlog now drops the story from the order (flai/internal/workitem/rules.go, RemoveFromOrder), board.md's order is empty, and 69 stories have started since with no recurrence. Reopen it with a new issue if flai check finds board.order again.
```

Source: I-0143.
