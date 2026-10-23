---
id: D-0348
title: The Components page cannot list, start, or stop anything through the dashboard, because the channel's project key and the components methods' sub-project share the one params field project
kind: cause
sources: [I-0163]
topics: [cli, dashboard, flai, go, flaiover, sveltekit, code, all]
paths: [flai/internal/hostapi/components.go, flai/internal/hostapi/components_test.go, flai/internal/channel/channel.go, flai/internal/channel/channel_test.go, flaiover/src/lib/server/agent.ts, flaiover/src/lib/server/agent.test.ts, flaiover/src/lib/server/repo.ts, flaiover/src/routes/api/components, design/system/component-library.md, design/system/flaiover-dashboard.md, docs/users/flai-reference.md, flai/internal/hostapi/params_test.go, design/system/dashboard-host-channel.md]
date: 2026-10-21
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0348 The Components page cannot list, start, or stop anything through the dashboard, because the channel's project key and the components methods' sub-project share the one params field project

I-0163 was closed on 2026-10-21.

The cause:

```text
The Components page cannot list, start, or stop anything through the dashboard, because the channel's project key and the components methods' sub-project share the one params field project
```

The fix:

```text
Fixed by S-0448, accepted 2026-10-18: the components methods take their sub-project as sub_project (382fab56f), and the dashboard's agent puts the channel's project key after a call's params (d94c114d4). Released in flai v1.63.1 and flaiover v0.52.5. S-0475 confirmed it: flai hostapi components.status '{"project":"sf"}' now answers the status, and the regression tests TestComponentsAnswerThroughTheChannelWithTheProjectKeyInTheirParams, TestNoMethodTakesAParameterNamedProject, and the S-0448 cases in flaiover's agent.test.ts and components.test.ts pass.
```

Source: I-0163.
