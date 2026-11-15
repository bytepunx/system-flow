---
id: D-0637
title: A test of a package moved onto a state's kind runs it on a workflow that renames in-progress and review with no alias and adds an active state between them, as serve's flow_test.go renamedFlow does, so that no comparison by name passes by accident
kind: decision
sources: [S-0743]
topics: [flai, go, code]
paths: [flai/internal/serve/flow_test.go, flai/internal/workitem/rules_test.go]
date: 2026-11-15
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0637 A test of a package moved onto a state's kind runs it on a workflow that renames in-progress and review with no alias and adds an active state between them, as serve's flow_test.go renamedFlow does, so that no comparison by name passes by accident

S-0743 moved flai serve's kluster deploys, forge pass, and merge queue off the names `review` and `done` and onto a state's kind and a transition's effect (ADR-0301). Its narrative's Decisions (`wip/archive/agents/S-0743.md`) say how its tests prove it:

> The test workflow renames in-progress and review with no alias and adds a state of kind active between them, so that no comparison by name can pass by accident; workitem's own fixture keeps review as an alias.

The fixture is `renamedFlow` in `flai/internal/serve/flow_test.go`, parsed by `renamedWorkflow(t)` for `Options.Workflow`: `doing` and `awaiting` stand for in-progress and review, and `testing`, of kind `active`, sits between them. A pass that still reads a state's name deploys, pushes, queues, or sends back the wrong story under it.

So a later E-0034 story that moves a package onto kinds tests it on a board of that shape, not only on the baseline, and not on `repoWithWorkflow`'s fixture in `flai/internal/workitem/rules_test.go` alone (D-0627), whose `review` alias would let a name comparison pass.

Source: S-0743's narrative, `wip/archive/agents/S-0743.md`, § Decisions.
