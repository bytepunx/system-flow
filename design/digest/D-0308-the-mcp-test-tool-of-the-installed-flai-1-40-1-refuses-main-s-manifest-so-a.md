---
id: D-0308
title: The MCP test tool of the installed flai 1.40.1 refuses main's manifest, so a story's agent runs flai test from the tree
kind: cause
sources: [I-0141]
topics: [release, cli, mcp, host, flai, go, code, all]
paths: [design/adrs/0153-flai-serve-upgrades-the-host-s-flai-when-the-project-s-history-has-a-newer-published-release-behind-auto-upgrade.md, design/adrs/README.md, design/system/flai-cli.md, flai/internal/serve/upgrade.go, flai/internal/serve/upgrade_test.go, flai/internal/serve/serve.go, flai/internal/serve/outdated_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/internal/verify/manifest.go, flai/internal/verify/manifest_test.go, flai/internal/mcpserver/test.go, flai/internal/mcpserver/test_test.go, flai/internal/mcpserver/verify.go, flai/internal/mcpserver/verify_test.go, docs/operators/index.md, docs/operators/settings.md, docs/operators/runbooks/update.md, docs/users/flai.md, design/issues/I-0141-the-mcp-test-tool-of-the-installed-flai-1-40-1-refuses-main-s-manifest-so-a-story-s-agent-runs-flai-test-from-the-tree.md, design/issues/summary.md, design/adrs/0158-flai-serve-checks-each-served-project-for-a-newer-flai-release-every-minute.md, flai/cmd/test.go, flai/cmd/test_test.go, flai/cmd/verify.go, flai/cmd/verify_test.go, flai/cmd/serve.go]
date: 2026-10-15
status: absorbed
absorbed_by: "design/system/flai-cli.md#versions-the-hosts-flai-and-the-tree"
asked_again: []
cited: []
---

# D-0308 The MCP test tool of the installed flai 1.40.1 refuses main's manifest, so a story's agent runs flai test from the tree

I-0141 was closed on 2026-10-15.

The cause:

```text
The MCP test tool of the installed flai 1.40.1 refuses main's manifest, so a story's agent runs flai test from the tree
```

The fix:

```text
Fixed by S-0387 (ADR-0158). flai serve now checks every minute whether its flai is older than the newest flai release in each served project's history. With the shell-only auto-upgrade host action on, it has flai host install that release and restart on it. Otherwise it asks the operator on one thread per project and release, so the host no longer lags unnoticed. flai test, flai verify, and the MCP tools test and verify now say an outdated flai refused the tiers and name the upgrade, instead of asking for the manifest to be fixed. ADR-0142's task_done refusal stays: it holds only until the host catches up.
```

Source: I-0141.
