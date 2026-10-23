---
id: D-0159
title: flai's wip markdown lint does not flag a space inside a code span, so a task body reached main and failed a story's close-out
kind: cause
sources: [I-0072]
topics: [cli, go, flai, code, all]
paths: [flai/internal/mdlint, design/system/flai-cli.md, docs/users/flai.md, design/issues]
date: 2026-10-05
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0159 flai's wip markdown lint does not flag a space inside a code span, so a task body reached main and failed a story's close-out

I-0072 was closed on 2026-10-05.

The cause:

```text
flai's wip markdown lint does not flag a space inside a code span, so a task body reached main and failed a story's close-out
```

The fix:

```text
S-0262 added MD038, spaces inside a code span, to flai's markdown lint (flai/internal/mdlint), agreeing with markdownlint-cli2 0.20.0 on a code-spans fixture that holds T-0822's line, so flai check reports such a span and the commands that write wip markdown refuse it
```

Source: I-0072.
