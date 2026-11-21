---
id: D-0348
title: The Components page cannot list, start, or stop anything through the dashboard, because the channel's project key and the components methods' sub-project share the one params field project
kind: cause
sources: [I-0163]
topics: [cli, dashboard, flai, go, flaiover, sveltekit, code, all]
paths: [flai/internal/hostapi/components.go, flai/internal/hostapi/components_test.go, flai/internal/channel/channel.go, flai/internal/channel/channel_test.go, flaiover/src/lib/server/agent.ts, flaiover/src/lib/server/agent.test.ts, flaiover/src/lib/server/repo.ts, flaiover/src/routes/api/components, design/system/component-library.md, design/system/flaiover-dashboard.md, docs/users/flai-reference.md, flai/internal/hostapi/params_test.go, design/system/dashboard-host-channel.md]
date: 2026-10-21
status: compacted
absorbed_by: ""
asked_again: []
cited: []
---

# D-0348 The Components page cannot list, start, or stop anything through the dashboard, because the channel's project key and the components methods' sub-project share the one params field project
