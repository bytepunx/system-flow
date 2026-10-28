---
id: I-0183
title: "flai check finds `item.chronology` outside the story at close-out"
class: efficiency
status: closed
count: 1
first_reported: 2026-10-21T12:18:49Z
last_reported: 2026-10-21T12:18:49Z
updated: 2026-10-28T14:39:55Z
---

# I-0183 flai check finds `item.chronology` outside the story at close-out

## Description
flai check finds `item.chronology` outside the story at close-out

## Instances

### 2026-10-21T12:18:49Z
Story: S-0473.
flai check found outside the story:
`wip/kanban/stories/S-0482-a-story-that-adds-a-flaiover-dependency-leaves-the-main-checkout-s-node-modules-without-it-so-every-other-close-out-fails-testeverylistedflaiovercomponentcompilesintheharness.md`: updated is before the last transition

## Remediation

Story S-0522 remediates this issue, created from it at 2026-10-28T12:53:36Z.
Closed 2026-10-28T14:39:55Z: Fixed by S-0522 (ADR-0208, refining ADR-0184). Repo.Save now raises a work item's updated to the later of its last transition's at and its created when the stamp a writer passes is earlier, under the items lock and after the merge, so no write flai makes leaves an item item.chronology reports. The replanner's forecast write (serve/replan.go), the merge queue's send-back (mergequeue/bisect.go), and the touches widening (itemedit/widen.go) now read their clock after the item they stamp. TestSaveKeepsUpdatedAtTheLastTransitionWhenTheWriterStampIsEarlier, TestSaveKeepsUpdatedAtTheLastTransitionAfterAMergeOfASameSecondMove, and TestSaveKeepsUpdatedAtCreated reproduce the S-0482 instance; TestAStoryMovedAfterTheLookBeganIsForecastNoEarlierThanItsMove, TestASendBackIsNeverStampedBeforeAMoveItRead, and TestWidenStoryReadsEachItemBeforeItsClock cover the writers. flai check --strict on the main checkout finds no item.chronology; one a hand edit makes is still reported.
