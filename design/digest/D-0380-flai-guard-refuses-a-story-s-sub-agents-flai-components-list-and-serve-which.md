---
id: D-0380
title: flai guard refuses a story's sub-agents flai components list and serve, which only read or run the project's own dev server
kind: cause
sources: [I-0156]
topics: [cli, components, flai, go, code, all]
paths: [flai/internal/guard/guard.go, flai/internal/guard/guard_test.go, flai/cmd/guard.go, flai/cmd/guard_test.go, docs/users/flai-reference.md, docs/users/flai.md, design/system/agent-adapters.md, design/system/component-library.md, design/issues/I-0165-flai-guard-refuses-a-story-s-sub-agent-flai-components-list-and-flai-components-serve-so-a-headless-pass-of-the-harness-falls-to-the-story-s-agent.md, design/issues/summary.md, design/adrs/0186-flai-guard-lets-every-agent-run-flai-components-list-and-a-story-s-task-sub.md, design/adrs/README.md, design/system/flai-cli.md, flai/internal/guard/call.go]
date: 2026-10-23
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0380 flai guard refuses a story's sub-agents flai components list and serve, which only read or run the project's own dev server

I-0156 was closed on 2026-10-23.

The cause:

```text
flai guard refuses a story's sub-agents flai components list and serve, which only read or run the project's own dev server
```

The fix:

```text
fixed by S-0477 (ADR-0186): flai guard's cliReads has a components entry, so flai components list is a read for every agent and a story's task sub-agent may run flai components serve; TestComponents-style cases in flai/internal/guard/guard_test.go reproduce the refusal I-0156 recorded, and pass at c4e2385c2
```

Source: I-0156.
