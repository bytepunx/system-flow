---
id: D-0350
title: A story that adds a flaiover dependency leaves the main checkout's node_modules without it, so every other close-out fails TestEveryListedFlaioverComponentCompilesInTheHarness
kind: cause
sources: [I-0178]
topics: [flai, cli, go, code, all]
paths: [scripts/integration.sh, scripts/README.md, design/system/devex.md, flai/internal/components/integration_test.go, design/issues/I-0178-a-story-that-adds-a-flaiover-dependency-leaves-the-main-checkout-s-node-modules-without-it-so-every-other-close-out-fails-testeverylistedflaiovercomponentcompilesintheharness.md, design/issues/summary.md]
date: 2026-10-21
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0350 A story that adds a flaiover dependency leaves the main checkout's node_modules without it, so every other close-out fails TestEveryListedFlaioverComponentCompilesInTheHarness

I-0178 was closed on 2026-10-21.

The cause:

```text
A story that adds a flaiover dependency leaves the main checkout's node_modules without it, so every other close-out fails TestEveryListedFlaioverComponentCompilesInTheHarness
```

The fix:

```text
S-0482: scripts/integration.sh runs flaiover-install.sh --if-needed in a worktree before go test, so a story's or the merge queue's integration tier runs on its own install from its own lockfile (verified in S-0482's worktree: the install took 4.3s and TestEveryListedFlaioverComponentCompilesInTheHarness passed on it). flaioverModules in flai/internal/components/integration_test.go borrows another checkout's node_modules only when its .pnpm/lock.yaml is this checkout's pnpm-lock.yaml, and otherwise fails naming the install; TestInstalledFromRefusesAnInstallFromAnotherLockfile reproduces the stale borrowed install.
```

Source: I-0178.
