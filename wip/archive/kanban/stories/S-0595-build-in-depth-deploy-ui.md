---
id: S-0595
type: story
nature: feature
title: Build in-depth Deploy UI
status: done
owner: alex
created: 2026-11-07T22:00:05Z
updated: 2026-11-07T23:06:27Z
transitions:
  - to: ready
    at: 2026-11-07T22:00:38Z
    by: alex
  - to: in-progress
    at: 2026-11-07T22:01:20Z
    by: agent-S-0595
  - to: review
    at: 2026-11-07T23:04:40Z
    by: agent-S-0595
  - to: done
    at: 2026-11-07T23:06:27Z
    by: orchestrator
tags: [dashboard, cli]
topics: [kluster, deployment, manifest]
touches: [flaiover/src, design/adrs/0263-a-kluster-block-may-list-components-each-built-as-its-own-step-with-its-own-log.md, design/adrs/README.md, flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, flai/internal/check/kluster_test.go, flai/internal/kluster/deploy.go, flai/internal/kluster/deploy_test.go, flai/internal/hostapi/kluster.go, flai/internal/hostapi/kluster_test.go, design/components/flaiover/src/lib/components/LogFeed.examples.json, design/components/flaiover/src/lib/components/DeploySettings.examples.json, design/system/kluster.md, design/system/flaiover-dashboard.md, design/system/project-manifest.md, docs/users/flaiover.md, docs/users/flai.md, docs/operators/index.md, docs/operators/settings.md, flai/internal/hostapi/settings.go, flai/internal/hostapi/settings_test.go, docs/users/flai-reference.md, design/issues/I-0308-flai-test-on-a-change-to-the-manifest-s-settings-catalog-does-not-run.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 3820
  priced_by: harness
  turns:
    - day: 2026-11-07
      ceremony: 4
      hand_edits: 2
      work: 92
  models:
    - model: claude-haiku-5-5
      input: 48
      output: 26415
      cache_read: 1301192
      cache_write: 173117
      cost: 0.0479
    - model: claude-opus-5-5
      input: 572
      output: 200946
      cache_read: 98223108
      cache_write: 775000
      cost: 28.7173
    - model: claude-sonnet-5-5
      input: 18
      output: 4846
      cache_read: 261666
      cache_write: 52011
      cost: 0.2047
  strategic:
    - kind: planner
      seconds: 198
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 82
          output: 429
          cache_read: 4398231
          cache_write: 168790
          cost: 0.3012
        - model: claude-opus-5-5
          input: 31
          output: 8123
          cache_read: 1454208
          cache_write: 109451
          cost: 1.0346
    - kind: orchestrator
      seconds: 143
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 27
          output: 237
          cache_read: 924302
          cache_write: 168015
          cost: 0.2773
        - model: claude-sonnet-5-5
          input: 6
          output: 47
          cache_read: 50522
          cache_write: 24704
          cost: 0.0548
cost_of_delay:
  inputs:
    revenue_per_week: 1250
    by: alex
    at: 2026-11-07T22:00:05Z
  value: 1250
  by: planner-S-0595
  at: 2026-11-07T22:12:42Z
forecast:
  duration: 7h
  delivery: 2026-11-08T05:00:00Z
  basis: "Raised from flai forecast's 11m (median 109 s per unit times size 6): the story's 9 tasks in 5 layers span the manifest schema, the kluster runner, the host API, flaiover's server, two components, two pages, and the docs, about 45m each."
  by: planner-S-0595
  at: 2026-11-07T22:12:42Z
---
# S-0595 Build in-depth Deploy UI

## Goal

Under Ops/Deployments, there is only a single button and when pressed, a new project shows the error:

```text
the kluster: block of system-flow.yaml names no deploy command; write kluster.deploy, the argument list that deploys into a namespace, and run flai check
```

The deployment screen should have the controls necessary to set up build and deploy machinery for a project's services and configuration.

## Acceptance criteria
- [x] The operator never has to edit system-flow.yaml directly to control or setup deployments
- [x] The UI/UX allow the operator to create or make any changes needed to control their projects deployments to the kluster managed kubernetes cluster
- [x] The Deployment page should use the vertical nav rail in page to group controls, livestream logs, and other status indicators for a given deployment
- [x] a Build sub-page should show both the status and an expandable live feed of each component's build. If there are 3 Docker images (for example), there should be 3 different component build panels with status indicators so that the operator can see which steps have completed, are still running, or have failed.

## Tasks
- T-2598 The kluster block takes components, checked by flai check and written by flai manifest set as a components setting
- T-2599 The host read kluster.log answers a step's log from an offset, and the dashboard requires it
- T-2600 A deploy builds each component as its own step and records every step as pending before the first starts
- T-2601 flaiover's server asks kluster.log through /api/kluster/log, and its types know pending steps and the components setting
- T-2602 LogFeed shows a step's status and follows its log live through /api/kluster/log
- T-2603 DeploySettings edits the kluster block's namespace, components, commands, and timeout, and asks an agent to write the scripts
- T-2604 The Deployments page groups Overview, Build, Logs, Access, and Setup in an in-page rail, with a panel per component build
- T-2605 The Klusters page's step 4 is the DeploySettings form
- T-2606 The kluster design, the dashboard design, and the user and operator guides describe components, pending steps, kluster.log, and the Deployments page's rail

## Notes

### Planning

The story's agent wrote ADR-0263 and the nine tasks after it started the story. This plan keeps them as written. It adds tags to each task, the topics `kluster`, `deployment`, and `manifest` to the story, and the touches and figures below.

Touches:

| Touch | From | Why |
|-------|------|-----|
| `flaiover/src` | declared | A folder touch, kept. The page panes may need routes or helpers no task names yet. The tasks' flaiover files replace it in the claim. |
| `design/adrs/0263-…-with-its-own-log.md`, `design/adrs/README.md` | declared | ADR-0263, written on the story branch. |
| `flai/internal/manifest/manifest.go`, `settings.go` and their tests, `flai/internal/check/kluster_test.go`, `docs/operators/settings.md` | declared by T-2598 | The `components` key, its checks, and its setting. |
| `flai/internal/kluster/deploy.go`, `deploy_test.go` | declared by T-2600 | One build step per component, and the pending plan in the record. |
| `flai/internal/hostapi/kluster.go`, `kluster_test.go` | declared by T-2599 | The `kluster.log` read. |
| `design/components/…/LogFeed.examples.json`, `DeploySettings.examples.json` | declared by T-2602, T-2603 | The new components' examples. |
| `design/system/kluster.md`, `flaiover-dashboard.md`, `project-manifest.md`, `docs/users/flaiover.md`, `docs/users/flai.md`, `docs/operators/index.md` | declared by T-2606 | The design and the guides. |

`flai touches suggest` listed `docs/users/flai-reference.md`, `design/system/flai-cli.md`, and `flai/internal/hostapi/writes.go`, from co-change. None is taken. No command or flag changes, so the generated reference does not. `kluster.log` registers in `klusterSpecs()` in `kluster.go`, not in `writes.go`. The declared `flai/cmd` was dropped by the story's agent: `flai kluster status` prints each step's name and outcome as it stands, so `build-<name>` and `pending` print with no change.

Figures:

- Forecast: 7h, raised from flai forecast's 11m. Its size of 6 counts only 4 criteria and 2 touches. The story has 9 tasks in 5 layers across Go and the dashboard, about 45m each. Delivery counts from the start at 2026-11-07T22:02:00Z.
- Cost of delay: 1250 USD a week, as `flai cod` gives it from the operator's revenue input. No penalty or time lost was given, so the figure stands.

No `after`: the story reads no output of an open story. It builds on S-0580 and S-0581, which are done.

### Accepted by the orchestrator

- Verified: f39e371e083ea07aba5e01bb0e9e2181f708dac6
- At: 2026-11-07T23:06:27Z

Verdict: pass at f39e371e083ea07aba5e01bb0e9e2181f708dac6
- 1: flaiover/src/lib/components/DeploySettings.svelte, flaiover/src/lib/settings.ts, flaiover/src/lib/server/repo.ts, flai/internal/manifest/manifest.go, flai/internal/manifest/settings.go, flai/internal/hostapi/settings.go, flaiover/src/routes/ops/deployments/+page.svelte
- 2: flaiover/src/lib/components/DeploySettings.svelte, flai/internal/kluster/deploy.go, flai/internal/kluster/deploy_test.go, flaiover/src/routes/ops/klusters/+page.svelte, flaiover/src/routes/ops/deployments/+page.svelte
- 3: flaiover/src/routes/ops/deployments/+page.svelte, flaiover/src/lib/components/LogFeed.svelte, flaiover/src/routes/api/kluster/log/+server.ts, flai/internal/hostapi/kluster.go
- 4: flaiover/src/routes/ops/deployments/+page.svelte, flaiover/src/lib/components/LogFeed.svelte, flai/internal/kluster/deploy.go, flaiover/src/routes/ops/deployments/deployments.svelte.test.ts, flaiover/src/lib/components/LogFeed.svelte.test.ts
