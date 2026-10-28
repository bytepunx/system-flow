---
id: D-0044
title: "S-0138's pack cannot fit 80 KB under ADR-0049's rules: which way: C - as long as the agent will read the full document if it determines from the brief that it is relevant"
kind: answer
sources: [TH-0032]
topics: [context, flai, cli, go, template, conventions, code, all]
paths: [flai/cmd, flai/internal/context, flai/internal/mcpserver/folder.go, design/adrs, design/system/agent-context.md, design/system/flai-cli.md, docs/users, CLAUDE.md, design/conventions/session-start.md, template]
date: 2026-09-29
status: absorbed
absorbed_by: "design/adrs/0050-a-document-a-story-names-only-by-its-path-written-out-is-briefed-when-it-is.md#decision"
asked_again: []
cited: []
---

# D-0044 S-0138's pack cannot fit 80 KB under ADR-0049's rules: which way: C - as long as the agent will read the full document if it determines from the brief that it is relevant

The operator, alex, answered TH-0032 at 2026-09-29T03:51:59Z.

What was asked:

```text
The budget, briefs, and ranked fill are built. Measured with `flai prime --story S-0138` on the S-0146 branch (header included):
```

The answer:

```text
C - as long as the agent will read the full document if it determines from the brief that it is relevant
```

Source: TH-0032.
