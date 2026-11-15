---
id: D-0628
title: Item.Validate still checks a status against the baseline workflow, aliases included, not the repository's workflow, and Repo.Save does not validate
kind: decision
sources: [S-0737, S-0741]
topics: [flai, go, code]
paths: [flai/internal/workitem/item.go, flai/internal/workitem/repo.go, flai/internal/check/check_items.go]
date: 2026-11-14
status: active
absorbed_by: ""
asked_again: []
cited: [TH-1206, TH-1208, TH-1209, TH-1210, TH-1211, TH-1212, TH-1214, TH-1218, TH-1219]
---

# D-0628 Item.Validate still checks a status against the baseline workflow, aliases included, not the repository's workflow, and Repo.Save does not validate

S-0737 moved the states and transition map behind `Repo.Move` onto the workflow a `Repo` runs, but left the status check in `Item.Validate` on the baseline. Its narrative's Decisions:

> Item.Validate still checks a status against the baseline, aliases included, not against the repository's workflow. Repo.Save does not validate, so the test workflow does not need it. A later story can move Validate to the repository's workflow.

So an item in a state that only a project's own workflow defines still fails `Item.Validate`, and every check built on it, until a story moves the check onto the repository's workflow. A test that saves an item in a test workflow's own state passes only because `Repo.Save` does not validate.

S-0741, E-0034's story for the check, moved `flai check`'s own item rules onto a state's kind. It left `workitem`'s helpers to S-0739, `Item.Validate`'s status check among them (`wip/archive/agents/S-0741.md`, § Decisions):

> What workitem gives the rules (Closed, EpicFollows, Item.Validate's status check, CancelledWith) still compares names and is S-0739's.

So the status check moves with S-0739, not S-0741. Until then, a project's own file (S-0755, S-0756) can rely on `workflow.item`, which judges a status against `Repo.Flow()`, and not on `Item.Validate`.

Sources: S-0737's narrative, `wip/archive/agents/S-0737.md`, and S-0741's narrative, `wip/archive/agents/S-0741.md`, § Decisions.
