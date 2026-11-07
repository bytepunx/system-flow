---
id: S-0586
type: story
nature: improvement
title: Provide first time wizard for kluster setup
status: done
owner: alex
created: 2026-11-07T08:24:44Z
updated: 2026-11-07T11:14:08Z
transitions:
  - to: ready
    at: 2026-11-07T08:24:52Z
    by: alex
  - to: in-progress
    at: 2026-11-07T08:34:45Z
    by: agent-S-0586
  - to: review
    at: 2026-11-07T11:12:30Z
    by: agent-S-0586
  - to: done
    at: 2026-11-07T11:14:08Z
    by: orchestrator
blocked:
  - from: 2026-11-07T09:47:29Z
    until: 2026-11-07T09:53:17Z
    reason: awaiting the operator on TH-0869
tags: [flai, flaiover, kluster]
topics: [kluster, deployment, manifest, dashboard]
touches: [design/adrs, design/adrs/README.md, design/system/kluster.md, design/system/project-manifest.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, design/tech/kluster.md, flai/internal/kluster/init.go, flai/internal/kluster/init_test.go, flai/internal/kluster/kluster.go, flai/internal/kluster/kluster_test.go, flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, flai/internal/check/kluster_test.go, flai/internal/docedit/file.go, flai/internal/docedit/file_test.go, flai/internal/serve/kluster.go, flai/internal/serve/kluster_test.go, flai/internal/hostapi/kluster.go, flai/internal/hostapi/kluster_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/cmd/kluster.go, flai/cmd/kluster_test.go, flaiover/src/lib/server/agent.ts, flaiover/src/lib/server/repo.ts, flaiover/src/routes/api/kluster/workloads/+server.ts, flaiover/src/routes/api/kluster/kluster.test.ts, flaiover/src/routes/ops/klusters/+page.svelte, flaiover/src/routes/ops/klusters/klusters.svelte.test.ts, flaiover/src/routes/ops/deployments/+page.svelte, flaiover/src/routes/ops/deployments/deployments.svelte.test.ts, flaiover/src/lib/components/SettingsPanel.svelte, flaiover/src/lib/components/SettingsPanel.svelte.test.ts, docs/users/flaiover.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/index.md, docs/operators/settings.md, flai/internal/hostapi/settings.go, flai/internal/hostapi/settings_test.go, flaiover/src/routes/api/kluster/config/+server.ts, flaiover/src/routes/api/kluster/config/config.test.ts, flaiover/e2e/klusters.e2e.ts, flaiover/src/lib/settings.ts, design/issues/I-0261-a-story-s-prime-pack-came-in-22-parts-most-of-each-part-the-same-topic-sources.md, design/issues/summary.md, design/issues/I-0265-a-story-s-sub-agents-write-and-monitor-are-refused-every-path-outside-its.md, flai/internal/hostapi/kluster_forward.go, flai/internal/hostapi/kluster_forward_test.go, flaiover/src/routes/api/kluster/forward/+server.ts, flaiover/src/routes/api/kluster/forward/forward.test.ts, design/issues/I-0269-repo-test-ts-s-whole-monorepo-read-times-out-at-vitest-s-5s-default-under-load.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 9132
  priced_by: harness
  turns:
    - day: 2026-11-07
      ceremony: 1
      test_runs: 1
      hand_edits: 5
      work: 208
  models:
    - model: claude-haiku-5-5
      input: 12
      output: 3229
      cache_read: 143120
      cache_write: 52289
      cost: 0.0096
    - model: claude-opus-5-5
      input: 1006
      output: 301496
      cache_read: 151354606
      cache_write: 1925021
      cost: 47.2794
    - model: claude-sonnet-5-5
      input: 14
      output: 5409
      cache_read: 190214
      cache_write: 54057
      cost: 0.2083
  strategic:
    - kind: planner
      seconds: 1715
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 4
          output: 1882
          cache_read: 4110
          cache_write: 5741
          cost: 0.0017
        - model: claude-opus-5-5
          input: 310
          output: 60455
          cache_read: 42214092
          cache_write: 603311
          cost: 14.4796
    - kind: orchestrator
      seconds: 238
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 57
          output: 384
          cache_read: 1722411
          cache_write: 365732
          cost: 0.5198
cost_of_delay:
  inputs:
    revenue_per_week: 1250
    by: alex
    at: 2026-11-07T08:24:44Z
  value: 1250
  by: planner-S-0586
  at: 2026-11-07T08:33:18Z
forecast:
  duration: 3h40m
  delivery: 2026-11-07T12:00:00Z
  basis: "The planner's 3h for six tasks, plus about 40 minutes for T-2569, the port-forward task the operator's answer 3 on TH-0861 added; T-2560 to T-2564 took about 70 minutes of agent time"
  by: agent-S-0586
  at: 2026-11-07T09:54:58Z
---
# S-0586 Provide first time wizard for kluster setup

## Goal

When enabling kluster in a project's feature settings, the following error appears at the bottom of the features page:

```text
this project's system-flow.yaml has no kluster: block, so flai runs no cluster for it; add one naming the kluster config and the deploy command, such as kluster: {config: integration/kluster.yaml, deploy: [scripts/kluster-deploy.sh]}, as design/system/kluster.md § The kluster: block describes, then run this again
```

After a refresh, the `kluster` feature is no longer enabled.

This is confusing and goes against the goal of making the dashboard the single place the operator needs access to. If there are hosted versions of flai in the future, we shouldn't assume shell access.

When enabling kluster, provide a scaffold for system-flow.yaml that is either functionally a no-op, or that starts with the empty kubernetes cluster.

The Klusters page should allow the user to manage all features completely without the need for CLI commands or file edits.

## Acceptance criteria
- [x] When an operator enables the kluster feature, provide a stub/no-op/defaults only block in system-flow.yaml
- [x] The Klusters page should have everything the user needs to set up, manage, observe, and control the kluster for that project
- [x] There are no operations, configuration, or commands the operator needs to issue outside of the dashboard to make full use of this feature

## Tasks
- T-2559 An ADR and kluster.md say what enabling kluster writes and what a block with no deploy means
- T-2560 flai serve enable kluster writes and commits a block and a config for a project with none, and flai manifest set changes the block's keys
- T-2561 flai serve's kluster pass leaves a project whose block names no deploy alone
- T-2562 kluster.status says what setup still lacks, and kluster.workloads lists the cluster's nodes and pods
- T-2563 The Klusters page leads an operator through first-time setup and then manages the cluster, with no shell or file edit
- T-2564 The Deployments page and the Features toggle say a block's state without naming a command, and a refused toggle shows its real state at once
- T-2565 The docs say a project's kluster is set up from the Klusters page
- T-2569 The Deployments page starts and stops each service's port-forward, which flai serve runs on the host bound to 0.0.0.0
- T-2585 Fix what the review of S-0586's diff found: a second port-forward start waits until kubectl forwards, a broken sentence on the Klusters page, and a misplaced bullet in the operators guide

## Notes

### Criterion 3: the two steps left on the host

The operator settled on TH-0861 (2026-11-07T09:51:27Z, D-0511) that two steps stay outside the dashboard, and `docs/operators/index.md` § The kluster host action (S-0577) names each with its reason:

| Step | Why it stays |
|------|--------------|
| Turn the `settings` host action on in a shell | ADR-0039 stands (answer 1): a dashboard switch would make the dashboard token worth a shell on the host |
| Install Docker, `kluster`, and `kubectl` | No Install button (answer 2): kluster installs to `/usr/local/bin`, `kluster setup` installs `kubectl` with `sudo`, and Docker needs root and a daemon; the Klusters page shows what is missing and how to install it |

Criterion 3 is ticked as met with these two exceptions. Everything else is on the Klusters and Deployments pages: turning kluster on writes the block, the profile, the config file, the block's keys, Up and Down, nodes and pods, deploys, and port-forwards. ADR-0259 records the decisions.

### Planning

Touches:

| Paths | Came from |
|-------|-----------|
| `flai/internal/hostapi/kluster.go`, `flai/internal/kluster/kluster.go`, `flaiover/src/routes/ops/klusters/+page.svelte`, `design/system/kluster.md` | Layout: the code behind the refusal the goal quotes and the page it names, read in the tree. `flai touches suggest` needs a starting path, since the story declared none, and these were the ones given it |
| `docs/users/flaiover.md`, `design/system/flai-cli.md`, `design/system/flaiover-dashboard.md`, `docs/operators/index.md`, `docs/users/flai-reference.md`, `docs/users/flai.md`, `flai/internal/hostapi/kluster_test.go`, `flaiover/src/lib/server/agent.ts` | Co-change: `flai touches suggest`, each 17% to 25% of 12 commits |
| `flai/internal/kluster/deploy.go` | Co-change, left out: its refusal for no `deploy` already says what this story needs |
| `flai/cmd/serve_actions.go` and its test | Layout: the refusal itself, `kluster.ErrNoBlock` returned from `flai serve enable` |
| `flai/internal/manifest/manifest.go`, `settings.go`, their tests, `flai/internal/check/kluster_test.go` | Design: `KlusterErrors` and the setting catalog `flai manifest set` writes |
| `flai/internal/kluster/init.go`, `init_test.go` | Layout: the new scaffold |
| `flai/internal/docedit/file.go` and its test | Design: ADR-0147's reach of `file.save` |
| `flai/internal/serve/kluster.go` and its test | Design: `design/system/kluster.md` § flai serve (S-0579) |
| `flai/cmd/kluster.go` and its test | Layout: the commands the `ErrNoBlock` sentence reaches |
| `flaiover/src/lib/server/repo.ts`, the `workloads` route, `kluster.test.ts`, the Deployments page, `SettingsPanel.svelte`, and their tests | Layout |
| `docs/operators/settings.md`, `design/system/project-manifest.md`, `design/tech/kluster.md` | Design: each describes the block by hand today |

`design/adrs` stays a folder touch because the story adds one ADR whose file name carries a number and a slug `flai adr new` gives it, so no task can name that file yet. `design/adrs/README.md` is the index flai generates from it.

Forecast: 3h. `flai forecast` gave 4m, from a median of 80 s per unit of size times size 3, counting 3 criteria and 0 touches. That band is small improvements, and this story is not one. It is adjusted to the agent time of the kluster stories it most resembles: S-0577 39m, for the manifest and commands; S-0580 51m and S-0581 38m, for the host methods and the pages; and part of S-0579's 51m, for the serve pass; about 2h50m with the docs. The plan thread's questions can add to it, and the figure is revised once they are answered. Delivery allows about an hour's wait behind a full in-progress column and S-0560.

Cost of delay: 1250 USD a week. It stands as `flai cod` works it out from the operator's revenue input of 1250 a week, with no penalty or time lost given.

No `after`: every kluster story this builds on, S-0577 to S-0581, is accepted.

### Accepted by the orchestrator

- Verified: 609704236e96c942b885c3c8ab58d9afa2578999
- At: 2026-11-07T11:14:08Z

Verdict: pass — all three criteria are met with the diff and the docs; the worktree is clean at 609704236e96c942b885c3c8ab58d9afa2578999 and the close-out passed.
- 1: flai/internal/kluster/init.go, flai/internal/kluster/kluster.go, flai/cmd/serve_actions.go, flai/internal/manifest/manifest.go, flai/internal/manifest/settings.go, flai/internal/serve/kluster.go, flai/internal/docedit/file.go, design/adrs/0259-enabling-kluster-on-a-project-with-no-kluster-block-writes-one-naming-only-its.md, design/system/kluster.md, design/system/project-manifest.md, flaiover/src/lib/components/SettingsPanel.svelte
- 2: flaiover/src/routes/ops/klusters/+page.svelte, flaiover/src/routes/ops/deployments/+page.svelte, flai/internal/hostapi/kluster.go, flai/internal/hostapi/kluster_forward.go, flaiover/src/lib/server/repo.ts, flaiover/src/lib/server/agent.ts, flaiover/src/routes/api/kluster/config/+server.ts, flaiover/src/routes/api/kluster/forward/+server.ts, flaiover/src/routes/api/kluster/workloads/+server.ts, flaiover/e2e/klusters.e2e.ts, docs/users/flaiover.md
- 3: docs/operators/index.md, docs/users/flaiover.md, docs/operators/settings.md, docs/users/flai.md, design/system/kluster.md, flai/internal/hostapi/kluster.go, flai/internal/hostapi/kluster_forward.go
