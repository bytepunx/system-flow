---
id: I-0250
title: The orchestrator's conventions sit within about 1.5KB of its 80KB pack budget, so a few more lines in work-management.md or strategic-agents.md drop every design brief from its pack
class: defect
status: open
count: 1
cost: 25m
first_reported: 2026-11-02T03:46:07Z
last_reported: 2026-11-02T03:46:07Z
updated: 2026-11-02T03:46:07Z
---

# I-0250 The orchestrator's conventions sit within about 1.5KB of its 80KB pack budget, so a few more lines in work-management.md or strategic-agents.md drop every design brief from its pack

## Description
The orchestrator's conventions sit within about 1.5KB of its 80KB pack budget, so a few more lines in work-management.md or strategic-agents.md drop every design brief from its pack

## Instances

### 2026-11-02T03:46:07Z
Story: S-0422.
S-0422's close-out stopped at integration: TestAThreadPackIsUnderAQuarterOfTheGeneralPack found the thread pack at 28.3% of the general orchestrate pack. Two new work-management.md rules (about 2.3KB) and a strategic-agents.md sentence took the orchestrator's conventions from 71504 to 74169 bytes, past the budget with the digest, so addDesign marked it exceeded conventions and loaded none of the 138 briefs; the general pack fell from 191973 to 155471 bytes. Shortening the rules to one of 903 bytes brought conventions to 72769 and the briefs back. The margin is now under about 1.5KB: the next story that adds a rule to either convention tips it again, and the orchestrator loses its design briefs with no warning but the test.

## Remediation
