---
id: I-0250
title: The orchestrator's conventions sit within about 1.5KB of its 80KB pack budget, so a few more lines in work-management.md or strategic-agents.md drop every design brief from its pack
class: defect
status: closed
count: 1
cost: 25m
first_reported: 2026-11-02T03:46:07Z
last_reported: 2026-11-02T03:46:07Z
updated: 2026-11-02T04:39:19Z
---

# I-0250 The orchestrator's conventions sit within about 1.5KB of its 80KB pack budget, so a few more lines in work-management.md or strategic-agents.md drop every design brief from its pack

## Description
The orchestrator's conventions sit within about 1.5KB of its 80KB pack budget, so a few more lines in work-management.md or strategic-agents.md drop every design brief from its pack

## Instances

### 2026-11-02T03:46:07Z
Story: S-0422.
S-0422's close-out stopped at integration: TestAThreadPackIsUnderAQuarterOfTheGeneralPack found the thread pack at 28.3% of the general orchestrate pack. Two new work-management.md rules (about 2.3KB) and a strategic-agents.md sentence took the orchestrator's conventions from 71504 to 74169 bytes, past the budget with the digest, so addDesign marked it exceeded conventions and loaded none of the 138 briefs; the general pack fell from 191973 to 155471 bytes. Shortening the rules to one of 903 bytes brought conventions to 72769 and the briefs back. The margin is now under about 1.5KB: the next story that adds a rule to either convention tips it again, and the orchestrator loses its design briefs with no warning but the test.

## Remediation
Closed 2026-11-02T04:39:19Z: S-0563 removed the cliff. On TH-0805 the operator chose option 1, recorded in ADR-0241, which refines ADR-0049. `addDesign` in `flai/internal/context/context.go` now loads what is named and every brief, and ranks nothing, when the conventions and the open-issue table alone take a pack over its budget. The pack keeps `exceeded: conventions`. So no number of open issues and no few lines of a convention drop the orchestrator's design briefs: at a 72KB budget its general pack is over on its conventions and still holds all 139 briefs. `TestTheHeaderSaysWhatAloneExceedsTheBudget` and `TestPrimeStoryDesign` reproduce it and fail on main. `TestAThreadPackIsUnderAQuarterOfTheGeneralPack` passes at 22.6%.
