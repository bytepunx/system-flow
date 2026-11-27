---
id: D-0653
title: workitem's EpicFollows, PlanOf, NewHolds, PlanCandidatesOf, PromotionCandidates, FinalizableLacks, Changes, and CancelledWith run the baseline workflow, and a caller outside workitem reads the repository's board only through their -In forms, which no caller outside workitem calls yet
kind: decision
sources: [S-0739, S-0843]
topics: [flai, go, code]
paths: [flai/internal/workitem, flai/internal/check/check_items.go, flai/internal/inbox, flai/internal/serve/replan.go, flai/internal/hostapi, flai/internal/mcpserver, flai/cmd, flai/internal/itemedit, flai/internal/messages, flai/internal/storygit, flai/internal/verify, flai/internal/metrics]
date: 2026-11-17
status: active
absorbed_by: ""
asked_again: []
cited: [TH-1269, TH-1270]
---

# D-0653 workitem's EpicFollows, PlanOf, NewHolds, PlanCandidatesOf, PromotionCandidates, FinalizableLacks, Changes, and CancelledWith run the baseline workflow, and a caller outside workitem reads the repository's board only through their -In forms, which no caller outside workitem calls yet

S-0739 moved `flai/internal/workitem`'s epic following, board view, task plan, holds, order, and promotion off state names and onto a state's kind (ADR-0301). Its narrative's Decisions (`wip/archive/agents/S-0739.md`) record how it kept every caller in other packages unchanged:

> Package functions that take no Repo (PlanOf, PullSequence, EpicFollows, NewHolds, NewBoardView and the like) read the workflow from what they are given (the Board, the Holds) or get a form that takes it; where an exported signature would change for callers in other packages, the existing form keeps its signature and runs the baseline, so those callers change in their own E-0034 stories (S-0742 to S-0749) and this story stays in `workitem`.
>
> Callers outside workitem keep the baseline forms (EpicFollows, PlanOf, NewHolds, PlanCandidatesOf, PromotionCandidates, FinalizableLacks, Changes, CancelledWith) and switch to the -In forms in their own stories: check/check_items.go (epicLags, CancelledWith, NewHolds) under S-0741; transition.go's `to == Cancelled` and rules.go's PlanOf under S-0744; inbox, serve/replan.go, hostapi, mcpserver, cmd, itemedit, messages, storygit, verify, metrics under S-0742, S-0746, S-0747, and S-0749.

The -In forms are `EpicFollowsIn`, `PlanOfIn`, `NewHoldsIn`, `PlanCandidatesIn`, `PromotionCandidatesIn`, `FinalizableLacksIn`, `ChangesIn`, `CancelledWithIn`, and `OrderedIn`. S-0741, S-0744, S-0746, S-0747, and S-0749 were accepted before S-0739 (2026-11-17), so none switched: on main at S-0739's acceptance every call outside `workitem` is still a baseline form, such as `workitem.NewHolds` in `check`, `cmd`, `itemedit`, `mcpserver`, `storygit`, `taskdone`, and `verify`, `workitem.Changes` in `inbox` and `serve/replan.go`, and `workitem.PlanOf` in `cmd`, `hostapi`, and `mcpserver`. A project whose workflow file renames a state gets the baseline's answer from those calls.

So a story that touches one of those callers, S-0742 among them, calls the `-In` form with `repo.Flow()` (D-0630) rather than the baseline form, and a story that finishes E-0034 checks that no baseline form is left outside `workitem` and its tests.

S-0843 made that switch for every caller with a repository at hand, so the title's "no caller outside workitem calls yet" held only until 2026-11-27. What is left on main is in [D-0676](D-0676-since-s-0843-every-caller-outside-workitem-that-has-a-repository-calls-the.md).

Source: S-0739's narrative, `wip/archive/agents/S-0739.md`, § Decisions.
