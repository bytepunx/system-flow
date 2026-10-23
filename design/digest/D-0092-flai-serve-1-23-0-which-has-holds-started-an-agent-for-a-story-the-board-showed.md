---
id: D-0092
title: flai serve 1.23.0, which has holds, started an agent for a story the board showed held
kind: cause
sources: [I-0050]
topics: [flai, cli, go, code, all]
paths: [flai/internal/serve, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, design/system/workflow.md, design/system/flai-cli.md, docs/users/flai.md, design/issues/I-0050-flai-serve-1-23-0-which-has-holds-started-an-agent-for-a-story-the-board-showed-held.md, design/issues/summary.md]
date: 2026-10-01
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0092 flai serve 1.23.0, which has holds, started an agent for a story the board showed held

I-0050 was closed on 2026-10-01.

The cause:

```text
flai serve 1.23.0, which has holds, started an agent for a story the board showed held
```

The fix:

```text
S-0182: the operator's Start agent started both held stories with the launcher's prompt; the prompt and the journal now say what an operator's start went past, and resume no longer passes a hold
```

Source: I-0050.
