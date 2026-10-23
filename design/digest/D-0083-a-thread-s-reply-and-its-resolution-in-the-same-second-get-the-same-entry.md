---
id: D-0083
title: A thread's reply and its resolution in the same second get the same entry heading, which fails the markdown lint
kind: cause
sources: [I-0043]
topics: [flai, cli, go, template, conventions, code, all]
paths: [flai/internal/threads, flai/internal/check, flai/internal/workitem, flai/internal/itemnew, flai/internal/itemedit, flai/internal/mcpserver/items_write.go, flai/internal/mdlint, scripts/mdlint-fixtures.sh, scripts/README.md, Makefile, design/system/flai-cli.md, design/adrs, design/conventions/tooling.md, template/root/design/conventions/tooling.md, template/template.yaml, template/CHANGELOG.md, docs/users/flai.md, design/issues, wip/archive/agents]
date: 2026-10-01
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0083 A thread's reply and its resolution in the same second get the same entry heading, which fails the markdown lint

I-0043 was closed on 2026-10-01.

The cause:

```text
A thread's reply and its resolution in the same second get the same entry heading, which fails the markdown lint
```

The fix:

```text
S-0179: threads.Reply and threads.Resolve join an entry by the same author in the same second to the last one under its heading, as narratives and issues do.
```

Source: I-0043.
