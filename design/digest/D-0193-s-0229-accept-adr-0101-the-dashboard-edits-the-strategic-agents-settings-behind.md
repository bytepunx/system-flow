---
id: D-0193
title: "S-0229: accept ADR-0101, the dashboard edits the strategic agents' settings behind the settings action: accept"
kind: answer
sources: [TH-0203]
topics: [orchestration, planning, analysis, flai, cli, go, flaiover, dashboard, sveltekit, template, conventions, code, all]
paths: [flaiover/src/routes, flai/internal/hostapi, flai/internal/manifest, design/system/flaiover-dashboard.md, docs/users/flaiover.md, flaiover/src/routes/settings, flaiover/src/lib/components, design/system/project-manifest.md, docs/operators/settings.md, flai/cmd/manifest.go, flai/cmd/manifest_test.go, flai/cmd/root.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/cmd/guard_permissions_change_test.go, flai/internal/guard, flai/internal/harness, flaiover/src/lib/settings.ts, flaiover/src/lib/server/agent.ts, design/system/strategic-agents.md, design/system/flai-cli.md, design/adrs, docs/operators/index.md, docs/users/flai.md, docs/users/flai-reference.md, design/conventions/strategic-agents.md, design/issues/I-0063-flai-adr-new-numbers-from-the-story-s-worktree-only-so-parallel-story-branches-take-the-same-adr-number.md, design/issues/summary.md, flaiover/src/lib/sitemenu.test.ts, flaiover/src/lib/sitemenu.ts, template/CHANGELOG.md, template/root/design/conventions/strategic-agents.md, template/template.yaml, design/issues/I-0073-flai-check-finds-threads-archived-outside-the-story-at-close-out.md, design/issues/I-0078-flai-check-finds-item-archive-outside-the-story-at-close-out.md]
date: 2026-10-06
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0193 S-0229: accept ADR-0101, the dashboard edits the strategic agents' settings behind the settings action: accept

The operator, alex, answered TH-0203 at 2026-10-06T22:51:48Z.

What was asked:

```text
ADR-0101 (proposed, on story/S-0229: `design/adrs/0101-with-the-settings-host-action-on-the-dashboard-edits-the-strategic-agents.md`) refines ADR-0039: with the `settings` host action on for a project, its dashboard may also change the manifest's `orchestration`, `planning` (not `currency`), and `analysis` keys through `flai manifest set`, committed with the dashboard's trailer. The consequence to weigh: a dashboard token can then grant the orchestrator `accept_reviews` and `publish`. Each permission's risk is shown beside it, and only the shell turns `settings` on or off. This is assumption…
```

The answer:

```text
accept
```

Source: TH-0203.
