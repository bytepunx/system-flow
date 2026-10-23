---
id: D-0022
title: "flai serve project remove cannot take a project below an import folder off the dashboard: add a removed list: A"
kind: answer
sources: [TH-0014]
topics: [flai, cli, go, flaiover, dashboard, sveltekit, code, all]
paths: [flai/cmd/serve_project.go, flai/cmd/serve.go, flai/cmd/dashboard_agent.go, flai/internal/serve/serve.go, flai/internal/serve/projects.go, flai/internal/channel/channeltest, flaiover/src/lib/server/agent.ts, docs/users/flai.md, docs/users/flai-reference.md, docs/operators, design/system/flai-cli.md, design/system/dashboard-host-channel.md, design/system/flaiover-dashboard.md]
date: 2026-09-26
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0022 flai serve project remove cannot take a project below an import folder off the dashboard: add a removed list: A

The operator, alex, answered TH-0014 at 2026-09-26T05:58:17Z.

What was asked:

```text
Since S-0120, `flai serve` serves every project below a folder named for import, registered or not. So `flai serve project remove blog` for a project under `~/git` (an import folder) unregisters it, but `flai serve` goes on serving it from the folder, and the switcher keeps it. Today `remove` says exactly that and names `flai serve import remove ~/git`, which also stops offering every other repository there. This question blocks moving S-0121 to review.
```

The answer:

```text
A
```

Source: TH-0014.
