---
id: D-0346
title: "The main checkout's flaiover/node_modules lacks @lucide/svelte since S-0458, so every close-out fails integration: Recommended: the operator runs `scripts/flaiover-install.sh --if-needed` in the main…"
kind: answer
sources: [TH-0563]
topics: [flai, cli, go, code, all]
paths: [flai/internal/host/host_test.go, design/issues/I-0173-testachildisstoppedonlyonceitsprocesshasended-fails-under-load-when-the-stubborn-child-dies-on-sigterm-before-it-ignores-it.md, design/issues/summary.md, design/issues/I-0175-a-story-s-in-progress-transition-from-flai-story-start-was-lost-when-a-concurrent-write-rewrote-the-item-file-leaving-it-ready-while-its-agent-worked-it.md, design/issues/I-0158-testasettledreplystartstheagentwaitingonitonce-in-flai-internal-serve-fails-now-and-then-in-the-close-out-s-integration-tier-on-a-branch-that-does-not-touch-serve.md, design/issues/I-0174-the-merge-queue-s-batch-runs-smoke-with-no-close-out-story-so-scripts-check-sh-runs-an-unscoped-flai-check-strict-and-fails-every-batch-on-findings-outside-its-stories.md, design/issues/I-0179-accepting-a-story-that-adds-a-flaiover-dependency-leaves-the-main-checkout-s-flaiover-node-modules-without-it-so-every-close-out-s-integration-tier-fails-testeverylistedflaiovercomponentcompilesintheharness-until-someone-installs.md]
date: 2026-10-21
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0346 The main checkout's flaiover/node_modules lacks @lucide/svelte since S-0458, so every close-out fails integration: Recommended: the operator runs `scripts/flaiover-install.sh --if-needed` in the main…

The operator, alex, answered TH-0563 at 2026-10-21T08:47:11Z.

What was asked:

```text
S-0463's close-out stops at the integration tier: `flai/internal/components/flaiover_test.go:93 src/lib/components/NavRail.svelte does not compile for the harness: 500`. The test links the main checkout's `flaiover/node_modules`, which has no `@lucide/svelte`. S-0458 added that package and was accepted at 08:32Z. This fails every story's close-out, not only S-0463's, which changes only `flai/internal/host/host_test.go` and `design/issues`. Recorded as I-0179.
```

The answer:

```text
Recommended: the operator runs `scripts/flaiover-install.sh --if-needed` in the main checkout. If it aborts with ERR_PNPM_ABORTED_REMOVE_MODULES_DIR_NO_TTY (I-0177), remove `flaiover/node_modules` there and run it again. S-0281's fix (I-0080) installs dependencies in story worktrees when their lockfile is newer. It does not reach the main checkout's `flaiover/node_modules`, which the integration test links, so S-0458's new @lucide/svelte is missing there. The agent should not install it itself: the main checkout is shared by every story's test runs and by the merge queue's batches, so rewriting it is a host change for you to time. Until it is done, every close-out stops at integration on NavRail.svelte. I-0179 records it; a lasting fix would have the integration test install what the main checkout's lockfile needs, as flaiover-install --if-needed does for a worktree.

Source: design/issues/I-0080-a-story-worktree-has-no-flaiover-node-modules-so-the-close-out-stops-at-the-flaiover-step-until-flaiover-install-runs-there.md § Remediation
```

Source: TH-0563.
