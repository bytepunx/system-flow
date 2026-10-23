---
id: D-0370
title: "I-0148's cause here is the close-out's config, not the rounds: pick the fix before building: Recommended for the operator: option 3. Do option 1 in S-0466 now, and record the queue-off case for later"
kind: answer
sources: [TH-0594]
topics: [cli, conventions, template, flai, go, code, all]
paths: [flai/cmd/verify.go, flai/cmd/verify_test.go, flai/internal/verify/story.go, design/system/flai-cli.md, docs/users/flai-reference.md, design/system/devex.md, docs/contributors/index.md, scripts/close-out.sh, design/issues/I-0148-a-close-out-whose-run-main-outruns-takes-two-rounds-of-the-integration-and-smoke-tiers-about-24-minutes-past-the-bash-tool-s-10-minute-limit-so-the-agent-s-call-is-killed-mid-round.md, design/issues/summary.md, design/issues/I-0182-criteria-tick-is-refused-on-another-open-story-s-narrative-state-finding-which-the-story-s-own-scoped-check-leaves-out.md, design/issues/I-0196-a-close-out-s-second-round-runs-go-test-again-when-the-rebase-regenerates-design-issues-summary-md-which-changes-the-branch-s-patch-id-so-a-close-out-main-outruns-still-takes-about-11-minutes-with-the-queue-on.md]
date: 2026-10-23
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0370 I-0148's cause here is the close-out's config, not the rounds: pick the fix before building: Recommended for the operator: option 3. Do option 1 in S-0466 now, and record the queue-off case for later

The operator, alex, answered TH-0594 at 2026-10-23T03:21:35Z.

What was asked:

```text
**Finding.** ADR-0174 already meant to close I-0148: while the merge queue can batch, flai verify leaves integration and smoke to the batch, and its Consequences say "I-0148 closes once this runs with the queue on". The queue is on for this project in flai serve's configuration (~/.flai/config.json: merge-queue, accept_reviews on, acceptance local). But scripts/close-out.sh runs the tree's flai through scripts/flai.sh, which reads .flai-cache/config.json, whose host_actions hold only push. So queueBatches (flai/cmd/verify.go:392) answers false in every close-out here, and integration and…
```

The answer:

```text
Recommended for the operator: option 3. Do option 1 in S-0466 now, and record the queue-off case for later.

- **The agent's evidence settles the cause here.** ADR-0174 already leaves integration and smoke to the merge queue, but the tree's `scripts/flai.sh` reads `.flai-cache/config.json`, which lacks the `merge-queue` host action. So every close-out in this repo runs 10 to 14 minutes of tiers the queue then runs again. The MCP verify of S-0483, which reads `~/.flai/config.json`, queued them and took 5 minutes.
- **Option 1 is small and closes I-0148 for this repo.** You run `scripts/flai.sh serve enable merge-queue`. flai verify names the condition that is off when it runs a `queue` tier. The dev docs say the dev config needs the host action.
- **The detached `flai close-out` still has value**, for projects whose tiers run long with the queue off, which is the template's default. But it is hours of work for a case this repo no longer hits. Keep it as a backlog story from I-0148, valued on its own.

This replaces my recommendation on TH-0596, which I made before reading this thread.

Source: design/adrs/0174-while-the-merge-queue-can-batch-flai-verify-and-the-close-out-leave-the-test.md
```

Source: TH-0594.
