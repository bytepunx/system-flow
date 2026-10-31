---
id: TH-0789
title: "S-0547's close-out fails only at smoke's OpenRouter gateway check: the key is out of credits (I-0237)"
anchor:
  path: wip/kanban/stories/S-0547-the-close-out-runs-the-queue-tiers-and-the-tests-of-the-packages-that-import-a.md
  item: S-0547
status: open
participants: [agent-S-0547, orchestrator]
created: 2026-10-31T10:11:06Z
updated: 2026-10-31T10:12:09Z
---

# TH-0789 S-0547's close-out fails only at smoke's OpenRouter gateway check: the key is out of credits (I-0237)

On wip/kanban/stories/S-0547-the-close-out-runs-the-queue-tiers-and-the-tests-of-the-packages-that-import-a.md.

## Entries

### 2026-10-31T10:11:06Z agent-S-0547
**Recommended:** raise the OpenRouter key's limit, or top it up. Then answer here and I'll run the close-out again.

## What happened

S-0547's close-out at 9794966 ran every tier. The new `when` worked:

| Tier | Result |
|------|--------|
| go-test | skipped, covered by integration |
| integration | passed, 9m41s, "ran here, its when flai/** matches flai/internal/check/tests_test.go" |
| smoke | failed, 2m55s |

Smoke failed only on `TestGatewayOpenRouter` (`scripts/gateway-smoke.sh`). OpenRouter answered:

```text
API Error: 402 This request requires more credits, or fewer max_tokens. You requested up to 32000 tokens, but can only afford 6139.
```

That is the key's limit, not S-0547's code. I recorded it as I-0237. The merge queue's batch runs smoke with the same keys, so every batch that selects smoke fails the same way until the key has credit. S-0548 moves the gateway checks out of smoke into a gateway tier of their own.

## Your choices

1. **Raise the key's limit (recommended).** I run the close-out again unchanged.
2. **Let me run the close-out without the OpenRouter key.** `gateway-smoke.sh` checks only the gateways whose key is set, so it skips OpenRouter. I record the skipped check in S-0547's notes.
3. **Wait for S-0548.** It takes the gateway checks out of smoke. I sync onto it once it's accepted and close out then.

### 2026-10-31T10:12:09Z orchestrator (recommendation)
**Recommendation:** choice 1. Top up the OpenRouter key, or raise its limit, then answer here so agent-S-0547 runs the close-out again unchanged.

**Why choice 1.**
- It is spend, so the decision is yours alone.
- D-0334 is the answer you gave last time to the same kind of failure (S-0450, TH-0530): fix the host's gateway setup rather than close out with the gateway check skipped. Under ADR-0093 I accept only on a verify run that passed at the branch head. A smoke run that skips a check the manifest selects would need a decision you have not recorded.
- The key also fails every merge-queue batch that selects smoke (I-0237), not only S-0547. One top-up fixes both.

**Why not choice 3 now.** S-0548 is in review but is not near landing. The queue lists it as ineligible: criteria 4 and 9 are unticked and TH-0788 is open. S-0547 would wait an unknown time.

**If you would rather not spend:** say choice 2 here. Skipping OpenRouter in the close-out is then your recorded decision, and S-0547 notes it.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/digest/D-0334-s-0450-s-close-out-stops-only-at-the-litellm-gateway-smoke-check-which-claude.md
