---
id: D-0347
title: TestASettledReplyStartsTheAgentWaitingOnItOnce in flai/internal/serve fails now and then in the close-out's integration tier on a branch that does not touch serve
kind: cause
sources: [I-0158]
topics: [flai, cli, go, code, all]
paths: [flai/internal/serve/agents_test.go, design/issues/I-0158-testasettledreplystartstheagentwaitingonitonce-in-flai-internal-serve-fails-now-and-then-in-the-close-out-s-integration-tier-on-a-branch-that-does-not-touch-serve.md, design/issues/summary.md, design/issues/I-0183-flai-check-finds-item-chronology-outside-the-story-at-close-out.md]
date: 2026-10-21
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0347 TestASettledReplyStartsTheAgentWaitingOnItOnce in flai/internal/serve fails now and then in the close-out's integration tier on a branch that does not touch serve

I-0158 was closed on 2026-10-21.

The cause:

```text
TestASettledReplyStartsTheAgentWaitingOnItOnce in flai/internal/serve fails now and then in the close-out's integration tier on a branch that does not touch serve
```

The fix:

```text
S-0473: the failure came from the test, not from flai serve. TestASettledReplyStartsTheAgentWaitingOnItOnce stamped the other story's run and the message from separate time.Now() calls, and messages.ToSince takes what is strictly after a run's start. So when the clock's second turned between the two putEnded calls, the message was new to the other story's agent and look started it, rightly by ADR-0154. The test now puts the waiting story's run, and so its message, a minute before the other run started. It also puts the other run a second earlier, which reproduced the failure on every run before the fix.
```

Source: I-0158.
