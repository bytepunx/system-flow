---
id: D-0448
title: "S-0547's close-out fails only at smoke's OpenRouter gateway check: the key is out of credits (I-0237): choice 2 - I have removed the OpenRouter key from my environment for now"
kind: recommendation
sources: [TH-0789, S-0547]
topics: [testing, cli, git, orchestration, planning, flai, go, code, all]
paths: [flai/internal/check, flai/internal/verify, flai/cmd, system-flow.yaml, design/system/flai-cli.md, design/conventions/testing.md, docs/users/flai.md, design/adrs, design/adrs/README.md, flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, flai/internal/mcpserver/verify.go, flai/internal/mcpserver/verify_test.go, docs/operators/settings.md, docs/users/flai-reference.md, design/system/project-manifest.md, design/conventions/code-quality.md, design/conventions/tooling.md, design/issues/I-0221-the-close-out-s-go-test-tier-runs-only-the-changed-packages-so-a-change-that-breaks-a-caller-package-s-tests-passes-before-review.md, design/issues/I-0203-a-story-the-operator-accepts-outside-the-merge-queue-lands-without-the-queue-tiers-so-a-stale-generated-file-on-main-fails-every-later-batch-at-integration.md, design/issues/summary.md, design/issues/I-0226-two-flai-serve-tests-check-before-what-they-wait-for-has-happened-and-fail-a-loaded-merge-queue-batch-on-timing.md, design/issues/I-0237-the-smoke-tier-s-openrouter-gateway-check-fails-on-a-402-because-the-operator-s.md]
date: 2026-11-01
status: active
absorbed_by: ""
asked_again: []
cited: [S-0685, TH-0989]
---

# D-0448 S-0547's close-out fails only at smoke's OpenRouter gateway check: the key is out of credits (I-0237): choice 2 - I have removed the OpenRouter key from my environment for now

The operator, alex, answered TH-0789 at 2026-11-01T23:29:54Z.

What was asked:

```text
**Recommended:** raise the OpenRouter key's limit, or top it up. Then answer here and I'll run the close-out again.
```

The recommendation it overruled:

```text
**Recommendation:** choice 1. Top up the OpenRouter key, or raise its limit, then answer here so agent-S-0547 runs the close-out again unchanged.

**Why choice 1.**
- It is spend, so the decision is yours alone.
- D-0334 is the answer you gave last time to the same kind of failure (S-0450, TH-0530): fix the host's gateway setup rather than close out with the gateway check skipped. Under ADR-0093 I accept only on a verify run that passed at the branch head. A smoke run that skips a check the manifest selects would need a decision you have not recorded.
- The key also fails every merge-queue batch that selects smoke (I-0237), not only S-0547. One top-up fixes both.

**Why not choice 3 now.** S-0548 is in review but is not near landing. The queue lists it as ineligible: criteria 4 and 9 are unticked and TH-0788 is open. S-0547 would wait an unknown time.

**If you would rather not spend:** say choice 2 here. Skipping OpenRouter in the close-out is then your recorded decision, and S-0547 notes it.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/digest/D-0334-s-0450-s-close-out-stops-only-at-the-litellm-gateway-smoke-check-which-claude.md
```

The operator's answer:

```text
choice 2 - I have removed the OpenRouter key from my environment for now.
```

Source: TH-0789.

What S-0547 did with it, from its narrative (`wip/archive/agents/S-0547.md`, § Decisions): "The operator chose choice 2 on TH-0789 (2026-11-01): close out without the OpenRouter key, which they removed from their environment. The close-out runs with OPENROUTER_API_KEY unset. S-0548, accepted meanwhile, moved the gateway checks out of smoke into a manual tier, so smoke no longer reaches them anyway." `scripts/gateway-smoke.sh` checks only the gateways whose key is set, so it skips OpenRouter without the key.

So the operator chose not to spend on the key for a close-out, against the orchestrator's recommendation, which cited D-0334. Since S-0548 the paid gateway checks run only as the manual `gateway` tier, `flai test gateway`, never in a close-out or a merge queue batch, so a close-out no longer stops on a gateway key's credit.

Sources: TH-0789; S-0547's narrative.
