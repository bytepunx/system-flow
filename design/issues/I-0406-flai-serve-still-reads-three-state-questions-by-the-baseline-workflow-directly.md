---
id: I-0406
title: flai serve still reads three state questions by the baseline workflow directly, not through a wrapper, so a project's own workflow file will not reach them
class: impression
status: open
count: 1
cost: 10m
first_reported: 2026-11-18T10:51:47Z
last_reported: 2026-11-18T10:51:47Z
updated: 2026-11-27T22:34:11Z
---

# I-0406 flai serve still reads three state questions by the baseline workflow directly, not through a wrapper, so a project's own workflow file will not reach them

## Description
flai serve still reads three state questions by the baseline workflow directly, not through a wrapper, so a project's own workflow file will not reach them

## Instances

### 2026-11-18T10:51:47Z
Story: S-0843.
S-0843 moved every caller of a baseline wrapper (NewHolds with items, HoldsMemo.NewHolds, NewPuller, PlanOf, Changes, CancelledWith, EpicFollows, PlanCandidatesOf, PromotionCandidates) that has a repository at hand to its In form. Its search also found three reads in flai/internal/serve that use the package's baselineFlow directly and so lie outside I-0379: kluster.go's send-back check (`baselineFlow.Gives(it.Type, ..., workflow.EffectSendBack)`, around line 295), components.go's StoryLeft (`wf := baselineFlow`, which its comment says is all a look reads of the project), and curate_schedule.go's acceptance() (`baselineFlow.WithEffect(TypeStory, EffectAccept)`, which its comment ties to roles' conditions naming the baseline's states until S-0748). Each may be deliberate or may lack a repository where it runs; none was changed. Under a board that renames review, in-progress, or done they would read the baseline's states. Fallbacks to the baseline when no project can be opened (serve/flow.go, hostapi moving, kluster.flowOf, guard.flow) are not this.

## Remediation

Story S-0851 remediates this issue, created from it at 2026-11-27T22:34:11Z.
