---
id: D-0638
title: flai serve takes its workflow from serve.Options.Workflow, nil for the baseline, and its passes ask flai/internal/workflow/review.go's AcceptFrom, AwaitsAcceptance, SendBack, Gives, Closed, and Underway where a story awaits acceptance, is sent back, or is closed
kind: decision
sources: [S-0743]
topics: [flai, go, code, kluster]
paths: [flai/internal/serve/serve.go, flai/internal/workflow/review.go, flai/internal/kluster/kluster.go]
date: 2026-11-15
status: absorbed
absorbed_by: "design/system/workflow.md#what-waits-on-an-acceptance-s-0296"
asked_again: []
cited: []
---

# D-0638 flai serve takes its workflow from serve.Options.Workflow, nil for the baseline, and its passes ask flai/internal/workflow/review.go's AcceptFrom, AwaitsAcceptance, SendBack, Gives, Closed, and Underway where a story awaits acceptance, is sent back, or is closed

S-0743 moved flai serve's kluster deploys, forge pass, and merge queue onto a state's kind and a transition's effect (ADR-0301). Its narrative's Decisions (`wip/archive/agents/S-0743.md`) record where serve gets its workflow and the helpers its three passes share:

> workitem.Open does not read a project's workflow file yet (that is another E-0034 story), so serve takes the workflow from Options.Workflow, nil for the baseline, as Repo.Workflow does. One shared set of methods in workflow/review.go, a new file, keeps the three from writing the same helpers three times and avoids S-0738's edits to workflow.go.

> T-3218 also touches flai/internal/kluster/kluster.go: Cluster gets the workflow there, so that DeployBy judges a story's state by its kind with no change to its callers in cmd and hostapi, which pass nil for the baseline until their own stories (S-0744, S-0747).

So a later story that needs, in flai serve, the state a story awaits acceptance in, its send-back transition, or whether a state is closed or underway, calls the methods on `workflow.Workflow` in `flai/internal/workflow/review.go` rather than writing its own, and reads the workflow from `Options.Workflow`. Until S-0755 makes flai read a project's own workflow file, that field is nil outside tests and serve runs the baseline (D-0627).

Source: S-0743's narrative, `wip/archive/agents/S-0743.md`, § Decisions.
