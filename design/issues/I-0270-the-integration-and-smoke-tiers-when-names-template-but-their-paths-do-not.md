---
id: I-0270
title: "The integration and smoke tiers' when names template/**, but their paths do not select a template-only change, so a story that changes only template/ runs neither at close-out, against what code-quality.md and tooling.md say"
class: defect
status: open
count: 1
cost: 5m
first_reported: 2026-11-07T04:41:00Z
last_reported: 2026-11-07T04:41:00Z
updated: 2026-11-10T08:54:31Z
---

# I-0270 The integration and smoke tiers' when names template/**, but their paths do not select a template-only change, so a story that changes only template/ runs neither at close-out, against what code-quality.md and tooling.md say

## Description
The integration and smoke tiers' when names template/**, but their paths do not select a template-only change, so a story that changes only template/ runs neither at close-out, against what code-quality.md and tooling.md say

## Instances

### 2026-11-07T04:41:00Z
Story: S-0558.
S-0558 changed only template/ and design/ paths, and its close-out at 56003eb ran rebase, sync, narrative, check, and the markdown and template tiers, with neither integration nor smoke run or listed as queued. SelectStoryLeavingQueue (flai/internal/verify/select.go) applies `when` only to the tiers selectStoryTiers already chose by `paths`, as ADR-0233 says ("a selected queue tier"). In system-flow.yaml, integration's paths are flai/**and scripts/integration.sh, and smoke's are flai/** and scripts/**/*.sh, so `when: template/**` decides nothing for a template-only branch. Yet code-quality.md's and tooling.md's project additions, and docs/users/flai.md, say a story that touches or changes template/ runs them. Either the tiers' paths should add template/**, or those documents and the when patterns should stop naming template/. No harm to S-0558, a one-clause prose change whose template tier passed; the risk is a template change that breaks what only the integration tier's repository tests catch, which then first fails in the merge queue's batch.

## Remediation

Story S-0696 remediates this issue, created from it at 2026-11-10T08:54:31Z.
