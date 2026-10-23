---
id: D-0080
title: scripts/flaiover-dev.sh points the dev server at a dashboard token file that flai no longer writes
kind: cause
sources: [I-0039]
topics: [flai, cli, go, flaiover, dashboard, sveltekit, code, all]
paths: [flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, flai/cmd/serve.go, flai/cmd/serve_test.go, flai/cmd/dashboard_token.go, flai/internal/host/host.go, flai/internal/host/host_test.go, scripts/env.sh, scripts/flai-test.sh, scripts/install-tools.sh, scripts/flaiover-dev.sh, flaiover/compose.yaml, docs/operators/index.md, docs/operators/settings.md, docs/users/flai-reference.md, design/system/flai-cli.md, design/system/devex.md, docs/contributors/index.md, design/issues/I-0001-golangci-lint-version-mismatch.md, design/issues/I-0039-scripts-flaiover-dev-sh-points-the-dev-server-at-a-dashboard-token-file-that-flai-no-longer-writes.md, design/issues/I-0044-agents-flai-serve-starts-inherit-the-host-s-address-and-token-so-a-flai-serve-an-agent-runs-by-hand-takes-over-the-operator-s-mcp-servers.md, design/issues/summary.md]
date: 2026-10-01
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0080 scripts/flaiover-dev.sh points the dev server at a dashboard token file that flai no longer writes

I-0039 was closed on 2026-10-01.

The cause:

```text
scripts/flaiover-dev.sh points the dev server at a dashboard token file that flai no longer writes
```

The fix:

```text
S-0183: scripts/flaiover-dev.sh takes the token file flai dashboard token --json names; flaiover/compose.yaml mounts ~/.flai/serve/dashboard.token (or DASHBOARD_TOKEN_FILE) and no longer mounts the repository or sets PROJECT_DIR.
```

Source: I-0039.
