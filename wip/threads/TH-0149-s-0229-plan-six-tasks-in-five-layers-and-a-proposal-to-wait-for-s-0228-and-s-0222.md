---
id: TH-0149
title: "S-0229 plan: six tasks in five layers, and a proposal to wait for S-0228 and S-0222"
anchor:
  path: wip/kanban/stories/S-0229-the-strategic-agents-settings-are-edited-in-the-dashboard-permissions-policy-release-policy-schedules-and-agents.md
  item: S-0229
status: resolved
participants: [planner-S-0229, alex]
created: 2026-10-05T05:48:20Z
updated: 2026-10-05T05:51:05Z
---

# TH-0149 S-0229 plan: six tasks in five layers, and a proposal to wait for S-0228 and S-0222

On wip/kanban/stories/S-0229-the-strategic-agents-settings-are-edited-in-the-dashboard-permissions-policy-release-policy-schedules-and-agents.md.

## Entries

### 2026-10-05T05:48:20Z planner-S-0229
Recommended: add S-0228 and S-0222 to S-0229's `after`. S-0228 makes the orchestrator and analyzer pages the panels go on, and S-0222 adds `orchestration.release.whole_epics`, which the first criterion names. Neither is in `after` today, directly or through S-0218. Both are ahead of S-0229 in the pull order, so the forecast does not move. I have not changed `after`: the story is finalized, so it is your call.

Tasks, by layer:

| Layer | Task | Waits for | Changes |
|-------|------|-----------|---------|
| 1 | T-0952 `flai manifest set` writes the strategic keys, validated, refusing with field and reason | nothing | `flai/internal/manifest/settings.go`, `flai/cmd/manifest.go`, `root.go`, the reference |
| 1 | T-0954 the guard holds a running orchestrator to changed permissions on its next call; its prompt says they may change | nothing (S-0218 is in the story's `after`) | `flai/cmd/guard_permissions_change_test.go`, `flai/internal/harness` |
| 2 | T-0961 hostapi `settings.manifest` runs the command; `settings.get` gives each key's value, default, sentence, and risk | T-0952 | `flai/internal/hostapi/settings.go`, `flai/cmd/serve_actions.go`, `REQUIRED_METHODS` |
| 3 | T-0963 a shared `StrategicSettings` panel: save, refusal on its field, read-only with the reason | T-0961 | `flaiover/src/lib/settings.ts`, `/api/settings`, the component |
| 4 | T-0964 the panel on the orchestrator, planner, and analyzer pages; the Settings page links to the Planner page | T-0963 | `flaiover/src/routes/workflow/*`, `SettingsPanel.svelte` |
| 5 | T-0966 an ADR refining ADR-0039, and the design and guides | T-0954, T-0964 | `design/adrs`, `design/system`, `docs` |

Assumptions:

1. ADR-0039 has every dashboard write run a flai command, so I planned a new `flai manifest set`. The name is my guess, and T-0952 lets the story's agent pick another.
2. `planning.currency` stays out of the dashboard, because changing it re-denominates every amount. The criterion does not list it.
3. Each setting's sentence, and each permission's risk, live in flai and reach the dashboard through `settings.get`. The command, `flai check`, and the panel then read one list, and the panel offers the values flai validates. That also covers the difference between S-0218's `throughput|cost_of_delay` and S-0217's `cod|wsjf|throughput|fifo`.
4. "Picks up changed permissions on its next decision" means the guard reads the manifest on each call, which S-0218's T-0887 already does, plus a sentence in the orchestrator's prompt. A save does not restart the orchestrator.
5. A save commits `system-flow.yaml` with the dashboard's trailer, as `settings.default_agent` does.
6. Widening what the `settings` action lets a dashboard token change needs an ADR. With it on, a token can grant the orchestrator `accept_reviews` and `publish`.

Figures: forecast 1h15m (flai's 1h raised to match S-0204 and S-0201); cost of delay 84.52 USD a week (its share of E-0016's inputs). The reasons are under `### Planning` in the story's Notes. I propose no task split, merge, or drop.

### 2026-10-05T05:51:05Z alex
Resolved.
