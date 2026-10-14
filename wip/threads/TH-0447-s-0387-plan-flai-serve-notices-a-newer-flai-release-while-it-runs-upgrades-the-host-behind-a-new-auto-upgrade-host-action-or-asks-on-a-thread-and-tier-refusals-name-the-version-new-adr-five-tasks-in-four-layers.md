---
id: TH-0447
title: "S-0387 plan: flai serve notices a newer flai release while it runs, upgrades the host behind a new auto-upgrade host action or asks on a thread, and tier refusals name the version (new ADR), five tasks in four layers"
anchor:
  path: wip/kanban/stories/S-0387-the-mcp-test-tool-of-the-installed-flai-1-40-1-refuses-main-s-manifest-so-a-story-s-agent-runs-flai-test-from-the-tree.md
  item: S-0387
status: open
participants: [planner-S-0387, orchestrator]
created: 2026-10-14T20:55:26Z
updated: 2026-10-14T20:55:43Z
---

# TH-0447 S-0387 plan: flai serve notices a newer flai release while it runs, upgrades the host behind a new auto-upgrade host action or asks on a thread, and tier refusals name the version (new ADR), five tasks in four layers

On wip/kanban/stories/S-0387-the-mcp-test-tool-of-the-installed-flai-1-40-1-refuses-main-s-manifest-so-a-story-s-agent-runs-flai-test-from-the-tree.md.

## Entries

### 2026-10-14T20:55:26Z planner-S-0387
**Plan:** S-0387 gets five tasks in four layers. The remedy needs a new ADR and a new host action, `auto-upgrade`, so it needs your approval before S-0387 is finalized.

## The cause

I-0141's two instances have one cause: the host's installed flai stays behind main after a flai release is published.

- 1.40.1 did not know `{importers}`, which S-0364 added to the manifest's tests.
- 1.48.0 was a minor version behind 1.49.1, so ADR-0142 refused `task_done`.
- Nothing upgrades the host. `flai serve` warns about an outdated flai only when it starts serving a project (`flai/internal/serve/serve.go`). A host that keeps running through a release never warns again.

## Proposed remedy

- **A periodic check.** `flai serve` checks `release.FlaiOutdated` for each project on its reconcile loop, at most once a minute.
- **With `auto-upgrade` on.** This new host action can be turned on only in a shell, like `auto-approve`. When it is on and the newer release is published and verifiable, `flai serve` asks the host to upgrade itself, as `flai host upgrade` does. Story agents are detached from the host, so they keep running. Their next session's `flai mcp` runs the new flai.
- **With it off.** `flai serve` opens one thread asking you to run `flai host upgrade`, as ADR-0108 asks when it reaches its restart limit.
- **Clearer refusals.** `flai test` and the MCP tools `test` and `verify` can refuse a manifest tier while the running flai is behind. They will then say the flai is older than the project and name the upgrade. Today they tell the agent to fix the manifest.

## Tasks

| Layer | Task | After | Touches |
|-------|------|-------|---------|
| 1 | T-1567 Record the remedy in a new ADR and in `flai-cli.md` § Versions | none | the ADR (predicted path), `design/adrs/README.md`, `design/system/flai-cli.md` |
| 2 | T-1568 The serve loop check, `auto-upgrade`, and the thread | T-1567 | `flai/internal/serve/upgrade.go` and its test, `serve.go`, `outdated_test.go`, `flai/cmd/serve_actions.go` and its test, `flai/internal/hostapi/writes.go` and its test |
| 2 | T-1569 Version-aware tier refusals | T-1567 | `flai/internal/verify/manifest.go`, `flai/internal/mcpserver/test.go` and `verify.go`, each with its test |
| 3 | T-1570 Operator and user docs | T-1568, T-1569 | `docs/operators/index.md`, `settings.md`, `runbooks/update.md`, `docs/users/flai.md` |
| 4 | T-1571 Close I-0141 | T-1570 | the I-0141 file, `design/issues/summary.md` |

T-1568 and T-1569 share no path, so they run together. Every touch names a file. No folder touch is kept.

## Figures

- **Forecast: 1h15m.** I raised it from flai's 34m. The plan has five tasks: an ADR, a host action, and a serve loop that drives the host's upgrade. S-0372 and S-0374 were raised to 1h for similar plans. Delivery is 2026-10-15T00:08:00Z.
- **Cost of delay: 10 USD a week.** This is flai cod's figure from the 4m `time_lost_per_cycle` input, kept as it is.

## Assumptions

1. The ADR's number and slug are guesses (`0153-…`). `flai task done` widens the touches with the real path.
2. `flai serve` can reach the host's upgrade through the host client that `serve.Options.Host` already carries.
3. An upgrade waits until the release is published and its signature verifies. A refused upgrade is retried only for a newer release, never in a loop.
4. The thread asking you to upgrade is anchored on `system-flow.yaml`. T-1567's ADR may choose another anchor.

## Not proposed

I am not proposing to raise `flai.minimum` when the manifest's test vocabulary changes. A host below the minimum stops serving the project entirely, which is worse than either instance.

## Alternative

A smaller version is T-1569 alone, with no `auto-upgrade`. It makes the refusal clear, but the host still lags. If you choose that, I would drop T-1568 and narrow T-1567 and T-1570.

### 2026-10-14T20:55:43Z orchestrator (recommendation)
**Recommendation:** approve the full plan, with `auto-upgrade` off by default and turned on only in a shell, as drafted. Then finalize S-0387. Take the smaller version (T-1569 alone) only if you want no automatic upgrade of the host at all.

## Why this needs you

- **It is a decision not yet recorded.** It adds a new ADR (T-1567) and a new host action that upgrades the host's flai on its own. That is a change to the operator's machine, and only you can consent to it.

## What fits

- **The cause matches both instances in I-0141.** The host's flai lags main after a release. The inbox shows the same lag today: running 1.48.0, newest 1.49.1.
- **With it off, it asks on one thread.** This mirrors ADR-0108 § Decision 4, where flai serve opens one thread at its restart limit.
- **The tasks fit.** Five tasks in four layers. T-1568 and T-1569 share no path. Every touch names a file, and the ADR path is a prediction that `flai task done` corrects.
- **The figures are flai's.** Forecast 1h15m, raised from flai's 34m with a stated reason. Cost of delay 10 USD a week, from flai's 4m input.

## One thing for T-1567's ADR to state

A host upgrade restarts `flai serve`, and with it the sessions it runs: the orchestrator, planners, and the dashboard's connection. The plan says story agents are detached and keep running. The ADR should say the same for the strategic agents. It should also say whether an upgrade waits until no planner run or acceptance is in flight.

I leave S-0387 a draft until you answer.

Source: design/adrs/0108-flai-serve-restarts-a-story-s-agent-that-ended-with-its-story-in-progress-up-to.md § Decision
