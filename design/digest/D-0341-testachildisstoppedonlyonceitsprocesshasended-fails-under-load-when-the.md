---
id: D-0341
title: TestAChildIsStoppedOnlyOnceItsProcessHasEnded fails under load when the stubborn child dies on SIGTERM before it ignores it
kind: cause
sources: [I-0173]
topics: [flai, cli, go, code, all]
paths: [flai/internal/host/host_test.go, design/issues/I-0173-testachildisstoppedonlyonceitsprocesshasended-fails-under-load-when-the-stubborn-child-dies-on-sigterm-before-it-ignores-it.md, design/issues/summary.md, design/issues/I-0175-a-story-s-in-progress-transition-from-flai-story-start-was-lost-when-a-concurrent-write-rewrote-the-item-file-leaving-it-ready-while-its-agent-worked-it.md, design/issues/I-0158-testasettledreplystartstheagentwaitingonitonce-in-flai-internal-serve-fails-now-and-then-in-the-close-out-s-integration-tier-on-a-branch-that-does-not-touch-serve.md, design/issues/I-0174-the-merge-queue-s-batch-runs-smoke-with-no-close-out-story-so-scripts-check-sh-runs-an-unscoped-flai-check-strict-and-fails-every-batch-on-findings-outside-its-stories.md, design/issues/I-0179-accepting-a-story-that-adds-a-flaiover-dependency-leaves-the-main-checkout-s-flaiover-node-modules-without-it-so-every-close-out-s-integration-tier-fails-testeverylistedflaiovercomponentcompilesintheharness-until-someone-installs.md]
date: 2026-10-21
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0341 TestAChildIsStoppedOnlyOnceItsProcessHasEnded fails under load when the stubborn child dies on SIGTERM before it ignores it

I-0173 was closed on 2026-10-21.

The cause:

```text
TestAChildIsStoppedOnlyOnceItsProcessHasEnded fails under load when the stubborn child dies on SIGTERM before it ignores it
```

The fix:

```text
Fixed by S-0463 (T-1977): the stubborn child of the host tests writes a file once it ignores SIGTERM, and TestAChildThatWillNotStopIsKilled and TestAChildIsStoppedOnlyOnceItsProcessHasEnded wait for that file instead of sleeping 100 ms. The child now waits 300 ms before it ignores SIGTERM, as a loaded runner might, and both tests check that the stop took at least the grace period, so a child that SIGTERM ended fails them every run.
```

Source: I-0173.
