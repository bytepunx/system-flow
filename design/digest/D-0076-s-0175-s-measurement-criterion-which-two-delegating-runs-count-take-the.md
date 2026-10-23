---
id: D-0076
title: "S-0175's measurement criterion: which two delegating runs count: take the recommended approach"
kind: answer
sources: [TH-0042]
topics: [conventions, flai, cli, go, template, code, all]
paths: [flai/internal/harness, flai/internal/context, flai/cmd/prime.go, flai/cmd/prime_test.go, flai/cmd/guard.go, flai/cmd/guard_test.go, flai/cmd/root.go, flai/internal/guard, flai/internal/mcpserver, flai/internal/conventions, template, design/conventions, design/system/agent-context.md, design/system/flai-cli.md, design/system/conventions.md, design/adrs, ".claude", docs/users, docs/operators/index.md, docs/operators/settings.md]
date: 2026-10-01
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0076 S-0175's measurement criterion: which two delegating runs count: take the recommended approach

The operator, alex, answered TH-0042 at 2026-10-01T08:05:17Z.

What was asked:

```text
The criterion asks for cost, cache reads, and turns on at least two stories that delegated, against comparable runs that did not. Stories only get the new prompt once a flai with it is released and installed on the host, which is after S-0175 is accepted. flai serve keeps 8 logs today; only S-0118's shows a sub-agent (one Explore).
```

The answer:

```text
take the recommended approach
```

Source: TH-0042.
