---
id: D-0082
title: Work items written in the main checkout reach CI without a markdown lint
kind: cause
sources: [I-0027]
topics: [flai, cli, go, template, conventions, code, all]
paths: [flai/internal/threads, flai/internal/check, flai/internal/workitem, flai/internal/itemnew, flai/internal/itemedit, flai/internal/mcpserver/items_write.go, flai/internal/mdlint, scripts/mdlint-fixtures.sh, scripts/README.md, Makefile, design/system/flai-cli.md, design/adrs, design/conventions/tooling.md, template/root/design/conventions/tooling.md, template/template.yaml, template/CHANGELOG.md, docs/users/flai.md, design/issues, wip/archive/agents]
date: 2026-10-01
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0082 Work items written in the main checkout reach CI without a markdown lint

I-0027 was closed on 2026-10-01.

The cause:

```text
Work items written in the main checkout reach CI without a markdown lint
```

The fix:

```text
S-0179 (ADR-0061): flai check lints wip markdown with the project's markdownlint configuration and warns on each finding, item creation, thread entries, and narrative log entries refuse what would bring one, titles lose MD026 punctuation, and the mirror block's removal leaves one blank line.
```

Source: I-0027.
