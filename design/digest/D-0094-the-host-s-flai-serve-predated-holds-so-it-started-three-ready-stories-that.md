---
id: D-0094
title: The host's flai serve predated holds, so it started three ready stories that open stories held
kind: cause
sources: [I-0049]
topics: [flai, cli, go, flaiover, dashboard, sveltekit, code, all]
paths: [flai/internal/workitem, flai/internal/threads, flai/internal/issues, flai/internal/manifest, flai/internal/serve, flai/internal/mcpserver, flai/internal/check, flai/internal/hostapi, flai/internal/release, flai/internal/buildinfo, flai/cmd, flaiover/src, design/system/flai-cli.md, design/system/work-hierarchy.md, design/system/project-manifest.md, design/issues, docs/operators, docs/users/flai.md, system-flow.yaml]
date: 2026-10-01
status: absorbed
absorbed_by: "design/system/work-hierarchy.md#fields-an-older-flai-does-not-know"
asked_again: []
cited: []
---

# D-0094 The host's flai serve predated holds, so it started three ready stories that open stories held

I-0049 was closed on 2026-10-01.

The cause:

```text
The host's flai serve predated holds, so it started three ready stories that open stories held
```

The fix:

```text
S-0181: flai serve logs, the MCP inbox reports (flai_outdated), and the dashboard's host badge shows when the host's flai is older than the newest flai/v* tag reachable from the project's HEAD, with the upgrade command; system-flow.yaml's flai.minimum stops an older flai before it reads any item, and publishing a flai release that changes the front-matter fields raises it.
```

Source: I-0049.
