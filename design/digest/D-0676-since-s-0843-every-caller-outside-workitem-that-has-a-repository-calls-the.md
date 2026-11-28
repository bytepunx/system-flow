---
id: D-0676
title: "Since S-0843 every caller outside workitem that has a repository calls the workitem -In forms with repo.Flow(); a NewHolds(nil, ...) claim of bare touches reads no task and keeps the baseline form, and serve's direct baselineFlow reads are I-0406"
kind: decision
sources: [S-0843, I-0379, I-0406]
topics: [flai, go, code, workflow]
paths: [flai/internal/workitem/hold.go, flai/internal/messages/share.go, flai/internal/check/check_items.go, flai/internal/storygit/sync.go, flai/internal/mcpserver/shared.go, flai/internal/verify/story.go, flai/internal/serve/replan.go, flai/internal/serve/flow.go, flai/internal/serve/kluster.go, flai/internal/serve/components.go, flai/internal/serve/curate_schedule.go]
date: 2026-11-27
status: active
absorbed_by: ""
asked_again: []
cited: [S-0851, TH-1333]
---

# D-0676 Since S-0843 every caller outside workitem that has a repository calls the workitem -In forms with repo.Flow(); a NewHolds(nil, ...) claim of bare touches reads no task and keeps the baseline form, and serve's direct baselineFlow reads are I-0406

S-0843 fixed I-0379, the callers that still passed the baseline workflow to `workitem`'s wrappers after S-0739 added the `-In` forms ([D-0653](D-0653-workitem-s-epicfollows-planof-newholds-plancandidatesof-promotioncandidates.md)). Its narrative's Decisions (`wip/archive/agents/S-0843.md`) record which callers moved and which stay:

> Added T-3609: a search for workitem.NewHolds with items found five more callers that read task states by the baseline through Holds.Claim, missed by the planner's table: messages/share.go (the issue's share hold), check/check_items.go (wip.overlap and componentTag), storygit/sync.go (outsideClaim), mcpserver/shared.go (storyClaim), and verify/story.go (storyRun.claim). Callers of NewHolds(nil, ...) with a bare touches item read no task and stay.
>
> Added T-3610 from T-3567's search: hostapi, mcpserver, and cmd/items read a story's task plan with PlanOf; serve/replan and inbox read changes with Changes; check reads the cancel cascade and epic follow with CancelledWith and EpicFollows; cmd/plan and inbox/starved read PlanCandidatesOf and PromotionCandidates; each has the repository at hand.
>
> I gave news a flow parameter and passed repo.Flow() from its three callers (replan.go, curate_schedule.go, kluster.go) and workflow.Baseline() from replan_test.go.

So on main the only baseline forms left outside `workitem` are `workitem.NewHolds(nil, ...)` calls that claim a bare `&workitem.Item{Touches: ...}`, in `itemedit`, `metrics`, `context`, `taskdone`, `preview`, `mcpserver`, and `cmd`. They read no task state, so they need no workflow. A new caller that claims an item with tasks, or reads a plan, changes, a cancel cascade, an epic follow, or candidates, calls the `-In` form with `repo.Flow()`, and its test builds a renamed workflow, as S-0843's tests do.

S-0843 left three reads in `flai/internal/serve` that use `baselineFlow` directly, and recorded them as I-0406 rather than widen the story: `kluster.go`'s send-back check, `components.go`'s `StoryLeft`, and `curate_schedule.go`'s `acceptance()`. Each may be deliberate or may lack a repository where it runs. The operator decides whether I-0406 becomes a story. Until then a story that touches one of them reads I-0406 first, and does not take the read for a fallback. The fallbacks when no project can be opened, `serve/flow.go`, hostapi, `kluster.flowOf`, and `guard.flow`, are not part of I-0406.

On main at S-0843's acceptance, `orchestrate_reasons.go` reads `baselineFlow` too, for the initial and waiting states. Neither S-0843 nor I-0406 names that read.

Source: S-0843's narrative, `wip/archive/agents/S-0843.md`, § Decisions; I-0379; I-0406.
