---
id: D-0306
title: "TestTheGateIsAskedOnEveryAttempt fails now and then: the state reads refused before the second refusal is logged"
kind: cause
sources: [I-0146]
topics: [flai, cli, go, code, all]
paths: [flai/internal/channel/channel.go, flai/internal/channel/channel_test.go, design/issues/I-0146-testthegateisaskedoneveryattempt-fails-now-and-then-the-state-reads-refused-before-the-second-refusal-is-logged.md, design/issues/summary.md, design/issues/I-0141-the-mcp-test-tool-of-the-installed-flai-1-40-1-refuses-main-s-manifest-so-a-story-s-agent-runs-flai-test-from-the-tree.md]
date: 2026-10-14
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0306 TestTheGateIsAskedOnEveryAttempt fails now and then: the state reads refused before the second refusal is logged

I-0146 was closed on 2026-10-14.

The cause:

```text
TestTheGateIsAskedOnEveryAttempt fails now and then: the state reads refused before the second refusal is logged
```

The fix:

```text
S-0389: Client.Run in flai/internal/channel/channel.go now logs each attempt's end before setState publishes it, so a reader that sees the state refused finds the refusal logged; TestTheStateIsPublishedAfterTheAttemptIsLogged logs through a slow writer and failed on the old order
```

Source: I-0146.
