---
id: D-0310
title: "S-0387 plan: flai serve notices a newer flai release while it runs, upgrades the host behind a new auto-upgrade host action or asks on a thread, and tier refusals name the version (new ADR), five…"
kind: recommendation
sources: [TH-0447]
topics: [release, cli, mcp, host, flai, go, code, all]
paths: [design/adrs/0153-flai-serve-upgrades-the-host-s-flai-when-the-project-s-history-has-a-newer-published-release-behind-auto-upgrade.md, design/adrs/README.md, design/system/flai-cli.md, flai/internal/serve/upgrade.go, flai/internal/serve/upgrade_test.go, flai/internal/serve/serve.go, flai/internal/serve/outdated_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/internal/verify/manifest.go, flai/internal/verify/manifest_test.go, flai/internal/mcpserver/test.go, flai/internal/mcpserver/test_test.go, flai/internal/mcpserver/verify.go, flai/internal/mcpserver/verify_test.go, docs/operators/index.md, docs/operators/settings.md, docs/operators/runbooks/update.md, docs/users/flai.md, design/issues/I-0141-the-mcp-test-tool-of-the-installed-flai-1-40-1-refuses-main-s-manifest-so-a-story-s-agent-runs-flai-test-from-the-tree.md, design/issues/summary.md, design/adrs/0158-flai-serve-checks-each-served-project-for-a-newer-flai-release-every-minute.md, flai/cmd/test.go, flai/cmd/test_test.go, flai/cmd/verify.go, flai/cmd/verify_test.go, flai/cmd/serve.go]
date: 2026-10-15
status: absorbed
absorbed_by: "design/system/flai-cli.md#versions-the-hosts-flai-and-the-tree"
asked_again: []
cited: []
---

# D-0310 S-0387 plan: flai serve notices a newer flai release while it runs, upgrades the host behind a new auto-upgrade host action or asks on a thread, and tier refusals name the version (new ADR), five…

The operator, alex, answered TH-0447 at 2026-10-15T02:28:10Z.

What was asked:

```text
**Plan:** S-0387 gets five tasks in four layers. The remedy needs a new ADR and a new host action, `auto-upgrade`, so it needs your approval before S-0387 is finalized.
```

The recommendation it overruled:

```text
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
```

The operator's answer:

```text
approve with recommendation
```

Source: TH-0447.
