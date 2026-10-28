---
id: D-0372
title: A close-out whose run main outruns takes two rounds of the integration and smoke tiers, about 24 minutes, past the Bash tool's 10-minute limit, so the agent's call is killed mid-round
kind: cause
sources: [I-0148]
topics: [cli, conventions, template, flai, go, code, all]
paths: [flai/cmd/verify.go, flai/cmd/verify_test.go, flai/internal/verify/story.go, design/system/flai-cli.md, docs/users/flai-reference.md, design/system/devex.md, docs/contributors/index.md, scripts/close-out.sh, design/issues/I-0148-a-close-out-whose-run-main-outruns-takes-two-rounds-of-the-integration-and-smoke-tiers-about-24-minutes-past-the-bash-tool-s-10-minute-limit-so-the-agent-s-call-is-killed-mid-round.md, design/issues/summary.md, design/issues/I-0182-criteria-tick-is-refused-on-another-open-story-s-narrative-state-finding-which-the-story-s-own-scoped-check-leaves-out.md, design/issues/I-0196-a-close-out-s-second-round-runs-go-test-again-when-the-rebase-regenerates-design-issues-summary-md-which-changes-the-branch-s-patch-id-so-a-close-out-main-outruns-still-takes-about-11-minutes-with-the-queue-on.md]
date: 2026-10-23
status: absorbed
absorbed_by: "design/adrs/0174-while-the-merge-queue-can-batch-flai-verify-and-the-close-out-leave-the-test.md#context"
asked_again: []
cited: []
---

# D-0372 A close-out whose run main outruns takes two rounds of the integration and smoke tiers, about 24 minutes, past the Bash tool's 10-minute limit, so the agent's call is killed mid-round

I-0148 was closed on 2026-10-23.

The cause:

```text
A close-out whose run main outruns takes two rounds of the integration and smoke tiers, about 24 minutes, past the Bash tool's 10-minute limit, so the agent's call is killed mid-round
```

The fix:

```text
S-0466: the close-out here runs the tree's flai through scripts/flai.sh, which reads .flai-cache/config.json, and that configuration lacked the merge-queue host action flai serve's has, so ADR-0174 never left integration and smoke to the queue and every round ran them (10 to 14 minutes). The operator turned the host action on in the dev configuration (TH-0601); flai verify now names the tiers it runs and the condition that is off when the queue cannot batch (queue_off, TestVerifyLeavesTheQueueTiersToTheMergeQueueWhileItCanBatch), and devex.md and docs/contributors say the dev configuration needs it. The detached close-out for projects with the queue off is left to a backlog story (TH-0594).
```

Source: I-0148.
