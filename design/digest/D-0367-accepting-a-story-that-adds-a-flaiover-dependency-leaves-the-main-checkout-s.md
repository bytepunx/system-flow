---
id: D-0367
title: "Accepting a story that adds a flaiover dependency leaves the main checkout's flaiover/node_modules without it, so every close-out's integration tier fails…"
kind: cause
sources: [I-0179]
topics: [release, testing, flai, cli, go, code, all]
paths: [design/adrs/README.md, design/system/project-manifest.md, design/system/flai-cli.md, flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, system-flow.yaml, docs/operators/settings.md, flai/cmd/accept.go, flai/cmd/accept_install.go, flai/cmd/accept_install_test.go, docs/users/flai.md, design/system/devex.md, design/issues/I-0179-accepting-a-story-that-adds-a-flaiover-dependency-leaves-the-main-checkout-s-flaiover-node-modules-without-it-so-every-close-out-s-integration-tier-fails-testeverylistedflaiovercomponentcompilesintheharness-until-someone-installs.md, design/issues/summary.md, design/adrs/0185-flai-accept-runs-a-sub-project-s-install-command-in-the-main-checkout-when-the.md, flai/internal/check/check.go, flai/internal/check/install_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/shared.go]
date: 2026-10-21
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0367 Accepting a story that adds a flaiover dependency leaves the main checkout's flaiover/node_modules without it, so every close-out's integration tier fails…

I-0179 was closed on 2026-10-21.

The cause:

```text
Accepting a story that adds a flaiover dependency leaves the main checkout's flaiover/node_modules without it, so every close-out's integration tier fails TestEveryListedFlaioverComponentCompilesInTheHarness until someone installs
```

The fix:

```text
S-0483: flai accept now runs a sub-project's projects[].install command in the main checkout after an acceptance that changed its install.paths (ADR-0185). flaiover's names scripts/flaiover-install.sh --if-needed on package.json and pnpm-lock.yaml, so accepting a story that adds a flaiover dependency leaves the main checkout's node_modules current. TestAcceptInstallsASubProjectWhoseLockTheStoryChanged reproduces the cause.
```

Source: I-0179.
