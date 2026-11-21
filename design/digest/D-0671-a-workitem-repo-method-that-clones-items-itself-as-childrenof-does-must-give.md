---
id: D-0671
title: "A workitem Repo method that clones items itself, as ChildrenOf does, must give each clone the repository's workflow, else Item.Closed answers by the baseline; TestCloneSharesNothing exempts Item.flow, which the items of one call share read-only"
kind: decision
sources: [S-0842, I-0377]
topics: [flai, go, code, workflow]
paths: [flai/internal/workitem/item.go, flai/internal/workitem/repo.go, flai/internal/workitem/view.go, flai/internal/workitem/store_test.go]
date: 2026-11-18
status: active
absorbed_by: ""
asked_again: []
cited: [S-0807, S-0839, S-0843, TH-1309, TH-1310, TH-1311, TH-1312, TH-1313, TH-1314, TH-1316]
---

# D-0671 A workitem Repo method that clones items itself, as ChildrenOf does, must give each clone the repository's workflow, else Item.Closed answers by the baseline; TestCloneSharesNothing exempts Item.flow, which the items of one call share read-only

S-0842 fixed I-0377 by giving each item a repository returns an unexported `flow`, the repository's workflow, so that `Item.Closed` reads the workflow's `accepted` and `cancelled` kinds rather than the baseline's `done` and `cancelled`. `design/system/workflow.md` § Transitions and who makes them states that behaviour. Its narrative's Decisions (`wip/archive/agents/S-0842.md`) record where it is set and what a later change must keep:

> Repo.List and Repo.Get set an unexported workflow on each item they return, after the parse cache's clone, so no shared item changes, and Item.Closed reads Workflow.Closed (the accepted and cancelled kinds). An item parsed alone keeps the baseline's names.
>
> T-3572's sub-agent stopped at Repo.ChildrenOf in view.go, outside its touches: ChildrenOf clones children itself, so without the workflow its items answered Closed by the baseline and TestChildrenOfIsChildrenOfTheList failed. [...] TestCloneSharesNothing exempts Item.flow, which is read-only and shared by the items of one call.

So a story that adds a `Repo` method in `flai/internal/workitem` that returns items it clones itself sets `flow` on each clone, as `view.go`'s `ChildrenOf` does with `r.Flow()`, or its items answer `Closed` by the baseline on a board that renames the closed states. A new pointer field on `Item` that is shared on purpose needs its own exemption in `store_test.go`'s `TestCloneSharesNothing`, beside `Item.flow`.

Source: S-0842's narrative, `wip/archive/agents/S-0842.md`, § Decisions; I-0377.
