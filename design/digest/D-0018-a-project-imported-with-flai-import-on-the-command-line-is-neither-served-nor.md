---
id: D-0018
title: A project imported with flai import on the command line is neither served nor offered, so the dashboard never shows it
kind: cause
sources: [I-0047]
topics: [flai, cli, go, code, all]
paths: [flai/cmd/import.go, flai/cmd/import_serve.go, flai/cmd/new.go, flai/cmd/repo_url.go, flai/cmd/serve.go, flai/cmd/serve_import.go, flai/cmd/serve_status_below.go, flai/cmd/dashboard_agent.go, flai/internal/serve, design/adrs, docs/users/flai.md, docs/operators, design/system/flai-cli.md]
date: 2026-09-26
status: absorbed
absorbed_by: "design/adrs/0045-a-system-flow-repository-under-a-folder-named-for-import-is-served-whether-or.md#decision"
asked_again: []
cited: []
---

# D-0018 A project imported with flai import on the command line is neither served nor offered, so the dashboard never shows it

I-0047 was closed on 2026-09-26.

The cause:

```text
A project imported with flai import on the command line is neither served nor offered, so the dashboard never shows it
```

The fix:

```text
S-0120: flai import registers with the host flai, flai serve serves the projects below import_roots, and flai serve status says why any is not served
```

Source: I-0047.
