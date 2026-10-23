---
id: D-0079
title: golangci-lint on the host is v1 but the config is v2
kind: cause
sources: [I-0001]
topics: [flai, cli, go, flaiover, dashboard, sveltekit, code, all]
paths: [flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, flai/cmd/serve.go, flai/cmd/serve_test.go, flai/cmd/dashboard_token.go, flai/internal/host/host.go, flai/internal/host/host_test.go, scripts/env.sh, scripts/flai-test.sh, scripts/install-tools.sh, scripts/flaiover-dev.sh, flaiover/compose.yaml, docs/operators/index.md, docs/operators/settings.md, docs/users/flai-reference.md, design/system/flai-cli.md, design/system/devex.md, docs/contributors/index.md, design/issues/I-0001-golangci-lint-version-mismatch.md, design/issues/I-0039-scripts-flaiover-dev-sh-points-the-dev-server-at-a-dashboard-token-file-that-flai-no-longer-writes.md, design/issues/I-0044-agents-flai-serve-starts-inherit-the-host-s-address-and-token-so-a-flai-serve-an-agent-runs-by-hand-takes-over-the-operator-s-mcp-servers.md, design/issues/summary.md]
date: 2026-10-01
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0079 golangci-lint on the host is v1 but the config is v2

I-0001 was closed on 2026-10-01.

The cause:

```text
The repository lints with golangci-lint v2 (`flai/.golangci.yaml`, `version: "2"`). The host has v1.64 in `~/go/bin`, which cannot read the config. Every session that lints has to install v2 into a scratch location first, and the scratch copy disappears between sessions.
```

The fix:

```text
S-0183: scripts/env.sh puts the main checkout's bin/ on PATH after the worktree's own, scripts/install-tools.sh installs the pinned golangci-lint v2 there from any tree, and scripts/flai-test.sh names the golangci-lint it runs and stops on anything but v2.
```

Source: I-0001.
