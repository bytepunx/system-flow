---
id: I-0237
title: The smoke tier's OpenRouter gateway check fails on a 402 because the operator's OpenRouter key is out of credits, so every close-out and batch that runs smoke fails
class: blocker
status: open
count: 1
cost: 15m
first_reported: 2026-10-31T10:10:59Z
last_reported: 2026-10-31T10:10:59Z
updated: 2026-11-10T08:37:11Z
---

# I-0237 The smoke tier's OpenRouter gateway check fails on a 402 because the operator's OpenRouter key is out of credits, so every close-out and batch that runs smoke fails

## Description
The smoke tier's OpenRouter gateway check fails on a 402 because the operator's OpenRouter key is out of credits, so every close-out and batch that runs smoke fails

## Instances

### 2026-10-31T10:10:59Z
Story: S-0547.
S-0547's close-out at 9794966 ran smoke in the worktree (ADR-0233), and smoke failed only on TestGatewayOpenRouter in flai/internal/serve (scripts/gateway-smoke.sh). Claude Code through OpenRouter answered "API Error: 402 This request requires more credits, or fewer max_tokens. You requested up to 32000 tokens, but can only afford 6139", so the check found no flai-check.md and no guard refusal. The key's limit is the operator's to raise. Every other tier passed, integration included (9m41s). The merge queue's batch runs smoke with the same keys, so it fails the same way until the key has credit or the gateway checks leave smoke (S-0548).

## Remediation

Story S-0685 remediates this issue, created from it at 2026-11-10T08:37:11Z.
