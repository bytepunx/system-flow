---
id: D-0059
title: Hand-edited work items drifted from the schema
kind: cause
sources: [I-0004]
topics: [flai, cli, go, flaiover, dashboard, sveltekit, code, all]
paths: [flai/internal/workitem, flai/internal/threads, flai/internal/issues, flai/internal/manifest, flai/internal/serve, flai/internal/mcpserver, flai/internal/check, flai/internal/hostapi, flai/internal/release, flai/internal/buildinfo, flai/cmd, flaiover/src, design/system/flai-cli.md, design/system/work-hierarchy.md, design/system/project-manifest.md, design/issues, docs/operators, docs/users/flai.md, system-flow.yaml]
date: 2026-10-01
status: absorbed
absorbed_by: "design/conventions/tooling.md#rules"
asked_again: []
cited: []
---

# D-0059 Hand-edited work items drifted from the schema

I-0004 was closed on 2026-10-01.

The cause:

```text
Before flai owned work items, front matter was written by hand and drifted: an unquoted colon in a title broke parsing, an archived story kept an unchecked criterion, and the board order listed backlog stories against the rule of the time.
```

The fix:

```text
Fixed: flai creates and moves items and flai check --strict gates them; no instance since 2026-09-15. Version skew on front matter is I-0051 (S-0181).
```

Source: I-0004.
