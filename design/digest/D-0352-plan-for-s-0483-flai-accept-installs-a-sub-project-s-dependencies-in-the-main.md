---
id: D-0352
title: "Plan for S-0483: flai accept installs a sub-project's dependencies in the main checkout when its lock file changed: Go ahead with 483 as planned and update the rule it is breaking to include the…"
kind: recommendation
sources: [TH-0578]
topics: [release, testing, flai, cli, go, code, all]
paths: [design/adrs/README.md, design/system/project-manifest.md, design/system/flai-cli.md, flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, system-flow.yaml, docs/operators/settings.md, flai/cmd/accept.go, flai/cmd/accept_install.go, flai/cmd/accept_install_test.go, docs/users/flai.md, design/system/devex.md, design/issues/I-0179-accepting-a-story-that-adds-a-flaiover-dependency-leaves-the-main-checkout-s-flaiover-node-modules-without-it-so-every-close-out-s-integration-tier-fails-testeverylistedflaiovercomponentcompilesintheharness-until-someone-installs.md, design/issues/summary.md, design/adrs/0185-flai-accept-runs-a-sub-project-s-install-command-in-the-main-checkout-when-the.md, flai/internal/check/check.go, flai/internal/check/install_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/shared.go]
date: 2026-10-21
status: absorbed
absorbed_by: "design/adrs/0185-flai-accept-runs-a-sub-project-s-install-command-in-the-main-checkout-when-the.md#decision"
asked_again: []
cited: []
---

# D-0352 Plan for S-0483: flai accept installs a sub-project's dependencies in the main checkout when its lock file changed: Go ahead with 483 as planned and update the rule it is breaking to include the…

The operator, alex, answered TH-0578 at 2026-10-21T12:36:20Z.

What was asked:

```text
S-0483 is planned: five tasks in four layers, 1h forecast, cost of delay 37.50 USD a week.
```

The recommendation it overruled:

```text
Recommended for the operator: go ahead with S-0482 now, and decide on S-0483 separately. My lean is to cancel S-0483, as on TH-0576, unless you want flai to install into the main checkout.

- The planner is right that the two fix different sides. S-0482 stops story worktrees from borrowing the main checkout's `node_modules`, so a close-out no longer fails when main's install is stale. S-0483 makes `flai accept` install into the main checkout when a lock file changed.
- Once S-0482 lands, nothing in the close-outs or the merge queue depends on main's install any more. What still uses it is the main checkout's own dev server and dashboard work, which you install by hand today.
- S-0483 reverses a recorded rule: "The main checkout is never installed into" (component-library.md). It adds a manifest key, `install`, and a step in every acceptance, which is why its plan starts with an ADR (T-1998). That is a decision not yet made, so it is yours.
- If you want main kept installed on its own, keep S-0483 as planned: five tasks, 1h, 37.50 USD a week. Otherwise cancel it, and S-0482 can close I-0179 as well.

S-0482 is in ready. I am keeping S-0483 in the backlog, not promoted, until you answer.

Source: design/system/component-library.md § A story's worktree (S-0442)
```

The operator's answer:

```text
Go ahead with 483 as planned and update the rule it is breaking to include the exception and reason for this edge case
```

Source: TH-0578.
