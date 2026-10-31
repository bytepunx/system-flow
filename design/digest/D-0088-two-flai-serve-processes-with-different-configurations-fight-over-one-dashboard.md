---
id: D-0088
title: Two flai serve processes with different configurations fight over one dashboard
kind: cause
sources: [I-0029, S-0141, S-0152, S-0214]
topics: [flai, cli, go, flaiover, dashboard, sveltekit, host, testing, code, all]
paths: [flai/cmd/dashboard_upgrade.go, flai/cmd/dashboard.go, flai/cmd/dashboard_watch.go, flai/cmd/host.go, flai/internal/host, flai/internal/channel, flai/internal/config, flaiover/src/lib/server/agent.ts, flaiover/src/lib/server/agent.test.ts, flaiover/Dockerfile, design/system/flaiover-dashboard.md, design/system/dashboard-host-channel.md, design/system/flai-cli.md, design/adrs/0062-flai-host-restarts-the-dashboard-container-flai-dashboard-recorded-when-it-is.md, design/adrs/0063-the-dashboard-refuses-a-second-flai-for-a-project-with-close-code-4409-while.md, design/adrs/README.md, docs/operators, docs/users/flai.md, docs/users/flai-reference.md, design/issues/summary.md, flai/cmd/dashboard_test.go, flai/cmd/dashboard_watch_test.go, design/issues/I-0025-a-pattern-kill-meant-for-a-scratch-server-stops-the-operator-s-dashboard-container-processes-match-pkill-on-the-host.md, design/issues/I-0029-two-flai-serve-processes-with-different-configurations-fight-over-one-dashboard.md]
date: 2026-10-01
status: compacted
absorbed_by: ""
asked_again: []
cited: []
---

# D-0088 Two flai serve processes with different configurations fight over one dashboard
