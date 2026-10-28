---
id: D-0088
title: Two flai serve processes with different configurations fight over one dashboard
kind: cause
sources: [I-0029, S-0141, S-0152, S-0214]
topics: [flai, cli, go, flaiover, dashboard, sveltekit, host, testing, code, all]
paths: [flai/cmd/dashboard_upgrade.go, flai/cmd/dashboard.go, flai/cmd/dashboard_watch.go, flai/cmd/host.go, flai/internal/host, flai/internal/channel, flai/internal/config, flaiover/src/lib/server/agent.ts, flaiover/src/lib/server/agent.test.ts, flaiover/Dockerfile, design/system/flaiover-dashboard.md, design/system/dashboard-host-channel.md, design/system/flai-cli.md, design/adrs/0062-flai-host-restarts-the-dashboard-container-flai-dashboard-recorded-when-it-is.md, design/adrs/0063-the-dashboard-refuses-a-second-flai-for-a-project-with-close-code-4409-while.md, design/adrs/README.md, docs/operators, docs/users/flai.md, docs/users/flai-reference.md, design/issues/summary.md, flai/cmd/dashboard_test.go, flai/cmd/dashboard_watch_test.go, design/issues/I-0025-a-pattern-kill-meant-for-a-scratch-server-stops-the-operator-s-dashboard-container-processes-match-pkill-on-the-host.md, design/issues/I-0029-two-flai-serve-processes-with-different-configurations-fight-over-one-dashboard.md]
date: 2026-10-01
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0088 Two flai serve processes with different configurations fight over one dashboard

I-0029 was closed on 2026-10-01.

The cause:

```text
Two flai serve processes with different configurations fight over one dashboard
```

The fix:

```text
S-0184: the dashboard refuses a second flai for a project with close code 4409 while the one it has answers a ping, replaces only one that does not, and flai waits a minute before it dials again (ADR-0063). One host per machine (ADR-0040) already refused a second host.
```

What stories do since, so that a measurement or a check never reaches the operator's host:

- S-0141 stubbed `/api/board` in the browser against a dev server with auth off, "instead of running a scratch flai host and serve, so nothing could reach the operator's host (I-0044)".
- S-0152 built its flai "to a scratch path" and left the operator's installed `flai serve` "not restarted or replaced".
- S-0214 fed a branch-built `flai stats` output to its checks, because "a second flai host on this machine would break the operator's".

Sources: I-0029, S-0141, S-0152, S-0214.
