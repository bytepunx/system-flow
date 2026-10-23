---
id: D-0023
title: "S-0121 conflicts with S-0120: accept S-0120 first, then I rebase S-0121: ok"
kind: answer
sources: [TH-0015]
topics: [flai, cli, go, flaiover, dashboard, sveltekit, code, all]
paths: [flai/cmd/serve_project.go, flai/cmd/serve.go, flai/cmd/dashboard_agent.go, flai/internal/serve/serve.go, flai/internal/serve/projects.go, flai/internal/channel/channeltest, flaiover/src/lib/server/agent.ts, docs/users/flai.md, docs/users/flai-reference.md, docs/operators, design/system/flai-cli.md, design/system/dashboard-host-channel.md, design/system/flaiover-dashboard.md]
date: 2026-09-26
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0023 S-0121 conflicts with S-0120: accept S-0120 first, then I rebase S-0121: ok

The operator, alex, answered TH-0015 at 2026-09-26T05:59:45Z.

What was asked:

```text
S-0121 is done: every criterion is ticked and verified, and flai-test and flaiover-test pass. But `story/S-0121` and `story/S-0120` (now in review) conflict in five files: `flai/internal/serve/serve.go`, `flai/cmd/serve.go`, `flai/cmd/dashboard_agent.go`, `docs/users/flai.md`, and `design/system/flai-cli.md`. Whichever is accepted second will not rebase cleanly. S-0120 also added `Status.ImportProjects` and `Status.Unserved` for the projects below the import folders, which `flai serve project list` here lists as "not served" from its own scan. After S-0120 lands, list should take them from…
```

The answer:

```text
ok
```

Source: TH-0015.
