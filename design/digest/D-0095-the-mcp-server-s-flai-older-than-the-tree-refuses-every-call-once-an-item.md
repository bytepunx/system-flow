---
id: D-0095
title: The MCP server's flai, older than the tree, refuses every call once an item carries a front matter field it does not know
kind: cause
sources: [I-0051]
topics: [flai, cli, go, flaiover, dashboard, sveltekit, code, all]
paths: [flai/internal/workitem, flai/internal/threads, flai/internal/issues, flai/internal/manifest, flai/internal/serve, flai/internal/mcpserver, flai/internal/check, flai/internal/hostapi, flai/internal/release, flai/internal/buildinfo, flai/cmd, flaiover/src, design/system/flai-cli.md, design/system/work-hierarchy.md, design/system/project-manifest.md, design/issues, docs/operators, docs/users/flai.md, system-flow.yaml]
date: 2026-10-01
status: absorbed
absorbed_by: "design/system/work-hierarchy.md#fields-an-older-flai-does-not-know"
asked_again: []
cited: []
---

# D-0095 The MCP server's flai, older than the tree, refuses every call once an item carries a front matter field it does not know

I-0051 was closed on 2026-10-01.

The cause:

```text
The MCP server's flai, older than the tree, refuses every call once an item carries a front matter field it does not know
```

The fix:

```text
S-0181: work items, threads, and issues are decoded leniently on the listing paths; an unknown front-matter field is kept and written back unchanged, logged once per file, and reported by flai check as item/threads/issues.unknown-field. flai.minimum in system-flow.yaml, raised by a flai release that changes the fields, makes an older flai say so instead.
```

Source: I-0051.
