---
id: D-0511
title: "S-0586 plan: six code and docs tasks plus an ADR in five layers, 3h, 1250 USD a week; three questions on what stays outside the dashboard: everything is good except I would prefer for #3 to change to…"
kind: recommendation
sources: [TH-0861]
topics: [kluster, deployment, manifest, dashboard, flai, cli, go, flaiover, sveltekit, code, all]
paths: [design/adrs, design/adrs/README.md, design/system/kluster.md, design/system/project-manifest.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, design/tech/kluster.md, flai/internal/kluster/init.go, flai/internal/kluster/init_test.go, flai/internal/kluster/kluster.go, flai/internal/kluster/kluster_test.go, flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, flai/internal/check/kluster_test.go, flai/internal/docedit/file.go, flai/internal/docedit/file_test.go, flai/internal/serve/kluster.go, flai/internal/serve/kluster_test.go, flai/internal/hostapi/kluster.go, flai/internal/hostapi/kluster_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/cmd/kluster.go, flai/cmd/kluster_test.go, flaiover/src/lib/server/agent.ts, flaiover/src/lib/server/repo.ts, flaiover/src/routes/api/kluster/workloads/+server.ts, flaiover/src/routes/api/kluster/kluster.test.ts, flaiover/src/routes/ops/klusters/+page.svelte, flaiover/src/routes/ops/klusters/klusters.svelte.test.ts, flaiover/src/routes/ops/deployments/+page.svelte, flaiover/src/routes/ops/deployments/deployments.svelte.test.ts, flaiover/src/lib/components/SettingsPanel.svelte, flaiover/src/lib/components/SettingsPanel.svelte.test.ts, docs/users/flaiover.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/index.md, docs/operators/settings.md, flai/internal/hostapi/settings.go, flai/internal/hostapi/settings_test.go, flaiover/src/routes/api/kluster/config/+server.ts, flaiover/src/routes/api/kluster/config/config.test.ts, flaiover/e2e/klusters.e2e.ts, flaiover/src/lib/settings.ts, design/issues/I-0261-a-story-s-prime-pack-came-in-22-parts-most-of-each-part-the-same-topic-sources.md, design/issues/summary.md, design/issues/I-0265-a-story-s-sub-agents-write-and-monitor-are-refused-every-path-outside-its.md]
date: 2026-11-07
status: active
absorbed_by: ""
asked_again: []
cited: [S-0586, T-2569, TH-0861, TH-0869]
---

# D-0511 S-0586 plan: six code and docs tasks plus an ADR in five layers, 3h, 1250 USD a week; three questions on what stays outside the dashboard: everything is good except I would prefer for #3 to change to…

The operator, alex, answered TH-0861 on S-0586 (the kluster first-time setup wizard) at 2026-11-07T09:51:27Z, taking the orchestrator's recommendation of 08:40:05Z except its answer 3.

What was settled, as the orchestrator summarised it on the thread at 09:52:54Z:

1. The `settings` host action stays the one step turned on in a shell (ADR-0039). A hosted flai turns it on at provisioning.
2. No Install button for `kluster` or `kubectl`. The Klusters page reports the missing tools and gives the install sentence, since kluster installs to `/usr/local/bin` and `kluster setup` installs `kubectl` with `sudo`; Docker stays a stated prerequisite.
3. The Deployments page starts and stops each service's port-forward, which flai serve runs bound to `0.0.0.0`, not `127.0.0.1`, because the operator means flai to run in the cloud. A forward on `0.0.0.0` has no auth of its own, so the docs say so.
4. The scaffold writes `profile: spire`, changeable on the Klusters page: kluster v1.3.1 has no empty cluster, and `spire` is its default.

The operator's answer:

```text
everything is good except I would prefer for #3 to change to binding to `0.0.0.0` (again, the future is running this in the cloud)
```

The recommendation it overruled in part:

```text
3. **Yes:** Start and Stop for each port-forward on the Deployments page, bound to `127.0.0.1`.
```

Source: TH-0861. T-2569 carries answer 3.
