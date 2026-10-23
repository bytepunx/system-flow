---
id: D-0315
title: The host's installed flai stays behind the project's releases after I-0147 closed, so the MCP task_done still refuses
kind: cause
sources: [I-0149]
topics: [release, host, mcp, flai, cli, go, code, all]
paths: [docs/operators/index.md, docs/operators/runbooks/update.md, design/system/flai-cli.md, design/issues/I-0149-the-host-s-installed-flai-stays-behind-the-project-s-releases-after-i-0147-closed-so-the-mcp-task-done-still-refuses.md, design/issues/summary.md, design/adrs/0161-flai-serve-s-upgrade-thread-names-a-flai-release-only-once-it-is-published-and.md, design/adrs/README.md, flai/cmd/serve.go, flai/cmd/serve_actions.go, flai/internal/serve/serve.go, flai/internal/serve/upgrade.go, flai/internal/serve/upgrade_test.go, docs/operators/runbooks/migrate.md]
date: 2026-10-15
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0315 The host's installed flai stays behind the project's releases after I-0147 closed, so the MCP task_done still refuses

I-0149 was closed on 2026-10-15.

The cause:

```text
The host's installed flai stays behind the project's releases after I-0147 closed, so the MCP task_done still refuses
```

The fix:

```text
S-0391. The host ran flai 1.48.0, which predates ADR-0158's minute check (1.50.1) and ADR-0160's close through the installed flai (1.51.1), so nothing upgraded it or asked. The operator ran flai host upgrade once (TH-0461). At 2026-10-15T22:54Z, flai host, flai serve, and the MCP servers run flai 1.52.0, the newest release in the project's history, and MCP task_done closed T-1601 with no refusal. From now on ADR-0158 keeps the host current. With auto-upgrade off, its thread now names a release only once it is published (T-1599, ADR-0161, TestTheUpgradeThreadWaitsForTheReleaseToBePublished), so the command it gives always installs. The operator docs name the one manual upgrade that a host older than 1.50.1 needs (T-1600).
```

Source: I-0149.
