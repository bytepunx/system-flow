---
id: D-0057
title: Go toolchain is not on PATH in the agent shell
kind: cause
sources: [I-0002]
topics: [flai, cli, go, flaiover, dashboard, sveltekit, code, all]
paths: [flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, flai/cmd/serve.go, flai/cmd/serve_test.go, flai/cmd/dashboard_token.go, flai/internal/host/host.go, flai/internal/host/host_test.go, scripts/env.sh, scripts/flai-test.sh, scripts/install-tools.sh, scripts/flaiover-dev.sh, flaiover/compose.yaml, docs/operators/index.md, docs/operators/settings.md, docs/users/flai-reference.md, design/system/flai-cli.md, design/system/devex.md, docs/contributors/index.md, design/issues/I-0001-golangci-lint-version-mismatch.md, design/issues/I-0039-scripts-flaiover-dev-sh-points-the-dev-server-at-a-dashboard-token-file-that-flai-no-longer-writes.md, design/issues/I-0044-agents-flai-serve-starts-inherit-the-host-s-address-and-token-so-a-flai-serve-an-agent-runs-by-hand-takes-over-the-operator-s-mcp-servers.md, design/issues/summary.md]
date: 2026-10-01
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0057 Go toolchain is not on PATH in the agent shell

I-0002 was closed on 2026-10-01.

The cause:

```text
`go` lives in `/usr/local/go/bin`, which is not on the PATH the agent shell inherits. Every Go command fails with "command not found" until PATH is prefixed by hand.
```

The fix:

```text
No longer applies: /usr/local/go is gone and go 1.26.4 is /opt/homebrew/bin/go on the normal PATH (triage 2026-10-01). The dead PATH entry in scripts/env.sh goes with S-0183.
```

Source: I-0002.
