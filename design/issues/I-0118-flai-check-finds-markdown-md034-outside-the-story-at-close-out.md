---
id: I-0118
title: "flai check finds `markdown.MD034` outside the story at close-out"
class: efficiency
status: closed
count: 5
first_reported: 2026-10-08T04:23:40Z
last_reported: 2026-10-08T05:31:18Z
updated: 2026-10-08T09:36:09Z
---

# I-0118 flai check finds `markdown.MD034` outside the story at close-out

## Description
flai check finds `markdown.MD034` outside the story at close-out

## Instances

### 2026-10-08T04:23:40Z
Story: S-0324.
flai check found outside the story:
`wip/agents/orchestrator.md`: MD034/no-bare-urls Bare URL used

### 2026-10-08T04:40:57Z
Story: S-0316.
flai check found outside the story:
`wip/agents/orchestrator.md`: MD034/no-bare-urls Bare URL used

### 2026-10-08T04:48:49Z
Story: S-0320.
flai check found outside the story:
`wip/agents/orchestrator.md`: MD034/no-bare-urls Bare URL used

### 2026-10-08T04:52:00Z
Story: S-0318.
flai check found outside the story:
`wip/agents/orchestrator.md`: MD034/no-bare-urls Bare URL used

### 2026-10-08T05:31:18Z
Story: S-0318.
flai check found outside the story:
`wip/threads/TH-0361-s-0320-s-close-out-fails-integration-on-a-bare-www-that-s-0316-s-acceptance-committed-to-main-s-orchestrator-log.md`: MD034/no-bare-urls Bare URL used

## Remediation

Story S-0346 remediates this issue, created from it at 2026-10-08T08:08:21Z.
Closed 2026-10-08T09:36:09Z: All five instances were bare `www.` literals that an installed flai older than 1.39.4 wrote in wip, before S-0324 gave flai's MD034 the bare `www.` rule (I-0110). Every write an agent makes to wip now passes through LintGuard, which refuses one. S-0346 fixed the one write flai makes from text no agent checks: a strategic run's end entry now puts the bare URLs in its summary in a code span (mdlint.QuoteBareURLs, called from logRunEndSaying). TestARunEndQuotesABareURLInItsSummary reproduces the refusal. flai check --strict on main finds no markdown.MD034.
